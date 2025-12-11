package reportes

import (
	"database/sql"
	"dbu-api/internal/models"
	"errors"

	"github.com/jmoiron/sqlx"
)

type ReportesRepository struct {
	DB *sqlx.DB
}

func NewReportesRepository(db *sqlx.DB) *ReportesRepository {
	return &ReportesRepository{db}
}

func (r *ReportesRepository) GetReportAttentionsDataStudents(fecha_inicio, fecha_fin string) ([]*models.ConsultationAttentionExcel, error) {
	var ms []*models.ConsultationAttentionExcel

	const sqlGetAll = `
		SELECT 
		cam.id as id, 
		cam.fecha_consulta as fecha_consulta, 
		p.tipo_persona as tipo_persona,
		p.escuela_profesional as escuela_profesional,
		p.sexo as sexo
	FROM consultas_areas_medicas cam
	JOIN pacientes p ON p.id = cam.paciente_id
	JOIN medicina_consulta_general mcg ON mcg.consulta_id = cam.id
	AND cam.created_at >= ?
	AND cam.created_at <= ?
	`

	err := r.DB.Select(&ms, sqlGetAll, fecha_inicio, fecha_fin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}

	return ms, nil

}

func (r *ReportesRepository) GetReportAttentionsDataTeachers(fecha_inicio, fecha_fin string) ([]*models.SRQMonthlySummary, error) {
	var ms []*models.SRQMonthlySummary

	const sqlGetAll = `
		SELECT
		p.tipo_participante,
		p.sexo,
		DATE_FORMAT(h.fecha_respuesta, '%Y-%m') AS mes,
		COUNT(*) AS total
	FROM historial h
	JOIN participantes p ON p.id_participante = h.id_participante
	WHERE h.fecha_respuesta BETWEEN ? AND ?
	  AND h.es_srq = 0
	GROUP BY p.tipo_participante, p.sexo, mes
	ORDER BY p.tipo_participante, p.sexo, mes;
	`

	err := r.DB.Select(&ms, sqlGetAll, fecha_inicio, fecha_fin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}

	return ms, nil
}

func (r *ReportesRepository) GetReportPatientsByDateRange(startDate, endDate string) ([]*models.PatientReportExcel, error) {
	query := `
        SELECT 
            p.dni,
            p.numero_atencion,
            p.tipo_participante,
            p.num_telefono,
            CONCAT(p.nombre, ' ', p.apellido) AS full_name,
            p.escuela,
            p.diagnostico,
            p.estado_evaluacion,
            MAX(h.created_date) AS created_date  -- Toma la fecha más reciente
        FROM historial h
        INNER JOIN participantes p ON h.id_participante = p.id_participante
        WHERE h.created_date BETWEEN ? AND ? AND h.es_srq = true
        GROUP BY p.id_participante, p.dni, p.numero_atencion, p.tipo_participante, 
                 p.num_telefono, p.nombre, p.apellido, p.escuela, p.diagnostico, p.estado_evaluacion
        ORDER BY MAX(h.created_date) DESC
    `

	var results []*models.PatientReportExcel
	err := r.DB.Select(&results, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return results, nil
}
