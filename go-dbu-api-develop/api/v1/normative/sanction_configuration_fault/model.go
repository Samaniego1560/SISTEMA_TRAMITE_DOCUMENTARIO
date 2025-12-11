package sanction_configuration_fault

import (
	"time"
)

// Estructura que representa la tabla de sanciones
type Sancion struct {
	ID              string    `json:"id"`
	ResolucionID    string    `json:"resolucion_id"`
	ArticuloID      string    `json:"articulo_id"`
	CapituloSancion string    `json:"capitulo_sancion"`
	ArticuloSancion string    `json:"articulo_sancion"`
	IncisoSancion   string    `json:"inciso_sancion"`
	DetalleSancion  string    `json:"detalle_sancion"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
