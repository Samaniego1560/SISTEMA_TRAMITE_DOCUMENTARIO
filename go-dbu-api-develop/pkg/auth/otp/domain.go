package otp

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// TokenOTPEstudiante representa un token OTP para autenticación de estudiantes
type TokenOTPEstudiante struct {
	ID               string     `json:"id" db:"id"`
	AlumnoID         int64      `json:"alumno_id" db:"alumno_id"`
	DNI              string     `json:"dni" db:"dni"`
	CodigoOTP        string     `json:"codigo_otp" db:"codigo_otp"`
	CorreoDestino    string     `json:"correo_destino" db:"correo_destino"`
	FechaGeneracion  time.Time  `json:"fecha_generacion" db:"fecha_generacion"`
	FechaExpiracion  time.Time  `json:"fecha_expiracion" db:"fecha_expiracion"`
	IntentosFallidos int        `json:"intentos_fallidos" db:"intentos_fallidos"`
	Estado           string     `json:"estado" db:"estado" valid:"in(pendiente|usado|expirado|bloqueado)"`
	DireccionIP      *string    `json:"direccion_ip,omitempty" db:"direccion_ip"`
	CreatedAt        *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// SesionEstudiante representa una sesión activa de un estudiante
type SesionEstudiante struct {
	ID                string     `json:"id" db:"id"`
	AlumnoID          int64      `json:"alumno_id" db:"alumno_id"`
	TokenJWT          string     `json:"token_jwt" db:"token_jwt"`
	DireccionIP       *string    `json:"direccion_ip,omitempty" db:"direccion_ip"`
	AgenteUsuario     *string    `json:"agente_usuario,omitempty" db:"agente_usuario"`
	FechaLogin        time.Time  `json:"fecha_login" db:"fecha_login"`
	FechaExpiracion   time.Time  `json:"fecha_expiracion" db:"fecha_expiracion"`
	FechaUltimoAcceso *time.Time `json:"fecha_ultimo_acceso,omitempty" db:"fecha_ultimo_acceso"`
	Activo            bool       `json:"activo" db:"activo"`
	CreatedAt         *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// AlumnoBasicInfo información básica del alumno para autenticación
type AlumnoBasicInfo struct {
	ID                  int64  `json:"id" db:"id"`
	CodigoEstudiante    string `json:"codigo_estudiante" db:"codigo_estudiante"`
	DNI                 string `json:"DNI" db:"DNI"`
	Nombres             string `json:"nombres" db:"nombres"`
	ApellidoPaterno     string `json:"apellido_paterno" db:"apellido_paterno"`
	ApellidoMaterno     string `json:"apellido_materno" db:"apellido_materno"`
	CorreoInstitucional string `json:"correo_institucional" db:"correo_institucional"`
}

func (m *TokenOTPEstudiante) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (m *SesionEstudiante) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
