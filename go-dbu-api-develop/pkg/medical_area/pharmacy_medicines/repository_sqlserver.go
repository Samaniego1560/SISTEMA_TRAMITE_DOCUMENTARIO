package pharmacy_medicines

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// sqlserver estructura de conexión a la BD
type sqlserver struct {
	DB   *sqlx.DB
	TxID string
}

func newMedicineSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

// existsByCode verifica si existe un medicamento con el código dado
func (s *sqlserver) existsByCode(code string) (bool, error) {
	var count int
	const sqlExists = `SELECT COUNT(1) FROM farmacia_medicamentos WHERE codigo = ? AND is_deleted = 0`
	err := s.DB.QueryRow(sqlExists, code).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// create inserta un nuevo medicamento
func (s *sqlserver) create(m *Medicine) error {
	const sqlInsert = `INSERT INTO farmacia_medicamentos 
		(id, codigo, nombre_generico, nombre_comercial, forma_farmaceutica, concentracion, 
		unidad_base, via_administracion, requiere_receta, controlado, descripcion, estado, 
		user_creator, created_at, updated_at) 
		VALUES (:id, :codigo, :nombre_generico, :nombre_comercial, :forma_farmaceutica, :concentracion, 
		:unidad_base, :via_administracion, :requiere_receta, :controlado, :descripcion, :estado, 
		:user_creator, :created_at, :updated_at)`

	rs, err := s.DB.NamedExec(sqlInsert, m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}

// update actualiza un medicamento existente
func (s *sqlserver) update(m *Medicine) error {
	const sqlUpdate = `UPDATE farmacia_medicamentos SET 
		codigo = :codigo,
		nombre_generico = :nombre_generico,
		nombre_comercial = :nombre_comercial,
		forma_farmaceutica = :forma_farmaceutica,
		concentracion = :concentracion,
		via_administracion = :via_administracion,
		requiere_receta = :requiere_receta,
		controlado = :controlado,
		descripcion = :descripcion,
		estado = :estado,
		updated_at = NOW()
		WHERE id = :id AND is_deleted = 0`

	rs, err := s.DB.NamedExec(sqlUpdate, m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}

// delete realiza soft delete de un medicamento
func (s *sqlserver) delete(id string) error {
	const sqlDelete = `UPDATE farmacia_medicamentos SET 
		is_deleted = 1, 
		deleted_at = NOW() 
		WHERE id = ? AND is_deleted = 0`

	rs, err := s.DB.Exec(sqlDelete, id)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}
	return nil
}

// getByID obtiene un medicamento por ID
func (s *sqlserver) getByID(id string) (*Medicine, error) {
	const sqlGetByID = `SELECT id, codigo, nombre_generico, nombre_comercial, forma_farmaceutica, 
		concentracion, unidad_base, via_administracion, requiere_receta, controlado, descripcion, 
		estado, user_creator, created_at, updated_at 
		FROM farmacia_medicamentos 
		WHERE id = ? AND is_deleted = 0`

	var m Medicine
	err := s.DB.Get(&m, sqlGetByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// getByCode obtiene un medicamento por código
func (s *sqlserver) getByCode(code string) (*Medicine, error) {
	const sqlGetByCode = `SELECT id, codigo, nombre_generico, nombre_comercial, forma_farmaceutica, 
		concentracion, unidad_base, via_administracion, requiere_receta, controlado, descripcion, 
		estado, user_creator, created_at, updated_at 
		FROM farmacia_medicamentos 
		WHERE codigo = ? AND is_deleted = 0`

	var m Medicine
	err := s.DB.Get(&m, sqlGetByCode, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// getAll obtiene todos los medicamentos activos
func (s *sqlserver) getAll() ([]*Medicine, error) {
	const sqlGetAll = `SELECT id, codigo, nombre_generico, nombre_comercial, forma_farmaceutica, 
		concentracion, unidad_base, via_administracion, requiere_receta, controlado, descripcion, 
		estado, user_creator, created_at, updated_at 
		FROM farmacia_medicamentos 
		WHERE is_deleted = 0 
		ORDER BY nombre_generico`

	var medicines []*Medicine
	err := s.DB.Select(&medicines, sqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return medicines, nil
}

// getAllWithStock obtiene todos los medicamentos con información de stock
func (s *sqlserver) getAllWithStock() ([]*MedicineWithStock, error) {
	const sqlGetAllWithStock = `
		SELECT 
			m.id, m.codigo, m.nombre_generico, m.nombre_comercial, m.forma_farmaceutica, 
			m.concentracion, m.unidad_base, m.via_administracion, m.requiere_receta, 
			m.controlado, m.descripcion, m.estado, m.user_creator, m.created_at, m.updated_at,
			COALESCE(SUM(l.cantidad_disponible), 0) as stock_total,
			COUNT(DISTINCT CASE WHEN l.cantidad_disponible > 0 AND l.fecha_vencimiento >= CURDATE() THEN l.id END) as lotes_activos
		FROM farmacia_medicamentos m
		LEFT JOIN farmacia_lotes l ON l.medicamento_id = m.id AND l.is_deleted = 0
		WHERE m.is_deleted = 0
		GROUP BY m.id
		ORDER BY m.nombre_generico`

	var medicines []*MedicineWithStock
	err := s.DB.Select(&medicines, sqlGetAllWithStock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return medicines, nil
}

// searchMedicines busca medicamentos con filtros
func (s *sqlserver) searchMedicines(search, estado string, limit, offset int64) ([]*MedicineWithStock, error) {
	query := `
		SELECT 
			m.id, m.codigo, m.nombre_generico, m.nombre_comercial, m.forma_farmaceutica, 
			m.concentracion, m.unidad_base, m.via_administracion, m.requiere_receta, 
			m.controlado, m.descripcion, m.estado, m.user_creator, m.created_at, m.updated_at,
			COALESCE(SUM(l.cantidad_disponible), 0) as stock_total,
			COUNT(DISTINCT CASE WHEN l.cantidad_disponible > 0 AND l.fecha_vencimiento >= CURDATE() THEN l.id END) as lotes_activos
		FROM farmacia_medicamentos m
		LEFT JOIN farmacia_lotes l ON l.medicamento_id = m.id AND l.is_deleted = 0
		WHERE m.is_deleted = 0`

	var conditions []string
	var params []interface{}

	if search != "" {
		conditions = append(conditions, "(m.codigo LIKE ? OR m.nombre_generico LIKE ? OR m.nombre_comercial LIKE ?)")
		searchParam := "%" + search + "%"
		params = append(params, searchParam, searchParam, searchParam)
	}

	if estado != "" {
		conditions = append(conditions, "m.estado = ?")
		params = append(params, estado)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " GROUP BY m.id ORDER BY m.nombre_generico LIMIT ? OFFSET ?"
	params = append(params, limit, offset)

	var medicines []*MedicineWithStock
	err := s.DB.Select(&medicines, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return medicines, nil
}

// countMedicines cuenta medicamentos con filtros
func (s *sqlserver) countMedicines(search, estado string) (int64, error) {
	query := `SELECT COUNT(DISTINCT m.id) FROM farmacia_medicamentos m WHERE m.is_deleted = 0`

	var conditions []string
	var params []interface{}

	if search != "" {
		conditions = append(conditions, "(m.codigo LIKE ? OR m.nombre_generico LIKE ? OR m.nombre_comercial LIKE ?)")
		searchParam := "%" + search + "%"
		params = append(params, searchParam, searchParam, searchParam)
	}

	if estado != "" {
		conditions = append(conditions, "m.estado = ?")
		params = append(params, estado)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var total int64
	err := s.DB.Get(&total, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return total, nil
}
