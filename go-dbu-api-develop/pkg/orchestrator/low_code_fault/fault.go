package low_code_fault

import (
	"dbu-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type FaultLowCode struct {
	db   *sqlx.DB
	user *models.User
	txID string
}

func NewFault(db *sqlx.DB, user *models.User, txID string) *FaultLowCode {
	return &FaultLowCode{
		db:   db,
		user: user,
		txID: txID,
	}
}

// Ejemplo de método: obtener faltas por alumno
func (f *FaultLowCode) GetFaultsByAlumnoID(alumnoID string) ([]models.Fault, int, error) {
	query := `
		SELECT 
			id, alumno_id,convocatoria_id, servicio_id, observacion,  fuente_informacion, 
			fecha_falta, estado
		FROM 
			faltas
		WHERE 
			alumno_id = @p1
	`

	var faults []models.Fault
	err := f.db.Select(&faults, query, alumnoID)
	if err != nil {
		return nil, 99, err
	}

	return faults, 29, nil
}
func (f *FaultLowCode) GetStudentsByFaultLowCode() ([]models.Student, int, error) {
	query := `
		SELECT DISTINCT a.id, a.dni, a.nombres, a.apellidos
		FROM faults fa
		JOIN alumnos a ON a.id = fa.alumno_id
	`

	var students []models.Student
	err := f.db.Select(&students, query)
	if err != nil {
		return nil, 99, err
	}

	return students, 29, nil
}
