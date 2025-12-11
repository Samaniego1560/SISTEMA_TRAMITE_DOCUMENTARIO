package asignaciones_cuartos

import (
	"database/sql"
	"dbu-api/internal/models"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

// sqlServer estructura de conexión a la BD de mssql
type sqlserver struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newAsignacionesCuartosSqlServerRepository(db *sqlx.DB, user *models.User, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *sqlserver) create(m *RoomAssignment) error {
	const sqlInsert = `INSERT INTO asignacion_cuartos (id, alumno_id, cuarto_id, convocatoria_id, 
                       fecha_asignacion, estado, observaciones, created_at, updated_at)
                       VALUES (:id, :alumno_id, :cuarto_id, :convocatoria_id,
                       :fecha_asignacion, :estado, :observaciones, :created_at, :updated_at)`
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
func (s *sqlserver) update(m *RoomAssignment) error {
	const sqlUpdate = `UPDATE asignacion_cuartos SET 
                      alumno_id = :alumno_id, 
                      cuarto_id = :cuarto_id,
                      convocatoria_id = :convocatoria_id,
                      fecha_asignacion = :fecha_asignacion,
                      estado = :estado,
                      observaciones = :observaciones,
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
func (s *sqlserver) delete(m *RoomAssignment) error {
	const sqlDelete = `UPDATE asignacion_cuartos SET estado = :estado, observaciones = :observaciones
    WHERE alumno_id = :alumno_id AND cuarto_id = :cuarto_id 
    AND convocatoria_id = :convocatoria_id`
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
func (s *sqlserver) getByID(id string) (*RoomAssignment, error) {
	const sqlGetByID = `SELECT id, alumno_id, cuarto_id, convocatoria_id, 
                       fecha_asignacion, estado, observaciones, created_at, updated_at 
                       FROM asignacion_cuartos WHERE id = ?`
	mdl := RoomAssignment{}
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
func (s *sqlserver) getAll() ([]*RoomAssignment, error) {
	var ms []*RoomAssignment
	const sqlGetAll = `SELECT id, alumno_id, cuarto_id, convocatoria_id, 
                      fecha_asignacion, estado, observaciones, created_at, updated_at 
                      FROM asignacion_cuartos`

	err := s.DB.Select(&ms, sqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) getRoomAssignmentByRoomIDSubmissionID(roomID string, submissionID int64) ([]*RoomAssignment, error) {
	var ms []*RoomAssignment
	const sqlGetByID = `SELECT id, alumno_id, cuarto_id, convocatoria_id, 
                       fecha_asignacion, estado, observaciones, created_at, updated_at 
                       FROM asignacion_cuartos WHERE cuarto_id = ? AND convocatoria_id = ?`
	err := s.DB.Select(&ms, sqlGetByID, roomID, submissionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *sqlserver) multiAssign(assignments []string) error {
	sqlInsert := fmt.Sprintf(`
        INSERT INTO asignacion_cuartos (
            id,
            alumno_id,
            cuarto_id,
            convocatoria_id,
            fecha_asignacion,
            estado,
            observaciones,
            created_at,
            updated_at
        ) VALUES %s;
    `, strings.Join(assignments, ", "))

	rs, err := s.DB.Exec(sqlInsert)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("rows affected error")
	}

	return nil
}

func (s *sqlserver) getAllRoomAssignmentsByStudentIDANDSubmissionID(studentID, callID int64) ([]*RoomAssignment, error) {
	var ms []*RoomAssignment
	const sqlGetAll = `SELECT id, alumno_id, cuarto_id, convocatoria_id,
                      fecha_asignacion, estado, observaciones, created_at, updated_at
                      FROM asignacion_cuartos WHERE alumno_id = ? AND convocatoria_id = ?`

	err := s.DB.Select(&ms, sqlGetAll, studentID, callID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getSubmissionsByStudentID obtiene las convocatorias donde el estudiante tiene asignación
func (s *sqlserver) getSubmissionsByStudentID(studentID int64) ([]*StudentSubmission, error) {
	var submissions []*StudentSubmission
	const query = `
		SELECT DISTINCT
			c.id as convocatoria_id,
			c.nombre as convocatoria_nombre,
			c.fecha_inicio,
			c.fecha_fin
		FROM convocatorias c
		INNER JOIN asignacion_cuartos ac ON ac.convocatoria_id = c.id
		WHERE ac.alumno_id = ?
		ORDER BY c.fecha_inicio DESC
	`
	err := s.DB.Select(&submissions, query, studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*StudentSubmission{}, nil
		}
		return nil, err
	}
	return submissions, nil
}

// getAssignmentDetailByStudentAndSubmission obtiene el detalle completo de la asignación
func (s *sqlserver) getAssignmentDetailByStudentAndSubmission(studentID, submissionID int64) (*StudentAssignmentDetail, error) {
	var detail StudentAssignmentDetail
	const query = `
		SELECT
			c.id as convocatoria_id,
			c.nombre as convocatoria_nombre,
			c.fecha_inicio,
			c.fecha_fin,
			ac.id as asignacion_id,
			DATE_FORMAT(ac.fecha_asignacion, '%Y-%m-%d') as fecha_asignacion,
			ac.estado as asignacion_estado,
			ac.observaciones as asignacion_observaciones,
			r.id as residencia_id,
			r.nombre as residencia_nombre,
			r.genero as residencia_genero,
			r.direccion as residencia_direccion,
			r.description as residencia_descripcion,
			cu.id as cuarto_id,
			cu.numero as cuarto_numero,
			cu.piso as cuarto_piso,
			cu.capacidad as cuarto_capacidad
		FROM asignacion_cuartos ac
		INNER JOIN convocatorias c ON c.id = ac.convocatoria_id
		INNER JOIN cuartos cu ON cu.id = ac.cuarto_id
		INNER JOIN residencias r ON r.id = cu.residencia_id
		WHERE ac.alumno_id = ? AND ac.convocatoria_id = ?
	`
	err := s.DB.Get(&detail, query, studentID, submissionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}

// getRoommatesByRoomAndSubmission obtiene los compañeros de cuarto
func (s *sqlserver) getRoommatesByRoomAndSubmission(roomID string, submissionID, excludeStudentID int64) ([]*RoommateInfo, error) {
	var roommates []*RoommateInfo
	const query = `
		SELECT
			a.id,
			a.nombres,
			a.apellido_paterno,
			a.apellido_materno,
			a.codigo_estudiante,
			a.facultad,
			a.escuela_profesional,
			a.correo_institucional
		FROM asignacion_cuartos ac
		INNER JOIN alumnos a ON a.id = ac.alumno_id
		WHERE ac.cuarto_id = ?
		  AND ac.convocatoria_id = ?
		  AND ac.alumno_id != ?
		  AND ac.estado = 'activo'
	`
	err := s.DB.Select(&roommates, query, roomID, submissionID, excludeStudentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*RoommateInfo{}, nil
		}
		return nil, err
	}
	return roommates, nil
}
