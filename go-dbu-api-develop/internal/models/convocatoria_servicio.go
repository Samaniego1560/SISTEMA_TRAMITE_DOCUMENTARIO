package models

import "time"

// ConvocatoriaServicio representa la relación entre convocatoria y servicio
type ConvocatoriaServicio struct {
	ID             uint64     `db:"id" json:"id"`
	ConvocatoriaID uint64     `db:"convocatoria_id" json:"convocatoria_id"`
	ServicioID     uint64     `db:"servicio_id" json:"servicio_id"`
	Cantidad       int        `db:"cantidad" json:"cantidad"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`
}
