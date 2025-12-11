package otp

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type sqlserver struct {
	DB   *sqlx.DB
	TxID string
}

func newOTPSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

// createOTPToken crea un nuevo token OTP
func (s *sqlserver) createOTPToken(m *TokenOTPEstudiante) error {
	const sqlInsert = `INSERT INTO tokens_otp_estudiante
		(id, alumno_id, dni, codigo_otp, correo_destino, fecha_generacion, fecha_expiracion,
		 intentos_fallidos, estado, direccion_ip, created_at, updated_at)
		VALUES (:id, :alumno_id, :dni, :codigo_otp, :correo_destino, :fecha_generacion, :fecha_expiracion,
		 :intentos_fallidos, :estado, :direccion_ip, :created_at, :updated_at)`

	_, err := s.DB.NamedExec(sqlInsert, m)
	return err
}

// getOTPTokenByDNIAndCode obtiene un token OTP por DNI y código
func (s *sqlserver) getOTPTokenByDNIAndCode(dni, code string) (*TokenOTPEstudiante, error) {
	const sqlSelect = `SELECT id, alumno_id, dni, codigo_otp, correo_destino, fecha_generacion,
		fecha_expiracion, intentos_fallidos, estado, direccion_ip, created_at, updated_at
		FROM tokens_otp_estudiante
		WHERE dni = ? AND codigo_otp = ? AND estado = 'pendiente'
		ORDER BY fecha_generacion DESC
		LIMIT 1`

	var token TokenOTPEstudiante
	err := s.DB.Get(&token, sqlSelect, dni, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// getLatestOTPByDNI obtiene el último token OTP generado para un DNI
func (s *sqlserver) getLatestOTPByDNI(dni string) (*TokenOTPEstudiante, error) {
	const sqlSelect = `SELECT id, alumno_id, dni, codigo_otp, correo_destino, fecha_generacion,
		fecha_expiracion, intentos_fallidos, estado, direccion_ip, created_at, updated_at
		FROM tokens_otp_estudiante
		WHERE dni = ?
		ORDER BY fecha_generacion DESC
		LIMIT 1`

	var token TokenOTPEstudiante
	err := s.DB.Get(&token, sqlSelect, dni)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// updateOTPToken actualiza un token OTP
func (s *sqlserver) updateOTPToken(m *TokenOTPEstudiante) error {
	const sqlUpdate = `UPDATE tokens_otp_estudiante
		SET intentos_fallidos = :intentos_fallidos,
			estado = :estado,
			updated_at = :updated_at
		WHERE id = :id`

	_, err := s.DB.NamedExec(sqlUpdate, m)
	return err
}

// countRecentOTPByDNI cuenta OTPs recientes por DNI (para rate limiting)
func (s *sqlserver) countRecentOTPByDNI(dni string, minutes int) (int, error) {
	const sqlCount = `SELECT COUNT(*)
		FROM tokens_otp_estudiante
		WHERE dni = ? AND fecha_generacion > DATE_SUB(NOW(), INTERVAL ? MINUTE)`

	var count int
	err := s.DB.Get(&count, sqlCount, dni, minutes)
	return count, err
}

// getAlumnoByDNI obtiene información básica del alumno por DNI
func (s *sqlserver) getAlumnoByDNI(dni string) (*AlumnoBasicInfo, error) {
	const sqlSelect = `SELECT id, codigo_estudiante, DNI, nombres, apellido_paterno,
		apellido_materno, correo_institucional
		FROM alumnos
		WHERE DNI = ?
		LIMIT 1`

	var alumno AlumnoBasicInfo
	err := s.DB.Get(&alumno, sqlSelect, dni)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &alumno, nil
}

// hasActiveRoomAssignment verifica si el alumno tiene una asignación de cuarto activa
func (s *sqlserver) hasActiveRoomAssignment(alumnoID int64) (bool, error) {
	const sqlCount = `SELECT COUNT(*)
		FROM asignacion_cuartos
		WHERE alumno_id = ?`

	var count int
	err := s.DB.Get(&count, sqlCount, alumnoID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// createSession crea una nueva sesión
func (s *sqlserver) createSession(m *SesionEstudiante) error {
	const sqlInsert = `INSERT INTO sesiones_estudiante
		(id, alumno_id, token_jwt, direccion_ip, agente_usuario, fecha_login, fecha_expiracion,
		 fecha_ultimo_acceso, activo, created_at, updated_at)
		VALUES (:id, :alumno_id, :token_jwt, :direccion_ip, :agente_usuario, :fecha_login,
		 :fecha_expiracion, :fecha_ultimo_acceso, :activo, :created_at, :updated_at)`

	_, err := s.DB.NamedExec(sqlInsert, m)
	return err
}

// getSessionByID obtiene una sesión por ID
func (s *sqlserver) getSessionByID(sessionID string) (*SesionEstudiante, error) {
	const sqlSelect = `SELECT id, alumno_id, token_jwt, direccion_ip, agente_usuario,
		fecha_login, fecha_expiracion, fecha_ultimo_acceso, activo, created_at, updated_at
		FROM sesiones_estudiante
		WHERE id = ?`

	var session SesionEstudiante
	err := s.DB.Get(&session, sqlSelect, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// updateSessionLastAccess actualiza el último acceso de una sesión
func (s *sqlserver) updateSessionLastAccess(sessionID string) error {
	const sqlUpdate = `UPDATE sesiones_estudiante
		SET fecha_ultimo_acceso = NOW(), updated_at = NOW()
		WHERE id = ?`

	_, err := s.DB.Exec(sqlUpdate, sessionID)
	return err
}

// deactivateSession desactiva una sesión
func (s *sqlserver) deactivateSession(sessionID string) error {
	const sqlUpdate = `UPDATE sesiones_estudiante
		SET activo = false, updated_at = NOW()
		WHERE id = ?`

	_, err := s.DB.Exec(sqlUpdate, sessionID)
	return err
}

// Exported methods for external use

// GetOTPTokenByDNIAndCode obtiene un token OTP por DNI y código (exported)
func (s *sqlserver) GetOTPTokenByDNIAndCode(dni, code string) (*TokenOTPEstudiante, error) {
	return s.getOTPTokenByDNIAndCode(dni, code)
}

// GetLatestOTPByDNI obtiene el último token OTP generado para un DNI (exported)
func (s *sqlserver) GetLatestOTPByDNI(dni string) (*TokenOTPEstudiante, error) {
	return s.getLatestOTPByDNI(dni)
}

// GetAlumnoByDNI obtiene información básica del alumno por DNI (exported)
func (s *sqlserver) GetAlumnoByDNI(dni string) (*AlumnoBasicInfo, error) {
	return s.getAlumnoByDNI(dni)
}

// GetSessionByID obtiene una sesión por ID (exported)
func (s *sqlserver) GetSessionByID(sessionID string) (*SesionEstudiante, error) {
	return s.getSessionByID(sessionID)
}

// UpdateSessionLastAccess actualiza el último acceso de una sesión (exported)
func (s *sqlserver) UpdateSessionLastAccess(sessionID string) error {
	return s.updateSessionLastAccess(sessionID)
}

// DeactivateSession desactiva una sesión (exported)
func (s *sqlserver) DeactivateSession(sessionID string) error {
	return s.deactivateSession(sessionID)
}
