package student_parent

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"dbu-api/internal/middleware"
	"dbu-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

type handlerSurvey struct {
	db   *sqlx.DB
	txID string
}

const maxMonthlyAmount = 99999999.99

type totalsStats struct {
	Total       int             `db:"total" json:"total"`
	HasChildren int             `db:"has_children" json:"has_children"`
	NoChildren  int             `db:"no_children" json:"no_children"`
	Gestation   int             `db:"gestation" json:"gestation"`
	Pregnant    int             `db:"pregnant" json:"pregnant"`
	AvgChildren sql.NullFloat64 `db:"avg_children" json:"avg_children"`
}

type schoolStats struct {
	School      string          `db:"school" json:"school"`
	Total       int             `db:"total" json:"total"`
	HasChildren int             `db:"has_children" json:"has_children"`
	NoChildren  int             `db:"no_children" json:"no_children"`
	Gestation   int             `db:"gestation" json:"gestation"`
	AvgChildren sql.NullFloat64 `db:"avg_children" json:"avg_children"`
}

type childrenCountStat struct {
	ChildrenCount sql.NullInt64 `db:"children_count" json:"children_count"`
	Total         int           `db:"total" json:"total"`
}

type totalsResponse struct {
	Total       int      `json:"total"`
	HasChildren int      `json:"has_children"`
	NoChildren  int      `json:"no_children"`
	Gestation   int      `json:"gestation"`
	Pregnant    int      `json:"pregnant"`
	AvgChildren *float64 `json:"avg_children"`
}

type schoolResponse struct {
	School      string   `json:"school"`
	Total       int      `json:"total"`
	HasChildren int      `json:"has_children"`
	NoChildren  int      `json:"no_children"`
	Gestation   int      `json:"gestation"`
	AvgChildren *float64 `json:"avg_children"`
}

type childrenCountResponse struct {
	ChildrenCount *int `json:"children_count"`
	Total         int  `json:"total"`
}

type verifyRequest struct {
	DNI string `json:"dni"`
}

type surveyRequest struct {
	DNI                    string   `json:"dni"`
	StudentCode            *string  `json:"student_code"`
	School                 *string  `json:"school"`
	AcademicCycle          *string  `json:"academic_cycle"`
	Age                    *int     `json:"age"`
	Gender                 *string  `json:"gender"`
	IsPregnant             *bool    `json:"is_pregnant"`
	HasChildren            *bool    `json:"has_children"`
	ChildrenCount          *int     `json:"children_count"`
	GestationSelected      *bool    `json:"gestation_selected"`
	Child1Age              *int     `json:"child_1_age"`
	Child2Age              *int     `json:"child_2_age"`
	Child3Age              *int     `json:"child_3_age"`
	Child4Age              *int     `json:"child_4_age"`
	ResidenceType          *string  `json:"residence_type"`
	ResidenceOther         *string  `json:"residence_other"`
	LivesWithPartner       *bool    `json:"lives_with_partner"`
	LivesWithChildren      *bool    `json:"lives_with_children"`
	CaregiverType          *string  `json:"caregiver_type"`
	CaregiverOther         *string  `json:"caregiver_other"`
	ParticipatesUpbringing *bool    `json:"participates_in_upbringing"`
	EconomicSupport        *bool    `json:"economic_support"`
	MonthlyAmount          *float64 `json:"monthly_amount"`
	HasAcademicSupport     *bool    `json:"has_academic_support"`
	SupportProvider        *string  `json:"support_provider"`
	SupportProviderOther   *string  `json:"support_provider_other"`
	InterruptedSemester    *bool    `json:"interrupted_semester"`
}

type academicResponse struct {
	Msg     string          `json:"msg"`
	Detalle json.RawMessage `json:"detalle"`
}

type academicDetail struct {
	CodSem     string `json:"codsem"`
	CodAlumno  int64  `json:"codalumno"`
	TDocumento string `json:"tdocumento"`
	Appaterno  string `json:"appaterno"`
	Apmaterno  string `json:"apmaterno"`
	Nombre     string `json:"nombre"`
	Sexo       string `json:"sexo"`
	FecNac     string `json:"fecnac"`
	NomEsp     string `json:"nomesp"`
	NumeSemCur int    `json:"nume_sem_cur"`
}

var allowedSchools = map[string]struct{}{
	"AGRONOMIA":                             {},
	"ZOOTECNIA":                             {},
	"INGENIERIA EN INDUSTRIAS ALIMENTARIAS": {},
	"INGENIERIA FORESTAL":                   {},
	"INGENIERIA EN CONSERVACION DE SUELOS Y AGUA": {},
	"ADMINISTRACION":                       {},
	"CONTABILIDAD":                         {},
	"ECONOMIA":                             {},
	"INGENIERIA EN INFORMATICA Y SISTEMAS": {},
	"INGENIERIA AMBIENTAL":                 {},
	"INGENIERIA EN RECURSOS NATURALES RENOVABLES": {},
	"INGENIERIA MECANICA ELECTRICA":               {},
	"INGENIERIA CIVIL":                            {},
	"TURISMO Y HOTELERIA":                         {},
	"INGENIERIA EN CIBERSEGURIDAD":                {},
}

func (h *handlerSurvey) Verify(c *fiber.Ctx) error {
	var req verifyRequest
	if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.DNI) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "dni es requerido",
		})
	}

	dni := strings.TrimSpace(req.DNI)
	exists, err := h.existsByDNI(dni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error interno"})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"alreadyAnswered": true})
	}

	student, err := fetchAcademicData(dni)
	if err != nil {
		if errors.Is(err, errNotEnrolled) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"exists":  false,
				"message": "Usted no se encuentra matriculado actualmente",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "error consultando datos académicos"})
	}

	age, err := calculateAge(student.FecNac)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "datos académicos incompletos"})
	}

	sexo := strings.ToUpper(student.Sexo)
	if sexo != "F" && sexo != "M" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "datos académicos incompletos"})
	}

	return c.JSON(fiber.Map{
		"exists":         true,
		"student_code":   fmt.Sprintf("%d", student.CodAlumno),
		"dni":            dni,
		"appaterno":      student.Appaterno,
		"apmaterno":      student.Apmaterno,
		"nombre":         student.Nombre,
		"school":         student.NomEsp,
		"academic_cycle": fmt.Sprintf("%d", student.NumeSemCur),
		"age":            age,
		"gender":         sexo,
	})
}

func (h *handlerSurvey) Store(c *fiber.Ctx) error {
	var req surveyRequest
	if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.DNI) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "dni es requerido"})
	}

	if req.IsPregnant == nil || req.HasChildren == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "is_pregnant y has_children son obligatorios"})
	}

	dni := strings.TrimSpace(req.DNI)

	exists, err := h.existsByDNI(dni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error interno"})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"alreadyAnswered": true})
	}

	student, err := fetchAcademicData(dni)
	if err != nil {
		if errors.Is(err, errNotEnrolled) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"exists":  false,
				"message": "Usted no se encuentra matriculado actualmente",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "error consultando datos académicos"})
	}

	age, err := calculateAge(student.FecNac)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "datos académicos incompletos"})
	}

	sexo := strings.ToUpper(student.Sexo)
	if sexo != "F" && sexo != "M" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "datos académicos incompletos"})
	}

	studentCode := fmt.Sprintf("%d", student.CodAlumno)
	if req.StudentCode != nil && strings.TrimSpace(*req.StudentCode) != "" {
		studentCode = strings.TrimSpace(*req.StudentCode)
	}

	school := strings.ToUpper(strings.TrimSpace(student.NomEsp))
	if req.School != nil && strings.TrimSpace(*req.School) != "" {
		school = strings.ToUpper(strings.TrimSpace(*req.School))
	}
	if _, ok := allowedSchools[school]; !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "school inválido"})
	}

	academicCycle := fmt.Sprintf("%d", student.NumeSemCur)
	if req.AcademicCycle != nil && strings.TrimSpace(*req.AcademicCycle) != "" {
		academicCycle = strings.TrimSpace(*req.AcademicCycle)
	}

	if req.Age != nil {
		age = *req.Age
	}
	if age < 15 || age > 90 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "age fuera de rango"})
	}

	if req.Gender != nil && strings.TrimSpace(*req.Gender) != "" {
		sexo = strings.ToUpper(strings.TrimSpace(*req.Gender))
	}
	if sexo != "F" && sexo != "M" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "gender inválido"})
	}

	payload := map[string]interface{}{
		"student_code":   studentCode,
		"dni":            dni,
		"appaterno":      student.Appaterno,
		"apmaterno":      student.Apmaterno,
		"nombre":         student.Nombre,
		"school":         school,
		"academic_cycle": academicCycle,
		"age":            age,
		"gender":         sexo,
		"is_pregnant":    *req.IsPregnant,
		"has_children":   *req.HasChildren,
		"gestation":      req.GestationSelected != nil && *req.GestationSelected,
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
	}

	if !*req.IsPregnant && !*req.HasChildren {
		h.nullifySections(payload)
		if err := h.insertSurvey(payload); err != nil {
			if errors.Is(err, errDuplicate) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"alreadyAnswered": true})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al guardar"})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"ok": true, "partial": true})
	}

	if req.ChildrenCount == nil {
		if req.GestationSelected == nil || !*req.GestationSelected {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "children_count es obligatorio"})
		}
	}
	if req.ChildrenCount != nil && *req.ChildrenCount < 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "children_count inválido"})
	}
	if (req.GestationSelected != nil && *req.GestationSelected) && !*req.IsPregnant {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "gestation_selected requiere is_pregnant = true"})
	}
	if req.ChildrenCount != nil && *req.ChildrenCount == 0 && !(req.GestationSelected != nil && *req.GestationSelected) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "children_count debe ser >= 1"})
	}

	if req.ChildrenCount == nil {
		payload["children_count"] = 0
	} else {
		payload["children_count"] = *req.ChildrenCount
	}
	if payload["children_count"].(int) == 0 {
		payload["child_1_age"] = 0
		payload["child_2_age"] = nil
		payload["child_3_age"] = nil
		payload["child_4_age"] = nil
	} else {
		if req.Child1Age == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "child_1_age es obligatorio"})
		}
		payload["child_1_age"] = *req.Child1Age
		if payload["children_count"].(int) >= 2 {
			if req.Child2Age == nil {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "child_2_age es obligatorio"})
			}
			payload["child_2_age"] = *req.Child2Age
		} else {
			payload["child_2_age"] = nil
		}
		if payload["children_count"].(int) >= 3 {
			if req.Child3Age == nil {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "child_3_age es obligatorio"})
			}
			payload["child_3_age"] = *req.Child3Age
		} else {
			payload["child_3_age"] = nil
		}
		if payload["children_count"].(int) >= 4 {
			if req.Child4Age == nil {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "child_4_age es obligatorio"})
			}
			payload["child_4_age"] = *req.Child4Age
		} else {
			payload["child_4_age"] = nil
		}
	}

	required := []struct {
		name  string
		value interface{}
	}{
		{"residence_type", req.ResidenceType},
		{"lives_with_partner", req.LivesWithPartner},
		{"lives_with_children", req.LivesWithChildren},
		{"participates_in_upbringing", req.ParticipatesUpbringing},
		{"economic_support", req.EconomicSupport},
		{"has_academic_support", req.HasAcademicSupport},
		{"interrupted_semester", req.InterruptedSemester},
	}
	for _, r := range required {
		if r.value == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": r.name + " es obligatorio"})
		}
	}

	payload["residence_type"] = *req.ResidenceType
	if *req.ResidenceType == "otro" {
		if req.ResidenceOther == nil || strings.TrimSpace(*req.ResidenceOther) == "" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "residence_other es obligatorio"})
		}
		payload["residence_other"] = strings.TrimSpace(*req.ResidenceOther)
	} else {
		payload["residence_other"] = nil
	}

	payload["lives_with_partner"] = *req.LivesWithPartner
	payload["lives_with_children"] = *req.LivesWithChildren

	if req.LivesWithChildren != nil && !*req.LivesWithChildren {
		if req.CaregiverType == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "caregiver_type es obligatorio"})
		}
		payload["caregiver_type"] = *req.CaregiverType
		if *req.CaregiverType == "otro" {
			if req.CaregiverOther == nil || strings.TrimSpace(*req.CaregiverOther) == "" {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "caregiver_other es obligatorio"})
			}
			payload["caregiver_other"] = strings.TrimSpace(*req.CaregiverOther)
		} else {
			payload["caregiver_other"] = nil
		}
	} else {
		payload["caregiver_type"] = nil
		payload["caregiver_other"] = nil
	}

	payload["participates_in_upbringing"] = *req.ParticipatesUpbringing
	payload["economic_support"] = *req.EconomicSupport
	if *req.EconomicSupport {
		if req.MonthlyAmount == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "monthly_amount es obligatorio"})
		}
		if math.IsNaN(*req.MonthlyAmount) || *req.MonthlyAmount < 0 || *req.MonthlyAmount > maxMonthlyAmount {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "monthly_amount fuera de rango"})
		}
		payload["monthly_amount"] = *req.MonthlyAmount
	} else {
		payload["monthly_amount"] = nil
	}

	payload["has_academic_support"] = *req.HasAcademicSupport
	if *req.HasAcademicSupport {
		if req.SupportProvider == nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "support_provider es obligatorio"})
		}
		payload["support_provider"] = *req.SupportProvider
		if *req.SupportProvider == "otro" {
			if req.SupportProviderOther == nil || strings.TrimSpace(*req.SupportProviderOther) == "" {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "support_provider_other es obligatorio"})
			}
			payload["support_provider_other"] = strings.TrimSpace(*req.SupportProviderOther)
		} else {
			payload["support_provider_other"] = nil
		}
	} else {
		payload["support_provider"] = nil
		payload["support_provider_other"] = nil
	}

	payload["interrupted_semester"] = *req.InterruptedSemester

	if err := h.insertSurvey(payload); err != nil {
		if errors.Is(err, errDuplicate) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"alreadyAnswered": true})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al guardar"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"ok": true})
}

func (h *handlerSurvey) Export(c *fiber.Ctx) error {
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{Error: true, Msg: "No autorizado"})
	}
	if user.IDLevelUser != 1 {
		return c.Status(fiber.StatusForbidden).JSON(models.Response{Error: true, Msg: "No autorizado"})
	}

	rows, err := h.fetchAllRows()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error al obtener datos: " + err.Error())
	}

	f := excelize.NewFile()
	sheet := "Encuesta"
	idx, _ := f.NewSheet(sheet)
	f.SetActiveSheet(idx)

	headers := []string{
		"CODIGO_ESTUDIANTE", "DNI", "AP_PATERNO", "AP_MATERNO", "NOMBRES", "NOMBRE_COMPLETO",
		"ESCUELA_PROFESIONAL", "CICLO_ACADEMICO", "EDAD", "SEXO", "GESTACION",
		"ESTA_EMBARAZADA_O_PAREJA", "TIENE_HIJOS", "NUM_HIJOS_MENORES_5",
		"EDAD_HIJO_1", "EDAD_HIJO_2", "EDAD_HIJO_3", "EDAD_HIJO_4",
		"TIPO_RESIDENCIA", "RESIDENCIA_OTRO", "VIVE_CON_PAREJA", "VIVE_CON_HIJOS",
		"QUIEN_CUIDA", "QUIEN_CUIDA_OTRO", "PARTICIPA_EN_CRIANZA",
		"APOYO_ECONOMICO", "MONTO_MENSUAL", "APOYO_ACADEMICO", "QUIEN_APOYA",
		"QUIEN_APOYA_OTRO", "INTERRUMPIO_SEMESTRE", "CREADO_EN",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for rowIdx, r := range rows {
		fullName := strings.TrimSpace(strings.Join([]string{r.Appaterno, r.Apmaterno, r.Nombre}, " "))
		values := []interface{}{
			r.StudentCode, r.DNI, r.Appaterno, r.Apmaterno, r.Nombre, fullName,
			r.School, r.AcademicCycle, r.Age, r.Gender, nullBoolToString(r.Gestation),
			boolToInt(r.IsPregnant), boolToInt(r.HasChildren), nullIntToString(r.ChildrenCount),
			nullIntToString(r.Child1Age), nullIntToString(r.Child2Age), nullIntToString(r.Child3Age), nullIntToString(r.Child4Age),
			nullString(r.ResidenceType), nullString(r.ResidenceOther),
			nullBoolToString(r.LivesWithPartner), nullBoolToString(r.LivesWithChildren),
			nullString(r.CaregiverType), nullString(r.CaregiverOther),
			nullBoolToString(r.ParticipatesInUpbringing),
			nullBoolToString(r.EconomicSupport), nullFloatToString(r.MonthlyAmount),
			nullBoolToString(r.HasAcademicSupport), nullString(r.SupportProvider),
			nullString(r.SupportProviderOther), nullBoolToString(r.InterruptedSemester),
			r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if idx, err := f.GetSheetIndex("Sheet1"); err == nil && idx != -1 {
		f.DeleteSheet("Sheet1")
	}

	fileName := "encuesta_estudiantes.xlsx"
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+fileName)
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	buf, err := f.WriteToBuffer()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generando Excel: " + err.Error())
	}
	return c.SendStream(buf)
}

func (h *handlerSurvey) Stats(c *fiber.Ctx) error {
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{Error: true, Msg: "No autorizado"})
	}
	if user.IDLevelUser != 1 {
		return c.Status(fiber.StatusForbidden).JSON(models.Response{Error: true, Msg: "No autorizado"})
	}

	var totals totalsStats
	totalsQuery := `
		SELECT
			COUNT(*) AS total,
			SUM(CASE WHEN has_children = 1 THEN 1 ELSE 0 END) AS has_children,
			SUM(CASE WHEN has_children = 0 THEN 1 ELSE 0 END) AS no_children,
			SUM(CASE WHEN gestation = 1 THEN 1 ELSE 0 END) AS gestation,
			SUM(CASE WHEN is_pregnant = 1 THEN 1 ELSE 0 END) AS pregnant,
			AVG(children_count) AS avg_children
		FROM student_parent_survey`
	if err := h.db.Get(&totals, totalsQuery); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error obteniendo totales"})
	}

	var bySchool []schoolStats
	bySchoolQuery := `
		SELECT
			COALESCE(NULLIF(TRIM(school), ''), 'SIN ESCUELA') AS school,
			COUNT(*) AS total,
			SUM(CASE WHEN has_children = 1 THEN 1 ELSE 0 END) AS has_children,
			SUM(CASE WHEN has_children = 0 THEN 1 ELSE 0 END) AS no_children,
			SUM(CASE WHEN gestation = 1 THEN 1 ELSE 0 END) AS gestation,
			AVG(children_count) AS avg_children
		FROM student_parent_survey
		GROUP BY COALESCE(NULLIF(TRIM(school), ''), 'SIN ESCUELA')
		ORDER BY total DESC`
	if err := h.db.Select(&bySchool, bySchoolQuery); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error obteniendo datos por escuela"})
	}

	var childrenCounts []childrenCountStat
	countQuery := `
		SELECT
			children_count,
			COUNT(*) AS total
		FROM student_parent_survey
		WHERE children_count IS NOT NULL
		GROUP BY children_count
		ORDER BY children_count`
	if err := h.db.Select(&childrenCounts, countQuery); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error obteniendo distribución de hijos"})
	}

	totalsResp := totalsResponse{
		Total:       totals.Total,
		HasChildren: totals.HasChildren,
		NoChildren:  totals.NoChildren,
		Gestation:   totals.Gestation,
		Pregnant:    totals.Pregnant,
		AvgChildren: nullFloatToPtr(totals.AvgChildren),
	}

	bySchoolResp := make([]schoolResponse, 0, len(bySchool))
	for _, s := range bySchool {
		bySchoolResp = append(bySchoolResp, schoolResponse{
			School:      s.School,
			Total:       s.Total,
			HasChildren: s.HasChildren,
			NoChildren:  s.NoChildren,
			Gestation:   s.Gestation,
			AvgChildren: nullFloatToPtr(s.AvgChildren),
		})
	}

	childrenCountResp := make([]childrenCountResponse, 0, len(childrenCounts))
	for _, ccount := range childrenCounts {
		childrenCountResp = append(childrenCountResp, childrenCountResponse{
			ChildrenCount: nullIntToPtr(ccount.ChildrenCount),
			Total:         ccount.Total,
		})
	}

	return c.JSON(fiber.Map{
		"totals":          totalsResp,
		"by_school":       bySchoolResp,
		"children_counts": childrenCountResp,
	})
}

type surveyRow struct {
	StudentCode              string          `db:"student_code"`
	DNI                      string          `db:"dni"`
	Appaterno                string          `db:"appaterno"`
	Apmaterno                string          `db:"apmaterno"`
	Nombre                   string          `db:"nombre"`
	School                   string          `db:"school"`
	AcademicCycle            string          `db:"academic_cycle"`
	Age                      int             `db:"age"`
	Gender                   string          `db:"gender"`
	Gestation                sql.NullBool    `db:"gestation"`
	IsPregnant               bool            `db:"is_pregnant"`
	HasChildren              bool            `db:"has_children"`
	ChildrenCount            sql.NullInt64   `db:"children_count"`
	Child1Age                sql.NullInt64   `db:"child_1_age"`
	Child2Age                sql.NullInt64   `db:"child_2_age"`
	Child3Age                sql.NullInt64   `db:"child_3_age"`
	Child4Age                sql.NullInt64   `db:"child_4_age"`
	ResidenceType            sql.NullString  `db:"residence_type"`
	ResidenceOther           sql.NullString  `db:"residence_other"`
	LivesWithPartner         sql.NullBool    `db:"lives_with_partner"`
	LivesWithChildren        sql.NullBool    `db:"lives_with_children"`
	CaregiverType            sql.NullString  `db:"caregiver_type"`
	CaregiverOther           sql.NullString  `db:"caregiver_other"`
	ParticipatesInUpbringing sql.NullBool    `db:"participates_in_upbringing"`
	EconomicSupport          sql.NullBool    `db:"economic_support"`
	MonthlyAmount            sql.NullFloat64 `db:"monthly_amount"`
	HasAcademicSupport       sql.NullBool    `db:"has_academic_support"`
	SupportProvider          sql.NullString  `db:"support_provider"`
	SupportProviderOther     sql.NullString  `db:"support_provider_other"`
	InterruptedSemester      sql.NullBool    `db:"interrupted_semester"`
	CreatedAt                time.Time       `db:"created_at"`
}

func (h *handlerSurvey) fetchAllRows() ([]surveyRow, error) {
	var rows []surveyRow
	query := `
		SELECT
			student_code, dni, appaterno, apmaterno, nombre, school, academic_cycle, age, gender, gestation,
			is_pregnant, has_children, children_count,
			child_1_age, child_2_age, child_3_age, child_4_age,
			residence_type, residence_other, lives_with_partner, lives_with_children,
			caregiver_type, caregiver_other, participates_in_upbringing,
			economic_support, monthly_amount, has_academic_support, support_provider,
			support_provider_other, interrupted_semester, created_at
		FROM student_parent_survey
		ORDER BY created_at DESC`
	if err := h.db.Select(&rows, query); err != nil {
		return nil, err
	}
	return rows, nil
}

func (h *handlerSurvey) existsByDNI(dni string) (bool, error) {
	var exists int
	err := h.db.Get(&exists, "SELECT 1 FROM student_parent_survey WHERE dni = ? LIMIT 1", dni)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

var errNotEnrolled = errors.New("not enrolled")
var errDuplicate = errors.New("duplicate")

func fetchAcademicData(dni string) (*academicDetail, error) {
	url := "https://bienestar.unas.edu.pe/backend/DatosAlumnoAcademico/show/" + dni
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var apiResp academicResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if len(apiResp.Detalle) == 0 || strings.TrimSpace(string(apiResp.Detalle)) == "[]" {
		return nil, errNotEnrolled
	}

	var detail academicDetail
	if err := json.Unmarshal(apiResp.Detalle, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

func calculateAge(dateStr string) (int, error) {
	if dateStr == "" {
		return 0, errors.New("empty")
	}
	dob, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, err
	}
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age, nil
}

func (h *handlerSurvey) insertSurvey(payload map[string]interface{}) error {
	query := `
		INSERT INTO student_parent_survey (
			student_code, dni, appaterno, apmaterno, nombre, school, academic_cycle, age, gender, gestation,
			is_pregnant, has_children, children_count, child_1_age, child_2_age, child_3_age, child_4_age,
			residence_type, residence_other, lives_with_partner, lives_with_children,
			caregiver_type, caregiver_other, participates_in_upbringing, economic_support, monthly_amount,
			has_academic_support, support_provider, support_provider_other, interrupted_semester,
			created_at, updated_at
		) VALUES (
			:student_code, :dni, :appaterno, :apmaterno, :nombre, :school, :academic_cycle, :age, :gender, :gestation,
			:is_pregnant, :has_children, :children_count, :child_1_age, :child_2_age, :child_3_age, :child_4_age,
			:residence_type, :residence_other, :lives_with_partner, :lives_with_children,
			:caregiver_type, :caregiver_other, :participates_in_upbringing, :economic_support, :monthly_amount,
			:has_academic_support, :support_provider, :support_provider_other, :interrupted_semester,
			:created_at, :updated_at
		)`
	_, err := h.db.NamedExec(query, payload)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return errDuplicate
		}
		return err
	}
	return nil
}

func (h *handlerSurvey) nullifySections(payload map[string]interface{}) {
	payload["children_count"] = nil
	payload["child_1_age"] = nil
	payload["child_2_age"] = nil
	payload["child_3_age"] = nil
	payload["child_4_age"] = nil
	payload["residence_type"] = nil
	payload["residence_other"] = nil
	payload["lives_with_partner"] = nil
	payload["lives_with_children"] = nil
	payload["caregiver_type"] = nil
	payload["caregiver_other"] = nil
	payload["participates_in_upbringing"] = nil
	payload["economic_support"] = nil
	payload["monthly_amount"] = nil
	payload["has_academic_support"] = nil
	payload["support_provider"] = nil
	payload["support_provider_other"] = nil
	payload["interrupted_semester"] = nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func nullString(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func nullIntToString(v sql.NullInt64) string {
	if v.Valid {
		return fmt.Sprintf("%d", v.Int64)
	}
	return ""
}

func nullBoolToString(v sql.NullBool) string {
	if v.Valid {
		if v.Bool {
			return "1"
		}
		return "0"
	}
	return ""
}

func nullFloatToString(v sql.NullFloat64) string {
	if v.Valid {
		return fmt.Sprintf("%.2f", v.Float64)
	}
	return ""
}

func nullFloatToPtr(v sql.NullFloat64) *float64 {
	if v.Valid {
		val := v.Float64
		return &val
	}
	return nil
}

func nullIntToPtr(v sql.NullInt64) *int {
	if v.Valid {
		val := int(v.Int64)
		return &val
	}
	return nil
}
