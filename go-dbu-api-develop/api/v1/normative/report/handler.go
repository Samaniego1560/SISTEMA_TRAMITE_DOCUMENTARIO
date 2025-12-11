package report

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

type ReportHandler struct {
	db *sqlx.DB
}

type StudentReportRow struct {
	DNI                     string     `db:"dni"`
	CodigoEstudiante        string     `db:"codigo_estudiante"`
	NombreApellido          string     `db:"nombre_apellido"`
	Sexo                    string     `db:"sexo"`
	EscuelaProfesional      string     `db:"escuela_profesional"`
	CorreoInstitucional     string     `db:"correo_institucional"`
	CelularEstudiante       string     `db:"celular_estudiante"`
	CelularPadre            string     `db:"celular_padre"`
	DepartamentoProcedencia string     `db:"departamento_procedencia"`
	ProvinciaProcedencia    string     `db:"provincia_procedencia"`
	DistritoProcedencia     string     `db:"distrito_procedencia"`
	Servicio                string     `db:"servicio"`
	FaltaCapitulo           string     `db:"falta_capitulo"`
	FaltaArticulo           string     `db:"falta_articulo"`
	GravedadArticulo        string     `db:"gravedad_articulo"`
	FaltaInciso             string     `db:"falta_inciso"`
	DetalleInciso           string     `db:"detalle_inciso"`
	Estado                  string     `db:"estado"`
	FechaFalta              *time.Time `db:"fecha_falta"`
	SancionCapitulo         string     `db:"sancion_capitulo"`
	SancionArticulo         string     `db:"sancion_articulo"`
	SancionInciso           string     `db:"sancion_inciso"`
	DetalleSancion          string     `db:"detalle_sancion"`
	FechaRegistro           *time.Time `db:"fecha_registro"`
	FechaInicio             *time.Time `db:"fecha_inicio"`
	FechaFin                *time.Time `db:"fecha_fin"`
	ApelacionMotivo         string     `db:"apelacion_motivo"`
	ApelacionObservaciones  string     `db:"apelacion_observaciones"`
	ApelacionVeredicto      string     `db:"apelacion_veredicto"`
}

func (h *ReportHandler) GetStudentReportExcel(c *fiber.Ctx) error {
	convocatoriaIDStr := c.Params("convocatoria_id")
	if convocatoriaIDStr == "" {
		return c.Status(400).SendString("convocatoria_id es requerido")
	}
	convocatoriaID, err := strconv.Atoi(convocatoriaIDStr)
	if err != nil {
		return c.Status(400).SendString("convocatoria_id inválido")
	}

	var convocatoriaNombre string
	err = h.db.Get(&convocatoriaNombre, "SELECT nombre FROM convocatorias WHERE id = ?", convocatoriaID)
	if err != nil || convocatoriaNombre == "" {
		convocatoriaNombre = "convocatoria"
	}
	safeConvocatoriaNombre := sanitizeFileName(convocatoriaNombre)

	//QUERY CORREGIDA: Usar ANY_VALUE() para columnas no agregadas o agregarlas al GROUP BY
	query := `
		SELECT
			al.dni,
			al.codigo_estudiante AS codigo_estudiante,
			CONCAT(al.nombres, ' ', al.apellido_paterno, ' ', al.apellido_materno) AS nombre_apellido,
			al.sexo AS sexo,
			al.escuela_profesional AS escuela_profesional,
			al.correo_institucional AS correo_institucional,
			-- Celulares agregados con MAX
			COALESCE(MAX(CASE WHEN req.nombre = 'celular de estudiante' THEN ds.respuesta_formulario END), '') AS celular_estudiante,
			COALESCE(MAX(CASE WHEN req.nombre = 'Celular padre' THEN ds.respuesta_formulario END), '') AS celular_padre,
			-- Departamento de procedencia
			COALESCE((
				SELECT d.name
				FROM solicitudes sol
				JOIN detalle_solicitudes ds2 ON ds2.solicitud_id = sol.id
				JOIN requisitos req2 ON req2.id = ds2.requisito_id
				JOIN departaments d ON d.id = ds2.opcion_seleccion
				WHERE sol.alumno_id = al.id AND req2.nombre = 'Departamento de procedencia'
				ORDER BY sol.created_at DESC
				LIMIT 1
			), '') AS departamento_procedencia,
			-- Provincia de procedencia
			COALESCE((
				SELECT p.name
				FROM solicitudes sol
				JOIN detalle_solicitudes ds2 ON ds2.solicitud_id = sol.id
				JOIN requisitos req2 ON req2.id = ds2.requisito_id
				JOIN provinces p ON p.id = ds2.opcion_seleccion
				WHERE sol.alumno_id = al.id AND req2.nombre = 'Provincia de procedencia'
				ORDER BY sol.created_at DESC
				LIMIT 1
			), '') AS provincia_procedencia,
			-- Distrito de procedencia
			COALESCE((
				SELECT dist.name
				FROM solicitudes sol
				JOIN detalle_solicitudes ds2 ON ds2.solicitud_id = sol.id
				JOIN requisitos req2 ON req2.id = ds2.requisito_id
				JOIN districts dist ON dist.id = ds2.opcion_seleccion
				WHERE sol.alumno_id = al.id AND req2.nombre = 'Distrito de procedencia'
				ORDER BY sol.created_at DESC
				LIMIT 1
			), '') AS distrito_procedencia,
			COALESCE(serv.nombre, '') AS servicio,
			COALESCE(cap.nombre, '') AS falta_capitulo,
			COALESCE(art.descripcion, '') AS falta_articulo,
			COALESCE(art.gravedad, '') AS gravedad_articulo,
			COALESCE(inc.nombre, '') AS falta_inciso,
			COALESCE(inc.descripcion, '') AS detalle_inciso,
			COALESCE(f.estado, '') AS estado,
			f.fecha_falta AS fecha_falta,
			COALESCE(saf.capitulo_sancion, '') AS sancion_capitulo,
			COALESCE(saf.articulo_sancion, '') AS sancion_articulo,
			COALESCE(saf.inciso_sancion, '') AS sancion_inciso,
			COALESCE(saf.detalle_sancion, '') AS detalle_sancion,
			MAX(sfa.fecha_asignacion) AS fecha_registro,
			MAX(sfa.fecha_inicio) AS fecha_inicio,
			MAX(sfa.fecha_fin) AS fecha_fin,
			COALESCE(MAX(ap.motivo), '') AS apelacion_motivo,
			COALESCE(MAX(ap.observaciones), '') AS apelacion_observaciones,
			COALESCE(MAX(ap.estado), '') AS apelacion_veredicto
		FROM alumnos al
		INNER JOIN faltas f ON f.alumno_id = al.id
		LEFT JOIN servicios serv ON serv.id = f.servicio_id
		LEFT JOIN faltas_articulos fa ON fa.falta_id = f.id
		LEFT JOIN articulos art ON art.id = fa.articulo_id
		LEFT JOIN capitulos cap ON cap.id = art.capitulo_id
		LEFT JOIN faltas_incisos fi ON fi.falta_id = f.id
		LEFT JOIN incisos inc ON inc.id = fi.inciso_id
		LEFT JOIN sancion_falta_asignada sfa ON sfa.falta_id = f.id
		LEFT JOIN sanciones_faltas_normativa saf ON saf.id = sfa.sancion_id
		LEFT JOIN apelaciones ap ON ap.sancion_falta_asignada_id = sfa.id
		LEFT JOIN solicitudes sol ON sol.alumno_id = al.id
		LEFT JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
		LEFT JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
		LEFT JOIN requisitos req ON req.id = ds.requisito_id
		WHERE f.convocatoria_id = ?
		GROUP BY
			al.id,
			al.dni,
			al.codigo_estudiante,
			al.nombres,
			al.apellido_paterno,
			al.apellido_materno,
			al.sexo,
			al.escuela_profesional,
			al.correo_institucional,
			f.id,
			f.fecha_falta,
			f.estado,
			serv.id,
			serv.nombre,
			cap.id,
			cap.nombre,
			art.id,
			art.descripcion,
			art.gravedad,
			inc.id,
			inc.nombre,
			inc.descripcion,
			saf.id,
			saf.capitulo_sancion,
			saf.articulo_sancion,
			saf.inciso_sancion,
			saf.detalle_sancion
		ORDER BY al.dni, f.fecha_falta DESC
	`

	var rows []StudentReportRow
	err = h.db.Select(&rows, query, convocatoriaID)
	if err != nil {
		return c.Status(500).SendString("Error al obtener datos: " + err.Error())
	}

	f := excelize.NewFile()
	sheet := "Reporte"
	idx, _ := f.NewSheet(sheet)
	f.SetActiveSheet(idx)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#305496"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	oddRowStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	evenRowStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	headers := []string{
		"DNI", "CODIGO_ESTUDIANTE", "NOMBRE Y APELLIDO", "SEXO", "ESCUELA_PROFESIONAL",
		"CORREO_INSTITUCIONAL", "CELULAR_ESTUDIANTE", "CELULAR_APODERADO",
		"DEPARTAMENTO", "PROVINCIA", "DISTRITO", "SERVICIO",
		"FALTA_CAPITULO", "FALTA_ARTICULO", "GRAVEDAD_ARTICULO", "FALTA_INCISO", "DETALLE_INCISO",
		"ESTADO", "FECHA_FALTA",
		"SANCION_CAPITULO", "SANCION_ARTICULO", "SANCION_INCISO", "DETALLE_SANCION",
		"FECHA_REGISTRO", "FECHA_INICIO", "FECHA_FIN",
		"APELACION_MOTIVO", "APELACION_OBSERVACIONES", "APELACION_VEREDICTO",
	}

	for i, h := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err == nil {
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
		}
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, col, col, 20)
	}

	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "",
		ActivePane:  "",
	})

	for rowIdx, s := range rows {
		values := []interface{}{
			s.DNI,
			s.CodigoEstudiante,
			s.NombreApellido,
			s.Sexo,
			s.EscuelaProfesional,
			s.CorreoInstitucional,
			s.CelularEstudiante,
			s.CelularPadre,
			s.DepartamentoProcedencia,
			s.ProvinciaProcedencia,
			s.DistritoProcedencia,
			s.Servicio,
			s.FaltaCapitulo,
			s.FaltaArticulo,
			s.GravedadArticulo,
			s.FaltaInciso,
			s.DetalleInciso,
			s.Estado,
			safeDate(s.FechaFalta),
			s.SancionCapitulo,
			s.SancionArticulo,
			s.SancionInciso,
			s.DetalleSancion,
			safeDate(s.FechaRegistro),
			safeDate(s.FechaInicio),
			safeDate(s.FechaFin),
			s.ApelacionMotivo,
			s.ApelacionObservaciones,
			s.ApelacionVeredicto,
		}
		rowStyle := evenRowStyle
		if rowIdx%2 == 0 {
			rowStyle = oddRowStyle
		}
		for colIdx, v := range values {
			cell, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			if err == nil {
				f.SetCellValue(sheet, cell, v)
				f.SetCellStyle(sheet, cell, cell, rowStyle)
			}
		}
	}

	if idx, err := f.GetSheetIndex("Sheet1"); err == nil && idx != -1 {
		f.DeleteSheet("Sheet1")
	}

	fileName := "reporte_sancionados_" + safeConvocatoriaNombre + ".xlsx"
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+fileName)
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	buf, err := f.WriteToBuffer()
	if err != nil {
		return c.Status(500).SendString("Error generando Excel: " + err.Error())
	}
	return c.SendStream(buf)
}

func sanitizeFileName(name string) string {
	safe := ""
	for _, r := range name {
		switch {
		case r == ' ':
			safe += "_"
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-':
			safe += string(r)
		}
	}
	if safe == "" {
		return "convocatoria"
	}
	return safe
}
