package chapters

import (
	"database/sql"
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type PortsServerChapter interface {
	CreateChapter(id string, nombre string, descripcion string, ResolucionId string) (*Chapter, int, error)
	UpdateChapter(id string, nombre string, descripcion string, ResolucionId string) (*Chapter, int, error)
	DeleteChapter(id string) (int, error)
	GetChapterByID(id string) (*Chapter, int, error)
	GetAllChapter() ([]*Chapter, error)
	GetChaptersByResolution(ResolucionId string) ([]*Chapter, int, error)
	UpdateOnlyCharacteristicsChapter(id string, nombre string, descripcion string) (*Chapter, int, error)
}

type service struct {
	repository ServiceChapterRepository
	user       *models.User
	txID       string
}

func NewChapterService(repository ServiceChapterRepository, user *models.User, TxID string) PortsServerChapter {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateChapter(id string, nombre string, descripcion string, Resolucion_id string) (*Chapter, int, error) {
	if id == "" {
		id = uuid.New().String()
	}

	m := NewChapter(id, strings.ToUpper(nombre), descripcion, Resolucion_id)
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

func (s *service) UpdateChapter(id string, nombre string, description string, resolucion_id string) (*Chapter, int, error) {
	m := NewChapter(id, strings.ToUpper(nombre), description, resolucion_id)
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

func (s *service) DeleteChapter(id string) (int, error) {
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

func (s *service) GetChapterByID(id string) (*Chapter, int, error) {
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

func (s *service) GetAllChapter() ([]*Chapter, error) {
	return s.repository.getAll()
}

func (s *service) GetChaptersByResolution(ResolucionId string) ([]*Chapter, int, error) {
	chapters, err := s.repository.GetByResolutionID(ResolucionId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logger.Error.Println(s.txID, " - no rows found for resolution: ", ResolucionId)
			return nil, 14, err
		default:
			logger.Error.Println(s.txID, " - couldn't get chapters: ", err)
			return nil, 13, err
		}
	}

	return chapters, 29, nil
}
func (s *service) UpdateOnlyCharacteristicsChapter(id string, nombre string, descripcion string) (*Chapter, int, error) {
	m := NewChapter(id, nombre, descripcion, "")
	if err := s.repository.updateOnlyCharacteristics(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Room :", err)
		return m, 18, err
	}
	return m, 29, nil
}
