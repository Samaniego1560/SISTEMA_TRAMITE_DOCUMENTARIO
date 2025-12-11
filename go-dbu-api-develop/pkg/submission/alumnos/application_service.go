package alumnos

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"
	"time"
)

type PortsServerAlumnos interface {
	GetStudentsAcceptedBySubmission(id int64) ([]*Alumno, error)
	GetStudentsAcceptedBySubmissionNewbie(id int64) ([]*Alumno, error)
	GetStudentsByResidenceANDBySubmission(residenceID string, submissionID int64, page, limit int, filter string) ([]*StudentInformation, error)
	GetTotalStudentsByResidenceANDBySubmission(residenceID string, submissionID int64, filter string) (int, error)
	GetStudentsBySubmission(submissionID int, page, limit int, gender string, statusService string, departmentRequirementID int, departmentRequired bool) ([]*StudentInformationSubmission, error)
	GetTotalStudentsBySubmission(submissionID int, gender string, statusService string, departmentRequirementID int, departmentRequired bool) (int, error)
	GetStudentsBySubmissionExcel(submissionID int) ([]*models.StudentExcel, error)
	GetStudentsByRooms(rooms []string, submissionID int64) ([]*StudentInformation, error)
	GetStudentAcceptedBySubmission(submissionID, studentID int64) (*Alumno, int, error)
	GetStudentProfile(studentID int64) (*Alumno, int, error)
	GetStudentByInstitutionalEmail(institutionalEmail string) (*Alumno, int, error)

	// Métodos de autenticación compartidos
	HasActiveRoomAssignment(alumnoID int64) (bool, error)
	CreateSession(session string, alumnoID int64, tokenJWT, ipAddress, userAgent string, active bool) (*SesionEstudiante, error)
	UpdateSessionStatus(sessionID string, active bool) error
	GetSessionByToken(tokenJWT string) (*SesionEstudiante, int, error)
}

type service struct {
	repository ServicesAlumnosRepository
	user       *models.User
	txID       string
}

func NewAlumnosService(repository ServicesAlumnosRepository, user *models.User, TxID string) PortsServerAlumnos {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) GetStudentsAcceptedBySubmission(id int64) ([]*Alumno, error) {
	return s.repository.getStudentsAcceptedBySubmission(id)
}

func (s *service) GetStudentsAcceptedBySubmissionNewbie(id int64) ([]*Alumno, error) {
	return s.repository.getStudentsAcceptedBySubmissionNewbie(id)
}

func (s *service) GetStudentsByResidenceANDBySubmission(residenceID string, submissionID int64, page, limit int, filter string) ([]*StudentInformation, error) {
	return s.repository.getStudentsByResidenceANDBySubmission(residenceID, submissionID, page, limit, filter)
}

func (s *service) GetTotalStudentsByResidenceANDBySubmission(residenceID string, submissionID int64, filter string) (int, error) {
	return s.repository.getTotalStudentsByResidenceANDBySubmission(residenceID, submissionID, filter)
}

func (s *service) GetStudentsBySubmission(submissionID int, page, limit int, gender string, statusService string, departmentRequirementID int, departmentRequired bool) ([]*StudentInformationSubmission, error) {
	return s.repository.getStudentsBySubmission(submissionID, page, limit, gender, statusService, departmentRequirementID, departmentRequired)
}

func (s *service) GetTotalStudentsBySubmission(submissionID int, gender string, statusService string, departmentRequirementID int, departmentRequired bool) (int, error) {
	return s.repository.getTotalStudentsBySubmission(submissionID, gender, statusService, departmentRequirementID, departmentRequired)
}

func (s *service) GetStudentsBySubmissionExcel(submissionID int) ([]*models.StudentExcel, error) {
	return s.repository.getStudentsBySubmissionExcel(submissionID)
}

func (s *service) GetStudentsByRooms(rooms []string, submissionID int64) ([]*StudentInformation, error) {
	return s.repository.getStudentsByRooms(rooms, submissionID)
}

func (s *service) GetStudentAcceptedBySubmission(submissionID, studentID int64) (*Alumno, int, error) {
	if submissionID < 1 || studentID < 1 {
		logger.Error.Println(s.txID, " - couldn't meet validations:")
		return nil, 15, fmt.Errorf("couldn't meet validations")
	}
	m, err := s.repository.getStudentAcceptedBySubmission(submissionID, studentID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetStudentProfile(studentID int64) (*Alumno, int, error) {
	if studentID < 1 {
		logger.Error.Println(s.txID, " - couldn't meet validations: invalid student ID")
		return nil, 15, fmt.Errorf("couldn't meet validations: invalid student ID")
	}

	m, err := s.repository.getByID(studentID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get student profile:", err)
		return nil, 3, err
	}

	if m == nil {
		logger.Error.Println(s.txID, " - student not found")
		return nil, 4, fmt.Errorf("student not found")
	}

	return m, 29, nil
}

func (s *service) GetStudentByInstitutionalEmail(institutionalEmail string) (*Alumno, int, error) {
	if institutionalEmail == "" {
		logger.Error.Println(s.txID, " - couldn't meet validations: institutional email is required")
		return nil, 15, fmt.Errorf("couldn't meet validations: institutional email is required")
	}

	m, err := s.repository.getByInstitutionalEmail(institutionalEmail)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get student by institutional email:", err)
		return nil, 22, err
	}

	if m == nil {
		logger.Error.Println(s.txID, " - student not found for institutional email:", institutionalEmail)
		return nil, 22, fmt.Errorf("student not found")
	}

	return m, 29, nil
}

// HasActiveRoomAssignment verifica si el alumno tiene una asignación de cuarto activa
func (s *service) HasActiveRoomAssignment(alumnoID int64) (bool, error) {
	if alumnoID < 1 {
		logger.Error.Println(s.txID, " - couldn't meet validations: invalid student ID")
		return false, fmt.Errorf("couldn't meet validations: invalid student ID")
	}

	hasAssignment, err := s.repository.hasActiveRoomAssignment(alumnoID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't verify room assignment:", err)
		return false, err
	}

	return hasAssignment, nil
}

// CreateSession crea una nueva sesión para un estudiante
func (s *service) CreateSession(sessionID string, alumnoID int64, tokenJWT, ipAddress, userAgent string, active bool) (*SesionEstudiante, error) {
	if alumnoID < 1 {
		logger.Error.Println(s.txID, " - couldn't meet validations: invalid student ID")
		return nil, fmt.Errorf("couldn't meet validations: invalid student ID")
	}

	if tokenJWT == "" {
		logger.Error.Println(s.txID, " - couldn't meet validations: token JWT is required")
		return nil, fmt.Errorf("couldn't meet validations: token JWT is required")
	}

	now := time.Now()

	session := &SesionEstudiante{
		ID:                sessionID,
		AlumnoID:          alumnoID,
		TokenJWT:          tokenJWT,
		DireccionIP:       &ipAddress,
		AgenteUsuario:     &userAgent,
		FechaLogin:        now,
		FechaExpiracion:   now.Add(24 * time.Hour), // 24 horas de expiración
		FechaUltimoAcceso: &now,
		Activo:            active,
		CreatedAt:         &now,
		UpdatedAt:         &now,
	}

	if err := s.repository.createSession(session); err != nil {
		logger.Error.Println(s.txID, " - couldn't create session:", err)
		return nil, err
	}

	logger.Info.Printf("%s - session created successfully for student ID: %d, session_id: %s", s.txID, alumnoID, session.ID)
	return session, nil
}

// UpdateSessionStatus actualiza el estado de una sesión
func (s *service) UpdateSessionStatus(sessionID string, active bool) error {
	if sessionID == "" {
		logger.Error.Println(s.txID, " - couldn't meet validations: session ID is required")
		return fmt.Errorf("couldn't meet validations: session ID is required")
	}

	// Obtener la sesión actual para verificar su estado
	currentSession, err := s.repository.getSessionByID(sessionID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get session by ID:", err)
		return err
	}

	if currentSession == nil {
		logger.Error.Println(s.txID, " - session not found:", sessionID)
		return fmt.Errorf("session not found")
	}

	// Verificar si el estado ya es el deseado
	if currentSession.Activo == active {
		logger.Info.Printf("%s - session status is already %v, no update needed: session_id=%s", s.txID, active, sessionID)
		return nil
	}

	// Actualizar el estado
	if err := s.repository.updateSessionStatus(sessionID, active); err != nil {
		logger.Error.Println(s.txID, " - couldn't update session status:", err)
		return err
	}

	logger.Info.Printf("%s - session status updated successfully: session_id=%s, active=%v", s.txID, sessionID, active)
	return nil
}

// GetSessionByToken obtiene una sesión por su token JWT
func (s *service) GetSessionByToken(tokenJWT string) (*SesionEstudiante, int, error) {
	if tokenJWT == "" {
		logger.Error.Println(s.txID, " - couldn't meet validations: token JWT is required")
		return nil, 15, fmt.Errorf("couldn't meet validations: token JWT is required")
	}

	session, err := s.repository.getSessionByToken(tokenJWT)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get session by token:", err)
		return nil, 22, err
	}

	if session == nil {
		logger.Error.Println(s.txID, " - session not found for token")
		return nil, 22, fmt.Errorf("session not found")
	}

	return session, 29, nil
}
