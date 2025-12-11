package faults

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Fault representa la entidad directamente mapeada a la tabla
type Fault struct {
	ID                 string     `db:"id" json:"id"`
	AlumnoID           int64      `db:"alumno_id" json:"alumno_id"`
	ConvocatoriaID     int64      `db:"convocatoria_id" json:"convocatoria_id"`
	ServicioID         int64      `db:"servicio_id" json:"servicio_id"`
	Observacion        string     `db:"observacion" json:"observacion"`
	FuenteInformacion  string     `db:"fuente_informacion" json:"fuente_informacion"`
	FechaFalta         time.Time  `db:"fecha_falta" json:"fecha_falta"`
	Estado             string     `db:"estado" json:"estado"`
	Apelable           bool       `db:"apelable" json:"apelable"`
	ApelacionDocumento *string    `db:"apelacion_documento" json:"apelacion_documento,omitempty"`
	MotivoResolucion   *string    `db:"motivo_resolucion" json:"motivo_resolucion,omitempty"`
	CreatedAt          *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// CreateFaultRequest estructura extendida para crear faltas con artículos e incisos
type CreateFaultRequest struct {
	AlumnoID          int64     `json:"alumno_id" valid:"required"`
	ServicioID        int64     `json:"servicio_id" valid:"required"`
	ConvocatoriaID    int64     `json:"convocatoria_id" valid:"required"`
	FuenteInformacion string    `json:"fuente_informacion" valid:"required"`
	FechaFalta        time.Time `json:"fecha_falta" valid:"required"`
	Estado            string    `json:"estado" valid:"in(registrada|notificada|apelada|resuelta)"`
	Observacion       string    `json:"observacion" valid:"required"`

	Articulos []string `json:"articulo_ids"`
	Incisos   []string `json:"inciso_ids"`
}

type FaultDocumento struct {
	ID        string    `db:"id" json:"id"`
	FaultID   string    `db:"falta_id" json:"falta_id"`
	URL       string    `db:"url" json:"url"`
	Archivo   []byte    `db:"archivo" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (m *CreateFaultRequest) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
