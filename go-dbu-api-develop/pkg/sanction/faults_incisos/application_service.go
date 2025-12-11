// pkg/sanction/faults_incisos/application_service.go
package faults_incisos

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"

	"github.com/asaskevich/govalidator"
)

type PortsServerFaultInciso interface {
	CreateFaultInciso(id string, faltaId string, incisoId string) (*FaultInciso, int, error)
	DeleteFaultInciso(id string) (int, error)
	GetAllFaultIncisos() ([]*FaultInciso, error)
}

type service struct {
	repository ServicesFaultIncisoRepository
	user       *models.User
	txID       string
}

func NewFaultIncisoService(repository ServicesFaultIncisoRepository, user *models.User, txID string) PortsServerFaultInciso {
	return &service{repository: repository, user: user, txID: txID}
}

func (s *service) CreateFaultInciso(id string, faltaId string, incisoId string) (*FaultInciso, int, error) {
	m := NewFaultInciso(id, faltaId, incisoId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create FaultInciso:", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) DeleteFaultInciso(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.delete(id); err != nil {
		logger.Error.Println(s.txID, " - couldn't delete FaultInciso:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetAllFaultIncisos() ([]*FaultInciso, error) {
	return s.repository.getAll()
}
