package chapters

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

func newChapterSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) create(m *Chapter) error {
	const sqlInsert = `INSERT INTO capitulos (id, nombre, descripcion, resolucion_id,  created_at, updated_at) 
              VALUES (:id, :nombre, :descripcion,  :resolucion_id, :created_at,  :updated_at)`
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
func (s *sqlserver) update(m *Chapter) error {
	const sqlUpdate = `UPDATE capitulos SET nombre = :nombre, descripcion = :descripcion,  resolucion_id = :resolucion_id, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM capitulos WHERE id = :id`
	m := Chapter{ID: id}
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
func (s *sqlserver) getByID(id string) (*Chapter, error) {
	const sqlGetByID = `SELECT id, nombre, descripcion, resolucion_id,created_at, updated_at FROM capitulos  WHERE id = ? `
	mdl := Chapter{}
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
func (s *sqlserver) getAll() ([]*Chapter, error) {
	var ms []*Chapter
	const sqlGetAll = `SELECT id, nombre, descripcion, resolucion_id, created_at, updated_at FROM capitulos`

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (r *sqlserver) GetByResolutionID(ResolucionId string) ([]*Chapter, error) {
	const sqlStatement = `SELECT id, nombre, descripcion, resolucion_id FROM capitulos WHERE resolucion_id = ?`
	rows, err := r.DB.Query(sqlStatement, ResolucionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []*Chapter
	for rows.Next() {
		var chapter Chapter
		err = rows.Scan(&chapter.ID, &chapter.Nombre, &chapter.Descripcion, &chapter.Resolucion_id)
		if err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}
	return chapters, nil
}

func (s *sqlserver) updateOnlyCharacteristics(m *Chapter) error {
	const sqlUpdate = `UPDATE cuartos SET capacidad = :capacidad, estado = :estado, updated_by = :updated_by, updated_at = :updated_at 
  						WHERE id = :id`
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}
