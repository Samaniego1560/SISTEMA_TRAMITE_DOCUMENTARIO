package configuration_sanction_fault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dbu-api/internal/logger"
	internalmodels "dbu-api/internal/models"
	"dbu-api/models"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

// Interface de puertos para el servicio de Sanciones a Faltas
// Interfaz para el repositorio de sanciones
type PortsSanctionFault interface {
	CreateSancion(id, resolucionID, articuloID, capituloSancion, articuloSancion, incisoSancion, detalleSancion string) (*Sancion, int, error)
	UpdateSancion(id, resolucionID, articuloID, capituloSancion, articuloSancion, incisoSancion, detalleSancion string) (*Sancion, int, error)
	DeleteSancion(id string) (int, error)
	GetSancionByID(id string) (*Sancion, int, error)
	GetAllSanciones() ([]*Sancion, error)
	AsignarSancionFalta(sfa *internalmodels.SancionFaltaAsignada) (int, error)
	RegistrarApelacion(ap *models.Apelacion) (int, error)
	GetSancionesAsignadasPorFalta(faltaID string) ([]*internalmodels.SancionAsignadaDetalle, error)
}

type service struct {
	repository SancionRepository
	user       *internalmodels.User
	txID       string
}

func (s *service) GetSancionesAsignadasPorFalta(faltaID string) ([]*internalmodels.SancionAsignadaDetalle, error) {
	return s.repository.GetSancionesAsignadasPorFalta(faltaID)
}

func NewSanctionFaultService(repository SancionRepository, user *internalmodels.User, txID string) PortsSanctionFault {
	return &service{repository: repository, user: user, txID: txID}
}

// Asignar una sanción a una falta concreta
func (s *service) AsignarSancionFalta(sfa *internalmodels.SancionFaltaAsignada) (int, error) {
	if sfa == nil || sfa.ID == "" || sfa.ResolucionID == "" || sfa.SancionID == "" {
		logger.Error.Println(s.txID, " - Faltan campos obligatorios para asignar sanción a falta")
		return 15, errors.New("faltan campos obligatorios")
	}
	if sfa.FechaAsignacion.IsZero() {
		sfa.FechaAsignacion = time.Now()
	}
	sfa.CreatedAt = time.Now()
	sfa.UpdatedAt = time.Now()
	// Se asegura que siempre se pasen fechas
	if sfa.FechaInicio == nil || sfa.FechaFin == nil {
		return 0, errors.New("fecha_inicio y fecha_fin son requeridas")
	}
	err := s.repository.AsignarSancionFalta(context.Background(), sfa)
	if err != nil {
		logger.Error.Println(s.txID, " - error asignando sanción a falta:", err)
		return 3, err
	}
	return 29, nil
}

// Registrar una apelación sobre una sanción asignada
func (s *service) RegistrarApelacion(ap *models.Apelacion) (int, error) {
	if ap == nil || ap.ID == "" || ap.SancionFaltaAsignadaID == "" || ap.Motivo == "" {
		logger.Error.Println(s.txID, " - Faltan campos obligatorios para registrar apelación")
		return 15, errors.New("faltan campos obligatorios")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if ap.FechaApelacion == nil || *ap.FechaApelacion == "" {
		ap.FechaApelacion = &now
	}
	ap.CreatedAt = &now
	ap.UpdatedAt = &now
	err := s.repository.RegistrarApelacion(context.Background(), ap)
	if err != nil {
		logger.Error.Println(s.txID, " - error registrando apelación:", err)
		return 3, err
	}
	return 29, nil
	// Crear una nueva sanción
}

// Crear una nueva sanción
func (s *service) CreateSancion(id, resolucionID, articuloID, capituloSancion, articuloSancion, incisoSancion, detalleSancion string) (*Sancion, int, error) {
	logger.Info.Println(s.txID, " - Creando sancion para resolucion:", resolucionID, " articulo:", articuloID)

	if resolucionID == "" || articuloID == "" || capituloSancion == "" || articuloSancion == "" || incisoSancion == "" || detalleSancion == "" {
		logger.Error.Println(s.txID, " - Faltan campos obligatorios")
		return nil, 15, errors.New("todos los campos son obligatorios")
	}

	if id == "" {
		id = uuid.New().String()
	}

	sancion := &Sancion{
		ID:              id,
		ResolucionID:    resolucionID,
		ArticuloID:      articuloID,
		CapituloSancion: capituloSancion,
		ArticuloSancion: articuloSancion,
		IncisoSancion:   incisoSancion,
		DetalleSancion:  detalleSancion,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if valid, err := sancion.Valid(); !valid {
		logger.Error.Println(s.txID, " - no cumple validaciones:", err)
		return sancion, 15, err
	}

	if err := s.repository.Create(sancion); err != nil {
		logger.Error.Println(s.txID, " - error al crear Sancion:", err)
		return sancion, 3, err
	}

	return sancion, 29, nil
}

// Actualizar sanción existente
func (s *service) UpdateSancion(id, resolucionID, articuloID, capituloSancion, articuloSancion, incisoSancion, detalleSancion string) (*Sancion, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - ID inválido:", id)
		return nil, 15, fmt.Errorf("id no es uuid")
	}

	existing, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - error obteniendo sancion:", err)
		return nil, 22, err
	}
	if existing == nil {
		logger.Error.Println(s.txID, " - sancion no existe:", id)
		return nil, 15, errors.New("la sanción no existe")
	}

	// Actualizar datos
	existing.ResolucionID = resolucionID
	existing.ArticuloID = articuloID
	existing.CapituloSancion = capituloSancion
	existing.ArticuloSancion = articuloSancion
	existing.IncisoSancion = incisoSancion
	existing.DetalleSancion = detalleSancion
	existing.UpdatedAt = time.Now()

	if valid, err := existing.Valid(); !valid {
		logger.Error.Println(s.txID, " - no cumple validaciones:", err)
		return existing, 15, err
	}

	if err := s.repository.Update(existing); err != nil {
		logger.Error.Println(s.txID, " - error actualizando sancion:", err)
		return existing, 18, err
	}
	return existing, 29, nil
}

// Eliminar sanción
func (s *service) DeleteSancion(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - ID inválido:", id)
		return 15, fmt.Errorf("id no es uuid")
	}
	existing, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - error obteniendo sancion:", err)
		return 22, err
	}
	if existing == nil {
		logger.Error.Println(s.txID, " - sancion no existe:", id)
		return 15, errors.New("la sanción no existe")
	}

	if err := s.repository.Delete(id); err != nil {
		logger.Error.Println(s.txID, " - error eliminando sancion:", err)
		return 20, err
	}
	return 28, nil
}

// Obtener por ID
func (s *service) GetSancionByID(id string) (*Sancion, int, error) {
	sancion, err := s.repository.GetByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - error obteniendo sancion por ID:", err)
		return nil, 22, err
	}
	if sancion == nil {
		return nil, 15, errors.New("la sanción no existe")
	}
	return sancion, 29, nil
}

// Obtener todas las sanciones
func (s *service) GetAllSanciones() ([]*Sancion, error) {
	list, err := s.repository.GetAll()
	if err != nil {
		logger.Error.Println(s.txID, " - error obteniendo todas las sanciones:", err)
		return nil, err
	}
	return list, nil
}
