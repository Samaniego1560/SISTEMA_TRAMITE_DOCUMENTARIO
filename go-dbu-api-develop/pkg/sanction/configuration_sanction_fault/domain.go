package configuration_sanction_fault

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Modelo de Sanción (para tabla sancionesAFaltas)
type Sancion struct {
	ID                  string    `db:"id" json:"id"`
	ResolucionID        string    `db:"resolucion_id" json:"resolucion_id"`
	ResolucionNombre    string    `db:"resolucion_nombre" json:"resolucion_nombre"`
	ArticuloID          string    `db:"articulo_id" json:"articulo_id"`
	ArticuloDescripcion string    `db:"articulo_descripcion" json:"articulo_descripcion"`
	Gravedad            string    `db:"gravedad" json:"gravedad"`
	CapituloSancion     string    `db:"capitulo_sancion" json:"capitulo_sancion"`
	CapituloNombre      string    `db:"capitulo_nombre" json:"capitulo_nombre"`
	ArticuloSancion     string    `db:"articulo_sancion" json:"articulo_sancion"`
	IncisoSancion       string    `db:"inciso_sancion" json:"inciso_sancion"`
	DetalleSancion      string    `db:"detalle_sancion" json:"detalle_sancion"`
	RequiereFechas      bool      `json:"requiereFechas"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

// Validación
func (m *Sancion) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
