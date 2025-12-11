package faults_articulos

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"

	"github.com/asaskevich/govalidator"
)

type PortsServerFaultArticulo interface {
	CreateFaultArticulo(id string, faltaId string, articuloId string) (*FaultArticulo, int, error)
	DeleteFaultArticulo(id string) (int, error)
	GetAllFaultArticulos() ([]*FaultArticulo, error)
}

type service struct {
	repository ServicesFaultArticuloRepository
	user       *models.User
	txID       string
}

func NewFaultArticuloService(repository ServicesFaultArticuloRepository, user *models.User, txID string) PortsServerFaultArticulo {
	return &service{repository: repository, user: user, txID: txID}
}

func (s *service) CreateFaultArticulo(id string, faltaId string, articuloId string) (*FaultArticulo, int, error) {
	m := NewFaultArticulo(id, faltaId, articuloId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create FaultArticulo:", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) DeleteFaultArticulo(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.delete(id); err != nil {
		logger.Error.Println(s.txID, " - couldn't delete FaultArticulo:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetAllFaultArticulos() ([]*FaultArticulo, error) {
	return s.repository.getAll()
}
