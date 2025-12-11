package incisos

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Inciso struct {
	ID          string     `db:"id" json:"id" valid:"-"`
	Nombre      string     `db:"nombre" json:"nombre" valid:"required"`
	Descripcion string     `db:"descripcion" json:"descripcion" valid:"required"`
	ArticuloId  string     `db:"articulo_id" json:"articulo_id" valid:"required"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
}

func NewInciso(id string, nombre string, descripcion string, articuloId string) *Inciso {
	now := time.Now()
	return &Inciso{
		ID:          id,
		Descripcion: descripcion,
		Nombre:      nombre,
		ArticuloId:  articuloId,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
}

func (m *Inciso) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
