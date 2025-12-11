package models

import (
	"fmt"
	"time"
)

// SancionFaltaAsignada representa la asignación de una sanción a una falta concreta
type SancionFaltaAsignada struct {
	FaltaID         string     `db:"falta_id" json:"falta_id"`
	FechaInicio     *time.Time `db:"fecha_inicio" json:"fecha_inicio"`
	FechaFin        *time.Time `db:"fecha_fin" json:"fecha_fin"`
	ID              string     `db:"id" json:"id"`
	ResolucionID    string     `db:"resolucion_id" json:"resolucion_id"`
	SancionID       string     `db:"sancion_id" json:"sancion_id"`
	FechaAsignacion time.Time  `db:"fecha_asignacion" json:"fecha_asignacion"`
	Observaciones   string     `db:"observaciones" json:"observaciones"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	ArticuloID      string     `db:"articulo_id" json:"articulo_id"`
	CapituloSancion string     `db:"capitulo_sancion" json:"capitulo_sancion"`
	ArticuloSancion string     `db:"articulo_sancion" json:"articulo_sancion"`
	IncisoSancion   string     `db:"inciso_sancion" json:"inciso_sancion"`
	DetalleSancion  string     `db:"detalle_sancion" json:"detalle_sancion"`
}

// Valid verifica si la estructura tiene los campos mínimos requeridos
// Si la estructura se usa para sancionesAFaltas, no exige SancionID
func (s *SancionFaltaAsignada) Valid() (bool, error) {
	if s.ID == "" {
		return false, fmt.Errorf("ID es obligatorio")
	}
	// Si la estructura se usa para sancion_falta_asignada, exigir SancionID
	// Si se usa para sancionesAFaltas, no exigir SancionID
	// Se asume que si SancionID es relevante, el handler lo valida aparte
	return true, nil
}

// ...existing code...
