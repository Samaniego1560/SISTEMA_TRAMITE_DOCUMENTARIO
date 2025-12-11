package resolutions

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type PortsServerResolution interface {
	CreateResolution(id string, nombre string, descripcion string, estado int, ServicioId string, RutaArchivo string) (*Resolution, int, error)
	UpdateResolution(id string, nombre string, descripcion string, estado int, ServicioId string, RutaArchivo string) (*Resolution, int, error)
	DeleteResolution(id string) (int, error)
	GetResolutionByID(id string) (*Resolution, int, error)
	GetAllResolution() ([]*Resolution, error)
}

type service struct {
	repository ServiceResolutionRepository
	user       *models.User
	txID       string
}

func NewResolutionService(repository ServiceResolutionRepository, user *models.User, TxID string) PortsServerResolution {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateResolution(id string, nombre string, descripcion string, estado int, ServicioId string, RutaArchivo string) (*Resolution, int, error) {
	m := NewResolution(id, strings.ToUpper(nombre), descripcion, estado, ServicioId, RutaArchivo)
	valid, err := m.valid()
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't meet validations:", err)
		return nil, 15, err
	}
	if !valid {
		logger.Error.Println(s.txID, " - don't meet validations:")
		return nil, 15, err
	}
	if err := s.repository.create(m); err != nil {
		if err.Error() == "rows affected error" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Resolution :", err)
		return nil, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateResolution(id string, nombre string, description string, estado int, ServicioId string, RutaArchivo string) (*Resolution, int, error) {
	m := NewResolution(id, strings.ToUpper(nombre), description, estado, ServicioId, RutaArchivo)
	valid, err := m.valid()
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't meet validations:", err)
		return nil, 15, err
	}
	if !valid {
		logger.Error.Println(s.txID, " - don't meet validations:")
		return nil, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Resolution :", err)
		return nil, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteResolution(id string) (int, error) {
	if err := uuid.Validate(id); err != nil {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "rows affected error" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetResolutionByID(id string) (*Resolution, int, error) {
	if err := uuid.Validate(id); err != nil {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllResolution() ([]*Resolution, error) {
	return s.repository.getAll()
}
