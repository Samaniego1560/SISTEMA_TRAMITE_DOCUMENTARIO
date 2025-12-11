package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Fault struct {
	ID                string    `json:"id" valid:"uuid,required"`
	AlumnoID          int64     `json:"alumno_id" valid:"required"`
	ServicioId        int64     `json:"servicio_id" valid:"required"`
	ConvocatoriaId    int64     `json:"convocatoria_id" valid:"required"`
	FuenteInformacion string    `json:"fuente_informacion" valid:"required"`
	FechaFalta        time.Time `json:"fecha_falta" valid:"required"`
	Estado            string    `json:"estado" valid:"required"`
	Observacion       string    `json:"observacion" valid:"-"`
	Articulos         []string  `json:"articulos" valid:"-"`
	Incisos           []string  `json:"incisos" valid:"-"`
	UrlDocumentos     []string  `json:"url_documentos" valid:"-"`
}

// ValidFault valida la estructura de Fault
func (m *Fault) ValidFault() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Student representa la información del estudiante asociado a una falta
type StudentSanctionsResponse struct {
	Student StudentInfoResponse  `json:"student"`
	Faults  []FaultWithSanctions `json:"faults"`
}

type StudentInfoResponse struct {
	ID                 int64  `json:"id"`
	FullName           string `json:"full_name"`
	Code               string `json:"code"`
	ProfessionalSchool string `json:"professional_school"`
	Room               string `json:"room"`
	AdmissionDate      string `json:"admission_date"`
}

type FaultWithSanctions struct {
	ID                string         `json:"id"`
	ServicioId        int64          `json:"servicio_id"`
	ConvocatoriaId    int64          `json:"convocatoria_id"`
	FuenteInformacion string         `json:"fuente_informacion"`
	FechaFalta        time.Time      `json:"fecha_falta"`
	Estado            string         `json:"estado"`
	Observacion       string         `json:"observacion,omitempty"`
	Articulos         []string       `json:"articulos,omitempty"`
	Incisos           []string       `json:"incisos,omitempty"`
	UrlDocumentos     []string       `json:"url_documentos,omitempty"`
	Sanctions         []SanctionInfo `json:"sanctions"`
}

type SanctionInfo struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Tipo        string    `json:"tipo"`
	Duracion    int       `json:"duracion"`
	FechaInicio time.Time `json:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"`
	Estado      string    `json:"estado"`
	Observacion string    `json:"observacion,omitempty"`
	Revisada    bool      `json:"revisada"`
}
type Studentf struct {
	ID        string `db:"id" json:"id"`
	DNI       int64  `db:"dni" json:"dni"`
	Nombres   string `db:"nombres" json:"nombres"`
	Apellidos string `db:"apellidos" json:"apellidos"`
}
