package otp

import (
	"crypto/rand"
	"dbu-api/internal/logger"
	"dbu-api/internal/middleware"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// SMTPService interfaz para evitar dependencia circular
type SMTPService interface {
	SendOTPCode(correo, nombreEstudiante, codigoOTP string, minutosExpiracion int) error
}

const (
	OTPExpirationMinutes   = 90
	MaxOTPAttemptsPerHour  = 3
	MaxFailedAttempts      = 3
	SessionExpirationHours = 24
)

type OTPService interface {
	RequestOTP(dni, ipAddress string) (*TokenOTPEstudiante, error)
	VerifyOTP(dni, code, ipAddress, userAgent string) (string, *AlumnoBasicInfo, error)
	Logout(sessionID string) error
}

type service struct {
	repository  OTPRepository
	smtpService SMTPService
	txID        string
}

func NewOTPService(repository OTPRepository, smtpService SMTPService, txID string) OTPService {
	return &service{
		repository:  repository,
		smtpService: smtpService,
		txID:        txID,
	}
}

// RequestOTP genera y envía un código OTP al estudiante
func (s *service) RequestOTP(dni, ipAddress string) (*TokenOTPEstudiante, error) {
	// 1. Verificar que el estudiante existe
	alumno, err := s.repository.getAlumnoByDNI(dni)
	if err != nil {
		logger.Error.Printf("%s - error al buscar alumno: %v", s.txID, err)
		return nil, fmt.Errorf("error al buscar estudiante")
	}
	if alumno == nil {
		logger.Warning.Printf("%s - estudiante no encontrado con DNI: %s", s.txID, dni)
		return nil, fmt.Errorf("estudiante no encontrado")
	}

	// 2. Verificar que el estudiante tiene una asignación de cuarto activa
	hasAssignment, err := s.repository.hasActiveRoomAssignment(alumno.ID)
	if err != nil {
		logger.Error.Printf("%s - error al verificar asignación de cuarto: %v", s.txID, err)
		return nil, fmt.Errorf("error al verificar asignación")
	}
	if !hasAssignment {
		logger.Warning.Printf("%s - estudiante sin asignación de cuarto con DNI: %s", s.txID, dni)
		return nil, fmt.Errorf("no tienes una asignación de cuarto activa")
	}

	// 3. Rate limiting: verificar que no haya solicitado muchos OTPs recientemente
	count, err := s.repository.countRecentOTPByDNI(dni, 60) // últimos 60 minutos
	if err != nil {
		logger.Error.Printf("%s - error al contar OTPs recientes: %v", s.txID, err)
		return nil, fmt.Errorf("error al procesar solicitud")
	}
	if count >= MaxOTPAttemptsPerHour {
		logger.Warning.Printf("%s - demasiados intentos de OTP para DNI: %s", s.txID, dni)
		return nil, fmt.Errorf("demasiados intentos. Intenta nuevamente en una hora")
	}

	// 4. Generar código OTP de 6 dígitos
	otpCode, err := generateOTPCode()
	if err != nil {
		logger.Error.Printf("%s - error al generar código OTP: %v", s.txID, err)
		return nil, fmt.Errorf("error al generar código")
	}

	// 5. Crear token OTP
	now := time.Now()
	token := &TokenOTPEstudiante{
		ID:               uuid.New().String(),
		AlumnoID:         alumno.ID,
		DNI:              dni,
		CodigoOTP:        otpCode,
		CorreoDestino:    alumno.CorreoInstitucional,
		FechaGeneracion:  now,
		FechaExpiracion:  now.Add(OTPExpirationMinutes * time.Second),
		IntentosFallidos: 0,
		Estado:           "pendiente",
		DireccionIP:      &ipAddress,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}

	// 6. Guardar en base de datos
	if err := s.repository.createOTPToken(token); err != nil {
		logger.Error.Printf("%s - error al guardar token OTP: %v", s.txID, err)
		return nil, fmt.Errorf("error al procesar solicitud")
	}

	logger.Info.Printf("%s - OTP generado para DNI: %s, alumno_id: %d", s.txID, dni, alumno.ID)

	token.CorreoDestino = maskEmail(token.CorreoDestino)
	// 7. Enviar email con código OTP
	nombreCompleto := alumno.Nombres + " " + alumno.ApellidoPaterno
	if err := s.smtpService.SendOTPCode(alumno.CorreoInstitucional, nombreCompleto, otpCode, OTPExpirationMinutes); err != nil {
		logger.Error.Printf("%s - error al enviar email OTP: %v", s.txID, err)
		// No retornamos error, el OTP ya fue guardado
	}

	// 8. Retornar token con correo parcialmente oculto
	return token, nil
}

// VerifyOTP verifica el código OTP y genera un JWT
func (s *service) VerifyOTP(dni, code, ipAddress, userAgent string) (string, *AlumnoBasicInfo, error) {
	// 1. Buscar token OTP
	token, err := s.repository.getOTPTokenByDNIAndCode(dni, code)
	if err != nil {
		logger.Error.Printf("%s - error al buscar token OTP: %v", s.txID, err)
		return "", nil, fmt.Errorf("error al verificar código")
	}
	if token == nil {
		logger.Warning.Printf("%s - código OTP no encontrado o inválido para DNI: %s", s.txID, dni)
		return "", nil, fmt.Errorf("código incorrecto o expirado")
	}

	// 2. Verificar que no esté bloqueado
	if token.Estado == "bloqueado" {
		logger.Warning.Printf("%s - token OTP bloqueado para DNI: %s", s.txID, dni)
		return "", nil, fmt.Errorf("código bloqueado por múltiples intentos fallidos")
	}

	// 3. Verificar expiración
	if time.Now().After(token.FechaExpiracion) {
		logger.Warning.Printf("%s - token OTP expirado para DNI: %s", s.txID, dni)
		// Marcar como expirado
		token.Estado = "expirado"
		now := time.Now()
		token.UpdatedAt = &now
		_ = s.repository.updateOTPToken(token)
		return "", nil, fmt.Errorf("código expirado. Solicita uno nuevo")
	}

	// 4. Verificar intentos fallidos
	if token.IntentosFallidos >= MaxFailedAttempts {
		logger.Warning.Printf("%s - demasiados intentos fallidos para DNI: %s", s.txID, dni)
		token.Estado = "bloqueado"
		now := time.Now()
		token.UpdatedAt = &now
		_ = s.repository.updateOTPToken(token)
		return "", nil, fmt.Errorf("código bloqueado por múltiples intentos fallidos")
	}

	// 5. Marcar token como usado
	token.Estado = "usado"
	now := time.Now()
	token.UpdatedAt = &now
	if err := s.repository.updateOTPToken(token); err != nil {
		logger.Error.Printf("%s - error al actualizar token OTP: %v", s.txID, err)
		return "", nil, fmt.Errorf("error al procesar verificación")
	}

	// 6. Obtener datos del alumno
	alumno, err := s.repository.getAlumnoByDNI(dni)
	if err != nil || alumno == nil {
		logger.Error.Printf("%s - error al obtener datos del alumno: %v", s.txID, err)
		return "", nil, fmt.Errorf("error al obtener datos del estudiante")
	}

	// 7. Generar JWT
	sessionID := uuid.New().String()
	jwtToken, err := middleware.CreateStudentJWT(alumno.ID, dni, sessionID)
	if err != nil {
		logger.Error.Printf("%s - error al crear JWT: %v", s.txID, err)
		return "", nil, fmt.Errorf("error al generar token de sesión")
	}

	// 8. Crear sesión
	session := &SesionEstudiante{
		ID:                sessionID,
		AlumnoID:          alumno.ID,
		TokenJWT:          jwtToken,
		DireccionIP:       &ipAddress,
		AgenteUsuario:     &userAgent,
		FechaLogin:        time.Now(),
		FechaExpiracion:   time.Now().Add(SessionExpirationHours * time.Hour),
		FechaUltimoAcceso: &now,
		Activo:            true,
		CreatedAt:         &now,
		UpdatedAt:         &now,
	}

	if err := s.repository.createSession(session); err != nil {
		logger.Error.Printf("%s - error al crear sesión: %v", s.txID, err)
		return "", nil, fmt.Errorf("error al crear sesión")
	}

	logger.Info.Printf("%s - sesión creada exitosamente para alumno_id: %d, session_id: %s", s.txID, alumno.ID, sessionID)

	return jwtToken, alumno, nil
}

// Logout desactiva una sesión
func (s *service) Logout(sessionID string) error {
	if err := s.repository.DeactivateSession(sessionID); err != nil {
		logger.Error.Printf("%s - error al cerrar sesión: %v", s.txID, err)
		return fmt.Errorf("error al cerrar sesión")
	}
	logger.Info.Printf("%s - sesión cerrada: %s", s.txID, sessionID)
	return nil
}

// generateOTPCode genera un código OTP de 6 dígitos
func generateOTPCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// maskEmail oculta parcialmente un correo electrónico
func maskEmail(email string) string {
	if len(email) < 3 {
		return "***"
	}

	atIndex := -1
	for i, char := range email {
		if char == '@' {
			atIndex = i
			break
		}
	}

	if atIndex == -1 {
		return "***"
	}

	visible := 1
	if atIndex > 3 {
		visible = 2
	}

	masked := email[:visible] + "***" + email[atIndex:]
	return masked
}
