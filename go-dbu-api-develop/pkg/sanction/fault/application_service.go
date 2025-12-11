package fault

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"dbu-api/internal/logger"
	"dbu-api/internal/models"

	"github.com/asaskevich/govalidator"

	"github.com/google/uuid"
)

// Puertos que expone el servicio Fault
type PortsServerFault interface {
	CreateFault(id string, alumnoID int64, servicioId int64, convocatoriaId int64, fuenteInformacion string, fechaFalta time.Time, estado string, observacion string, articulos []string, incisos []string) (*Fault, int, error)
	UpdateFault(id string, alumnoID int64, servicioId int64, convocatoriaId int64, fuenteInformacion string, fechaFalta time.Time, estado string, observacion string) (*Fault, int, error)
	DeleteFault(id string) (int, error)
	GetFaultByID(id string) (*Fault, int, error)
	GetJSONParaNotificacion(faltaID string) (string, error)
	GetAllFault() ([]*FaultWithStudent, error)
	CreateDocumentoFalta(faltaId string, url string) error
	SubirDocumentoFalta(faltaID, fileName string, data []byte) error
	DescargarDocumento(docID string) (*FaultDocumento, error)
	GetDetalleFaltaAgrupado(faltaID string) (DetalleFaltaAgrupado, error)
	// Nuevo método para exponer el resumen de incisos
	GetResumenIncisosAlumno(alumnoID int64, faltaActualID string) (IncisosAlumnoResumen, error)
}

// Implementación concreta del servicio
type service struct {
	repository ServicesFaultRepository
	user       *models.User
	txID       string
}

func NewFaultService(repository ServicesFaultRepository, user *models.User, TxID string) PortsServerFault {
	return &service{repository: repository, user: user, txID: TxID}
}

// Resumen de incisos cometidos por un alumno
type IncisosAlumnoResumen struct {
	IncisosFaltaActual int                   `json:"incisos_falta_actual"`
	IncisosTotales     int                   `json:"incisos_totales"`
	LevesTotales       int                   `json:"leves_totales"`
	GravesTotales      int                   `json:"graves_totales"`
	LevesFaltaActual   int                   `json:"leves_falta_actual"`
	GravesFaltaActual  int                   `json:"graves_falta_actual"`
	PrimerIncisoLeve   bool                  `json:"primer_inciso_leve"`
	PuedeSuspender     bool                  `json:"puede_suspender"`
	PuedeExpulsion     bool                  `json:"puede_expulsion"`
	Detalles           []*FaultIncisoDetalle `json:"detalles"`
}

// Devuelve el resumen de incisos cometidos por un alumno
func (s *service) GetResumenIncisosAlumno(alumnoID int64, faltaActualID string) (IncisosAlumnoResumen, error) {
	detalles, err := s.repository.GetAllIncisosByAlumnoID(alumnoID)
	if err != nil {
		return IncisosAlumnoResumen{}, err
	}
	var incisosFaltaActual, levesTotales, gravesTotales, levesFaltaActual, gravesFaltaActual int
	for _, d := range detalles {
		if d.FaultID == faltaActualID {
			incisosFaltaActual++
			if d.Gravedad == "leve" {
				levesFaltaActual++
			} else if d.Gravedad == "grave" {
				gravesFaltaActual++
			}
		}
		if d.Gravedad == "leve" {
			levesTotales++
		} else if d.Gravedad == "grave" {
			gravesTotales++
		}
	}
	primerIncisoLeve := levesTotales == 1
	puedeSuspender := levesTotales >= 2
	puedeExpulsion := gravesTotales > 0 || levesTotales >= 2
	return IncisosAlumnoResumen{
		IncisosFaltaActual: incisosFaltaActual,
		IncisosTotales:     len(detalles),
		LevesTotales:       levesTotales,
		GravesTotales:      gravesTotales,
		LevesFaltaActual:   levesFaltaActual,
		GravesFaltaActual:  gravesFaltaActual,
		PrimerIncisoLeve:   primerIncisoLeve,
		PuedeSuspender:     puedeSuspender,
		PuedeExpulsion:     puedeExpulsion,
		Detalles:           detalles,
	}, nil
}

// Crear Falta
func (s *service) CreateFault(
	id string,
	alumnoID int64,
	servicioId int64,
	convocatoriaId int64,
	fuenteInformacion string,
	fechaFalta time.Time,
	estado string,
	observacion string,
	articulos []string,
	incisos []string,
) (*Fault, int, error) {
	// ✅ Si estado viene vacío, poner "registrada"
	if strings.TrimSpace(estado) == "" {
		estado = "registrada"
	}
	// ✅ Ajustar a la nueva firma de NewFault
	m := NewFault(id, alumnoID, servicioId, convocatoriaId, fuenteInformacion, fechaFalta, estado, true, "", "", observacion)

	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if strings.Contains(err.Error(), "ecatch:108") {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Fault:", err)
		return m, 3, err
	}

	for _, articuloID := range articulos {
		fa := NewFaultArticulo(uuid.New().String(), id, articuloID)
		if err := s.repository.createFaultArticulo(fa); err != nil {
			logger.Error.Println(s.txID, " - couldn't create FaultArticulo:", err)
			return m, 3, err
		}
	}

	for _, incisoID := range incisos {
		fi := NewFaultInciso(uuid.New().String(), id, incisoID)
		if err := s.repository.createFaultInciso(fi); err != nil {
			logger.Error.Println(s.txID, " - couldn't create FaultInciso:", err)
			return m, 3, err
		}
	}

	return m, 29, nil
}

// Obtener Detalle de Falta Agrupado
func (s *service) GetDetalleFaltaAgrupado(faltaID string) (DetalleFaltaAgrupado, error) {
	if !govalidator.IsUUID(faltaID) {
		return DetalleFaltaAgrupado{}, fmt.Errorf("ID de falta inválido")
	}

	detalles, err := s.repository.GetDetalleFalta(faltaID)
	if err != nil {
		return DetalleFaltaAgrupado{}, err
	}
	if len(detalles) == 0 {
		return DetalleFaltaAgrupado{}, fmt.Errorf("no existe la falta")
	}
	nombreServicio, _ := s.repository.GetServicioNombreByID(detalles[0].ServicioID)
	agrupado := AgruparDetalleFalta(detalles, nombreServicio)
	return agrupado, nil
}

// Crear Documento asociado a Falta
func (s *service) CreateDocumentoFalta(faultID string, url string) error {
	if !govalidator.IsUUID(faultID) {
		return fmt.Errorf("falta_id inválido")
	}

	doc := NewFaultDocumento(uuid.New().String(), faultID, url)
	if err := s.repository.CreateFaultDocumento(doc); err != nil {
		logger.Error.Println(s.txID, " - error al guardar documento:", err)
		return err
	}

	return nil
}

// Subir documento BLOB
func (s *service) SubirDocumentoFalta(faltaID, fileName string, data []byte) error {
	if !govalidator.IsUUID(faltaID) {
		return fmt.Errorf("ID de falta inválido")
	}
	if len(data) == 0 {
		return fmt.Errorf("Archivo vacío")
	}
	doc := &FaultDocumento{
		ID:        uuid.New().String(),
		FaultID:   faltaID,
		URL:       fileName,
		Archivo:   data,
		CreatedAt: time.Now(),
	}
	return s.repository.CreateFaultDocumento(doc)
}

// Descargar Documento BLOB
func (s *service) DescargarDocumento(docID string) (*FaultDocumento, error) {
	if !govalidator.IsUUID(docID) {
		return nil, fmt.Errorf("ID de documento inválido")
	}
	return s.repository.GetFaultDocumentoByID(docID)
}

// Actualizar Falta
func (s *service) UpdateFault(id string, alumnoID int64, servicioId int64, convocatoriaId int64, fuenteInformacion string, fechaFalta time.Time, estado string, observacion string) (*Fault, int, error) {
	if strings.TrimSpace(estado) == "" {
		estado = "registrada"
	}
	// ✅ Ajustar a la nueva firma de NewFault
	m := NewFault(id, alumnoID, servicioId, convocatoriaId, fuenteInformacion, fechaFalta, estado, false, "", "", observacion)

	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Fault:", err)
		return m, 18, err
	}
	return m, 29, nil
}
func (s *service) DeleteFault(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}
	if err := s.repository.delete(id); err != nil {
		if strings.Contains(err.Error(), "ecatch:108") {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't delete Fault:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetFaultByID(id string) (*Fault, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllFault() ([]*FaultWithStudent, error) {
	allFaults, err := s.repository.getAll()
	if err != nil {
		return nil, err
	}

	for i := range allFaults {
		gravedades := strings.Split(allFaults[i].Gravedades, ",")
		if len(gravedades) > 0 {
			if contains(gravedades, "grave") {
				allFaults[i].Gravedad = "grave"
			} else if contains(gravedades, "leve") {
				allFaults[i].Gravedad = "leve"
			}
		}
	}

	return allFaults, nil
}

// Utilidad para buscar valores en slice
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// Agrupar detalles de la Falta
func AgruparDetalleFalta(detalles []*FaultDetalle, nombreServicio string) DetalleFaltaAgrupado {
	if len(detalles) == 0 {
		return DetalleFaltaAgrupado{}
	}

	agrupado := DetalleFaltaAgrupado{
		FaltaID:           detalles[0].FaltaID,
		FuenteInformacion: detalles[0].FuenteInformacion,
		Alumno: Alumno{
			DNI:                 detalles[0].DNI,
			Nombres:             detalles[0].Nombres,
			ApellidoPaterno:     detalles[0].ApellidoPaterno,
			ApellidoMaterno:     detalles[0].ApellidoMaterno,
			Sexo:                detalles[0].Sexo,
			Facultad:            detalles[0].Facultad,
			EscuelaProfesional:  detalles[0].EscuelaProfesional,
			Edad:                detalles[0].Edad,
			CorreoInstitucional: detalles[0].CorreoInstitucional,
			Direccion:           detalles[0].Direccion,
			LugarProcedencia:    detalles[0].LugarProcedencia,
			CelularEstudiante:   detalles[0].CelularEstudiante,
		},
		Resolucion: ResolucionDetalle{
			ResolucionID:     detalles[0].ResolucionID,
			ResolucionNombre: detalles[0].ResolucionNombre,
		},
		FechaFalta: detalles[0].FechaFalta.Format("2006-01-02"),
		Servicio:   nombreServicio,
	}

	documentoSet := make(map[string]struct{})

	for _, d := range detalles {
		if d.DocumentoURL != "" {
			documentoSet[d.DocumentoURL] = struct{}{}
		}

		var capitulo *CapituloDetalle
		for i := range agrupado.Resolucion.Capitulos {
			if agrupado.Resolucion.Capitulos[i].CapituloID == d.CapituloID {
				capitulo = &agrupado.Resolucion.Capitulos[i]
				break
			}
		}

		if capitulo == nil && d.CapituloID != "" {
			newCap := CapituloDetalle{
				CapituloID:     d.CapituloID,
				CapituloNombre: d.CapituloNombre,
			}
			agrupado.Resolucion.Capitulos = append(agrupado.Resolucion.Capitulos, newCap)
			capitulo = &agrupado.Resolucion.Capitulos[len(agrupado.Resolucion.Capitulos)-1]
		}

		if capitulo != nil {
			var articulo *ArticuloDetalle
			for i := range capitulo.Articulos {
				if capitulo.Articulos[i].ArticuloID == d.ArticuloID {
					articulo = &capitulo.Articulos[i]
					break
				}
			}

			if articulo == nil && d.ArticuloID != "" {
				newArt := ArticuloDetalle{
					ArticuloID:          d.ArticuloID,
					ArticuloDescripcion: d.ArticuloDescripcion,
					ArticuloGravedad:    d.ArticuloGravedad,
				}
				capitulo.Articulos = append(capitulo.Articulos, newArt)
				articulo = &capitulo.Articulos[len(capitulo.Articulos)-1]
			}

			if articulo != nil && d.IncisoID != "" {
				existsInciso := false
				for _, i := range articulo.Incisos {
					if i.IncisoID == d.IncisoID {
						existsInciso = true
						break
					}
				}
				if !existsInciso {
					articulo.Incisos = append(articulo.Incisos, IncisoDetalle{
						IncisoID:          d.IncisoID,
						IncisoNombre:      d.IncisoNombre,
						IncisoDescripcion: d.IncisoDescripcion,
					})
				}
			}
		}
	}

	for doc := range documentoSet {
		agrupado.Documentos = append(agrupado.Documentos, doc)
	}

	return agrupado
}

// Obtener JSON para notificación
func (s *service) GetJSONParaNotificacion(faltaID string) (string, error) {
	if !govalidator.IsUUID(faltaID) {
		return "", fmt.Errorf("ID de falta inválido")
	}

	detalles, err := s.repository.GetDetalleFalta(faltaID)
	if err != nil {
		return "", fmt.Errorf("error al obtener detalles de la falta: %v", err)
	}
	if len(detalles) == 0 {
		return "", fmt.Errorf("no existe la falta")
	}

	nombreServicio, err := s.repository.GetServicioNombreByID(detalles[0].ServicioID)
	if err != nil {
		nombreServicio = ""
	}

	agrupado := AgruparDetalleFalta(detalles, nombreServicio)

	jsonBytes, err := json.MarshalIndent(agrupado, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error al convertir a JSON: %v", err)
	}

	return string(jsonBytes), nil
}
