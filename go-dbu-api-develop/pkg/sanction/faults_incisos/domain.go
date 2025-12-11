package faults_incisos

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type FaultInciso struct {
	ID        string    `json:"id" valid:"uuid,required" db:"id"`
	FaultID   string    `json:"fault_id" valid:"uuid,required" db:"falta_id"`
	IncisoID  string    `json:"inciso_id" valid:"uuid,required" db:"inciso_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewFaultInciso(id, faultID, incisoID string) *FaultInciso {
	return &FaultInciso{
		ID:       id,
		FaultID:  faultID,
		IncisoID: incisoID,
	}
}

func (m *FaultInciso) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
