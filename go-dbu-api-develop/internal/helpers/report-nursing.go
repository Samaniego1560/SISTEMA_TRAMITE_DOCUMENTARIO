package helpers

import (
	"dbu-api/internal/models"
)

func CreateHeaderReportNursingRua(consultation *models.ReportNursingRua, row int) models.ExcelPageRow {
	if row == 1 {
		return models.ExcelPageRow{
			Row: 1,
			Columns: []models.ExcelPageColumn{
				{Column: "A", Value: "FECHA"},
				{Column: "B", Value: "Nº HCL/ CODIGO"},
				{Column: "C", Value: "DNI"},
				{Column: "D", Value: "APELLIDOS Y NOMBRES"},
				{Column: "E", Value: "FECHA DE NACIMIENTO"},
				{Column: "F", Value: "SEXO (M)(F)"},
				{Column: "G", Value: "EDAD"},
				{Column: "H", Value: "DOMICILIO CON REFERENCIA"},
				{Column: "I", Value: "LUGAR DE PROCEDENCIA"},
				{Column: "J", Value: "COND. EE.SS"},
				{Column: "K", Value: "DOSIS VACUNA"},
				{Column: "L", Value: "VACUNA"},
				{Column: "M", Value: "MINSA"},
				{Column: "N", Value: "DIAGNOSTICO Y/O PROCEDIMIENTO"},
				{Column: "O", Value: "TRATAMIENTO"},
				{Column: "P", Value: "SERVICIO"},
				{Column: "Q", Value: "RESPONSABLE DE LA ATENCIÓN"},
				{Column: "R", Value: "TRABAJADOR ADM/DOC/OTROS/ESTUDIANTE"},
				{Column: "S", Value: "ESCUELA  ACADEMICA PROFESIONAL"},
				{Column: "T", Value: "Nº DE CELULAR DEL FAMILIAR"},
				{Column: "U", Value: "OBSERVACIONES"},
				{Column: "V", Value: "N° RECIBO"},
				{Column: "W", Value: "MONTO"},
			},
		}
	}

	return models.ExcelPageRow{
		Row: row,
		Columns: []models.ExcelPageColumn{
			{Column: "A", Value: consultation.Fecha},
			{Column: "B", Value: consultation.Codigo},
			{Column: "C", Value: consultation.DNI},
			{Column: "D", Value: consultation.ApellidosNombres},
			{Column: "E", Value: consultation.FechaNacimiento},
			{Column: "F", Value: consultation.Sexo},
			{Column: "G", Value: consultation.Edad},
			{Column: "H", Value: consultation.Domicilio},
			{Column: "I", Value: consultation.Procedencia},
			{Column: "J", Value: consultation.CondicionSalud},
			{Column: "K", Value: consultation.DosisVacuna},
			{Column: "L", Value: consultation.Vacuna},
			{Column: "M", Value: consultation.Minsa},
			{Column: "N", Value: consultation.Diagnostico},
			{Column: "O", Value: consultation.Tratamiento},
			{Column: "P", Value: consultation.Servicio},
			{Column: "Q", Value: consultation.Responsable},
			{Column: "R", Value: consultation.TipoPaciente},
			{Column: "S", Value: consultation.Escuela},
			{Column: "T", Value: consultation.Celular},
			{Column: "U", Value: consultation.Observaciones},
			{Column: "V", Value: consultation.Recibo},
			{Column: "W", Value: consultation.Monto},
		},
	}
}

func CreateHeaderReport1Admin(consultation *models.Report1Admin, row int) models.ExcelPageRow {
	if row == 1 {
		return models.ExcelPageRow{
			Row: 1,
			Columns: []models.ExcelPageColumn{
				{Column: "A", Value: "FECHA"},
				{Column: "B", Value: "TIPO PERSONA"},
				{Column: "C", Value: "CÓDIGO"},
				{Column: "D", Value: "DNI"},
				{Column: "E", Value: "APELLIDOS Y NOMBRES"},
				{Column: "F", Value: "FECHA DE NACIMIENTO"},
				{Column: "G", Value: "SEXO"},
				{Column: "H", Value: "EDAD"},
				{Column: "I", Value: "ESCUELA PROFESIONAL"},
				{Column: "J", Value: "OCUPACIÓN"},
				{Column: "K", Value: "AREA MEDICA"},
			},
		}
	}

	return models.ExcelPageRow{
		Row: row,
		Columns: []models.ExcelPageColumn{
			{Column: "A", Value: consultation.Fecha},
			{Column: "B", Value: consultation.TipoPersona},
			{Column: "C", Value: consultation.Codigo},
			{Column: "D", Value: consultation.DNI},
			{Column: "E", Value: consultation.ApellidosNombres},
			{Column: "F", Value: consultation.FechaNacimiento},
			{Column: "G", Value: consultation.Sexo},
			{Column: "H", Value: consultation.Edad},
			{Column: "I", Value: consultation.EscuelaProfesional},
			{Column: "J", Value: consultation.Ocupacion},
			{Column: "K", Value: consultation.AreaMedica},
		},
	}
}

func CreateHeaderReport2Admin(consultation *models.Report2Admin, row int) models.ExcelPageRow {
	if row == 1 {
		return models.ExcelPageRow{
			Row: 1,
			Columns: []models.ExcelPageColumn{
				{Column: "A", Value: "Total enfermería"},
				{Column: "B", Value: "Total medicina"},
				{Column: "C", Value: "Total odontología"},
			},
		}
	}

	return models.ExcelPageRow{
		Row: row,
		Columns: []models.ExcelPageColumn{
			{Column: "A", Value: consultation.TotalEnfermeria},
			{Column: "B", Value: consultation.TotalMedicina},
			{Column: "C", Value: consultation.TotalOdontologia},
		},
	}
}
