package chapters

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Chapter struct {
	ID            string     `db:"id" json:"id" valid:"-"`
	Nombre        string     `db:"nombre" json:"nombre" valid:"required"`
	Descripcion   string     `db:"descripcion" json:"descripcion" valid:"required"`
	Resolucion_id string     `db:"resolucion_id" json:"servicio_id" valid:"required"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

func NewChapter(id string, nombre string, descripcion string, resolucion_id string) *Chapter {
	now := time.Now()
	return &Chapter{
		ID:            id,
		Nombre:        nombre,
		Descripcion:   descripcion,
		Resolucion_id: resolucion_id,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
}

func (m *Chapter) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
