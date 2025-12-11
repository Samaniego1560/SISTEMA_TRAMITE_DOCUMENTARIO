package asignaciones_cuartos

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"time"

	"dbu-api/internal/logger"
	"dbu-api/internal/models"
)

type RoomAssignmentServer interface {
	CreateRoomAssignment(studentID int64, roomID string, callID int64, assignmentDate time.Time, status string) (*RoomAssignment, int, error)
	UpdateRoomAssignment(id string, studentID int64, roomID string, callID int64, assignmentDate time.Time, status string) (*RoomAssignment, int, error)
	DeleteRoomAssignment(studentID int64, roomID string, callID int64, status, observation string) (int, error)
	GetRoomAssignmentByID(id string) (*RoomAssignment, int, error)
	GetAllRoomAssignments() ([]*RoomAssignment, error)
	GetAllRoomAssignmentsByStudentIDANDSubmissionID(studentID, callID int64) ([]*RoomAssignment, error)
	GetRoomAssignmentByRoomIDSubmissionID(roomID string, submissionID int64) ([]*RoomAssignment, error)
	MultiAssignRoom(submission int64, assignments []models.RoomGPT) (int, error)

	// Métodos para consulta de estudiantes
	GetSubmissionsByStudent(studentID int64) ([]*StudentSubmission, int, error)
	GetAssignmentDetailForStudent(studentID, submissionID int64) (*StudentAssignmentDetail, []*RoommateInfo, int, error)
}

type service struct {
	repository RoomAssignmentRepository
	user       *models.User
	txID       string
}

func NewRoomAssignmentService(repository RoomAssignmentRepository, user *models.User, txID string) RoomAssignmentServer {
	return &service{repository: repository, user: user, txID: txID}
}

func (s *service) CreateRoomAssignment(studentID int64, roomID string, callID int64, assignmentDate time.Time, status string) (*RoomAssignment, int, error) {
	m := NewRoomAssignment(uuid.New().String(), studentID, roomID, callID, assignmentDate, status, "")
	if valid, err := m.Valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.create(m); err != nil {
		if err.Error() == "rows affected error" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create RoomAssignment:", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateRoomAssignment(id string, studentID int64, roomID string, callID int64, assignmentDate time.Time, status string) (*RoomAssignment, int, error) {
	m := NewRoomAssignment(id, studentID, roomID, callID, assignmentDate, status, "")
	if id == "" {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.Valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update RoomAssignment:", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteRoomAssignment(studentID int64, roomID string, callID int64, status, observation string) (int, error) {
	m := NewRoomAssignment(uuid.New().String(), studentID, roomID, callID, time.Now(), status, observation)
	if valid, err := m.Valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return 15, err
	}

	if err := s.repository.delete(m); err != nil {
		if err.Error() == "rows affected error" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetRoomAssignmentByID(id string) (*RoomAssignment, int, error) {
	if id == "" {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllRoomAssignments() ([]*RoomAssignment, error) {
	return s.repository.getAll()
}

func (s *service) GetRoomAssignmentByRoomIDSubmissionID(roomID string, submissionID int64) ([]*RoomAssignment, error) {
	if !govalidator.IsUUID(roomID) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, fmt.Errorf("id is required")
	}
	m, err := s.repository.getRoomAssignmentByRoomIDSubmissionID(roomID, submissionID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID row:", err)
		return nil, err
	}
	if m == nil {
		return []*RoomAssignment{}, nil
	}
	return m, nil
}

func (s *service) MultiAssignRoom(submission int64, assignments []models.RoomGPT) (int, error) {
	var values []string
	for _, assignment := range assignments {
		for _, student := range assignment.Students {
			values = append(values, fmt.Sprintf(
				"('%s', %d, '%s', %d, NOW(), 'activo', NULL, NOW(), NOW())",
				uuid.New().String(),
				student.StudentId,
				assignment.RoomId,
				submission,
			))
		}
	}

	if err := s.repository.multiAssign(values); err != nil {
		if err.Error() == "rows affected error" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetAllRoomAssignmentsByStudentIDANDSubmissionID(studentID, callID int64) ([]*RoomAssignment, error) {
	return s.repository.getAllRoomAssignmentsByStudentIDANDSubmissionID(studentID, callID)
}

// GetSubmissionsByStudent obtiene las convocatorias donde el estudiante tiene asignaciones
func (s *service) GetSubmissionsByStudent(studentID int64) ([]*StudentSubmission, int, error) {
	submissions, err := s.repository.getSubmissionsByStudentID(studentID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get submissions:", err)
		return nil, 15, err
	}
	if submissions == nil || len(submissions) == 0 {
		return []*StudentSubmission{}, 214, nil
	}
	return submissions, 214, nil
}

// GetAssignmentDetailForStudent obtiene el detalle completo de la asignación del estudiante
func (s *service) GetAssignmentDetailForStudent(studentID, submissionID int64) (*StudentAssignmentDetail, []*RoommateInfo, int, error) {
	detail, err := s.repository.getAssignmentDetailByStudentAndSubmission(studentID, submissionID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get assignment detail:", err)
		return nil, nil, 15, err
	}
	if detail == nil {
		logger.Warning.Println(s.txID, " - no assignment found for student:", studentID, "submission:", submissionID)
		return nil, nil, 4, fmt.Errorf("no assignment found")
	}

	// Obtener compañeros de cuarto
	roommates, err := s.repository.getRoommatesByRoomAndSubmission(detail.CuartoID, submissionID, studentID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get roommates:", err)
		// No retornamos error crítico, solo array vacío
		return detail, []*RoommateInfo{}, 214, nil
	}

	return detail, roommates, 214, nil
}
