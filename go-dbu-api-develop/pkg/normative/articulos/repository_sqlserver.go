package articles

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

func newArticleSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) create(m *Article) error {
	const sqlInsert = `INSERT INTO articulos (id, descripcion, gravedad, capitulo_id,  created_at, updated_at) 
              VALUES (:id, :descripcion, :gravedad, :capitulo_id, :created_at,  :updated_at)`
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
func (s *sqlserver) update(m *Article) error {
	const sqlUpdate = `UPDATE articulos SET  descripcion = :descripcion, gravedad = :gravedad, capitulo_id = :capitulo_id, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM articulos WHERE id = :id`
	m := Article{ID: id}
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
func (s *sqlserver) getByID(id string) (*Article, error) {
	const sqlGetByID = `SELECT id, gravedad, descripcion, capitulo_id,created_at, updated_at FROM articulos  WHERE id = ? `
	mdl := Article{}
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
func (s *sqlserver) getAll() ([]*Article, error) {
	var ms []*Article
	// Solo artículos cuyos capítulos pertenecen a una resolución activa (estado = 1)
	const sqlGetAll = `
	       SELECT a.id, a.gravedad, a.descripcion, a.capitulo_id, a.created_at, a.updated_at, c.resolucion_id
	       FROM articulos a
	       JOIN capitulos c ON a.capitulo_id = c.id
	       JOIN resoluciones r ON c.resolucion_id = r.id
	       WHERE r.estado = 1`

	var rows []Article
	err := s.DB.Select(&rows, sqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convertir a []*Article y agregar resolucion_id como campo adicional usando map[string]interface{} si es necesario
	for _, row := range rows {
		art := &Article{
			ID:           row.ID,
			Descripcion:  row.Descripcion,
			Gravedad:     row.Gravedad,
			Capitulo_id:  row.Capitulo_id,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
			ResolucionID: row.ResolucionID,
		}
		ms = append(ms, art)
	}
	// Si quieres devolver resolucion_id directamente en el modelo, agrega el campo a Article y asígnalo aquí
	return ms, nil
}

func (r *sqlserver) GetByChapterID(CapituloId int64) ([]*Article, error) {
	const sqlStatement = `SELECT id, nombre, descripcion, capitulo_id, gravedad FROM articulos WHERE capitulo_id = $1`
	rows, err := r.DB.Query(sqlStatement, CapituloId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*Article
	for rows.Next() {
		var article Article
		err = rows.Scan(&article.ID, &article.Descripcion, &article.Capitulo_id)
		if err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}
	return articles, nil
}
func (s *sqlserver) updateOnlyCharacteristics(m *Article) error {
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
