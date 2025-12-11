package faults

import (
	"database/sql"
	"dbu-api/internal/authorization"
	"dbu-api/internal/logger"
	"dbu-api/internal/middleware"
	"dbu-api/internal/models"
	"dbu-api/pkg/sanction"
	"dbu-api/pkg/sanction/fault"
	_ "encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// GetResumenIncisosAlumno godoc
// @Summary Obtener resumen de incisos cometidos por un alumno
// @Description Devuelve el resumen de incisos (leve/grave, totales, etc.) para un alumno y una falta actual
// @tags Faltas
// @Accept json
// @Produce json
// @Param alumno_id path int true "ID del alumno"
// @Param falta_id path string true "ID de la falta actual"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/Faltas/resumen-incisos/{alumno_id}/{falta_id} [get]

// GetAlumnoProfileByCodigo godoc
// @Summary Obtiene el perfil completo del alumno por código de estudiante
// @Description Retorna información completa del alumno desde la tabla alumnos
// @Tags Faltas
// @Accept json
// @Produce json
// @Param codigo_estudiante path string true "Código del alumno"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /v1/alumnos/perfil-codigo/{codigo_estudiante} [GET]
func (h *handlerFaults) GetAlumnoProfileByCodigo(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	codigo := c.Params("codigo_estudiante")

	if codigo == "" {
		res.Code, res.Type, res.Msg = 1, "error", "Código de estudiante requerido"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Validamos permisos si aplica
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	err = authorization.ValidPermissions(user, h.db, c)
	if err != nil {
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Consultamos al alumno con validación de cuarto asignado y datos adicionales
	type AlumnoDetalle struct {
		AlumnoID                int64  `db:"alumno_id" json:"alumno_id"`
		DNI                     string `db:"DNI" json:"dni"`
		FullName                string `db:"full_name" json:"full_name"`
		Code                    string `db:"code" json:"code"`
		ProfessionalSchool      string `db:"professional_school" json:"professional_school"`
		Faculty                 string `db:"faculty" json:"faculty"`
		RoomNumber              string `db:"room_number" json:"room_number"`
		ResidenceName           string `db:"residence_name" json:"residence_name"`
		AdmissionDate           string `db:"admission_date" json:"admission_date"`
		CelularEstudiante       string `db:"celular_estudiante" json:"celular_estudiante"`
		CelularPadre            string `db:"celular_padre" json:"celular_padre"`
		DepartamentoProcedencia string `db:"departamento_procedencia" json:"departamento_procedencia"` // Added
		ProvinciaProcedencia    string `db:"provincia_procedencia" json:"provincia_procedencia"`       // Added
		DistritoProcedencia     string `db:"distrito_procedencia" json:"distrito_procedencia"`         // Added
		LugarProcedencia        string `db:"lugar_procedencia" json:"lugar_procedencia"`               // Added
		Direccion               string `db:"direccion" json:"direccion"`
	}

	query := `
		SELECT 
			a.id AS alumno_id,
			a.DNI,
			CONCAT(a.nombres, ' ', a.apellido_paterno, ' ', a.apellido_materno) AS full_name,
			a.codigo_estudiante AS code,
			a.escuela_profesional AS professional_school,
			a.facultad AS faculty,
			-- Cuarto y residencia
			COALESCE((SELECT CAST(c.numero AS CHAR)
			 FROM asignacion_cuartos ac
			 JOIN cuartos c ON c.id = ac.cuarto_id
			 WHERE ac.alumno_id = a.id AND ac.estado = 'activo'
			 LIMIT 1), '') AS room_number,
			COALESCE((SELECT r.nombre
			 FROM asignacion_cuartos ac
			 JOIN cuartos c ON c.id = ac.cuarto_id
			 JOIN residencias r ON r.id = c.residencia_id
			 WHERE ac.alumno_id = a.id AND ac.estado = 'activo'
			 LIMIT 1), '') AS residence_name,
			-- Fecha de asignación o updated_at
			COALESCE((SELECT ac.fecha_asignacion FROM asignacion_cuartos ac WHERE ac.alumno_id = a.id AND ac.estado = 'activo' LIMIT 1), '') AS admission_date,
			-- Celular estudiante
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'celular de estudiante'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS celular_estudiante,
			-- Celular padre
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Celular padre'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS celular_padre,
			-- Departamento de procedencia (lookup por opcion_seleccion -> departaments.name)
			COALESCE((
				SELECT d.name
				FROM solicitudes sol
				JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
				JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
				JOIN requisitos req ON req.id = ds.requisito_id
				JOIN departaments d ON d.id = ds.opcion_seleccion
				WHERE sol.alumno_id = a.id AND req.nombre = 'Departamento de procedencia'
				ORDER BY sol.created_at DESC
				LIMIT 1
			), '') AS departamento_procedencia,
			-- Provincia de procedencia (lookup por opcion_seleccion -> provinces.name)
			COALESCE((
				SELECT p.name
				FROM solicitudes sol
				JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
				JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
				JOIN requisitos req ON req.id = ds.requisito_id
				JOIN provinces p ON p.id = ds.opcion_seleccion
				WHERE sol.alumno_id = a.id AND req.nombre = 'Provincia de procedencia'
				ORDER BY sol.created_at DESC
				LIMIT 1
			), '') AS provincia_procedencia,
			-- Distrito de procedencia (lookup por opcion_seleccion -> districts.name)
			COALESCE((
				SELECT dist.name
				FROM solicitudes sol
				JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
				JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
				JOIN requisitos req ON req.id = ds.requisito_id
				JOIN districts dist ON dist.id = ds.opcion_seleccion
				WHERE sol.alumno_id = a.id AND req.nombre = 'Distrito de procedencia'
				ORDER BY sol.created_at DESC
				LIMIT 1
			), '') AS distrito_procedencia,
			-- Lugar de procedencia concatenado (departamento/provincia/distrito)
			CONCAT(
				COALESCE((
					SELECT d.name
					FROM solicitudes sol
					JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
					JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
					JOIN requisitos req ON req.id = ds.requisito_id
					JOIN departaments d ON d.id = ds.opcion_seleccion
					WHERE sol.alumno_id = a.id AND req.nombre = 'Departamento de procedencia'
					ORDER BY sol.created_at DESC
					LIMIT 1
				), ''),
				'/',
				COALESCE((
					SELECT p.name
					FROM solicitudes sol
					JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
					JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
					JOIN requisitos req ON req.id = ds.requisito_id
					JOIN provinces p ON p.id = ds.opcion_seleccion
					WHERE sol.alumno_id = a.id AND req.nombre = 'Provincia de procedencia'
					ORDER BY sol.created_at DESC
					LIMIT 1
				), ''),
				'/',
				COALESCE((
					SELECT dist.name
					FROM solicitudes sol
					JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
					JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
					JOIN requisitos req ON req.id = ds.requisito_id
					JOIN districts dist ON dist.id = ds.opcion_seleccion
					WHERE sol.alumno_id = a.id AND req.nombre = 'Distrito de procedencia'
					ORDER BY sol.created_at DESC
					LIMIT 1
				), '')
			) AS lugar_procedencia,
			-- Dirección: residencia/cuarto si tiene, blanco si no
			CASE WHEN EXISTS (
				SELECT 1 FROM asignacion_cuartos ac WHERE ac.alumno_id = a.id AND ac.estado = 'activo'
			) THEN CONCAT(
				COALESCE((SELECT r.nombre FROM asignacion_cuartos ac JOIN cuartos c ON c.id = ac.cuarto_id JOIN residencias r ON r.id = c.residencia_id WHERE ac.alumno_id = a.id AND ac.estado = 'activo' LIMIT 1), ''),
				' / Cuarto: ',
				COALESCE((SELECT CAST(c.numero AS CHAR) FROM asignacion_cuartos ac JOIN cuartos c ON c.id = ac.cuarto_id WHERE ac.alumno_id = a.id AND ac.estado = 'activo' LIMIT 1), '')
			) ELSE '' END AS direccion
		FROM alumnos a
		WHERE a.codigo_estudiante = ?
		LIMIT 1
	`

	var alumno AlumnoDetalle
	err = h.db.Get(&alumno, query, codigo)
	if err != nil {
		logger.Error.Printf("No se encontró el alumno con código %s, error: %v", codigo, err)
		res.Code, res.Type, res.Msg = 1, "error", "No se encontró alumno con ese código"
		return c.Status(fiber.StatusNotFound).JSON(res)
	}

	res.Data = alumno
	res.Error = false
	return c.JSON(res)
}
func (h *handlerFaults) GetResumenIncisosAlumno(c *fiber.Ctx) error {
	alumnoIDStr := c.Params("alumno_id")
	faltaID := c.Params("falta_id")
	if alumnoIDStr == "" || faltaID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "alumno_id y falta_id son requeridos"})
	}
	var alumnoID int64
	_, err := fmt.Sscan(alumnoIDStr, &alumnoID)
	if err != nil || alumnoID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "alumno_id debe ser un número válido"})
	}
	user, err := middleware.GetUser(c.Get("Authorization"), h.db)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "No autorizado"})
	}
	srv := sanction.NewServerSanction(h.db, user, h.txID)
	resumen, err := srv.SrvFault.GetResumenIncisosAlumno(alumnoID, faltaID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(resumen)
}

type handlerFaults struct {
	db   *sqlx.DB
	txID string
}

// UpdateEstadoFalta actualiza el estado de una falta
// @Summary Actualiza el estado de una falta
// @Description Permite actualizar el campo estado de una falta
// @tags Faltas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la falta"
// @Param body body struct{Estado string `json:"estado"`} true "Nuevo estado"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /v1/faltas/{id}/estado [patch]
func (h *handlerFaults) UpdateEstadoFalta(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	id := c.Params("id")
	if id == "" {
		res.Code, res.Type, res.Msg = 1, "error", "ID de falta requerido"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	var body struct {
		Estado string `json:"estado"`
	}
	if err := c.BodyParser(&body); err != nil {
		res.Code, res.Type, res.Msg = 1, "error", "Error al leer el cuerpo"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	if body.Estado == "" {
		res.Code, res.Type, res.Msg = 1, "error", "Estado requerido"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	// Actualizar en BD
	err := h.updateEstadoFaltaDB(id, body.Estado)
	if err != nil {
		res.Code, res.Type, res.Msg = 2, "error", err.Error()
		return c.Status(http.StatusInternalServerError).JSON(res)
	}
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Estado actualizado"
	return c.Status(http.StatusOK).JSON(res)
}

// updateEstadoFaltaDB actualiza el estado en la base de datos
func (h *handlerFaults) updateEstadoFaltaDB(id string, estado string) error {
	query := "UPDATE faltas SET estado = ? WHERE id = ?"
	_, err := h.db.Exec(query, estado, id)
	return err
}

// CreateFaults godoc
// @Summary Crear una instancia de Faltas
// @Description Método que permite crear una instancia del objeto Faltas en la base de datos
// @tags Faltas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param models.Fault body models.Fault true "Datos para crear Faltas"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 202 {object} models.Response
// @Router /v1/Faltas [POST]
func (h *handlerFaults) CreateFault(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	req := models.Fault{}

	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	err = authorization.ValidPermissions(user, h.db, c)
	if err != nil {
		logger.Error.Printf("User does not have permission to call the api, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "error", "invalid request body"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	isValid, err := req.ValidFault()
	if err != nil || !isValid {
		logger.Error.Printf("couldn't validate body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "error", "validation failed"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	srv := sanction.NewServerSanction(h.db, user, h.txID)

	// Enviamos el DNI directamente desde el campo alumno.dni
	fault, code, err := srv.SrvFault.CreateFault(
		req.ID,
		req.AlumnoID,
		req.ServicioId,
		req.ConvocatoriaId,
		req.FuenteInformacion,
		req.FechaFalta,
		req.Estado,
		req.Observacion,
		req.Articulos,
		req.Incisos,
	)
	if err != nil {
		logger.Error.Printf("couldn't create Faults, error: %v", err)
		res.Code, res.Type, res.Msg = code, "error", err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = fault
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Falta creada correctamente"
	return c.Status(http.StatusCreated).JSON(res)
}

// UpdateFaults godoc
// @Summary Actualiza una instancia de Faltas
// @Description Método que permite Actualiza una instancia del objeto Faltas en la base de datos
// @tags Faltas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param models.Fault body models.Fault true "Datos para actualizar Faltas"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 202 {object} models.Response
// @Router /v1/Faltas [PUT]
func (h *handlerFaults) UpdateFault(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	req := models.Fault{}

	bearer := c.Get("Authorization")

	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	err = authorization.ValidPermissions(user, h.db, c)
	if err != nil {
		logger.Error.Printf("User does not have permission to call the api, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "", "Error al leer el cuerpo de la solicitud"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	isValid, err := req.ValidFault()
	if err != nil || !isValid {
		logger.Error.Printf("invalid request data, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "", "Datos inválidos para la falta"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	srv := sanction.NewServerSanction(h.db, user, h.txID)

	// Aquí simplemente pasamos el alumno_id como int64 directamente
	fault, code, err := srv.SrvFault.UpdateFault(
		req.ID,
		req.AlumnoID, // 👈🏻 ya debe ser int64 en el modelo
		req.ServicioId,
		req.ConvocatoriaId,
		req.FuenteInformacion,
		req.FechaFalta,
		req.Estado,
		req.Observacion,
	)
	if err != nil {
		logger.Error.Printf("couldn't update Faltas, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", "No se pudo actualizar la falta"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = fault
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", "Falta actualizada correctamente"
	return c.Status(http.StatusOK).JSON(res)
}

// DeleteFaults godoc
// @Summary Elimina una instancia de Faltas
// @Description Método que permite eliminar una instancia del objeto Faltas en la base de datos
// @tags Faltas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param	id	path string true "Faltas ID" format(uuid)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 202 {object} models.Response
// @Router /v1/Faltas [DELETE]
func (h *handlerFaults) DeleteFaults(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	bearer := c.Get("Authorization")

	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	err = authorization.ValidPermissions(user, h.db, c)
	if err != nil {
		logger.Error.Printf("User does not have permission to call the api, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	idStr := c.Params("id")
	if idStr == "" {
		logger.Error.Println("couldn't parse id request")
		res.Code, res.Type, res.Msg = 1, "", ""
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	srv := sanction.NewServerSanction(h.db, user, h.txID)
	code, err := srv.SrvFault.DeleteFault(idStr)
	if err != nil {
		logger.Error.Printf("couldn't delete Fault, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

// GetFaltasByID godoc
// @Summary Obtiene una instancia de Faltas por su id
// @Description Método que permite obtener una instancia del objeto Faltas en la base de datos por su id
// @tags Faltas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param	id	path string true "Faltas ID" format(uuid)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 202 {object} models.Response
// @Router /v1/Faltas/:{id} [GET]
func (h *handlerFaults) GetFaultsByID(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	bearer := c.Get("Authorization")

	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	err = authorization.ValidPermissions(user, h.db, c)
	if err != nil {
		logger.Error.Printf("User does not have permission to call the api, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	idStr := c.Params("id")
	if idStr == "" {
		logger.Error.Println("couldn't parse id request")
		res.Code, res.Type, res.Msg = 1, "", ""
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	srv := sanction.NewServerSanction(h.db, user, h.txID)
	data, code, err := srv.SrvFault.GetFaultByID(idStr)
	if err != nil {
		logger.Error.Printf("couldn't get faults by id, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = data
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

// / GetAllFaults godoc
// @Summary Obtiene todas las instancias de Faltas
// @Description Método que permite obtener todas las instancias del objeto Faltas en la base de datos
// @tags Faltas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 202 {object} models.Response
// @Router /v1/Faltas [GET]
func (h *handlerFaults) GetAllFaults(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	bearer := c.Get("Authorization")

	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	srv := sanction.NewServerSanction(h.db, user, h.txID)

	// Traemos todas las faltas con datos del alumno
	allFaults, err := srv.SrvFault.GetAllFault()
	if err != nil {
		logger.Error.Printf("couldn't get all faults, error: %v", err)
		res.Code, res.Type, res.Msg = 99, "error", err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = allFaults
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Faltas obtenidas correctamente"
	return c.Status(http.StatusOK).JSON(res)
}

// Función auxiliar para comprobar si un slice contiene un valor
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// GetAllStudentsByFault godoc
// @Summary Obtiene todos los alumnos según la Residencia
// @Description Método que permite obtener todas las instancias de alumnos según la Residencia en la base de datos
// @tags Faltas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param	id	path string true "Faltas ID" format(uuid)
// @Success 200 {object} models.Response
// @Failure 202 {object} models.Response
// @Router /v1/Faltas/:id/alumnos [GET]
func (h *handlerFaults) GetAllStudentsByFault(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	bearer := c.Get("Authorization")

	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	err = authorization.ValidPermissions(user, h.db, c)
	if err != nil {
		logger.Error.Printf("User does not have permission to call the api, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Listado de alumnos con lógica de cuarto, celular y procedencia
	type AlumnoDetalle struct {
		AlumnoID                int64  `db:"alumno_id" json:"alumno_id"`
		DNI                     string `db:"DNI" json:"dni"`
		FullName                string `db:"full_name" json:"full_name"`
		Code                    string `db:"code" json:"code"`
		ProfessionalSchool      string `db:"professional_school" json:"professional_school"`
		Faculty                 string `db:"faculty" json:"faculty"`
		RoomNumber              string `db:"room_number" json:"room_number"`
		ResidenceName           string `db:"residence_name" json:"residence_name"`
		AdmissionDate           string `db:"admission_date" json:"admission_date"`
		CelularEstudiante       string `db:"celular_estudiante" json:"celular_estudiante"`
		CelularPadre            string `db:"celular_padre" json:"celular_padre"`
		DepartamentoProcedencia string `db:"departamento_procedencia" json:"departamento_procedencia"`
		ProvinciaProcedencia    string `db:"provincia_procedencia" json:"provincia_procedencia"`
		DistritoProcedencia     string `db:"distrito_procedencia" json:"distrito_procedencia"`
	}

	// Puedes ajustar el WHERE según el filtro que uses (por ejemplo, por residencia, convocatoria, etc.)
	query := `
		SELECT 
			a.id AS alumno_id,
			a.DNI,
			CONCAT(a.nombres, ' ', a.apellido_paterno, ' ', a.apellido_materno) AS full_name,
			a.codigo_estudiante AS code,
			a.escuela_profesional AS professional_school,
			a.facultad AS faculty,
			-- Cuarto y residencia
			COALESCE((SELECT CAST(c.numero AS CHAR)
			 FROM asignacion_cuartos ac
			 JOIN cuartos c ON c.id = ac.cuarto_id
			 WHERE ac.alumno_id = a.id AND ac.estado = 'activo'
			 LIMIT 1), '') AS room_number,
			COALESCE((SELECT r.nombre
			 FROM asignacion_cuartos ac
			 JOIN cuartos c ON c.id = ac.cuarto_id
			 JOIN residencias r ON r.id = c.residencia_id
			 WHERE ac.alumno_id = a.id AND ac.estado = 'activo'
			 LIMIT 1), '') AS residence_name,
			-- Fecha de asignación o updated_at
			COALESCE((SELECT ac.fecha_asignacion FROM asignacion_cuartos ac WHERE ac.alumno_id = a.id AND ac.estado = 'activo' LIMIT 1), '') AS admission_date,
			-- Celular estudiante
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'celular de estudiante'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS celular_estudiante,
			-- Celular padre
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Celular padre'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS celular_padre,
			-- Departamento de procedencia
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Departamento de procedencia'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS departamento_procedencia,
			-- Provincia de procedencia
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Provincia de procedencia'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS provincia_procedencia,
			-- Distrito de procedencia
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Distrito de procedencia'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS distrito_procedencia
		FROM alumnos a
		WHERE EXISTS (
			SELECT 1 FROM asignacion_cuartos ac
			JOIN cuartos c ON c.id = ac.cuarto_id
			JOIN residencias r ON r.id = c.residencia_id
			WHERE r.id = ? AND ac.alumno_id = a.id
		)
	`

	residenciaID := c.Params("id")
	var alumnos []AlumnoDetalle
	err = h.db.Select(&alumnos, query, residenciaID)
	if err != nil {
		logger.Error.Printf("No se pudo obtener alumnos por residencia %s, error: %v", residenciaID, err)
		res.Code, res.Type, res.Msg = 99, "error", err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = alumnos
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

// GetAlumnoProfileByDNI godoc
// @Summary Obtiene el perfil completo del alumno por DNI
// @Description Retorna información completa del alumno desde la tabla alumnos
// @Tags Faltas
// @Accept json
// @Produce json
// @Param dni path string true "DNI del alumno"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /v1/alumnos/perfil/{dni} [GET]
func (h *handlerFaults) GetAlumnoProfileByDNI(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	dni := c.Params("dni")

	if dni == "" {
		res.Code, res.Type, res.Msg = 1, "error", "DNI requerido"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Validamos permisos si aplica
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	err = authorization.ValidPermissions(user, h.db, c)
	if err != nil {
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Consultamos al alumno con validación de cuarto asignado y datos adicionales
	type AlumnoDetalle struct {
		AlumnoID                int64  `db:"alumno_id" json:"alumno_id"`
		DNI                     string `db:"DNI" json:"dni"`
		FullName                string `db:"full_name" json:"full_name"`
		Code                    string `db:"code" json:"code"`
		ProfessionalSchool      string `db:"professional_school" json:"professional_school"`
		Faculty                 string `db:"faculty" json:"faculty"`
		RoomNumber              string `db:"room_number" json:"room_number"`
		ResidenceName           string `db:"residence_name" json:"residence_name"`
		AdmissionDate           string `db:"admission_date" json:"admission_date"`
		CelularEstudiante       string `db:"celular_estudiante" json:"celular_estudiante"`
		CelularPadre            string `db:"celular_padre" json:"celular_padre"`
		DepartamentoProcedencia string `db:"departamento_procedencia" json:"departamento_procedencia"`
		ProvinciaProcedencia    string `db:"provincia_procedencia" json:"provincia_procedencia"`
		DistritoProcedencia     string `db:"distrito_procedencia" json:"distrito_procedencia"`
	}

	query := `
		SELECT 
			a.id AS alumno_id,
			a.DNI,
			CONCAT(a.nombres, ' ', a.apellido_paterno, ' ', a.apellido_materno) AS full_name,
			a.codigo_estudiante AS code,
			a.escuela_profesional AS professional_school,
			a.facultad AS faculty,
			-- Cuarto y residencia
			COALESCE((SELECT CAST(c.numero AS CHAR)
			 FROM asignacion_cuartos ac
			 JOIN cuartos c ON c.id = ac.cuarto_id
			 WHERE ac.alumno_id = a.id AND ac.estado = 'activo'
			 LIMIT 1), '') AS room_number,
			COALESCE((SELECT r.nombre
			 FROM asignacion_cuartos ac
			 JOIN cuartos c ON c.id = ac.cuarto_id
			 JOIN residencias r ON r.id = c.residencia_id
			 WHERE ac.alumno_id = a.id AND ac.estado = 'activo'
			 LIMIT 1), '') AS residence_name,
			-- Fecha de asignación o updated_at
			COALESCE((SELECT ac.fecha_asignacion FROM asignacion_cuartos ac WHERE ac.alumno_id = a.id AND ac.estado = 'activo' LIMIT 1), '') AS admission_date,
			-- Celular estudiante
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'celular de estudiante'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS celular_estudiante,
			-- Celular padre
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Celular padre'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS celular_padre,
			-- Departamento de procedencia
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Departamento de procedencia'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS departamento_procedencia,
			-- Provincia de procedencia
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Provincia de procedencia'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS provincia_procedencia,
			-- Distrito de procedencia
			COALESCE((SELECT ds.respuesta_formulario
			 FROM solicitudes sol
			 JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id
			 JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
			 JOIN requisitos req ON req.id = ds.requisito_id
			 WHERE sol.alumno_id = a.id AND req.nombre = 'Distrito de procedencia'
			 ORDER BY sol.created_at DESC
			 LIMIT 1), '') AS distrito_procedencia
		FROM alumnos a
		WHERE a.DNI = ?
		LIMIT 1
	`

	var alumno AlumnoDetalle
	err = h.db.Get(&alumno, query, dni)
	if err != nil {
		logger.Error.Printf("No se encontró el alumno con DNI %s, error: %v", dni, err)
		res.Code, res.Type, res.Msg = 1, "error", "No se encontró alumno con ese DNI"
		return c.Status(fiber.StatusNotFound).JSON(res)
	}

	res.Data = alumno
	res.Error = false
	res.Code = 29
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *handlerFaults) CrearDocumentoFalta(c *fiber.Ctx) error {
	faltaId := c.Params("id")
	var payload struct {
		URL string `json:"url"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body inválido"})
	}
	fmt.Println("➡️ Recibiendo documento:", payload.URL, "para falta ID:", faltaId)
	user, err := middleware.GetUser(c.Get("Authorization"), h.db)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "No autorizado"})
	}
	srv := sanction.NewServerSanction(h.db, user, h.txID)
	if err := srv.SrvFault.CreateDocumentoFalta(faltaId, payload.URL); err != nil {
		fmt.Println("💥 Error al guardar documento:", err)
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo guardar documento"})
	}
	return c.JSON(fiber.Map{"message": "Documento guardado"})
}

// UploadDocumento godoc
// @Summary Sube un archivo y devuelve su URL
// @Description Carga un archivo y devuelve su ruta local simulada
// @Tags Faltas
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Archivo"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /v1/faltas/upload [post]
func (h *handlerFaults) UploadDocumento(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Archivo requerido"})
	}
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo abrir el archivo"})
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo leer el archivo"})
	}

	faltaID := c.FormValue("falta_id")
	if faltaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "falta_id requerido"})
	}

	user, err := middleware.GetUser(c.Get("Authorization"), h.db)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "No autorizado"})
	}
	srv := sanction.NewServerSanction(h.db, user, h.txID)
	if err := srv.SrvFault.SubirDocumentoFalta(faltaID, fileHeader.Filename, data); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Documento guardado en base de datos"})
}
func (h *handlerFaults) DescargarDocumento(c *fiber.Ctx) error {
	docID := c.Params("id")
	user, err := middleware.GetUser(c.Get("Authorization"), h.db)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "No autorizado"})
	}
	srv := sanction.NewServerSanction(h.db, user, h.txID)
	doc, err := srv.SrvFault.DescargarDocumento(docID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Documento no encontrado"})
	}

	// 2. Determinar tipo de archivo (mime) usando la extensión del nombre (campo url)
	ext := filepath.Ext(doc.URL)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 3. Soportar descarga o visualización según parámetro
	disposition := "inline"
	if c.Query("download") == "1" {
		disposition = "attachment"
	}

	c.Set("Content-Type", mimeType)
	c.Set("Content-Disposition", fmt.Sprintf(`%s; filename="%s"`, disposition, doc.URL))
	return c.Send(doc.Archivo)
}
func (h *handlerFaults) GetDocumentosPorFalta(c *fiber.Ctx) error {
	faltaID := c.Params("id")
	// Validación (opcional)
	if !govalidator.IsUUID(faltaID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de falta inválido"})
	}

	var docs []FaultDocumento
	err := h.db.Select(&docs, "SELECT id, url FROM faltas_documentos WHERE falta_id = ?", faltaID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudieron obtener los documentos"})
	}
	return c.JSON(docs)
}

func (h *handlerFaults) GetDetalleFalta(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	faltaID := c.Params("id")

	if !govalidator.IsUUID(faltaID) {
		res.Code, res.Type, res.Msg = 1, "error", "ID de falta inválido"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	var detalles []models.DetalleFalta // Define este struct en tu modelo

	const query = `...` // Usa la consulta SQL de arriba

	if err := h.db.Select(&detalles, query, sql.Named("falta_id", faltaID)); err != nil {
		res.Code, res.Type, res.Msg = 99, "error", err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res.Data = detalles
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Detalle obtenido"
	return c.Status(fiber.StatusOK).JSON(res)
}
func (h *handlerFaults) GetDetalleFaltaAgrupado(c *fiber.Ctx) error {
	faltaID := c.Params("id")
	if !govalidator.IsUUID(faltaID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de falta inválido"})
	}

	repository := fault.FactoryStorage(h.db, h.txID)
	detalles, err := repository.GetDetalleFalta(faltaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener los detalles", "details": err.Error()})
	}

	if len(detalles) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No existen detalles para la falta"})
	}

	// ¡Aquí obtienes el nombre del servicio!
	nombreServicio, err := repository.GetServicioNombreByID(detalles[0].ServicioID)
	if err != nil {
		nombreServicio = "" // O puedes retornar error si prefieres
	}

	agrupado := fault.AgruparDetalleFalta(detalles, nombreServicio)
	return c.Status(fiber.StatusOK).JSON(agrupado)
}

// ValidarPostulacionPorDNI godoc
// @Summary Valida si el estudiante puede postular a un servicio según faltas y sanciones vigentes
// @Description Valida por DNI y servicio si el estudiante tiene faltas vigentes en el servicio solicitado o sanción en residencia
// @Tags Faltas
// @Accept json
// @Produce json
// @Param dni path string true "DNI del alumno"
// @Param servicio_id path int true "ID del servicio a postular"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /v1/alumnos/validar-postulacion/{dni}/{servicio_id} [get]
// ...existing code...

// ValidarPostulacionPorDNI godoc
// @Summary Valida si el estudiante puede postular a un servicio según faltas y sanciones vigentes
// @Description Valida por DNI y servicio si el estudiante tiene faltas vigentes en el servicio solicitado o sanción en residencia
// @Tags Faltas
// @Accept json
// @Produce json
// @Param dni path string true "DNI del alumno"
// @Param servicio_id path int true "ID del servicio a postular"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /v1/alumnos/validar-postulacion/{dni}/{servicio_id} [get]
type FaltaDetalle struct {
	FaltaID      string `db:"falta_id"`
	ServicioID   int    `db:"servicio_id"`
	Observacion  string `db:"observacion"`
	Estado       string `db:"estado"`
	SancionID    string `db:"sancion_id"`
	DocumentoID  string `db:"documento_id"`
	DocumentoURL string `db:"documento_url"`
}

func (h *handlerFaults) ValidarPostulacionPorDNI(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	dni := c.Params("dni")
	servicioID := c.Params("servicio_id")
	if dni == "" || servicioID == "" {
		res.Code, res.Type, res.Msg = 1, "error", "DNI y servicio_id requeridos"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Obtener alumno
	var alumnoID int64
	err := h.db.Get(&alumnoID, "SELECT id FROM alumnos WHERE DNI = ? LIMIT 1", dni)
	if err != nil || alumnoID == 0 {
		res.Code, res.Type, res.Msg = 2, "error", "Alumno no encontrado"
		return c.Status(fiber.StatusNotFound).JSON(res)
	}

	hoy := "2025-10-15" // Fecha actual, puede usarse time.Now().Format("2006-01-02")
	var faltas []FaltaDetalle
	query := `SELECT f.id as falta_id, f.servicio_id, f.observacion, f.estado,
		sfa.fecha_asignacion, sfa.fecha_inicio, sfa.fecha_fin, sfa.sancion_id,
		d.id as documento_id, d.url as documento_url
		FROM faltas f
		LEFT JOIN sancion_falta_asignada sfa ON sfa.falta_id = f.id
		LEFT JOIN faltas_documentos d ON d.falta_id = f.id
		WHERE f.alumno_id = ?
		AND sfa.fecha_inicio <= ? AND sfa.fecha_fin >= ?
		AND f.estado IN ('registrada','sancionada','notificada','apelada','resuelta')`
	err = h.db.Select(&faltas, query, alumnoID, hoy, hoy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error consultando faltas"})
	}

	// Convertir servicioID a int para comparar
	sid, err := strconv.Atoi(servicioID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "servicio_id inválido"})
	}

	// Validar el servicio solicitado
	resultado := fiber.Map{"puede_postular": true, "motivo": "Puede postular a este servicio", "detalle": nil}

	for _, f := range faltas {
		if f.ServicioID == sid {
			// Usar método GetDetalleFalta para obtener detalles normativos
			repo := fault.FactoryStorage(h.db, h.txID)
			detalles, err := repo.GetDetalleFalta(f.FaltaID)
			var detalleMap interface{} = nil
			if err == nil && len(detalles) > 0 {
				detalleMap = detalles[0]
			}

			// Validar sanción vigente
			var sancionAsignada struct {
				FechaInicio      string
				FechaFin         string
				SancionID        string
				Observaciones    string
				CapituloSancion  string
				ArticuloSancion  string
				IncisoSancion    string
				DetalleSancion   string
				DocumentoSancion string
			}
			h.db.Get(&sancionAsignada, `SELECT sfa.fecha_inicio, sfa.fecha_fin, sfa.sancion_id, sfa.observaciones, 
				saf.capitulo_sancion, saf.articulo_sancion, saf.inciso_sancion, saf.detalle_sancion, 
				ds.url as documento_sancion 
				FROM sancion_falta_asignada sfa 
				LEFT JOIN sanciones_faltas_normativa saf ON saf.falta_id = sfa.falta_id 
				LEFT JOIN faltas_documentos ds ON ds.falta_id = sfa.falta_id 
				WHERE sfa.falta_id = ? LIMIT 1`, f.FaltaID)

			// Si la falta está apelada, permite postular aunque la sanción esté vigente
			if f.Estado == "apelada" {
				resultado["puede_postular"] = true
				resultado["motivo"] = "Puede postular, falta apelada"
			} else if sancionAsignada.SancionID != "" && sancionAsignada.FechaInicio != "" &&
				sancionAsignada.FechaFin != "" && sancionAsignada.FechaInicio <= hoy &&
				sancionAsignada.FechaFin >= hoy {
				resultado["puede_postular"] = false
				resultado["motivo"] = fmt.Sprintf("No puede postular a este servicio (ID %d) por sanción vigente", sid)
			}

			resultado["detalle"] = fiber.Map{
				"normativa":             detalleMap,
				"fecha_inicio_sancion":  sancionAsignada.FechaInicio,
				"fecha_fin_sancion":     sancionAsignada.FechaFin,
				"sancion_id":            sancionAsignada.SancionID,
				"observaciones_sancion": sancionAsignada.Observaciones,
				"capitulo_sancion":      sancionAsignada.CapituloSancion,
				"articulo_sancion":      sancionAsignada.ArticuloSancion,
				"inciso_sancion":        sancionAsignada.IncisoSancion,
				"detalle_sancion":       sancionAsignada.DetalleSancion,
				"documento_sancion":     sancionAsignada.DocumentoSancion,
				"falta_cometida":        f.Observacion,
			}
			break
		}
	}

	return c.Status(200).JSON(resultado)
}

// ValidarPostulacionMultiple godoc
// @Summary Valida si el estudiante puede postular a múltiples servicios
// @Description Valida por DNI si el estudiante tiene faltas vigentes en los servicios solicitados
// @Tags Faltas
// @Accept json
// @Produce json
// @Param dni path string true "DNI del alumno"
// @Param body body struct{ServicioIDs []int `json:"servicio_ids"`} true "IDs de servicios"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/alumnos/validar-postulacion-multiple/{dni} [post]
func (h *handlerFaults) ValidarPostulacionMultiple(c *fiber.Ctx) error {
	dni := c.Params("dni")
	if dni == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "DNI requerido"})
	}

	var req struct {
		ServicioIDs []int `json:"servicio_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Body inválido", "details": err.Error()})
	}

	if len(req.ServicioIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "servicio_ids requeridos"})
	}

	// Obtener alumno
	var alumnoID int64
	err := h.db.Get(&alumnoID, "SELECT id FROM alumnos WHERE DNI = ? LIMIT 1", dni)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Alumno no encontrado"})
	}

	// Fecha actual desde la base de datos
	var hoy string
	err = h.db.Get(&hoy, "SELECT CURRENT_DATE")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo obtener la fecha actual del servidor", "details": err.Error()})
	}

	// Validar cada servicio solicitado
	resultados := make(map[string]interface{})

	for _, sid := range req.ServicioIDs {

		// Consultar nombre del servicio
		var nombreServicio string
		errNombre := h.db.Get(&nombreServicio, "SELECT nombre FROM servicios WHERE id = ? LIMIT 1", sid)
		if errNombre != nil {
			nombreServicio = ""
		}

		resultado := fiber.Map{
			"puede_postular":  true,
			"motivo":          "Puede postular a este servicio",
			"detalle":         nil,
			"servicio_nombre": nombreServicio,
		}

		// Primero: Buscar TODAS las faltas de este alumno en este servicio
		type FaltaSimple struct {
			FaltaID    string `db:"falta_id"`
			Estado     string `db:"estado"`
			ServicioID int    `db:"servicio_id"`
		}

		var faltasServicio []FaltaSimple
		queryFaltas := `
			SELECT 
				f.id as falta_id,
				f.estado,
				f.servicio_id
			FROM faltas f
			WHERE f.alumno_id = ?
			AND f.servicio_id = ?
		`

		_ = h.db.Select(&faltasServicio, queryFaltas, alumnoID, sid)
		// No logging

		// Segundo: Buscar falta CON sanción vigente
		type FaltaConSancion struct {
			FaltaID         string `db:"falta_id"`
			ServicioID      int    `db:"servicio_id"`
			Observacion     string `db:"observacion"`
			Estado          string `db:"estado"`
			FechaFalta      string `db:"fecha_falta"`
			FechaInicio     string `db:"fecha_inicio"`
			FechaFin        string `db:"fecha_fin"`
			SancionID       string `db:"sancion_id"`
			Observaciones   string `db:"observaciones"`
			CapituloSancion string `db:"capitulo_sancion"`
			ArticuloSancion string `db:"articulo_sancion"`
			IncisoSancion   string `db:"inciso_sancion"`
			DetalleSancion  string `db:"detalle_sancion"`
		}

		var falta FaltaConSancion
		queryCompleta := `
			SELECT 
    f.id as falta_id,
    f.servicio_id,
    COALESCE(f.observacion, '') as observacion,
    f.estado,
    COALESCE(f.fecha_falta, '') as fecha_falta,
    COALESCE(sfa.fecha_inicio, '') as fecha_inicio,
    COALESCE(sfa.fecha_fin, '') as fecha_fin,
    COALESCE(sfa.sancion_id, '') as sancion_id,
    COALESCE(sfa.observaciones, '') as observaciones,
    COALESCE(saf.capitulo_sancion, '') as capitulo_sancion,
    COALESCE(saf.articulo_sancion, '') as articulo_sancion,
    COALESCE(saf.inciso_sancion, '') as inciso_sancion,
    COALESCE(saf.detalle_sancion, '') as detalle_sancion
FROM faltas f
LEFT JOIN sancion_falta_asignada sfa ON sfa.falta_id = f.id
LEFT JOIN sanciones_faltas_normativa saf ON saf.id = sfa.sancion_id
WHERE f.alumno_id = ?
AND f.servicio_id = ?
AND f.estado IN ('registrada','sancionada','notificada','apelada','resuelta')
ORDER BY f.created_at DESC
LIMIT 1;
		`

		err = h.db.Get(&falta, queryCompleta, alumnoID, sid)

		if err != nil {
			// No logging
		} else {
			// Validar sanción
			if falta.SancionID != "" {
				fechaInicio := falta.FechaInicio
				fechaFin := falta.FechaFin

				// VALIDACIÓN DE SANCIÓN VIGENTE
				sancionVigente := fechaInicio != "" && fechaFin != "" &&
					fechaInicio <= hoy && fechaFin >= hoy

				// REGLA 1: Si está apelada, puede postular
				if falta.Estado == "apelada" {
					resultado["puede_postular"] = true
					resultado["motivo"] = "Puede postular, falta apelada"
				} else if sancionVigente {
					// REGLA 2: Sanción vigente = NO puede postular
					resultado["puede_postular"] = false
					resultado["motivo"] = fmt.Sprintf("No puede postular a %s por sanción vigente hasta %s", nombreServicio, fechaFin)
				} else {
					resultado["puede_postular"] = true
					resultado["motivo"] = "Sanción expirada o no vigente"
				}

				// Obtener detalles normativos
				repo := fault.FactoryStorage(h.db, h.txID)
				detalles, errDetalle := repo.GetDetalleFalta(falta.FaltaID)
				var detalleMap interface{} = nil
				if errDetalle == nil && len(detalles) > 0 {
					detalleMap = detalles[0]
				}

				// Buscar documento
				var documentoURL string
				h.db.Get(&documentoURL, "SELECT COALESCE(url, '') FROM faltas_documentos WHERE falta_id = ? LIMIT 1", falta.FaltaID)

				resultado["detalle"] = fiber.Map{
					"normativa":             detalleMap,
					"fecha_inicio_sancion":  fechaInicio,
					"fecha_fin_sancion":     fechaFin,
					"sancion_id":            falta.SancionID,
					"observaciones_sancion": falta.Observaciones,
					"capitulo_sancion":      falta.CapituloSancion,
					"articulo_sancion":      falta.ArticuloSancion,
					"inciso_sancion":        falta.IncisoSancion,
					"detalle_sancion":       falta.DetalleSancion,
					"documento_sancion":     documentoURL,
					"falta_cometida":        falta.Observacion,
					"fecha_falta":           falta.FechaFalta,
					"estado_falta":          falta.Estado,
					"servicio_nombre":       nombreServicio,
				}
			} else {
			}
		}

		resultados[fmt.Sprintf("servicio_%d", sid)] = resultado
	}

	return c.Status(200).JSON(resultados)
}

// GetAllServicios godoc
// @Summary Obtiene todos los servicios
// @Description Devuelve el listado de servicios (id y nombre)
// @Tags Servicios
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /v1/servicios [get]
func (h *handlerFaults) GetAllServicios(c *fiber.Ctx) error {
	var servicios []struct {
		ID     int    `db:"id" json:"id"`
		Nombre string `db:"nombre" json:"nombre"`
	}
	err := h.db.Select(&servicios, "SELECT id, nombre FROM servicios ORDER BY nombre ASC")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudieron obtener los servicios", "details": err.Error()})
	}
	return c.JSON(servicios)
}
