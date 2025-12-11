package submissions

import "dbu-api/internal/models"

// ResponseStudentsSubmission representa la respuesta de alumnos por convocatoria
type ResponseStudentsSubmission struct {
	Total    int               `json:"total"`    // Total de estudiantes encontrados
	Students []*models.Student `json:"students"` // Lista de estudiantes
}
