package faults_incisos

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type sqlserver struct {
	DB   *sqlx.DB
	TxID string
}

func newFaultIncisoSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

func (s *sqlserver) create(m *FaultInciso) error {
	const sqlInsert = `INSERT INTO faltas_incisos (id, falta_id, inciso_id, created_at, updated_at)
	                   VALUES (:id, :falta_id, :inciso_id, :created_at, :updated_at)`
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	rs, err := s.DB.NamedExec(sqlInsert, m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *sqlserver) delete(id string) error {
	const sqlDelete = `DELETE FROM faltas_incisos WHERE id = :id`
	m := FaultInciso{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *sqlserver) getAll() ([]*FaultInciso, error) {
	var ms []*FaultInciso
	const sqlGetAll = `SELECT id, falta_id, inciso_id, created_at, updated_at FROM faltas_incisos`
	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
