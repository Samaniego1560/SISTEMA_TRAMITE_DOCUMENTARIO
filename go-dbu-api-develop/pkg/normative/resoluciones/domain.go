package resolutions

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Resolution struct {
	ID           string     `db:"id" json:"id" valid:"-"`
	Nombre       string     `db:"nombre" json:"nombre" valid:"required"`
	Descripcion  string     `db:"descripcion" json:"descripcion" valid:"required"`
	Servicio_id  string     `db:"servicio_id" json:"servicio_id" valid:"required"`
	Ruta_archivo string     `db:"ruta_archivo" json:"ruta_archivo" valid:"required"`
	Estado       int        `db:"estado" json:"estado" valid:"-"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at"`
}

func NewResolution(id string, nombre string, descripcion string, estado int, servicio_id string, ruta_archivo string) *Resolution {
	now := time.Now()
	return &Resolution{
		ID:           id,
		Nombre:       nombre,
		Descripcion:  descripcion,
		Estado:       estado,
		Servicio_id:  servicio_id,
		Ruta_archivo: ruta_archivo,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
}

func (m *Resolution) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
