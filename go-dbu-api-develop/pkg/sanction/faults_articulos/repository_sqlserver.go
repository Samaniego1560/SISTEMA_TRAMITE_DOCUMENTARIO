// pkg/sanction/faults_articulos/repository_sqlserver.go
package faults_articulos

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

func newFaultArticuloSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

func (s *sqlserver) create(m *FaultArticulo) error {
	const sqlInsert = `INSERT INTO faltas_articulos (id, falta_id, articulo_id, created_at, updated_at)
	                   VALUES (:id, :falta_id, :articulo_id, :created_at, :updated_at)`
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
	const sqlDelete = `DELETE FROM faltas_articulos WHERE id = :id`
	m := FaultArticulo{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *sqlserver) getAll() ([]*FaultArticulo, error) {
	var ms []*FaultArticulo
	const sqlGetAll = `SELECT id, falta_id, articulo_id, created_at, updated_at FROM faltas_articulos`
	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
