package low_code_medical_area

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

var months = map[string]string{
	"1":  "ENERO",
	"2":  "FEBRERO",
	"3":  "MARZO",
	"4":  "ABRIL",
	"5":  "MAYO",
	"6":  "JUNIO",
	"7":  "JULIO",
	"8":  "AGOSTO",
	"9":  "SETIEMBRE",
	"10": "OCTUBRE",
	"11": "NOVIEMBRE",
	"12": "DICIEMBRE",
}

var schools = []string{
	"AGRONOMIA",
	"ZOOTECNIA",
	"INGENIERIA EN INDUSTRIAS ALIMENTARIAS",
	"INGENIERIA FORESTAL",
	"INGENIERIA EN CONSERVACION DE SUELOS Y AGUA",
	"INGENIERIA AMBIENTAL",
	"INGENIERIA EN RECURSOS NATURALES RENOVABLES",
	"ECONOMIA",
	"ADMINISTRACION",
	"CONTABILIDAD",
	"INGENIERIA EN INFORMATICA Y SISTEMAS",
	"INGENIERIA MECANICA ELECTRICA",
	"INGENIERIA EN CIBERSEGURIDAD",
	"INGENIERIA CIVIL",
	"TURISMO Y HOTELERIA",
}

var nursingProceduresHeader = []string{
	"INYECT.",
	"C.F.V",
	"LAVADO DE OÍDO",
	"CIRUG. MENOR",
	"VENOCLISIS",
	"OTROS",
}

var nursingProceduresComparation = []string{
	"INYECTABLES",
	"CONTROL DE FUNCIONES VITALES",
	"LAVADO DE OIDO",
	"CIRUGIA MENOR",
	"VENOCLISIS",
	"OTROS",
}

var nursingProcedures = []string{
	"INYECTABLES",
	"CONTROL DE FUNCIONES VITALES",
	"LAVADO DE OIDO",
	"BAÑO RAYOS INFRARROJOS",
	"EXTRACCION DE UÑERO",
	"CIRUGIA MENOR",
	"VENOCLISIS",
	"OTROS",
}

var nursingStaffs = []string{
	"ADMINISTRATIVOS",
	"DOCENTES",
	"OTROS",
}

var nursingStaffsComparation = []string{
	"Administrativo",
	"Docente",
	"Otros",
}

var dentistryProcedures = []string{
	"PERIODONCIA PROFILAXIS",
	"OPERATORIA RESINA",
	"CIRUGÍA EXODONCIA",
	"DIAGNÓSTICO EXAMEN BUCAL",
	"CONSULTA DENTAL",
	"PREVENCIÓN FLUORIZACIÓN",
	"PREVENCIÓN: COLUTORIO CH2+CPC 0.12%",
	"PREVENCIÓN: IHO - TCA CEPILLADO",
	"RADIOGRAFÍA DENTAL",
}

var dentistryStaffs = []string{
	"ADMINISTRATIVOS",
	"DOCENTES",
}

var dentistryStaffsComparation = []string{
	"Administrativo",
	"Docente",
}

func getTotalNameMonthsByRange(monthStart, monthEnd string) []string {
	var totalMonths []string
	monthStartNumber, _ := strconv.Atoi(monthStart)
	monthEndNumber, _ := strconv.Atoi(monthEnd)
	for i := monthStartNumber; i <= monthEndNumber; i++ {
		month := months[strconv.Itoa(i)]
		totalMonths = append(totalMonths, month)
	}

	return totalMonths
}

func getTotalNumberMonthsByRange(monthStart, monthEnd string) []int {
	var numberMonths []int
	monthStartNumber, _ := strconv.Atoi(monthStart)
	monthEndNumber, _ := strconv.Atoi(monthEnd)
	for i := monthStartNumber; i <= monthEndNumber; i++ {
		numberMonths = append(numberMonths, i)
	}

	return numberMonths
}

func obtenerRangoTrimestre(trimestre string, year string) (string, string) {
	yearRange, _ := strconv.Atoi(year)

	var startDate, endDate time.Time
	switch trimestre {
	case "1,2,3":
		startDate = time.Date(yearRange, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(yearRange, 4, 0, 23, 59, 59, 0, time.UTC)
	case "4,5,6":
		startDate = time.Date(yearRange, 4, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(yearRange, 7, 0, 23, 59, 59, 0, time.UTC)
	case "7,8,9":
		startDate = time.Date(yearRange, 7, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(yearRange, 10, 0, 23, 59, 59, 0, time.UTC)
	case "10,11,12":
		startDate = time.Date(yearRange, 10, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(yearRange, 13, 0, 23, 59, 59, 0, time.UTC)
	}
	return startDate.Format("2006-01-02"), endDate.Format("2006-01-02")
}

func getRangeMonthsOfYear(year string) []struct {
	Start time.Time
	End   time.Time
} {

	yearRange, _ := strconv.Atoi(year)
	var rangos []struct {
		Start time.Time
		End   time.Time
	}

	mesInicial := time.January

	for i := 0; i < 12; i++ {
		mes := mesInicial + time.Month(i)
		inicio := time.Date(yearRange, mes, 1, 0, 0, 0, 0, time.UTC)
		fin := time.Date(yearRange, mes+1, 0, 23, 59, 59, 0, time.UTC)

		rangos = append(rangos, struct {
			Start time.Time
			End   time.Time
		}{Start: inicio, End: fin})
	}

	return rangos
}

func getRangeMonthsOfTrimestre(trimestre string, year string) []struct {
	Start time.Time
	End   time.Time
} {

	yearRange, _ := strconv.Atoi(year)
	var rangos []struct {
		Start time.Time
		End   time.Time
	}

	var mesInicial time.Month
	switch trimestre {
	case "1,2,3":
		mesInicial = time.January
	case "4,5,6":
		mesInicial = time.April
	case "7,8,9":
		mesInicial = time.July
	case "10,11,12":
		mesInicial = time.October
	}

	for i := 0; i < 3; i++ {
		mes := mesInicial + time.Month(i)
		inicio := time.Date(yearRange, mes, 1, 0, 0, 0, 0, time.UTC)
		fin := time.Date(yearRange, mes+1, 0, 23, 59, 59, 0, time.UTC)

		rangos = append(rangos, struct {
			Start time.Time
			End   time.Time
		}{Start: inicio, End: fin})
	}

	return rangos
}

func getRangeAllTrimestres(year string) []struct {
	Start time.Time
	End   time.Time
} {
	yearRange, _ := strconv.Atoi(year)
	rangos := make([]struct {
		Start time.Time
		End   time.Time
	}, 4)

	trimestres := []time.Month{
		time.January,
		time.April,
		time.July,
		time.October,
	}

	for i, mesInicial := range trimestres {
		inicio := time.Date(yearRange, mesInicial, 1, 0, 0, 0, 0, time.UTC)

		if i == 3 {
			fin := time.Date(yearRange+1, time.January, 0, 23, 59, 59, 0, time.UTC)
			rangos[i] = struct {
				Start time.Time
				End   time.Time
			}{Start: inicio, End: fin}
		} else {
			fin := time.Date(yearRange, trimestres[i+1], 0, 23, 59, 59, 0, time.UTC)
			rangos[i] = struct {
				Start time.Time
				End   time.Time
			}{Start: inicio, End: fin}
		}
	}

	return rangos
}

func styleExcel(f *excelize.File) (int, int) {
	styleHeaderID, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
			Size: 8,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#AFDF9F"},
			Pattern: 1,
		},
	})
	if err != nil {
		logger.Error.Println("Error creating style:", err)
	}

	styleBorderID, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		logger.Error.Println("Error creating style:", err)
	}

	return styleHeaderID, styleBorderID
}

func createDentistryExcelRow(consultation *models.ConsultationPatientsMedicalAreaExcel, row int) models.ExcelPageRow {
	if row == 1 {
		return models.ExcelPageRow{
			Row: 1,
			Columns: []models.ExcelPageColumn{
				{Column: "A", Value: "FECHA"},
				{Column: "B", Value: "DNI"},
				{Column: "C", Value: "CODIGO"},
				{Column: "D", Value: "APELLIDOS Y NOMBRES"},
				{Column: "E", Value: "SEXO"},
				{Column: "F", Value: "FECHA DE NACIMIENTO"},
				{Column: "G", Value: "EDAD"},
				{Column: "H", Value: "DOMICILIO"},
				{Column: "I", Value: "TELEFONO"},
				{Column: "J", Value: "ESCUELA PROFESIONAL"},
				{Column: "K", Value: "OCUPACION"},
				{Column: "L", Value: "STUD/ADMVO/CAT"},
				{Column: "M", Value: "PROCEDIMIENTO"},
				{Column: "N", Value: "Pieza dental"},
				{Column: "O", Value: "TIPO PROCEDIMIENTO"},
				{Column: "P", Value: "Recibo"},
				{Column: "Q", Value: "Costo"},
				{Column: "R", Value: "Fecha de pago"},
			},
		}
	}

	return models.ExcelPageRow{
		Row: row,
		Columns: []models.ExcelPageColumn{
			{Column: "A", Value: consultation.FechaConsulta},
			{Column: "B", Value: consultation.DNI},
			{Column: "C", Value: consultation.CodigoSga},
			{Column: "D", Value: consultation.NombreCompleto},
			{Column: "E", Value: consultation.Sexo},
			{Column: "F", Value: consultation.FechaNacimiento},
			{Column: "G", Value: consultation.Edad},
			{Column: "H", Value: consultation.Domicilio},
			{Column: "I", Value: consultation.NumeroCelular},
			{Column: "J", Value: consultation.EscuelaProfesional},
			{Column: "K", Value: consultation.Ocupacion},
			{Column: "L", Value: consultation.TipoPersona},
			{Column: "M", Value: consultation.Servicios},
			{Column: "N", Value: consultation.PiezaDental},
			{Column: "O", Value: consultation.TipoProcedimiento},
			{Column: "P", Value: consultation.Recibo},
			{Column: "Q", Value: consultation.Costo},
			{Column: "R", Value: consultation.FechaPago},
		},
	}
}
func createMedicalExcelRow(consultation *models.ConsultationPatientsMedicalAreaExcel, row int) models.ExcelPageRow {
	if row == 1 {
		return models.ExcelPageRow{
			Row: 1,
			Columns: []models.ExcelPageColumn{
				{Column: "A", Value: "Fecha consulta"},
				{Column: "B", Value: "Tipo de persona"},
				{Column: "C", Value: "DNI"},
				{Column: "D", Value: "Nombre Completo"},
				{Column: "E", Value: "Sexo"},
				{Column: "F", Value: "Fecha de nacimiento"},
				{Column: "G", Value: "Número de teléfono"},
				{Column: "H", Value: "Servicios"},
			},
		}
	}

	return models.ExcelPageRow{
		Row: row,
		Columns: []models.ExcelPageColumn{
			{Column: "A", Value: consultation.FechaConsulta},
			{Column: "B", Value: consultation.TipoPersona},
			{Column: "C", Value: consultation.DNI},
			{Column: "D", Value: consultation.NombreCompleto},
			{Column: "E", Value: consultation.Sexo},
			{Column: "F", Value: consultation.FechaNacimiento},
			{Column: "G", Value: consultation.NumeroCelular},
			{Column: "H", Value: consultation.Servicios},
		},
	}
}
