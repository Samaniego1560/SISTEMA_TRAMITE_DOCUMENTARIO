package sanctions

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
)

type PortsServerSanction interface {
	CreateSanction(faultID string, tipoSanction string, duracion int, fechaInicio time.Time, fechaFin time.Time, estado string, observacion string) (*Sanction, int, error)
	UpdateSanction(id string, tipoSanction string, duracion int, fechaInicio time.Time, fechaFin time.Time, estado string, observacion string) (*Sanction, int, error)
	DeleteSanction(id string) (int, error)
	GetSanctionByID(id string) (*Sanction, int, error)
	GetAllSanctions() ([]*Sanction, error)
}

type service struct {
	repository ServicesSanctionRepository
	user       *models.User
	txID       string
}

func NewSanctionService(repository ServicesSanctionRepository, user *models.User, TxID string) PortsServerSanction {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateSanction(faultID string, tipoSanction string, duracion int, fechaInicio time.Time, fechaFin time.Time, estado string, observacion string) (*Sanction, int, error) {
	m := NewSanction("", faultID, tipoSanction, duracion, fechaInicio, fechaFin, estado, observacion)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't create Sanction :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateSanction(id string, tipoSanction string, duracion int, fechaInicio time.Time, fechaFin time.Time, estado string, observacion string) (*Sanction, int, error) {
	m := NewSanction(id, "", tipoSanction, duracion, fechaInicio, fechaFin, estado, observacion)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Sanction :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteSanction(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.delete(id); err != nil {
		logger.Error.Println(s.txID, " - couldn't delete Sanction:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetSanctionByID(id string) (*Sanction, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't getByID Sanction:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllSanctions() ([]*Sanction, error) {
	return s.repository.getAll()
}
