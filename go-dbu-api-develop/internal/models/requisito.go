package models

import "time"

// Requisito representa un requisito dentro de una sección
type Requisito struct {
	ID              uint64     `db:"id" json:"id"`
	SeccionID       uint64     `db:"seccion_id" json:"seccion_id"`
	Nombre          string     `db:"nombre" json:"nombre"`
	Descripcion     *string    `db:"descripcion" json:"descripcion"`
	UrlGuia         *string    `db:"url_guia" json:"url_guia"`
	UrlPlantilla    *string    `db:"url_plantilla" json:"url_plantilla"`
	Opciones        *string    `db:"opciones" json:"opciones"`
	TipoRequisitoID uint64     `db:"tipo_requisito_id" json:"tipo_requisito_id"`
	UserID          uint64     `db:"user_id" json:"user_id"`
	CreatedAt       *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at"`
}
