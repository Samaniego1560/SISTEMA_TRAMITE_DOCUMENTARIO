package nursing_consultation_performed_procedures

import (
	"database/sql"
	"dbu-api/internal/models"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type sqlserver struct {
	DB   *sqlx.DB
	TxID string
}

func newPerformedProceduresSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

func (s sqlserver) create(m *PerformedProcedures) error {
	const sqlInsertConsulta = `INSERT INTO enfermeria_procedimientos_realizados (id, consulta_enfermeria_id, procedimiento, numero_recibo, comentarios, costo, fecha_pago, area_solicitante, especialista_solicitante, is_deleted, user_deleted, deleted_at, user_creator, created_at, updated_at) VALUES (:id, :consulta_enfermeria_id, :procedimiento, :numero_recibo, :comentarios, :costo, :fecha_pago, :area_solicitante, :especialista_solicitante, :is_deleted, :user_deleted, :deleted_at, :user_creator, :created_at, :updated_at)`

	rs, err := s.DB.NamedExec(sqlInsertConsulta, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}

func (s sqlserver) update(m *PerformedProcedures) error {
	const sqlUpdate = `UPDATE enfermeria_procedimientos_realizados SET consulta_enfermeria_id = :consulta_enfermeria_id, procedimiento = :procedimiento, numero_recibo = :numero_recibo, comentarios = :comentarios, costo = :costo, fecha_pago = :fecha_pago, area_solicitante = :area_solicitante, especialista_solicitante = :especialista_solicitante, is_deleted = :is_deleted, user_deleted = :user_deleted, deleted_at = :deleted_at, user_creator = :user_creator, created_at = :created_at, updated_at = :updated_at WHERE id = :id`
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}

func (s sqlserver) delete(id string) error {
	const psqlDelete = `DELETE FROM enfermeria_procedimientos_realizados WHERE id = :id`
	m := PerformedProcedures{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}

	return nil
}

func (s sqlserver) deleteByIDConsultation(consulta_enfermeria_id string) error {
	const psqlDelete = `DELETE FROM enfermeria_procedimientos_realizados WHERE consulta_enfermeria_id = :consulta_enfermeria_id`
	m := PerformedProcedures{IDConsultaEnfermeria: consulta_enfermeria_id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}

	return nil
}

func (s sqlserver) getByID(id string) (*PerformedProcedures, error) {
	const sqlGetByID = `SELECT id, consulta_enfermeria_id, procedimiento, numero_recibo, comentarios, costo, fecha_pago, area_solicitante, especialista_solicitante, is_deleted, user_deleted, deleted_at, user_creator, created_at, updated_at FROM enfermeria_procedimientos_realizados WHERE id = ?`

	mdl := PerformedProcedures{}
	err := s.DB.Get(&mdl, sqlGetByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &mdl, nil
}

func (s sqlserver) getAll() ([]*PerformedProcedures, error) {
	var ms []*PerformedProcedures
	const sqlGetAll = `SELECT id, consulta_enfermeria_id, procedimiento, numero_recibo, comentarios, costo, fecha_pago, area_solicitante, especialista_solicitante, is_deleted, user_deleted, deleted_at, user_creator, created_at, updated_at FROM enfermeria_procedimientos_realizados`

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}

	return ms, nil
}

func (s sqlserver) GetByIDConsultation(id string) (*PerformedProcedures, error) {
	mdl := PerformedProcedures{}
	const sqlGetByIDConsultation = `SELECT id, consulta_enfermeria_id, procedimiento, numero_recibo, comentarios, costo, fecha_pago, area_solicitante, especialista_solicitante, is_deleted, user_deleted, deleted_at, user_creator, created_at, updated_at FROM enfermeria_procedimientos_realizados WHERE consulta_enfermeria_id = ?`

	err := s.DB.Get(&mdl, sqlGetByIDConsultation, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &mdl, err
	}

	return &mdl, nil
}

func (s sqlserver) getAllByDateExcel(fecha_inicio, fecha_fin string) ([]*models.PerformedProceduresExcel, error) {
	var ms []*models.PerformedProceduresExcel
	const sqlGetAll = `
		SELECT 
		cam.id as id, 
		cam.fecha_consulta as fecha_consulta, 
		p.tipo_persona as tipo_persona,
		p.escuela_profesional as escuela_profesional,
		p.sexo as sexo,
		epr.procedimiento as tipo_procedimiento
	FROM consultas_areas_medicas cam
	JOIN pacientes p ON p.id = cam.paciente_id
	JOIN enfermeria_procedimientos_realizados epr ON epr.consulta_enfermeria_id = cam.id
	WHERE STR_TO_DATE(cam.fecha_consulta, '%Y-%m-%d') >= ? AND STR_TO_DATE(cam.fecha_consulta, '%Y-%m-%d') <= ?
	`

	err := s.DB.Select(&ms, sqlGetAll, fecha_inicio, fecha_fin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}

	return ms, nil
}

func (s sqlserver) getByYearAndMonths(year int, months []int) ([]*models.PerformedProceduresExcel, error) {
	var ms []*models.PerformedProceduresExcel

	whereClause, arguments := generateWhereClauseQuery(year, months)
	query := fmt.Sprintf(`
		SELECT 
		cam.id as id, 
		cam.fecha_consulta as fecha_consulta, 
		p.tipo_persona as tipo_persona,
		p.escuela_profesional as escuela_profesional,
		p.sexo as sexo,
		epr.procedimiento as tipo_procedimiento
	FROM consultas_areas_medicas cam
	JOIN pacientes p ON p.id = cam.paciente_id
	JOIN enfermeria_procedimientos_realizados epr ON epr.consulta_enfermeria_id = cam.id
	WHERE %s`, whereClause)

	err := s.DB.Select(&ms, query, arguments...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}

	return ms, nil
}

func generateWhereClauseQuery(year int, months []int) (string, []interface{}) {
	var conditions []string
	var arguments []interface{}

	for _, mes := range months {
		fechaInicio := time.Date(year, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)

		nextMonth := fechaInicio.AddDate(0, 1, 0)
		dateEnd := nextMonth.Add(-24 * time.Hour)

		condition := "(fecha_consulta BETWEEN ? AND ?)"
		conditions = append(conditions, condition)
		arguments = append(arguments, fechaInicio.Format("2006-01-02"))
		arguments = append(arguments, dateEnd.Format("2006-01-02"))
	}

	whereClause := strings.Join(conditions, " OR ")

	return whereClause, arguments
}
