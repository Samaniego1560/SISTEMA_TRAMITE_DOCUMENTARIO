package models

import "time"

// SancionAsignadaDetalle representa la sanción asignada a una falta con detalle
// Incluye los datos de la asignación y los campos de la sanción
// Se usa para la respuesta del endpoint /asignadas-por-falta/:falta_id
//
type SancionAsignadaDetalle struct {
	ID              string     `db:"id" json:"id"`
	FaltaID         string     `db:"falta_id" json:"falta_id"`
	ResolucionID    string     `db:"resolucion_id" json:"resolucion_id"`
	SancionID       string     `db:"sancion_id" json:"sancion_id"`
	FechaAsignacion time.Time  `db:"fecha_asignacion" json:"fecha_asignacion"`
	FechaInicio     *time.Time `db:"fecha_inicio" json:"fecha_inicio"`
	FechaFin        *time.Time `db:"fecha_fin" json:"fecha_fin"`
	Observaciones   string     `db:"observaciones" json:"observaciones"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	CapituloSancion string     `db:"capitulo_sancion" json:"capitulo_sancion"`
	ArticuloSancion string     `db:"articulo_sancion" json:"articulo_sancion"`
	IncisoSancion   string     `db:"inciso_sancion" json:"inciso_sancion"`
	DetalleSancion  string     `db:"detalle_sancion" json:"detalle_sancion"`
	Estado          string     `db:"estado" json:"estado"`
}
