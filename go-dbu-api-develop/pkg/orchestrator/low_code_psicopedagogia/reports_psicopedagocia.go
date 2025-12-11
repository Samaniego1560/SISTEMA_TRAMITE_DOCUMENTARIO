package low_code_psicopedagogia

import (
	"bytes"
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"dbu-api/pkg/psicopedagogia"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

type ReportsPsicopedagogiaService struct {
	db   *sqlx.DB
	usr  *models.User
	txID string
}

type PortsPsicopedagogia interface {
	GetReportPsicopedagogiaLowCode(numeroCuadro, month, year, startDate, endDate string) (string, int, error)
}

func NewReportsPsicopedagogia(db *sqlx.DB, usr *models.User, txID string) PortsPsicopedagogia {
	return &ReportsPsicopedagogiaService{db: db, usr: usr, txID: txID}
}

var teachersAndExterns = []string{
	"DOCENTES",
	"NO DOCENTES",
	"EXTERNOS",
}

var teachersAndExternsComparation = []string{
	"DOCENTES",
	"NO DOCENTES",
	"EXTERNOS",
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

var rangeDateTrimester = map[string][2]string{
	"1,2,3":    {"01-01", "03-31"},
	"4,5,6":    {"04-01", "06-30"},
	"7,8,9":    {"07-01", "09-30"},
	"10,11,12": {"10-01", "12-31"},
}

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

type HeaderAttentionExcel struct {
	Frame string `json:"frame"`
	Title string `json:"title"`
	Area  string `json:"area"`
}

type MergeRange struct {
	Column1 string
	Column2 string
}
type DateRange struct {
	Start time.Time
	End   time.Time
}

func (s *ReportsPsicopedagogiaService) GetReportPsicopedagogiaLowCode(numeroCuadro, month, year, startDate, endDate string) (string, int, error) {

	if numeroCuadro == "" {
		return "", 1, errors.New("missing numeroCuadro Psicopedagogia")
	}

	if month == "" {
		month = "1,2,3"
	}

	if year == "" {
		year = time.Now().Format("2006")
	}

	if numeroCuadro == "1" {
		base64Str, code, err := s.GetDataAttentionsStudents(year, month)
		if err != nil {
			logger.Error.Println("Error generating report medical frame 3:", err)
			return "", code, err
		}
		return base64Str, code, err
	} else if numeroCuadro == "2" {
		base64Str, code, err := s.GetDataAttentionsTeachers(year, month)
		if err != nil {
			logger.Error.Println("Error generating report medical frame 3:", err)
			return "", code, err
		}
		return base64Str, code, err
	} else if numeroCuadro == "3" {
		base64Str, code, err := s.GetReportByDateRange(startDate, endDate)
		if err != nil {
			logger.Error.Println("Error generating report medical frame 3:", err)
			return "", code, err
		}
		return base64Str, code, err
	}

	return "", 1, errors.New("invalid numeroCuadro")
}

func (s *ReportsPsicopedagogiaService) GetDataAttentionsStudents(year, month string) (string, int, error) {

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	styleHeaderID, styleBorderID := styleExcel(f)

	trimester := getTrimester(month)

	headerAttention := HeaderAttentionExcel{
		Frame: "CUADRO",
		Title: fmt.Sprintf("ATENCIONES PSICOLOGICAS A ESTUDIANTES DEL %s, - AÑO %s", trimester, year),
		Area:  "ATENCIONES",
	}

	f = HeaderNursingConsultingTrimesterExcel(f, sheet, headerAttention, styleHeaderID, month, trimester, year)

	for i, school := range schools {
		cell := fmt.Sprintf("A%d", 6+i)
		cellEnd := fmt.Sprintf("U%d", 6+i)
		_ = f.SetCellValue(sheet, cell, school)
		_ = f.SetCellStyle(sheet, cell, cellEnd, styleBorderID)
	}

	_ = f.SetColWidth(sheet, "A", "A", 50)

	dateStart, dateEnd := getDateRange(month, year)

	srv := psicopedagogia.NewServerpsicopedagogia(s.db, s.usr, s.txID)
	dataAttentions, err := srv.SrvReportes.GetReportAttentionsDataStudents(dateStart, dateEnd)

	if err != nil {
		logger.Error.Println("", " error al consultar las atenciones:", err)
		return "", 15, err
	}

	startRow := 6
	totalRow := 5

	startTime, _ := time.Parse("2006-01-02", dateStart)
	endTime, _ := time.Parse("2006-01-02", dateEnd)
	nextStart := startTime.AddDate(0, 1, 0)
	prevEnd := endTime.AddDate(0, -1, 0)

	dateRanges := []DateRange{
		{Start: startTime, End: nextStart},
		{Start: nextStart, End: prevEnd},
		{Start: prevEnd, End: endTime},
	}

	columnTotalTrimester := []string{"B", "C", "D", "E", "F"}
	columnsTotalMonthStart := []string{"G", "L", "Q"}

	shouldCount := func(attn *models.ConsultationAttentionExcel, school string, sexo string) bool {
		if attn.TipoPersona != "Estudiante" {
			return false
		}
		if attn.EscuelaProfesional != school {
			return false
		}
		if sexo == "" {
			return true
		}
		return attn.Sexo == sexo
	}
	for i, sexo := range []string{"", "M", "F", "E.S", "C.P"} {
		totalCount := 0

		for j, school := range schools {
			row := startRow + j
			count := 0

			for _, attn := range dataAttentions {
				if shouldCount(attn, school, sexo) {
					count++
					totalCount++
				}
			}
			_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", columnTotalTrimester[i], row), count)
		}
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", columnTotalTrimester[i], totalRow), totalCount)
	}

	for monthIdx, dateRange := range dateRanges {
		for i, sexo := range []string{"", "M", "F", "E.S", "C.P"} {
			colNumStart, _ := excelize.ColumnNameToNumber(columnsTotalMonthStart[monthIdx])
			col, _ := excelize.ColumnNumberToName(colNumStart + i)
			totalCount := 0

			for j, school := range schools {
				row := startRow + j
				count := 0

				for _, attn := range dataAttentions {
					date, err := time.Parse("2006-01-02", attn.FechaConsulta[:10])
					if err != nil {
						logger.Error.Println("", " - couldn't parse date:", err)
						return "", 15, err
					}

					if (date.Equal(dateRange.Start) || date.After(dateRange.Start)) && date.Before(dateRange.End) {
						if shouldCount(attn, school, sexo) {
							count++
							totalCount++
						}
					}
				}
				_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), count)
			}
			_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, totalRow), totalCount)
		}
	}

	base64Str, code, err := SaveExcelAndReturnBase64(f, "Atenciones Psicologicas")

	return base64Str, code, err

}

func HeaderTeachersConsultingTrimesterExcel(f *excelize.File, sheet string, headerMedicalArea HeaderAttentionExcel, styleHeaderID int, month, trimester, year string) *excelize.File {
	var mergeColumns = []MergeRange{
		{Column1: "A1", Column2: "Q1"},
		{Column1: "A2", Column2: "Q2"},
		{Column1: "A3", Column2: "A3"},
		{Column1: "B3", Column2: "E3"},
		{Column1: "F3", Column2: "I3"},
		{Column1: "J3", Column2: "M3"},
		{Column1: "N3", Column2: "Q3"},
	}

	for _, merge := range mergeColumns {
		_ = f.MergeCell(sheet, merge.Column1, merge.Column2)
		_ = f.SetCellStyle(sheet, merge.Column1, merge.Column2, styleHeaderID)
	}

	cellValues := map[string]string{
		"A1": headerMedicalArea.Frame,
		"A2": headerMedicalArea.Title,
		"A3": headerMedicalArea.Area,
		"A4": "",
		"B3": fmt.Sprintf("TOTAL %s %s", trimester, year),
	}

	cols := []string{"G3", "L3", "Q3"}
	monthCodes := strings.Split(month, ",")
	for i, code := range monthCodes {
		if i < len(cols) {
			cellValues[cols[i]] = months[code]
		}
	}

	for cell, value := range cellValues {
		_ = f.SetCellValue(sheet, cell, value)
		_ = f.SetCellStyle(sheet, cell, cell, styleHeaderID)
	}

	columns := []string{
		"B", "C", "D", "E", "F",
		"G", "H", "I", "J", "K",
		"L", "M", "N", "O", "P",
		"Q",
	}

	values := []string{"T", "M", "F", "C.P"}
	for i, col := range columns {
		value := values[i%len(values)]
		cell := col + "4"
		_ = f.SetCellStyle(sheet, cell, cell, styleHeaderID)
		_ = f.SetCellValue(sheet, cell, value)
		_ = f.SetColWidth(sheet, col, col, 6)
	}

	_ = f.SetRowHeight(sheet, 3, 38)
	_ = f.SetRowHeight(sheet, 4, 24)

	return f
}

func (s *ReportsPsicopedagogiaService) GetDataAttentionsTeachers(year, month string) (string, int, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	styleHeaderID, _ := styleExcel(f)

	trimester := getTrimester(month)
	headerAttention := HeaderAttentionExcel{
		Frame: "CUADRO",
		Title: fmt.Sprintf("ATENCIONES PSICOLOGICAS A DOCENTES, NO DOCENTES Y EXTERNOS DEL %s, - AÑO %s", trimester, year),
		Area:  "PSICOPEDAGOGIA",
	}

	f = HeaderTeachersConsultingTrimesterExcel(f, sheet, headerAttention, styleHeaderID, month, trimester, year)

	for i, sName := range teachersAndExterns {
		row := 5 + i
		start, end := fmt.Sprintf("A%d", row), fmt.Sprintf("Q%d", row)
		_ = f.SetCellValue(sheet, start, sName)
		_ = f.SetCellStyle(sheet, start, end, styleHeaderID)
	}

	_ = f.SetCellValue(sheet, "A8", "TOTAL")
	_ = f.SetCellStyle(sheet, "A8", "Q8", styleHeaderID)
	_ = f.SetColWidth(sheet, "A", "A", 50)

	startDate, endDate := getDateRange(month, year)
	srv := psicopedagogia.NewServerpsicopedagogia(s.db, s.usr, s.txID)
	data, err := srv.SrvReportes.GetReportAttentionsDataTeachers(startDate, endDate)
	if err != nil {
		return "", 15, err
	}

	rowMap := map[string]int{
		"DOCENTE":    5,
		"NO DOCENTE": 6,
		"EXTERNO":    7,
		"ALUMNO":     6,
		"ESTUDIANTE": 6,
	}

	monthColumns := map[string]map[string]string{
		"04": {"M": "G", "F": "H", "CP": "I", "T": "F"},
		"05": {"M": "K", "F": "L", "CP": "M", "T": "J"},
		"06": {"M": "O", "F": "P", "CP": "Q", "T": "N"},
	}

	for _, r := range data {
		mes := strings.Split(r.Mes, "-")[1]
		row, ok := rowMap[strings.ToUpper(r.TipoParticipante)]
		if !ok {
			continue
		}
		colSexo := monthColumns[mes][r.Sexo]
		colCP := monthColumns[mes]["CP"]
		colT := monthColumns[mes]["T"]

		cellSexo := fmt.Sprintf("%s%d", colSexo, row)
		val, _ := f.GetCellValue(sheet, cellSexo)
		actual, _ := strconv.Atoi(val)
		_ = f.SetCellValue(sheet, cellSexo, actual+r.Total)

		cellCP := fmt.Sprintf("%s%d", colCP, row)
		valCP, _ := f.GetCellValue(sheet, cellCP)
		actualCP, _ := strconv.Atoi(valCP)
		_ = f.SetCellValue(sheet, cellCP, actualCP+r.Total)

		cellT := fmt.Sprintf("%s%d", colT, row)
		valT, _ := f.GetCellValue(sheet, cellT)
		actualT, _ := strconv.Atoi(valT)
		_ = f.SetCellValue(sheet, cellT, actualT+r.Total)

		cellSexoTot := fmt.Sprintf("%s8", colSexo)
		v1, _ := f.GetCellValue(sheet, cellSexoTot)
		a1, _ := strconv.Atoi(v1)
		_ = f.SetCellValue(sheet, cellSexoTot, a1+r.Total)

		cellCPTot := fmt.Sprintf("%s8", colCP)
		v2, _ := f.GetCellValue(sheet, cellCPTot)
		a2, _ := strconv.Atoi(v2)
		_ = f.SetCellValue(sheet, cellCPTot, a2+r.Total)

		cellTTot := fmt.Sprintf("%s8", colT)
		v3, _ := f.GetCellValue(sheet, cellTTot)
		a3, _ := strconv.Atoi(v3)
		_ = f.SetCellValue(sheet, cellTTot, a3+r.Total)
	}

	return SaveExcelAndReturnBase64(f, "Atenciones Psicologicas a profesores y externos")
}

func getTrimester(month string) string {
	trimester := ""
	if month == "1,2,3" {
		trimester = "I TRIMESTRE"
	}
	if month == "4,5,6" {
		trimester = "II TRIMESTRE"
	}
	if month == "7,8,9" {
		trimester = "III TRIMESTRE"
	}
	if month == "10,11,12" {
		trimester = "IV TRIMESTRE"
	}
	return trimester
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

func HeaderNursingConsultingTrimesterExcel(f *excelize.File, sheet string, headerMedicalArea HeaderAttentionExcel, styleHeaderID int, month, trimester, year string) *excelize.File {

	var mergeColumns = []MergeRange{
		{Column1: "A1", Column2: "u1"},
		{Column1: "A2", Column2: "u2"},
		{Column1: "A3", Column2: "A4"},
		{Column1: "B3", Column2: "F3"},
		{Column1: "G3", Column2: "K3"},
		{Column1: "L3", Column2: "P3"},
		{Column1: "Q3", Column2: "U3"},
	}

	for _, merge := range mergeColumns {
		_ = f.MergeCell(sheet, merge.Column1, merge.Column2)
		_ = f.SetCellStyle(sheet, merge.Column1, merge.Column2, styleHeaderID)
	}

	cellValues := map[string]string{
		"A1": headerMedicalArea.Frame,
		"A2": headerMedicalArea.Title,
		"A3": headerMedicalArea.Area,
		"A5": "TOTAL",
		"B3": fmt.Sprintf("TOTAL %s %s", trimester, year),
	}

	for col := 'B'; col <= 'U'; col++ {
		cell := fmt.Sprintf("%c5", col)
		cellValues[cell] = ""
	}

	cols := []string{"G3", "L3", "Q3"}
	monthCodes := strings.Split(month, ",")
	for i, code := range monthCodes {
		if i < len(cols) {
			cellValues[cols[i]] = months[code]
		}
	}

	for cell, value := range cellValues {
		_ = f.SetCellValue(sheet, cell, value)
		_ = f.SetCellStyle(sheet, cell, cell, styleHeaderID)
	}

	columns := []string{
		"B", "C", "D", "E", "F",
		"G", "H", "I", "J", "K",
		"L", "M", "N", "O", "P",
		"Q", "R", "S", "T", "U",
	}

	values := []string{"T", "M", "F", "E.S", "C.P"}
	for i, col := range columns {
		value := values[i%len(values)]
		cell := col + "4"
		_ = f.SetCellStyle(sheet, cell, cell, styleHeaderID)
		_ = f.SetCellValue(sheet, cell, value)
		_ = f.SetColWidth(sheet, col, col, 6)
	}

	_ = f.SetRowHeight(sheet, 3, 38)

	return f
}

func getDateRange(month, year string) (string, string) {

	rangeDates := rangeDateTrimester[month]

	dateStart := fmt.Sprintf("%s-%s", year, rangeDates[0])
	dateEnd := fmt.Sprintf("%s-%s", year, rangeDates[1])

	tEnd, _ := time.Parse("2006-01-02", dateEnd)
	tEnd = tEnd.AddDate(0, 0, 1)
	dateEnd = tEnd.Format("2006-01-02")

	return dateStart, dateEnd
}

func SaveExcelAndReturnBase64(f *excelize.File, fileName string) (string, int, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.xlsx", fileName, timestamp)
	outputPath := fmt.Sprintf("./reports/medical_area/%s", filename)

	if err := f.SaveAs(outputPath); err != nil {
		return "", 0, fmt.Errorf("error saving Excel file: %w", err)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return "", 0, fmt.Errorf("error writing Excel to buffer: %w", err)
	}
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Str, 0, nil
}

func (s *ReportsPsicopedagogiaService) GetReportByDateRange(startDate, endDate string) (string, int, error) {
	// Create Excel file
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	// Header style
	styleHeaderID, _ := styleExcel(f)

	// Column headers in Spanish (as shown in your screenshot)
	headers := []string{"DNI", "NUM. ATENCIONES", "TIPO PACIENTE", "NUM. CELULAR", "NOMBRE Y APELLIDOS", "ESCUELA", "DIAGNOSTICO", "ESTADO", "FECHA DE ENCUESTA"}
	for i, header := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s1", col), header)
		_ = f.SetCellStyle(sheet, fmt.Sprintf("%s1", col), fmt.Sprintf("%s1", col), styleHeaderID)
		_ = f.SetColWidth(sheet, col, col, 25)
	}

	// Get data from service
	srv := psicopedagogia.NewServerpsicopedagogia(s.db, s.usr, s.txID)
	data, err := srv.SrvReportes.GetReportPatientsByDateRange(startDate, endDate)
	if err != nil {
		return "", 15, err
	}

	// Fill data rows
	for i, row := range data {
		r := i + 2
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", r), row.DNI)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.NumAttentions)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row.PatientType)
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", r), row.PhoneNumber)
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", r), row.FullName)
		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", r), row.School)
		_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", r), row.Diagnosis)
		_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", r), row.Status)
		if fechaParsed, err := time.Parse(time.RFC3339, row.FechaRegistro); err == nil {
			_ = f.SetCellValue(sheet, fmt.Sprintf("I%d", r), fechaParsed.Format("2006-01-02"))
		} else {
			_ = f.SetCellValue(sheet, fmt.Sprintf("I%d", r), row.FechaRegistro)
		}
	}

	return SaveExcelAndReturnBase64(f, "Reporte_Psicopedagogia_Rango_Fechas")
}
