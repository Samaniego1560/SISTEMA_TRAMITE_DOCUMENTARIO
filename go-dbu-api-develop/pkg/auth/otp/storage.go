package otp

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

// OTPRepository interface para operaciones de tokens OTP
type OTPRepository interface {
	// Tokens OTP
	createOTPToken(m *TokenOTPEstudiante) error
	getOTPTokenByDNIAndCode(dni, code string) (*TokenOTPEstudiante, error)
	GetOTPTokenByDNIAndCode(dni, code string) (*TokenOTPEstudiante, error)
	getLatestOTPByDNI(dni string) (*TokenOTPEstudiante, error)
	GetLatestOTPByDNI(dni string) (*TokenOTPEstudiante, error)
	updateOTPToken(m *TokenOTPEstudiante) error
	countRecentOTPByDNI(dni string, minutes int) (int, error)

	// Alumnos
	getAlumnoByDNI(dni string) (*AlumnoBasicInfo, error)
	GetAlumnoByDNI(dni string) (*AlumnoBasicInfo, error)
	hasActiveRoomAssignment(alumnoID int64) (bool, error)

	// Sesiones
	createSession(m *SesionEstudiante) error
	GetSessionByID(sessionID string) (*SesionEstudiante, error)
	UpdateSessionLastAccess(sessionID string) error
	DeactivateSession(sessionID string) error
}

func FactoryStorage(db *sqlx.DB, txID string) OTPRepository {
	return newOTPSqlServerRepository(db, txID)
}
