// pkg/sanction/faults_articulos/domain.go
package faults_articulos

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type FaultArticulo struct {
	ID         string    `json:"id" valid:"uuid,required" db:"id"`
	FaultID    string    `json:"fault_id" valid:"uuid,required" db:"falta_id"`
	ArticuloID string    `json:"articulo_id" valid:"uuid,required" db:"articulo_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func NewFaultArticulo(id, faultID, articuloID string) *FaultArticulo {
	return &FaultArticulo{
		ID:         id,
		FaultID:    faultID,
		ArticuloID: articuloID,
	}
}

func (m *FaultArticulo) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
