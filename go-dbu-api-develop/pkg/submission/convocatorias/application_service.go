package convocatorias

import (
	"fmt"
	"time"

	"dbu-api/internal/logger"
	"dbu-api/internal/models"
)

type PortsServerConvocatorias interface {
	CreateSubmissions(fechaInicio, fechaFin time.Time, nombre string, userId int64, creditoMinimo, notaMinima int) (*Convocatorias, int, error)
	UpdateSubmissions(id int64, fechaInicio, fechaFin time.Time, nombre string, userId int64, creditoMinimo, notaMinima int) (*Convocatorias, int, error)
	DeleteSubmissions(id int64) (int, error)
	GetSubmissionsByID(id int64) (*Convocatorias, int, error)
	GetAllSubmissions() ([]*Convocatorias, error)
	GetAllSubmissionsByService(id int64) ([]*Convocatorias, error)
	GetActiveSubmissions() (*Convocatorias, int, error)
	GetLastSubmission() (*Convocatorias, int, error)
	// Nuevos métodos para manejar convocatorias con relaciones
	CreateConvocatoriaWithRelations(req *CreateConvocatoriaRequest) (*ConvocatoriaResponse, int, error)
	UpdateConvocatoriaWithRelations(id int64, req *CreateConvocatoriaRequest) (*ConvocatoriaResponse, int, error)
	GetConvocatoriaWithRelations(id int64) (*ConvocatoriaResponse, int, error)
}

type service struct {
	repository ServicesConvocatoriasRepository
	user       *models.User
	txID       string
}

func NewConvocatoriasService(repository ServicesConvocatoriasRepository, user *models.User, TxID string) PortsServerConvocatorias {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateSubmissions(fechaInicio, fechaFin time.Time, nombre string, userId int64, creditoMinimo, notaMinima int) (*Convocatorias, int, error) {
	credMin := creditoMinimo
	notaMin := notaMinima
	m := NewCreateSubmissions(&fechaInicio, &fechaFin, nombre, userId, &credMin, &notaMin)

	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "rows affected error" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Convocatorias :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateSubmissions(id int64, fechaInicio, fechaFin time.Time, nombre string, userId int64, creditoMinimo, notaMinima int) (*Convocatorias, int, error) {
	credMin := creditoMinimo
	notaMin := notaMinima
	m := NewSubmissions(id, &fechaInicio, &fechaFin, nombre, userId, &credMin, &notaMin)

	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}

	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Convocatorias :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteSubmissions(id int64) (int, error) {
	if id == 0 {
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

func (s *service) GetSubmissionsByID(id int64) (*Convocatorias, int, error) {
	if id == 0 {
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

func (s *service) GetAllSubmissions() ([]*Convocatorias, error) {
	return s.repository.getAll()
}

func (s *service) GetAllSubmissionsByService(id int64) ([]*Convocatorias, error) {
	return s.repository.getAllByService(id)
}

func (s *service) GetActiveSubmissions() (*Convocatorias, int, error) {
	m, err := s.repository.getActive()
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
func (s *service) GetLastSubmission() (*Convocatorias, int, error) {
	m, err := s.repository.getLast()
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getLast row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

// parseDate intenta parsear una fecha con múltiples formatos
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",       // Formato del frontend
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05Z",      // RFC3339 sin timezone
		"2006-01-02T15:04:05",       // ISO8601 sin timezone
		"2006-01-02",                // Solo fecha
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", dateStr)
}

// CreateConvocatoriaWithRelations crea una convocatoria con todas sus relaciones en una transacción
func (s *service) CreateConvocatoriaWithRelations(req *CreateConvocatoriaRequest) (*ConvocatoriaResponse, int, error) {
	// Parsear fechas
	fechaInicio, err := parseDate(req.FechaInicio)
	if err != nil {
		logger.Error.Println(s.txID, " - error parsing fecha_inicio:", err)
		return nil, 15, err
	}

	fechaFin, err := parseDate(req.FechaFin)
	if err != nil {
		logger.Error.Println(s.txID, " - error parsing fecha_fin:", err)
		return nil, 15, err
	}

	// Validar que fecha_fin >= fecha_inicio
	if fechaFin.Before(fechaInicio) {
		logger.Error.Println(s.txID, " - fecha_fin must be after fecha_inicio")
		return nil, 15, fmt.Errorf("fecha_fin debe ser mayor o igual a fecha_inicio")
	}

	// Verificar que no haya convocatorias solapadas
	hasOverlap, err := s.repository.checkOverlappingConvocatorias(fechaInicio, nil)
	if err != nil {
		logger.Error.Println(s.txID, " - error checking overlapping convocatorias:", err)
		return nil, 3, err
	}
	if hasOverlap {
		logger.Error.Println(s.txID, " - overlapping convocatorias found")
		return nil, 15, fmt.Errorf("no se puede crear la convocatoria durante una fecha establecida entre convocatorias existentes")
	}

	// Crear convocatoria principal
	convocatoria := &Convocatorias{
		FechaInicio:   &fechaInicio,
		FechaFin:      &fechaFin,
		Nombre:        req.Nombre,
		UserId:        s.user.ID,
		CreditoMinimo: req.CreditoMinimo,
		NotaMinima:    req.NotaMinima,
	}

	if err := s.repository.create(convocatoria); err != nil {
		if err.Error() == "rows affected error" {
			return nil, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create convocatoria:", err)
		return nil, 3, err
	}

	// Crear servicios de convocatoria
	var servicios []ConvocatoriaServicio
	for _, csReq := range req.ConvocatoriaServicio {
		cs := &ConvocatoriaServicio{
			ConvocatoriaID: convocatoria.ID,
			ServicioID:     csReq.ServicioID,
			Cantidad:       csReq.Cantidad,
		}
		if err := s.repository.createConvocatoriaServicio(cs); err != nil {
			logger.Error.Println(s.txID, " - couldn't create convocatoria_servicio:", err)
			return nil, 3, err
		}
		servicios = append(servicios, *cs)
	}

	// Crear secciones con sus requisitos
	var secciones []SeccionResponse
	for _, secReq := range req.Secciones {
		seccion := &Seccion{
			ConvocatoriaID: convocatoria.ID,
			Descripcion:    secReq.Descripcion,
		}
		if err := s.repository.createSeccion(seccion); err != nil {
			logger.Error.Println(s.txID, " - couldn't create seccion:", err)
			return nil, 3, err
		}

		// Crear requisitos para cada sección
		var requisitos []Requisito
		for _, reqReq := range secReq.Requisitos {
			requisito := &Requisito{
				SeccionID:       seccion.ID,
				Nombre:          reqReq.Nombre,
				Descripcion:     reqReq.Descripcion,
				UrlGuia:         reqReq.UrlGuia,
				UrlPlantilla:    reqReq.UrlPlantilla,
				Opciones:        reqReq.Opciones,
				TipoRequisitoID: reqReq.TipoRequisitoID,
				UserID:          s.user.ID,
			}
			if err := s.repository.createRequisito(requisito); err != nil {
				logger.Error.Println(s.txID, " - couldn't create requisito:", err)
				return nil, 3, err
			}
			requisitos = append(requisitos, *requisito)
		}

		secciones = append(secciones, SeccionResponse{
			ID:             seccion.ID,
			ConvocatoriaID: seccion.ConvocatoriaID,
			Descripcion:    seccion.Descripcion,
			Requisitos:     requisitos,
			CreatedAt:      seccion.CreatedAt,
			UpdatedAt:      seccion.UpdatedAt,
		})
	}

	// Preparar respuesta
	response := &ConvocatoriaResponse{
		ID:                   convocatoria.ID,
		FechaInicio:          convocatoria.FechaInicio,
		FechaFin:             convocatoria.FechaFin,
		Nombre:               convocatoria.Nombre,
		UserID:               convocatoria.UserId,
		CreditoMinimo:        convocatoria.CreditoMinimo,
		NotaMinima:           convocatoria.NotaMinima,
		ConvocatoriaServicio: servicios,
		Secciones:            secciones,
		CreatedAt:            convocatoria.CreatedAt,
		UpdatedAt:            convocatoria.UpdatedAt,
	}

	return response, 29, nil
}

// UpdateConvocatoriaWithRelations actualiza una convocatoria y todas sus relaciones
func (s *service) UpdateConvocatoriaWithRelations(id int64, req *CreateConvocatoriaRequest) (*ConvocatoriaResponse, int, error) {
	// Validar request
	if valid, err := req.Valid(); !valid {
		logger.Error.Println(s.txID, " - request validation failed:", err)
		return nil, 15, err
	}

	// Verificar que la convocatoria existe
	existingConv, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get convocatoria:", err)
		return nil, 22, err
	}
	if existingConv == nil {
		logger.Error.Println(s.txID, " - convocatoria not found")
		return nil, 22, fmt.Errorf("convocatoria no encontrada")
	}

	// Parsear fechas
	fechaInicio, err := parseDate(req.FechaInicio)
	if err != nil {
		logger.Error.Println(s.txID, " - error parsing fecha_inicio:", err)
		return nil, 15, err
	}

	fechaFin, err := parseDate(req.FechaFin)
	if err != nil {
		logger.Error.Println(s.txID, " - error parsing fecha_fin:", err)
		return nil, 15, err
	}

	// Validar que fecha_fin >= fecha_inicio
	if fechaFin.Before(fechaInicio) {
		logger.Error.Println(s.txID, " - fecha_fin must be after fecha_inicio")
		return nil, 15, fmt.Errorf("fecha_fin debe ser mayor o igual a fecha_inicio")
	}

	// Actualizar convocatoria principal
	convocatoria := &Convocatorias{
		ID:            id,
		FechaInicio:   &fechaInicio,
		FechaFin:      &fechaFin,
		Nombre:        req.Nombre,
		UserId:        existingConv.UserId,
		CreditoMinimo: req.CreditoMinimo,
		NotaMinima:    req.NotaMinima,
	}

	if err := s.repository.update(convocatoria); err != nil {
		logger.Error.Println(s.txID, " - couldn't update convocatoria:", err)
		return nil, 18, err
	}

	// Eliminar servicios existentes
	if err := s.repository.deleteConvocatoriaServicios(id); err != nil {
		logger.Error.Println(s.txID, " - couldn't delete old servicios:", err)
		return nil, 3, err
	}

	// Crear nuevos servicios
	var servicios []ConvocatoriaServicio
	for _, csReq := range req.ConvocatoriaServicio {
		cs := &ConvocatoriaServicio{
			ConvocatoriaID: id,
			ServicioID:     csReq.ServicioID,
			Cantidad:       csReq.Cantidad,
		}
		if err := s.repository.createConvocatoriaServicio(cs); err != nil {
			logger.Error.Println(s.txID, " - couldn't create convocatoria_servicio:", err)
			return nil, 3, err
		}
		servicios = append(servicios, *cs)
	}

	// Eliminar secciones existentes (y sus requisitos en cascada)
	if err := s.repository.deleteSecciones(id); err != nil {
		logger.Error.Println(s.txID, " - couldn't delete old secciones:", err)
		return nil, 3, err
	}

	// Crear nuevas secciones con sus requisitos
	var secciones []SeccionResponse
	for _, secReq := range req.Secciones {
		seccion := &Seccion{
			ConvocatoriaID: id,
			Descripcion:    secReq.Descripcion,
		}
		if err := s.repository.createSeccion(seccion); err != nil {
			logger.Error.Println(s.txID, " - couldn't create seccion:", err)
			return nil, 3, err
		}

		var requisitos []Requisito
		for _, reqReq := range secReq.Requisitos {
			requisito := &Requisito{
				SeccionID:       seccion.ID,
				Nombre:          reqReq.Nombre,
				Descripcion:     reqReq.Descripcion,
				UrlGuia:         reqReq.UrlGuia,
				UrlPlantilla:    reqReq.UrlPlantilla,
				Opciones:        reqReq.Opciones,
				TipoRequisitoID: reqReq.TipoRequisitoID,
				UserID:          s.user.ID,
			}
			if err := s.repository.createRequisito(requisito); err != nil {
				logger.Error.Println(s.txID, " - couldn't create requisito:", err)
				return nil, 3, err
			}
			requisitos = append(requisitos, *requisito)
		}

		secciones = append(secciones, SeccionResponse{
			ID:             seccion.ID,
			ConvocatoriaID: seccion.ConvocatoriaID,
			Descripcion:    seccion.Descripcion,
			Requisitos:     requisitos,
			CreatedAt:      seccion.CreatedAt,
			UpdatedAt:      seccion.UpdatedAt,
		})
	}

	// Preparar respuesta
	response := &ConvocatoriaResponse{
		ID:                   id,
		FechaInicio:          convocatoria.FechaInicio,
		FechaFin:             convocatoria.FechaFin,
		Nombre:               convocatoria.Nombre,
		UserID:               convocatoria.UserId,
		CreditoMinimo:        convocatoria.CreditoMinimo,
		NotaMinima:           convocatoria.NotaMinima,
		ConvocatoriaServicio: servicios,
		Secciones:            secciones,
		CreatedAt:            existingConv.CreatedAt,
		UpdatedAt:            convocatoria.UpdatedAt,
	}

	return response, 29, nil
}

// GetConvocatoriaWithRelations obtiene una convocatoria con todas sus relaciones
func (s *service) GetConvocatoriaWithRelations(id int64) (*ConvocatoriaResponse, int, error) {
	// Obtener convocatoria principal
	convocatoria, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get convocatoria:", err)
		return nil, 22, err
	}
	if convocatoria == nil {
		return nil, 22, fmt.Errorf("convocatoria no encontrada")
	}

	// Obtener servicios
	servicios, err := s.repository.getConvocatoriaServiciosByConvocatoriaID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get servicios:", err)
		return nil, 3, err
	}

	// Obtener secciones
	secciones, err := s.repository.getSeccionesByConvocatoriaID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't get secciones:", err)
		return nil, 3, err
	}

	// Obtener requisitos para cada sección
	var seccionesResponse []SeccionResponse
	for _, seccion := range secciones {
		requisitos, err := s.repository.getRequisitosBySeccionID(seccion.ID)
		if err != nil {
			logger.Error.Println(s.txID, " - couldn't get requisitos:", err)
			return nil, 3, err
		}

		seccionesResponse = append(seccionesResponse, SeccionResponse{
			ID:             seccion.ID,
			ConvocatoriaID: seccion.ConvocatoriaID,
			Descripcion:    seccion.Descripcion,
			Requisitos:     requisitos,
			CreatedAt:      seccion.CreatedAt,
			UpdatedAt:      seccion.UpdatedAt,
		})
	}

	// Preparar respuesta
	response := &ConvocatoriaResponse{
		ID:                   convocatoria.ID,
		FechaInicio:          convocatoria.FechaInicio,
		FechaFin:             convocatoria.FechaFin,
		Nombre:               convocatoria.Nombre,
		UserID:               convocatoria.UserId,
		CreditoMinimo:        convocatoria.CreditoMinimo,
		NotaMinima:           convocatoria.NotaMinima,
		ConvocatoriaServicio: servicios,
		Secciones:            seccionesResponse,
		CreatedAt:            convocatoria.CreatedAt,
		UpdatedAt:            convocatoria.UpdatedAt,
	}

	return response, 29, nil
}
