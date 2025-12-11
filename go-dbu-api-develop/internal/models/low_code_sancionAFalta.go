package models

import (
	"time"
)

type Sancion struct {
	ID             string    `json:"id" gorm:"primaryKey;size:36"`
	ResolucionID   string    `json:"resolucion_id" gorm:"size:36;not null"`
	ArticuloID     string    `json:"articulo_id" gorm:"size:36;not null"`
	IncisoSancion  string    `json:"inciso_sancion" gorm:"size:10;not null"`
	DetalleSancion string    `json:"detalle_sancion" gorm:"type:text;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
