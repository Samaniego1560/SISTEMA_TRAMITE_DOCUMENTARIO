package models

import "time"

// Convocatoria representa la tabla 'convocatorias' en la base de datos
// Puedes ajustar los tipos según tus necesidades
//
type Convocatoria struct {
	ID            uint64     `db:"id" json:"id"`
	FechaInicio   *time.Time `db:"fecha_inicio" json:"fecha_inicio"`
	FechaFin      *time.Time `db:"fecha_fin" json:"fecha_fin"`
	Nombre        string     `db:"nombre" json:"nombre"`
	UserID        uint64     `db:"user_id" json:"user_id"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
	CreditoMinimo *int       `db:"credito_minimo" json:"credito_minimo"`
	NotaMinima    *int       `db:"nota_minima" json:"nota_minima"`
}
