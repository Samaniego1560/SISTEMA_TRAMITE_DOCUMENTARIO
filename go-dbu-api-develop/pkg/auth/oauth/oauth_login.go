package auth_orchestrator

import (
	"context"
	"dbu-api/internal/logger"
	"dbu-api/internal/middleware"
	oauth "dbu-api/pkg/auth/oauth/providers"
	"dbu-api/pkg/submission/alumnos"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// OAuthLoginService maneja la lógica de negocio para autenticación OAuth
type OAuthLoginService struct {
	db       *sqlx.DB
	txID     string
	provider oauth.OAuthProvider
}

// PortsServerOAuthLogin define la interfaz para el servicio de OAuth login
type PortsServerOAuthLogin interface {
	// GetAuthorizationURL genera la URL de autorización del proveedor OAuth
	GetAuthorizationURL(state string) string

	// ProcessCallback procesa el callback de OAuth y retorna los tokens JWT
	ProcessCallback(ctx context.Context, code string, baseURL string, ipAddress string, userAgent string) (token string, code_response int, err error)

	// GetProviderName retorna el nombre del proveedor OAuth
	GetProviderName() string

	// UpdateSessionStatus actualiza el estado de una sesión
	UpdateSessionStatus(sessionID string, active bool) error

	// GetSessionByToken obtiene una sesión por su token JWT
	GetSessionByToken(tokenJWT string) (*alumnos.SesionEstudiante, int, error)
}

// NewOAuthLogin crea una nueva instancia del servicio OAuth login
func NewOAuthLogin(db *sqlx.DB, provider oauth.OAuthProvider, txID string) PortsServerOAuthLogin {
	return &OAuthLoginService{
		db:       db,
		txID:     txID,
		provider: provider,
	}
}

// GetAuthorizationURL genera la URL de autorización del proveedor OAuth
func (s *OAuthLoginService) GetAuthorizationURL(state string) string {
	return s.provider.GetAuthURL(state)
}

// ProcessCallback procesa el callback de OAuth
// 1. Intercambia el código de autorización por información del usuario
// 2. Busca el alumno en la base de datos por correo institucional
// 3. Verifica que tenga asignación de cuarto activa
// 4. Genera JWT y refresh token
// 5. Crea sesión de estudiante
// 6. Retorna los tokens
func (s *OAuthLoginService) ProcessCallback(ctx context.Context, code string, baseURL string, ipAddress string, userAgent string) (string, int, error) {
	// Intercambiar código por información del usuario
	userInfo, err := s.provider.ExchangeCode(ctx, code)
	if err != nil {
		logger.Error.Printf("%s - failed to exchange OAuth code: %v", s.txID, err)
		return "", 54, fmt.Errorf("failed to authenticate with %s: %w", s.provider.GetProviderName(), err)
	}

	// Log de información del usuario recibida
	logger.Info.Printf("%s - OAuth user authenticated: institutional_email=%s, email=%s, name=%s, provider=%s",
		s.txID, userInfo.InstitutionalEmail, userInfo.Email, userInfo.Name, userInfo.Provider)

	// Buscar alumno en la base de datos por correo institucional
	alumnosRepo := alumnos.FactoryStorage(s.db, nil, s.txID)
	alumnosService := alumnos.NewAlumnosService(alumnosRepo, nil, s.txID)

	// Normalizar correo institucional: trim y lowercase
	normalizedEmail := strings.ToLower(strings.TrimSpace(userInfo.InstitutionalEmail))
	student, _, err := alumnosService.GetStudentByInstitutionalEmail(normalizedEmail)
	if err != nil {
		logger.Error.Printf("%s - couldn't get student by institutional email: %v", s.txID, err)
		return "", 99, fmt.Errorf("student not found in database")
	}

	if student == nil {
		logger.Error.Printf("%s - student not found for institutional email: %s", s.txID, userInfo.InstitutionalEmail)
		return "", 99, fmt.Errorf("student with institutional email %s is not registered in the system", userInfo.InstitutionalEmail)
	}

	// Verificar que el estudiante tenga una asignación de cuarto activa
	hasAssignment, err := alumnosService.HasActiveRoomAssignment(student.ID)
	if err != nil {
		logger.Error.Printf("%s - error verifying room assignment: %v", s.txID, err)
		return "", 98, fmt.Errorf("error verifying room assignment")
	}
	if !hasAssignment {
		logger.Warning.Printf("%s - student without active room assignment: %s (ID: %d)", s.txID, student.CodigoEstudiante, student.ID)
		return "", 97, fmt.Errorf("you don't have an active room assignment")
	}

	sessionID := uuid.New().String()

	// Generar JWT Token con el ID del estudiante
	token, err := middleware.CreateStudentJWT(student.ID, student.DNI, sessionID)
	if err != nil {
		logger.Error.Printf("%s - error generating JWT token: %v", s.txID, err)
		return "", 54, fmt.Errorf("error generating token: %w", err)
	}

	// Crear sesión de estudiante
	_, err = alumnosService.CreateSession(sessionID, student.ID, token, ipAddress, userAgent, false)
	if err != nil {
		logger.Error.Printf("%s - error creating session: %v", s.txID, err)
		// No retornamos error aquí, la sesión es opcional
	}

	logger.Info.Printf("%s - OAuth login successful for student: %s (ID: %d)", s.txID, student.CodigoEstudiante, student.ID)

	return token, 220, nil
}

// GetProviderName retorna el nombre del proveedor OAuth
func (s *OAuthLoginService) GetProviderName() string {
	return s.provider.GetProviderName()
}

// UpdateSessionStatus actualiza el estado de una sesión
func (s *OAuthLoginService) UpdateSessionStatus(sessionID string, active bool) error {
	alumnosRepo := alumnos.FactoryStorage(s.db, nil, s.txID)
	alumnosService := alumnos.NewAlumnosService(alumnosRepo, nil, s.txID)

	return alumnosService.UpdateSessionStatus(sessionID, active)
}

// GetSessionByToken obtiene una sesión por su token JWT
func (s *OAuthLoginService) GetSessionByToken(tokenJWT string) (*alumnos.SesionEstudiante, int, error) {
	alumnosRepo := alumnos.FactoryStorage(s.db, nil, s.txID)
	alumnosService := alumnos.NewAlumnosService(alumnosRepo, nil, s.txID)

	return alumnosService.GetSessionByToken(tokenJWT)
}
