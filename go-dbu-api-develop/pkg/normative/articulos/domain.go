package articles

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Article struct {
	ID           string     `db:"id" json:"id" valid:"-"`
	Descripcion  string     `db:"descripcion" json:"descripcion" valid:"required"`
	Gravedad     string     `db:"gravedad" json:"gravedad" valid:"required"`
	Capitulo_id  string     `db:"capitulo_id" json:"capitulo_id" valid:"required"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at"`
	ResolucionID string     `db:"resolucion_id" json:"resolucion_id"`
}

func NewArticle(id string, descripcion string, gravedad string, capitulo_id string) *Article {
	now := time.Now()
	return &Article{
		ID:          id,
		Descripcion: descripcion,
		Gravedad:    gravedad,
		Capitulo_id: capitulo_id,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
}

func (m *Article) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
