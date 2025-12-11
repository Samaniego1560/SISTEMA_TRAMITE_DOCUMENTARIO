package nursing_consultation_vaccine

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"
	"github.com/google/uuid"
)

type PortsServerVaccine interface {
	CreateVaccine(id, consulta_enfermeria_id, tipo_vacuna, fecha_dosis string, minsa bool, tipoAtencion, indicaciones, observaciones string) (*Vaccine, int, error)
	UpdateVaccine(id, consulta_enfermeria_id, tipo_vacuna, fecha_dosis string, minsa bool, tipoAtencion, indicaciones, observaciones string) (*Vaccine, int, error)
	DeleteVaccine(id string) (int, error)
	DeleteVaccineByIDConsultation(id string) (int, error)
	GetVaccineByID(id string) ([]*Vaccine, int, error)
	GetVaccineByIDPatient(id string) ([]*Vaccine, int, error)
	GetAllVaccine() ([]*Vaccine, error)
	GetAllTypesVaccines() ([]*VaccineType, error)
	GetAllVaccinesByPatientDni(dni string) ([]*Vaccine, error)
	GetAllRequiredVaccineTypes() ([]*VaccineType, error)
	GetAllNotRequiredVaccineTypes() ([]*VaccineType, error)
	GetVaccineIntervalsByVaccineType(tipoVacunaID string) ([]*VaccineInterval, error)
	GetNextVaccineIntervalsByVaccineType(tipoVacunaID string, currentDosis int) (*VaccineInterval, int, error)
}

type service struct {
	repository ServicesVaccineRepository
	user       *models.User
	txID       string
}

func NewVaccineService(repository ServicesVaccineRepository, user *models.User, TxID string) PortsServerVaccine {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateVaccine(id, consulta_enfermeria_id, tipo_vacuna, fecha_dosis string, minsa bool, tipoAtencion, indicaciones, observaciones string) (*Vaccine, int, error) {

	m := NewVaccine(id, consulta_enfermeria_id, tipo_vacuna, fecha_dosis, minsa, tipoAtencion, indicaciones, observaciones)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.create(m); err != nil {
		if err.Error() == "rows affected error" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create vaccine :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateVaccine(id, consulta_enfermeria_id, tipo_vacuna, fecha_dosis string, minsa bool, tipoAtencion, indicaciones, observaciones string) (*Vaccine, int, error) {
	m := NewVaccine(id, consulta_enfermeria_id, tipo_vacuna, fecha_dosis, minsa, tipoAtencion, indicaciones, observaciones)
	valid, err := m.valid()
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't meet validations:", err)
		return nil, 15, err
	}
	if !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return nil, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update vaccine :", err)
		return nil, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteVaccine(id string) (int, error) {
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

func (s *service) DeleteVaccineByIDConsultation(id string) (int, error) {
	if err := uuid.Validate(id); err != nil {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
	}

	if err := s.repository.deleteByIDConsultation(id); err != nil {
		if err.Error() == "rows affected error" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetVaccineByID(id string) ([]*Vaccine, int, error) {
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

func (s *service) GetVaccineByIDPatient(id string) ([]*Vaccine, int, error) {
	if err := uuid.Validate(id); err != nil {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByIDPatient(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllVaccine() ([]*Vaccine, error) {
	return s.repository.getAll()
}

func (s *service) GetAllTypesVaccines() ([]*VaccineType, error) {
	return s.repository.getAllTypesVaccines()
}

func (s *service) GetAllVaccinesByPatientDni(dni string) ([]*Vaccine, error) {
	return s.repository.getAllVaccinesByPatientDni(dni)
}

func (s *service) GetAllRequiredVaccineTypes() ([]*VaccineType, error) {
	return s.repository.getAllRequiredVaccineTypes()
}

func (s *service) GetAllNotRequiredVaccineTypes() ([]*VaccineType, error) {
	return s.repository.getAllNotRequiredVaccineTypes()
}

func (s *service) GetVaccineIntervalsByVaccineType(tipoVacunaID string) ([]*VaccineInterval, error) {
	return s.repository.getVaccineIntervalsByVaccineType(tipoVacunaID)
}

func (s *service) GetNextVaccineIntervalsByVaccineType(tipoVacunaID string, currentDosis int) (*VaccineInterval, int, error) {
	if err := uuid.Validate(tipoVacunaID); err != nil {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getNextVaccineIntervalsByVaccineType(tipoVacunaID, currentDosis)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
