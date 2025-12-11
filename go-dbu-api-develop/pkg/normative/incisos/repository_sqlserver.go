package incisos

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

func newIncisoSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) create(m *Inciso) error {
	const sqlInsert = `INSERT INTO incisos (id, nombre, descripcion, articulo_id,  created_at, updated_at) 
              VALUES (:id, :nombre, :descripcion, :articulo_id, :created_at,  :updated_at)`
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
func (s *sqlserver) update(m *Inciso) error {
	const sqlUpdate = `UPDATE incisos SET  nombre = :nombre, descripcion = :descripcion, articulo_id = :articulo_id, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM incisos WHERE id = :id`
	m := Inciso{ID: id}
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
func (s *sqlserver) getByID(id string) (*Inciso, error) {
	const sqlGetByID = `SELECT id, nombre, descripcion, articulo_id,created_at, updated_at FROM incisos  WHERE id = ? `
	mdl := Inciso{}
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
func (s *sqlserver) getAll() ([]*Inciso, error) {
	var ms []*Inciso
	const sqlGetAll = `SELECT id, nombre, descripcion, articulo_id, created_at, updated_at FROM incisos`

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (r *sqlserver) GetByarticleID(articuloId string) ([]*Inciso, error) {
	const sqlStatement = `SELECT id, nombre, descripcion, articulo_id FROM incisos WHERE articulo_id = ?`

	rows, err := r.DB.Query(sqlStatement, articuloId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incisos []*Inciso
	for rows.Next() {
		var inciso Inciso
		err = rows.Scan(&inciso.ID, &inciso.Nombre, &inciso.Descripcion, &inciso.ArticuloId)
		if err != nil {
			return nil, err
		}
		incisos = append(incisos, &inciso)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return incisos, nil
}
