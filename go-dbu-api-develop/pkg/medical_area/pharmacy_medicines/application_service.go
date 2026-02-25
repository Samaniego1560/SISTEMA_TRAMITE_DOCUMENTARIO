package pharmacy_medicines

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// PortsServerMedicine define la interfaz del servicio
type PortsServerMedicine interface {
	CreateMedicine(id, codigo, nombreGenerico string, nombreComercial *string, formaFarmaceutica, concentracion string, viaAdministracion *string, requiereReceta, controlado bool, descripcion *string, estado string) (*Medicine, int, error)
	UpdateMedicine(medicine *Medicine) (int, error)
	DeleteMedicine(id string) (int, error)
	GetMedicineByID(id string) (*Medicine, int, error)
	GetMedicineByCode(code string) (*Medicine, int, error)
	GetAllMedicines() ([]*Medicine, error)
	GetAllMedicinesWithStock() ([]*MedicineWithStock, error)
	SearchMedicines(search, estado string, limit, offset int64) ([]*MedicineWithStock, int64, error)
}

type service struct {
	repository ServicesMedicineRepository
	user       *models.User
	txID       string
}

// NewMedicineService crea una nueva instancia del servicio
func NewMedicineService(repository ServicesMedicineRepository, user *models.User, txID string) PortsServerMedicine {
	return &service{
		repository: repository,
		user:       user,
		txID:       txID,
	}
}

// CreateMedicine crea un nuevo medicamento
func (s *service) CreateMedicine(id, codigo, nombreGenerico string, nombreComercial *string, formaFarmaceutica, concentracion string, viaAdministracion *string, requiereReceta, controlado bool, descripcion *string, estado string) (*Medicine, int, error) {
	// Validar que el código no exista
	exists, err := s.repository.existsByCode(codigo)
	if err != nil {
		logger.Error.Println(s.txID, " - error checking if medicine exists by code:", err)
		return nil, 3, err
	}
	if exists {
		logger.Error.Println(s.txID, " - medicine with code already exists")
		return nil, 16, fmt.Errorf("medicine with code already exists")
	}

	// Crear medicamento
	m := NewMedicine(id, codigo, nombreGenerico, nombreComercial, formaFarmaceutica, concentracion, viaAdministracion, requiereReceta, controlado, descripcion, estado, s.user.Username)

	// Validar estructura
	if valid, err := m.Valid(); !valid {
		logger.Error.Println(s.txID, " - validation error:", err)
		return m, 15, err
	}

	// Insertar en BD
	if err := s.repository.create(m); err != nil {
		if err.Error() == "rows affected error" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create medicine:", err)
		return m, 3, err
	}

	return m, 29, nil
}

// UpdateMedicine actualiza un medicamento existente
func (s *service) UpdateMedicine(medicine *Medicine) (int, error) {
	// Validar estructura
	valid, err := medicine.Valid()
	if err != nil {
		logger.Error.Println(s.txID, " - validation error:", err)
		return 15, err
	}
	if !valid {
		logger.Error.Println(s.txID, " - validation failed")
		return 15, errors.New("validation failed")
	}

	// Actualizar en BD
	if err := s.repository.update(medicine); err != nil {
		if err.Error() == "rows affected error" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update medicine:", err)
		return 18, err
	}

	return 29, nil
}

// DeleteMedicine elimina (soft delete) un medicamento
func (s *service) DeleteMedicine(id string) (int, error) {
	// Validar UUID
	if err := uuid.Validate(id); err != nil {
		logger.Error.Println(s.txID, " - invalid UUID:", err)
		return 15, fmt.Errorf("invalid UUID")
	}

	// Eliminar (soft delete)
	if err := s.repository.delete(id); err != nil {
		if err.Error() == "rows affected error" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't delete medicine:", err)
		return 20, err
	}

	return 28, nil
}

// GetMedicineByID obtiene un medicamento por ID
func (s *service) GetMedicineByID(id string) (*Medicine, int, error) {
	// Validar UUID
	if err := uuid.Validate(id); err != nil {
		logger.Error.Println(s.txID, " - invalid UUID:", err)
		return nil, 15, fmt.Errorf("invalid UUID")
	}

	// Obtener de BD
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get medicine by ID:", err)
		return nil, 22, err
	}

	if m == nil {
		return nil, 22, fmt.Errorf("medicine not found")
	}

	return m, 29, nil
}

// GetMedicineByCode obtiene un medicamento por código
func (s *service) GetMedicineByCode(code string) (*Medicine, int, error) {
	if code == "" {
		logger.Error.Println(s.txID, " - code is required")
		return nil, 15, fmt.Errorf("code is required")
	}

	// Obtener de BD
	m, err := s.repository.getByCode(code)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get medicine by code:", err)
		return nil, 22, err
	}

	if m == nil {
		return nil, 22, fmt.Errorf("medicine not found")
	}

	return m, 29, nil
}

// GetAllMedicines obtiene todos los medicamentos
func (s *service) GetAllMedicines() ([]*Medicine, error) {
	return s.repository.getAll()
}

// GetAllMedicinesWithStock obtiene todos los medicamentos con información de stock
func (s *service) GetAllMedicinesWithStock() ([]*MedicineWithStock, error) {
	return s.repository.getAllWithStock()
}

// SearchMedicines busca medicamentos con filtros y paginación
func (s *service) SearchMedicines(search, estado string, limit, offset int64) ([]*MedicineWithStock, int64, error) {
	// Obtener total
	total, err := s.repository.countMedicines(search, estado)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't count medicines:", err)
		return nil, 0, err
	}

	// Obtener resultados
	medicines, err := s.repository.searchMedicines(search, estado, limit, offset)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't search medicines:", err)
		return nil, 0, err
	}

	return medicines, total, nil
}
