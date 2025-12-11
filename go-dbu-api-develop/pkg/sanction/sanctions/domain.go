package sanctions

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
)

type Sanction struct {
	ID          string    `json:"id" valid:"uuid,required" db:"id"`
	FaultID     string    `json:"fault_id" valid:"required" db:"fault_id"`
	Tipo        string    `json:"tipo" valid:"required" db:"tipo_sancion"`
	Duracion    int       `json:"duracion" valid:"required" db:"duracion"`
	FechaInicio time.Time `json:"fecha_inicio" valid:"required" db:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin" valid:"required" db:"fecha_fin"`
	Estado      string    `json:"estado" valid:"required" db:"estado"`
	Observacion string    `json:"observacion" db:"observacion"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewSanction(id string, faultID string, tipo string, duracion int, fechaInicio time.Time, fechaFin time.Time, estado string, observacion string) *Sanction {
	return &Sanction{
		ID:          id,
		FaultID:     faultID,
		Tipo:        tipo,
		Duracion:    duracion,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
		Estado:      estado,
		Observacion: observacion,
	}
}

func (s *Sanction) valid() (bool, error) {
	if !govalidator.IsUUID(s.ID) {
		return false, fmt.Errorf("invalid UUID format")
	}
	return true, nil
}
