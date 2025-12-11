package assignment

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"dbu-api/pkg/residence/asignaciones_cuartos"
	"dbu-api/pkg/submission/alumnos"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type handlerStudentAssignment struct {
	db                    *sqlx.DB
	roomAssingmentService asignaciones_cuartos.RoomAssignmentServer
	studentService        alumnos.PortsServerAlumnos
}

func NewStudentAssignment(db *sqlx.DB) *handlerStudentAssignment {
	roomRepo := asignaciones_cuartos.FactoryStorage(db, nil, "")
	roomService := asignaciones_cuartos.NewRoomAssignmentService(roomRepo, nil, "")

	studentRepo := alumnos.FactoryStorage(db, nil, "")
	studentService := alumnos.NewAlumnosService(studentRepo, nil, "")

	return &handlerStudentAssignment{
		db:                    db,
		roomAssingmentService: roomService,
		studentService:        studentService,
	}
}

// GetSubmissions obtiene todas las convocatorias con detalle completo del estudiante
// @Summary Listar convocatorias con detalle del estudiante
// @Description Obtiene todas las convocatorias donde el estudiante tiene asignaciones, incluyendo información completa de residencia, cuarto, compañeros y objetos
// @Tags Student Assignment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]SubmissionWithDetailResponse}
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/student/submissions [get]
func (h *handlerStudentAssignment) GetSubmissions(c *fiber.Ctx) error {
	// Obtener alumno_id del contexto (viene del middleware)
	alumnoID, ok := c.Locals("alumno_id").(int64)
	if !ok || alumnoID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Error: true,
			Msg:   "No se pudo identificar al estudiante",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Obtener convocatorias del estudiante
	submissions, code, err := h.roomAssingmentService.GetSubmissionsByStudent(alumnoID)
	if err != nil {
		logger.Error.Printf("error getting submissions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Error: true,
			Msg:   "Error al obtener las convocatorias",
			Data:  nil,
			Code:  code,
			Type:  "error",
		})
	}

	// Si no hay convocatorias, retornar array vacío
	if submissions == nil || len(submissions) == 0 {
		return c.Status(fiber.StatusOK).JSON(models.Response{
			Error: false,
			Msg:   "No tienes asignaciones registradas",
			Data:  []SubmissionWithDetailResponse{},
			Code:  200,
			Type:  "success",
		})
	}

	// Obtener detalle completo de cada convocatoria
	response := make([]SubmissionWithDetailResponse, 0, len(submissions))
	for _, sub := range submissions {
		// Obtener detalle de asignación y compañeros
		detail, roommates, detailCode, err := h.roomAssingmentService.GetAssignmentDetailForStudent(alumnoID, sub.ConvocatoriaID)
		if err != nil {
			logger.Error.Printf("error getting assignment detail for submission %d: %v", sub.ConvocatoriaID, err)
			// Si no se encuentra el detalle, continuar con la siguiente convocatoria
			if detailCode == 4 {
				continue
			}
			// Para otros errores, retornar error
			return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
				Error: true,
				Msg:   "Error al obtener detalle de asignación",
				Data:  nil,
				Code:  detailCode,
				Type:  "error",
			})
		}

		// Mapear compañeros
		companeros := make([]CompaneroInfo, 0)
		if roommates != nil {
			for _, rm := range roommates {
				nombreCompleto := rm.Nombres + " " + rm.ApellidoPaterno + " " + rm.ApellidoMaterno
				carrera := rm.EscuelaProfesional
				if rm.Facultad != "" {
					carrera = rm.Facultad + " - " + rm.EscuelaProfesional
				}
				companeros = append(companeros, CompaneroInfo{
					CodigoEstudiante:    rm.CodigoEstudiante,
					NombreCompleto:      nombreCompleto,
					Carrera:             carrera,
					CorreoInstitucional: rm.CorreoInstitucional,
				})
			}
		}

		// Helper para convertir punteros a strings
		residenciaDesc := ""
		if detail.ResidenciaDescripcion != nil {
			residenciaDesc = *detail.ResidenciaDescripcion
		}

		// Agregar a respuesta
		response = append(response, SubmissionWithDetailResponse{
			Convocatoria: ConvocatoriaInfo{
				Nombre: detail.ConvocatoriaNombre,
			},
			Residencia: ResidenciaInfo{
				Nombre:      detail.ResidenciaNombre,
				Direccion:   detail.ResidenciaDireccion,
				Descripcion: residenciaDesc,
			},
			Cuarto: CuartoInfo{
				Numero:       detail.CuartoNumero,
				Piso:         detail.CuartoPiso,
				Capacidad:    detail.CuartoCapacidad,
				FechaIngreso: detail.FechaAsignacion,
			},
			Companeros: companeros,
			Objetos:    []ObjetoInfo{}, // Por ahora vacío
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Error: false,
		Msg:   "success",
		Data:  response,
		Code:  200,
		Type:  "success",
	})
}

// GetAssignmentDetail obtiene el detalle completo de la asignación del estudiante para una convocatoria específica
// @Summary Obtener detalle de asignación
// @Description Retorna información completa de la asignación: residencia, cuarto, compañeros y objetos asignados
// @Tags Student Assignment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param convocatoria_id path int true "ID de la convocatoria"
// @Success 200 {object} models.Response{data=AssignmentDetailResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/student/assignment/{convocatoria_id} [get]
func (h *handlerStudentAssignment) GetAssignmentDetail(c *fiber.Ctx) error {
	// Obtener alumno_id del contexto
	alumnoID, ok := c.Locals("alumno_id").(int64)
	if !ok || alumnoID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Error: true,
			Msg:   "No se pudo identificar al estudiante",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Obtener y validar convocatoria_id del path
	convocatoriaIDStr := c.Params("convocatoria_id")
	if convocatoriaIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Error: true,
			Msg:   "ID de convocatoria es requerido",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	convocatoriaID, err := strconv.ParseInt(convocatoriaIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Error: true,
			Msg:   "ID de convocatoria inválido",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Obtener detalle de asignación y compañeros
	detail, roommates, code, err := h.roomAssingmentService.GetAssignmentDetailForStudent(alumnoID, convocatoriaID)
	if err != nil {
		logger.Error.Printf("error getting assignment detail: %v", err)

		if code == 4 {
			return c.Status(fiber.StatusNotFound).JSON(models.Response{
				Error: true,
				Msg:   "No se encontró asignación para esta convocatoria",
				Data:  nil,
				Code:  code,
				Type:  "error",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Error: true,
			Msg:   "Error al obtener detalle de asignación",
			Data:  nil,
			Code:  code,
			Type:  "error",
		})
	}

	// Mapear compañeros
	companeros := make([]CompaneroInfo, 0)
	if roommates != nil {
		for _, rm := range roommates {
			nombreCompleto := rm.Nombres + " " + rm.ApellidoPaterno + " " + rm.ApellidoMaterno
			carrera := rm.EscuelaProfesional
			if rm.Facultad != "" {
				carrera = rm.Facultad + " - " + rm.EscuelaProfesional
			}
			companeros = append(companeros, CompaneroInfo{
				CodigoEstudiante:    rm.CodigoEstudiante,
				NombreCompleto:      nombreCompleto,
				Carrera:             carrera,
				CorreoInstitucional: rm.CorreoInstitucional,
			})
		}
	}

	// Helper para convertir punteros a strings
	residenciaDesc := ""
	if detail.ResidenciaDescripcion != nil {
		residenciaDesc = *detail.ResidenciaDescripcion
	}

	// Preparar response
	response := AssignmentDetailResponse{
		Convocatoria: ConvocatoriaInfo{
			Nombre: detail.ConvocatoriaNombre,
		},
		Residencia: ResidenciaInfo{
			Nombre:      detail.ResidenciaNombre,
			Direccion:   detail.ResidenciaDireccion,
			Descripcion: residenciaDesc,
		},
		Cuarto: CuartoInfo{
			Numero:       detail.CuartoNumero,
			Piso:         detail.CuartoPiso,
			Capacidad:    detail.CuartoCapacidad,
			FechaIngreso: detail.FechaAsignacion,
		},
		Companeros: companeros,
		Objetos:    []ObjetoInfo{}, // Por ahora vacío, se implementará después
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Error: false,
		Msg:   "success",
		Data:  response,
		Code:  200,
		Type:  "success",
	})
}

// GetProfile obtiene los datos personales básicos del estudiante autenticado
// @Summary Obtener perfil del estudiante
// @Description Retorna los datos personales básicos del estudiante autenticado (nombre, género, edad, DNI, código)
// @Tags Student Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{data=StudentProfileResponse}
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/student/profile [get]
func (h *handlerStudentAssignment) GetProfile(c *fiber.Ctx) error {
	// Obtener alumno_id del contexto (viene del middleware)
	alumnoID, ok := c.Locals("alumno_id").(int64)
	if !ok || alumnoID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Error: true,
			Msg:   "No se pudo identificar al estudiante",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Obtener perfil del estudiante usando el servicio
	student, code, err := h.studentService.GetStudentProfile(alumnoID)
	if err != nil {
		logger.Error.Printf("error getting student profile: %v", err)

		// Si el estudiante no existe, retornar 404
		if code == 4 {
			return c.Status(fiber.StatusNotFound).JSON(models.Response{
				Error: true,
				Msg:   "Estudiante no encontrado",
				Data:  nil,
				Code:  code,
				Type:  "error",
			})
		}

		// Error de validación
		if code == 15 {
			return c.Status(fiber.StatusBadRequest).JSON(models.Response{
				Error: true,
				Msg:   "Datos de entrada inválidos",
				Data:  nil,
				Code:  code,
				Type:  "error",
			})
		}

		// Error de base de datos u otro error interno
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Error: true,
			Msg:   "Error al obtener el perfil del estudiante",
			Data:  nil,
			Code:  code,
			Type:  "error",
		})
	}

	// Construir nombre completo
	nombreCompleto := student.Nombres + " " + student.ApellidoPaterno + " " + student.ApellidoMaterno

	// Mapear género (M/F a Masculino/Femenino)
	genero := "Masculino"
	if student.Sexo == "F" || student.Sexo == "femenino" {
		genero = "Femenino"
	}

	// Preparar response
	response := StudentProfileResponse{
		NombreCompleto:   nombreCompleto,
		Genero:           genero,
		Edad:             student.Edad,
		DNI:              student.DNI,
		CodigoEstudiante: student.CodigoEstudiante,
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Error: false,
		Msg:   "success",
		Data:  response,
		Code:  200,
		Type:  "success",
	})
}
