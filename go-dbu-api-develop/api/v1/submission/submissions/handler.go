package submissions

import (
	"dbu-api/internal/authorization"
	"dbu-api/internal/logger"
	"dbu-api/internal/middleware"
	"dbu-api/internal/models"
	"dbu-api/pkg/orchestrator/low_code_submissions"
	"dbu-api/pkg/orchestrator/response_messages"
	"dbu-api/pkg/submission/convocatorias"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type handlerResidencias struct {
	db   *sqlx.DB
	txID string
	msg  response_messages.Message
}

// StudentsBySubmissions godoc
// @Summary Obtener alumnos por convocatoria
// @Description Método que permite obtener la lista de alumnos asociados a una convocatoria
// @Tags Convocatorias
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "ID de la convocatoria"
// @Param page query integer false "Número de página" minimum(0) default(0)
// @Param limit query integer false "Límite de registros por página" minimum(0) maximum(100) default(0)
// @Param gender query string false "Género del estudiante (masculino/femenino)" Enums(masculino,femenino)
// @Success 200 {object} models.Response{error=boolean,data=ResponseStudentsSubmission,code=integer,type=string,msg=string} "Lista de alumnos obtenida exitosamente"
// @Failure 400 {object} models.Response "Error en la solicitud"
// @Failure 401 {object} models.Response "Error de autenticación"
// @Failure 500 {object} models.Response "Error interno del servidor"
// @Router /v1/convocatorias/{id}/alumnos-aceptados [GET]
func (h *handlerResidencias) StudentsBySubmissions(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	// Autenticación
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Autorización
	if err := authorization.ValidPermissions(user, h.db, c); err != nil {
		logger.Error.Printf("User does not have permission to call the api, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Parámetros de consulta
	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 0)
	gender := c.Query("gender")
	departmentRequired := c.QueryBool("department_required", false)

	if page < 1 {
		limit = 0
	}

	if page >= 1 && (limit < 1 || limit > 100) {
		logger.Error.Printf("invalid limit number: %d", limit)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(2)
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	if gender != "" && (gender != "masculino" && gender != "femenino") {
		logger.Error.Printf("invalid gender: %s", gender)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(2)
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// ID de convocatoria
	submission, err := c.ParamsInt("id")
	if err != nil {
		logger.Error.Printf("missing submission_id parameter")
		res.Code, res.Type, res.Msg = h.msg.GetByCode(2)
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Servicio
	srv := low_code_submissions.NewSubmission(h.db, user, h.txID)
	students, code, total, err := srv.GetStudentsBySubmissionsLowCode(submission, page, limit, gender, departmentRequired)
	if err != nil {
		logger.Error.Printf("error getting students: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(code)
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	// Respuesta
	res.Data = ResponseStudentsSubmission{
		Total:    total,
		Students: students,
	}

	res.Error = false
	res.Code, res.Type, res.Msg = h.msg.GetByCode(210)
	return c.Status(http.StatusOK).JSON(res)
}

// StudentReportsBySubmissions godoc
// @Summary Generar reporte Excel de alumnos por convocatoria
// @Description Método que permite generar un reporte Excel en formato base64 de los alumnos asociados a una convocatoria
// @Tags Convocatorias
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "ID de la convocatoria"
// @Success 200 {object} models.Response{error=boolean,data=string,code=integer,type=string,msg=string} "Reporte Excel generado exitosamente en formato base64"
// @Failure 400 {object} models.Response "Error en la solicitud"
// @Failure 401 {object} models.Response "Error de autenticación"
// @Failure 500 {object} models.Response "Error interno del servidor"
// @Router /v1/convocatorias/{id}/reporte-residencias [GET]
func (h *handlerResidencias) StudentReportsBySubmissions(c *fiber.Ctx) error {
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

	submission, err := c.ParamsInt("id")
	if err != nil {
		logger.Error.Printf("missing submission_id parameter")
		res.Code, res.Type, res.Msg = h.msg.GetByCode(2)
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	srv := low_code_submissions.NewSubmission(h.db, user, h.txID)
	base64, code, err := srv.GetReportBySubmissionsLowCode(submission)
	if err != nil {
		logger.Error.Printf("error generate report: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(code)
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res.Data = base64
	res.Error = false
	res.Code, res.Type, res.Msg = h.msg.GetByCode(223)
	return c.Status(http.StatusOK).JSON(res)
}

// CreateConvocatoria godoc
// @Summary Crear una convocatoria con sus relaciones
// @Description Método que permite crear una convocatoria completa con servicios, secciones y requisitos
// @Tags Convocatorias
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param convocatoria body convocatorias.CreateConvocatoriaRequest true "Datos de la convocatoria"
// @Success 201 {object} models.Response{error=boolean,data=convocatorias.ConvocatoriaResponse,code=integer,type=string,msg=string} "Convocatoria creada exitosamente"
// @Failure 400 {object} models.Response "Error en la solicitud"
// @Failure 401 {object} models.Response "Error de autenticación"
// @Failure 422 {object} models.Response "Error de validación"
// @Router /v1/convocatoria/create [POST]
func (h *handlerResidencias) CreateConvocatoria(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	// Autenticación
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(53)
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Parsear request
	var req convocatorias.CreateConvocatoriaRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(1)
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	// Validar request
	isValid, err := req.Valid()
	if err != nil || !isValid {
		logger.Error.Printf("validation error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(2)
		res.Data = err.Error()
		return c.Status(fiber.StatusUnprocessableEntity).JSON(res)
	}

	// Crear servicio
	srv := convocatorias.NewConvocatoriasService(
		convocatorias.FactoryStorage(h.db, user, h.txID),
		user,
		h.txID,
	)

	// Crear convocatoria con relaciones
	convocatoria, code, err := srv.CreateConvocatoriaWithRelations(&req)
	if err != nil {
		logger.Error.Printf("couldn't create convocatoria: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(code)
		res.Data = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	// Respuesta exitosa
	res.Data = convocatoria
	res.Error = false
	res.Code, res.Type, res.Msg = h.msg.GetByCode(211)
	return c.Status(http.StatusCreated).JSON(res)
}

// UpdateConvocatoria godoc
// @Summary Actualizar una convocatoria con sus relaciones
// @Description Método que permite actualizar una convocatoria completa con sus servicios, secciones y requisitos
// @Tags Convocatorias
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "ID de la convocatoria"
// @Param convocatoria body convocatorias.CreateConvocatoriaRequest true "Datos actualizados de la convocatoria"
// @Success 200 {object} models.Response{error=boolean,data=convocatorias.ConvocatoriaResponse,code=integer,type=string,msg=string} "Convocatoria actualizada exitosamente"
// @Failure 400 {object} models.Response "Error en la solicitud"
// @Failure 401 {object} models.Response "Error de autenticación"
// @Failure 404 {object} models.Response "Convocatoria no encontrada"
// @Failure 422 {object} models.Response "Error de validación"
// @Router /v1/convocatoria/update/{id} [PUT]
func (h *handlerResidencias) UpdateConvocatoria(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	// Autenticación
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(53)
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Obtener ID de la convocatoria
	convocatoriaID, err := c.ParamsInt("id")
	if err != nil {
		logger.Error.Printf("invalid convocatoria ID, error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(2)
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	// Parsear request
	var req convocatorias.CreateConvocatoriaRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(1)
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	// Validar request
	isValid, err := req.Valid()
	if err != nil || !isValid {
		logger.Error.Printf("validation error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(2)
		res.Data = err.Error()
		return c.Status(fiber.StatusUnprocessableEntity).JSON(res)
	}

	// Crear servicio
	srv := convocatorias.NewConvocatoriasService(
		convocatorias.FactoryStorage(h.db, user, h.txID),
		user,
		h.txID,
	)

	// Actualizar convocatoria con relaciones
	convocatoria, code, err := srv.UpdateConvocatoriaWithRelations(int64(convocatoriaID), &req)
	if err != nil {
		logger.Error.Printf("couldn't update convocatoria: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(code)
		res.Data = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	// Respuesta exitosa
	res.Data = convocatoria
	res.Error = false
	res.Code, res.Type, res.Msg = h.msg.GetByCode(212)
	return c.Status(http.StatusOK).JSON(res)
}

// GetConvocatoriaWithRelations godoc
// @Summary Obtener una convocatoria con todas sus relaciones
// @Description Método que permite obtener una convocatoria completa con servicios, secciones y requisitos
// @Tags Convocatorias
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "ID de la convocatoria"
// @Success 200 {object} models.Response{error=boolean,data=convocatorias.ConvocatoriaResponse,code=integer,type=string,msg=string} "Convocatoria obtenida exitosamente"
// @Failure 400 {object} models.Response "Error en la solicitud"
// @Failure 401 {object} models.Response "Error de autenticación"
// @Failure 404 {object} models.Response "Convocatoria no encontrada"
// @Router /v1/convocatoria/show/{id} [GET]
func (h *handlerResidencias) GetConvocatoriaWithRelations(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	// Autenticación
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(53)
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Autorización
	if err := authorization.ValidPermissions(user, h.db, c); err != nil {
		logger.Error.Printf("User does not have permission, error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(10)
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Obtener ID de la convocatoria
	convocatoriaID, err := c.ParamsInt("id")
	if err != nil {
		logger.Error.Printf("invalid convocatoria ID, error: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(2)
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	// Crear servicio
	srv := convocatorias.NewConvocatoriasService(
		convocatorias.FactoryStorage(h.db, user, h.txID),
		user,
		h.txID,
	)

	// Obtener convocatoria con relaciones
	convocatoria, code, err := srv.GetConvocatoriaWithRelations(int64(convocatoriaID))
	if err != nil {
		logger.Error.Printf("couldn't get convocatoria: %v", err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	// Respuesta exitosa
	res.Data = convocatoria
	res.Error = false
	res.Code, res.Type, res.Msg = h.msg.GetByCode(29)
	return c.Status(http.StatusOK).JSON(res)
}
