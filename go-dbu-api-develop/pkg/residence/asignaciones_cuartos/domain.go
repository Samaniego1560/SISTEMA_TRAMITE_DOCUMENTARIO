package asignaciones_cuartos

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type RoomAssignment struct {
	ID             string     `json:"id" db:"id" valid:"-"`
	StudentID      int64      `json:"alumno_id" db:"alumno_id" valid:"required"`
	RoomID         string     `json:"cuarto_id" db:"cuarto_id" valid:"required"`
	CallID         int64      `json:"convocatoria_id" db:"convocatoria_id" valid:"required"`
	AssignmentDate time.Time  `json:"fecha_asignacion" db:"fecha_asignacion" valid:"required"`
	Status         string     `json:"estado" db:"estado" valid:"in(activo|desocupado|suspendido|cancelado),required"`
	Observations   *string    `json:"observaciones" db:"observaciones"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
}

func NewRoomAssignment(id string, studentID int64, roomID string, callID int64, assignmentDate time.Time, status, observation string) *RoomAssignment {
	now := time.Now()
	return &RoomAssignment{
		ID:             id,
		StudentID:      studentID,
		RoomID:         roomID,
		CallID:         callID,
		AssignmentDate: assignmentDate,
		Status:         status,
		Observations:   &observation,
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}
}

func (m *RoomAssignment) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

// StudentSubmission representa una convocatoria donde el estudiante tiene asignación
type StudentSubmission struct {
	ConvocatoriaID     int64      `json:"convocatoria_id" db:"convocatoria_id"`
	ConvocatoriaNombre string     `json:"convocatoria_nombre" db:"convocatoria_nombre"`
	FechaInicio        *time.Time `json:"fecha_inicio" db:"fecha_inicio"`
	FechaFin           *time.Time `json:"fecha_fin" db:"fecha_fin"`
}

// StudentAssignmentDetail detalle completo de la asignación del estudiante
type StudentAssignmentDetail struct {
	// Convocatoria
	ConvocatoriaID     int64      `json:"convocatoria_id" db:"convocatoria_id"`
	ConvocatoriaNombre string     `json:"convocatoria_nombre" db:"convocatoria_nombre"`
	FechaInicio        *time.Time `json:"fecha_inicio" db:"fecha_inicio"`
	FechaFin           *time.Time `json:"fecha_fin" db:"fecha_fin"`

	// Asignación
	AsignacionID            string  `json:"asignacion_id" db:"asignacion_id"`
	FechaAsignacion         string  `json:"fecha_asignacion" db:"fecha_asignacion"`
	AsignacionEstado        string  `json:"asignacion_estado" db:"asignacion_estado"`
	AsignacionObservaciones *string `json:"asignacion_observaciones" db:"asignacion_observaciones"`

	// Residencia
	ResidenciaID          string  `json:"residencia_id" db:"residencia_id"`
	ResidenciaNombre      string  `json:"residencia_nombre" db:"residencia_nombre"`
	ResidenciaGenero      string  `json:"residencia_genero" db:"residencia_genero"`
	ResidenciaDireccion   string  `json:"residencia_direccion" db:"residencia_direccion"`
	ResidenciaDescripcion *string `json:"residencia_descripcion" db:"residencia_descripcion"`

	// Cuarto
	CuartoID        string `json:"cuarto_id" db:"cuarto_id"`
	CuartoNumero    int    `json:"cuarto_numero" db:"cuarto_numero"`
	CuartoPiso      int    `json:"cuarto_piso" db:"cuarto_piso"`
	CuartoCapacidad int    `json:"cuarto_capacidad" db:"cuarto_capacidad"`
}

// RoommateInfo información de un compañero de cuarto
type RoommateInfo struct {
	ID                  int64  `json:"id" db:"id"`
	Nombres             string `json:"nombres" db:"nombres"`
	ApellidoPaterno     string `json:"apellido_paterno" db:"apellido_paterno"`
	ApellidoMaterno     string `json:"apellido_materno" db:"apellido_materno"`
	CodigoEstudiante    string `json:"codigo_estudiante" db:"codigo_estudiante"`
	Facultad            string `json:"facultad" db:"facultad"`
	EscuelaProfesional  string `json:"escuela_profesional" db:"escuela_profesional"`
	CorreoInstitucional string `json:"correo_institucional" db:"correo_institucional"`
}
