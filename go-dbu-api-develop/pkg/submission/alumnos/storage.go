package alumnos

import (
	"github.com/jmoiron/sqlx"

	"dbu-api/internal/models"
)

type ServicesAlumnosRepository interface {
	getStudentsAcceptedBySubmission(id int64) ([]*Alumno, error)
	getStudentsAcceptedBySubmissionNewbie(id int64) ([]*Alumno, error)
	getStudentsByResidenceANDBySubmission(residenceID string, submissionID int64, page, limit int, filter string) ([]*StudentInformation, error)
	getTotalStudentsByResidenceANDBySubmission(residenceID string, submissionID int64, filter string) (int, error)
	getStudentsBySubmission(submissionID int, page, limit int, gender string, statusService string, departmentRequirementID int, departmentRequired bool) ([]*StudentInformationSubmission, error)
	getTotalStudentsBySubmission(submissionID int, gender string, statusService string, departmentRequirementID int, departmentRequired bool) (int, error)
	getStudentsBySubmissionExcel(submissionID int) ([]*models.StudentExcel, error)
	getStudentsByRooms(rooms []string, submissionID int64) ([]*StudentInformation, error)
	getStudentAcceptedBySubmission(submissionID, studentID int64) (*Alumno, error)
	getByID(id int64) (*Alumno, error)
	getByInstitutionalEmail(institutionalEmail string) (*Alumno, error)
	hasActiveRoomAssignment(alumnoID int64) (bool, error)
	createSession(m *SesionEstudiante) error
	updateSessionStatus(sessionID string, active bool) error
	getSessionByToken(tokenJWT string) (*SesionEstudiante, error)
	getSessionByID(sessionID string) (*SesionEstudiante, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesAlumnosRepository {
	return newConvocatoriasSqlServerRepository(db, user, txID)
}
