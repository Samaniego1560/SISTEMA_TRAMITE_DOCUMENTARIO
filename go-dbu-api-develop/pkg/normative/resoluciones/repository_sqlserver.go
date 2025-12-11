package resolutions

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// sqlServer estructura de conexión a la BD de mssql
type sqlserver struct {
	DB   *sqlx.DB
	TxID string
}

func newResolutionSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) create(m *Resolution) error {
	const sqlInsert = `INSERT INTO resoluciones (id, nombre, descripcion, servicio_id, ruta_archivo,  created_at, updated_at) 
              VALUES (:id, :nombre, :descripcion, :servicio_id, :ruta_archivo, :created_at,  :updated_at)`
	rs, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) update(m *Resolution) error {
	const sqlUpdate = `UPDATE resoluciones SET nombre = :nombre, descripcion = :descripcion, estado = :estado, servicio_id = :servicio_id, ruta_archivo = :ruta_archivo, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) delete(id string) error {

	// Physical delete
	const psqlDelete = `DELETE FROM resoluciones WHERE id = :id`
	m := Resolution{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}

	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) getByID(id string) (*Resolution, error) {
	const sqlGetByID = `SELECT id, nombre, descripcion, servicio_id, ruta_archivo,created_at, updated_at FROM resoluciones  WHERE id = ? `
	mdl := Resolution{}
	err := s.DB.Get(&mdl, sqlGetByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *sqlserver) getAll() ([]*Resolution, error) {
	var ms []*Resolution
	const sqlGetAll = `
	SELECT id, nombre, descripcion, servicio_id, ruta_archivo, estado, created_at, updated_at 
	FROM resoluciones
	ORDER BY created_at DESC
	`

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
