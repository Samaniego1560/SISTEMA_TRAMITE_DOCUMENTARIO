package models

import "time"

// Seccion representa una sección dentro de una convocatoria
type Seccion struct {
	ID             uint64     `db:"id" json:"id"`
	ConvocatoriaID uint64     `db:"convocatoria_id" json:"convocatoria_id"`
	Descripcion    string     `db:"descripcion" json:"descripcion"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`
}
