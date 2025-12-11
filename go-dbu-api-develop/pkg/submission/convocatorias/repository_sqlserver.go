package convocatorias

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"dbu-api/internal/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// sqlServer estructura de conexión a la BD de mssql
type sqlserver struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newConvocatoriasSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

func (s *sqlserver) create(m *Convocatorias) error {
	date := time.Now()
	m.UpdatedAt = &date
	m.CreatedAt = &date

	const sqlInsert = `INSERT INTO convocatorias 
        (fecha_inicio, fecha_fin, nombre, user_id, credito_minimo, nota_minima, 
        created_at, updated_at) 
        VALUES 
        (:fecha_inicio, :fecha_fin, :nombre, :user_id, :credito_minimo, :nota_minima, 
        :created_at, :updated_at)`

	rs, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		return err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return err
	}

	if id == 0 {
		return fmt.Errorf("rows affected error")
	}

	m.ID = id // Cambiado a int64 para coincidir con bigint unsigned
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) update(m *Convocatorias) error {
	date := time.Now()
	m.UpdatedAt = &date

	const sqlUpdate = `UPDATE convocatorias 
        SET fecha_inicio = :fecha_inicio, 
            fecha_fin = :fecha_fin, 
            nombre = :nombre, 
            user_id = :user_id, 
            credito_minimo = :credito_minimo, 
            nota_minima = :nota_minima, 
            updated_at = :updated_at 
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

// Delete elimina un registro de la BD
func (s *sqlserver) delete(id int64) error {
	const sqlDelete = `DELETE FROM convocatorias WHERE id = :id`
	m := Convocatorias{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) getByID(id int64) (*Convocatorias, error) {
	const sqlGetByID = `SELECT id, fecha_inicio, fecha_fin, nombre, user_id, 
        credito_minimo, nota_minima, created_at, updated_at 
        FROM convocatorias WHERE id = ?`

	mdl := Convocatorias{}
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
func (s *sqlserver) getAll() ([]*Convocatorias, error) {
	var ms []*Convocatorias
	const sqlGetAll = `SELECT id, fecha_inicio, fecha_fin, nombre, user_id, 
        credito_minimo, nota_minima, created_at, updated_at 
        FROM convocatorias`

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) getAllByService(id int64) ([]*Convocatorias, error) {
	var ms []*Convocatorias
	const sqlGetAll = `SELECT c.id, c.fecha_inicio, c.fecha_fin, c.nombre, c.user_id, c.credito_minimo, c.nota_minima, c.created_at, c.updated_at 
FROM convocatorias c
LEFT JOIN convocatoria_servicio cs on (cs.convocatoria_id = c.id)
WHERE cs.servicio_id = ?`

	err := s.DB.Select(&ms, sqlGetAll, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) getActive() (*Convocatorias, error) {
	const sqlGetByID = `SELECT 
            id, 
            fecha_inicio, 
            fecha_fin, 
            nombre,
            user_id,
            created_at,
            updated_at,
            credito_minimo,
            nota_minima
        FROM convocatorias
        WHERE fecha_inicio <= CURRENT_TIMESTAMP()
        AND fecha_fin >= CURRENT_TIMESTAMP() LIMIT 1;`

	mdl := Convocatorias{}
	err := s.DB.Get(&mdl, sqlGetByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
func (s *sqlserver) getLast() (*Convocatorias, error) {
	const sqlGetLast = `
        SELECT id, fecha_inicio, fecha_fin, nombre, user_id, credito_minimo, nota_minima, created_at, updated_at
        FROM convocatorias
        ORDER BY id DESC
        OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY;`

	mdl := Convocatorias{}
	err := s.DB.Get(&mdl, sqlGetLast)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// createConvocatoriaServicio crea la relación convocatoria-servicio
func (s *sqlserver) createConvocatoriaServicio(cs *ConvocatoriaServicio) error {
	date := time.Now()
	cs.CreatedAt = &date
	cs.UpdatedAt = &date

	const sqlInsert = `INSERT INTO convocatoria_servicio
		(convocatoria_id, servicio_id, cantidad, created_at, updated_at)
		VALUES (:convocatoria_id, :servicio_id, :cantidad, :created_at, :updated_at)`

	rs, err := s.DB.NamedExec(sqlInsert, cs)
	if err != nil {
		return err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return err
	}

	if id == 0 {
		return fmt.Errorf("rows affected error")
	}

	cs.ID = id
	return nil
}

// createSeccion crea una sección en la base de datos
func (s *sqlserver) createSeccion(sec *Seccion) error {
	date := time.Now()
	sec.CreatedAt = &date
	sec.UpdatedAt = &date

	const sqlInsert = `INSERT INTO secciones
		(convocatoria_id, descripcion, created_at, updated_at)
		VALUES (:convocatoria_id, :descripcion, :created_at, :updated_at)`

	rs, err := s.DB.NamedExec(sqlInsert, sec)
	if err != nil {
		return err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return err
	}

	if id == 0 {
		return fmt.Errorf("rows affected error")
	}

	sec.ID = id
	return nil
}

// createRequisito crea un requisito en la base de datos
func (s *sqlserver) createRequisito(req *Requisito) error {
	date := time.Now()
	req.CreatedAt = &date
	req.UpdatedAt = &date

	const sqlInsert = `INSERT INTO requisitos
		(seccion_id, nombre, descripcion, url_guia, url_plantilla, opciones, tipo_requisito_id, user_id, created_at, updated_at)
		VALUES (:seccion_id, :nombre, :descripcion, :url_guia, :url_plantilla, :opciones, :tipo_requisito_id, :user_id, :created_at, :updated_at)`

	rs, err := s.DB.NamedExec(sqlInsert, req)
	if err != nil {
		return err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return err
	}

	if id == 0 {
		return fmt.Errorf("rows affected error")
	}

	req.ID = id
	return nil
}

// getConvocatoriaServiciosByConvocatoriaID obtiene los servicios de una convocatoria
func (s *sqlserver) getConvocatoriaServiciosByConvocatoriaID(convocatoriaID int64) ([]ConvocatoriaServicio, error) {
	const sqlGet = `SELECT id, convocatoria_id, servicio_id, cantidad, created_at, updated_at
		FROM convocatoria_servicio WHERE convocatoria_id = ?`

	var servicios []ConvocatoriaServicio
	err := s.DB.Select(&servicios, sqlGet, convocatoriaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []ConvocatoriaServicio{}, nil
		}
		return nil, err
	}
	return servicios, nil
}

// getSeccionesByConvocatoriaID obtiene las secciones de una convocatoria
func (s *sqlserver) getSeccionesByConvocatoriaID(convocatoriaID int64) ([]Seccion, error) {
	const sqlGet = `SELECT id, convocatoria_id, descripcion, created_at, updated_at
		FROM secciones WHERE convocatoria_id = ?`

	var secciones []Seccion
	err := s.DB.Select(&secciones, sqlGet, convocatoriaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Seccion{}, nil
		}
		return nil, err
	}
	return secciones, nil
}

// getRequisitosBySeccionID obtiene los requisitos de una sección
func (s *sqlserver) getRequisitosBySeccionID(seccionID int64) ([]Requisito, error) {
	const sqlGet = `SELECT id, seccion_id, nombre, descripcion, url_guia, url_plantilla, opciones, tipo_requisito_id, user_id, created_at, updated_at
		FROM requisitos WHERE seccion_id = ?`

	var requisitos []Requisito
	err := s.DB.Select(&requisitos, sqlGet, seccionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Requisito{}, nil
		}
		return nil, err
	}
	return requisitos, nil
}

// deleteConvocatoriaServicios elimina todos los servicios de una convocatoria
func (s *sqlserver) deleteConvocatoriaServicios(convocatoriaID int64) error {
	const sqlDelete = `DELETE FROM convocatoria_servicio WHERE convocatoria_id = ?`
	_, err := s.DB.Exec(sqlDelete, convocatoriaID)
	return err
}

// deleteRequisitosBySeccionID elimina todos los requisitos de una sección
func (s *sqlserver) deleteRequisitosBySeccionID(seccionID int64) error {
	const sqlDelete = `DELETE FROM requisitos WHERE seccion_id = ?`
	_, err := s.DB.Exec(sqlDelete, seccionID)
	return err
}

// deleteSecciones elimina todas las secciones de una convocatoria
func (s *sqlserver) deleteSecciones(convocatoriaID int64) error {
	// Primero obtener todas las secciones
	secciones, err := s.getSeccionesByConvocatoriaID(convocatoriaID)
	if err != nil {
		return err
	}

	// Eliminar requisitos de cada sección
	for _, seccion := range secciones {
		if err := s.deleteRequisitosBySeccionID(seccion.ID); err != nil {
			return err
		}
	}

	// Eliminar las secciones
	const sqlDelete = `DELETE FROM secciones WHERE convocatoria_id = ?`
	_, err = s.DB.Exec(sqlDelete, convocatoriaID)
	return err
}

// checkOverlappingConvocatorias verifica si hay convocatorias solapadas
func (s *sqlserver) checkOverlappingConvocatorias(fechaInicio time.Time, excludeID *int64) (bool, error) {
	var count int
	var sqlQuery string
	var args []interface{}

	if excludeID != nil {
		sqlQuery = `SELECT COUNT(*) FROM convocatorias WHERE fecha_fin >= ? AND id != ?`
		args = []interface{}{fechaInicio, *excludeID}
	} else {
		sqlQuery = `SELECT COUNT(*) FROM convocatorias WHERE fecha_fin >= ?`
		args = []interface{}{fechaInicio}
	}

	err := s.DB.Get(&count, sqlQuery, args...)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
