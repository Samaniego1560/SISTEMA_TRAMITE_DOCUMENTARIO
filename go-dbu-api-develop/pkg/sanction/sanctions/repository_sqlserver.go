// repository_sqlserver.go

package sanctions

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

func newSanctionSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

func (s *sqlserver) create(m *Sanction) error {
	date := time.Now()
	m.CreatedAt = date
	m.UpdatedAt = date
	const sqlInsert = `INSERT INTO sanciones (id, fault_id, tipo_sancion, duracion, fecha_inicio, fecha_fin, estado, observacion, created_at, updated_at) 
					   VALUES (:id, :fault_id, :tipo_sancion, :duracion, :fecha_inicio, :fecha_fin, :estado, :observacion, :created_at, :updated_at)`
	rs, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *sqlserver) update(m *Sanction) error {
	date := time.Now()
	m.UpdatedAt = date
	const sqlUpdate = `UPDATE sanciones SET fault_id = :fault_id, tipo_sancion = :tipo_sancion, duracion = :duracion, fecha_inicio = :fecha_inicio, 
					   fecha_fin = :fecha_fin, estado = :estado, observacion = :observacion, updated_at = :updated_at WHERE id = :id`
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *sqlserver) delete(id string) error {
	const sqlDelete = `DELETE FROM sanciones WHERE id = :id`
	m := Sanction{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *sqlserver) getByID(id string) (*Sanction, error) {
	const sqlGetByID = `SELECT id, fault_id, tipo_sancion, duracion, fecha_inicio, fecha_fin, estado, observacion, created_at, updated_at 
						FROM sanciones WHERE id = @id`
	mdl := Sanction{}
	err := s.DB.Get(&mdl, sqlGetByID, sql.Named("id", id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *sqlserver) getAll() ([]*Sanction, error) {
	var ms []*Sanction
	const sqlGetAll = `SELECT id, fault_id, tipo_sancion, duracion, fecha_inicio, fecha_fin, estado, observacion, created_at, updated_at 
					   FROM sanciones`
	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
