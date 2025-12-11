package incisos

import (
	"database/sql"
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"

	"github.com/google/uuid"
)

type PortsServerInciso interface {
	CreateInciso(id string, nombre string, descripcion string, articuloId string) (*Inciso, int, error)
	UpdateInciso(id string, nombre string, descripcion string, articuloId string) (*Inciso, int, error)
	DeleteInciso(id string) (int, error)
	GetIncisoByID(id string) (*Inciso, int, error)
	GetAllInciso() ([]*Inciso, error)
	GetIncisosByArticleId(articleId string) ([]*Inciso, int, error)
}

type service struct {
	repository ServiceIncisoRepository
	user       *models.User
	txID       string
}

func NewIncisoService(repository ServiceIncisoRepository, user *models.User, TxID string) PortsServerInciso {
	return &service{repository: repository, user: user, txID: TxID}
}
func (s *service) CreateInciso(id string, nombre string, descripcion string, articuloId string) (*Inciso, int, error) {
	m := NewInciso(id, nombre, descripcion, articuloId)
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
		logger.Error.Println(s.txID, " - couldn't create Inciso :", err)
		return nil, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateInciso(id string, nombre string, descripcion string, articuloId string) (*Inciso, int, error) {
	m := NewInciso(id, nombre, descripcion, articuloId)
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
		logger.Error.Println(s.txID, " - couldn't update Inciso :", err)
		return nil, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteInciso(id string) (int, error) {
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

func (s *service) GetIncisoByID(id string) (*Inciso, int, error) {
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

func (s *service) GetAllInciso() ([]*Inciso, error) {
	return s.repository.getAll()
}

func (s *service) GetIncisosByArticleId(articleId string) ([]*Inciso, int, error) {
	Incisos, err := s.repository.GetByarticleID(articleId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logger.Error.Println(s.txID, " - no rows found for Chapter: ", articleId)
			return nil, 14, err
		default:
			logger.Error.Println(s.txID, " - couldn't get Incisos: ", err)
			return nil, 13, err
		}
	}

	return Incisos, 29, nil
}
