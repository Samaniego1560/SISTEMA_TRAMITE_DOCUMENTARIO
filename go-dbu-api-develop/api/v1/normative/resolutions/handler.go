package resolutions

import (
	"dbu-api/internal/authorization"
	"dbu-api/internal/logger"
	"dbu-api/internal/middleware"
	"dbu-api/internal/models"
	"dbu-api/pkg/normative"
	"dbu-api/pkg/orchestrator/low_code_normative"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type handlerResolutions struct {
	db   *sqlx.DB
	txID string
}

// CreateResoluciones godoc
// @Summary Crear una instancia de Resoluciones
// @Description Método que permite crear una instancia del objeto Resoluciones en la base de datos
// @tags Resoluciones
// @Accept json
// @Produce json
// @Param RequestResoluciones body RequestResoluciones true "Datos para crear Resoluciones"
// @Success 201 {object} ResponseResoluciones
// @Failure 400 {object} ResponseResoluciones
// @Failure 202 {object} ResponseResoluciones
// @Router /api/v1/Resoluciones [POST]
func (h *handlerResolutions) CreateResoluciones(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	req := models.Normative{}
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
		res.Code, res.Type, res.Msg = 1, "", ""
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	isValid, err := req.Valid()
	if err != nil {
		logger.Error.Printf("couldn't validate body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "", ""
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	if !isValid {
		logger.Error.Println("couldn't validate body request")
		res.Code, res.Type, res.Msg = 1, "", ""
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	srv := low_code_normative.NewResolution(h.db, user, h.txID)
	code, err := srv.CreateResolutionLowCode(&req)
	if err != nil {
		logger.Error.Printf("couldn't create Resolution, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = req
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusCreated).JSON(res)
}

// UpdateResoluciones godoc
// @Summary Actualiza una instancia de Resoluciones
// @Description Método que permite Actualiza una instancia del objeto Resoluciones en la base de datos
// @tags Resoluciones
// @Accept json
// @Produce json
// @Param RequestResoluciones body RequestResoluciones true "Datos para actualizar Resoluciones"
// @Success 200 {object} ResponseResoluciones
// @Failure 400 {object} ResponseResoluciones
// @Failure 202 {object} ResponseResoluciones
// @Router /api/v1/Resoluciones [PUT]
func (h *handlerResolutions) UpdateResoluciones(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	req := models.Normative{}

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
		res.Code, res.Type, res.Msg = 1, "", ""
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	isValid, err := req.Valid()
	if err != nil {
		logger.Error.Printf("couldn't validate body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "", ""
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	if !isValid {
		logger.Error.Println("couldn't validate body request")
		res.Code, res.Type, res.Msg = 1, "", ""
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	srv := normative.NewServerResolution(h.db, user, h.txID)
	_, code, err := srv.SrvResolution.UpdateResolution(req.ID, req.Nombre, req.Descripcion, req.Estado, req.ServicioId, req.RutaArchivo)
	if err != nil {
		logger.Error.Printf("couldn't update Resolution, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = req
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

// DeleteResoluciones godoc
// @Summary Elimina una instancia de Resoluciones
// @Description Método que permite eliminar una instancia del objeto Resoluciones en la base de datos
// @tags Resoluciones
// @Accept json
// @Produce json
// @Param	id	path int true "Resoluciones ID"
// @Success 200 {object} ResponseResoluciones
// @Failure 400 {object} ResponseResoluciones
// @Failure 202 {object} ResponseResoluciones
// @Router /api/v1/Resoluciones [DELETE]
func (h *handlerResolutions) DeleteResoluciones(c *fiber.Ctx) error {
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
	srv := normative.NewServerResolution(h.db, nil, h.txID)
	code, err := srv.SrvResolution.DeleteResolution(idStr)
	if err != nil {
		logger.Error.Printf("couldn't delete Resolution, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

// GetResolucionesByID godoc
// @Summary Obtiene una instancia de Resoluciones por su id
// @Description Método que permite obtener una instancia del objeto Resoluciones en la base de datos por su id
// @tags Resoluciones
// @Accept json
// @Produce json
// @Param	id	path int true "Resoluciones ID"
// @Success 200 {object} ResponseResoluciones
// @Failure 400 {object} ResponseResoluciones
// @Failure 202 {object} ResponseResoluciones
// @Router /api/v1/Resoluciones/:id [GET]
func (h *handlerResolutions) GetResolucionesByID(c *fiber.Ctx) error {
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
	srv := normative.NewServerResolution(h.db, nil, h.txID)
	data, code, err := srv.SrvResolution.GetResolutionByID(idStr)
	if err != nil {
		logger.Error.Printf("couldn't get Resolution by id, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = data
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

// GetResolucionesByID godoc
// @Summary Obtiene todas las instancias de Resoluciones
// @Description Método que permite obtener todas las instancias del objeto Resoluciones en la base de datos
// @tags Resoluciones
// @Accept json
// @Produce json
// @Success 200 {object} ResponseResoluciones
// @Failure 202 {object} ResponseResoluciones
// @Router /api/v1/Resoluciones/all [GET]
func (h *handlerResolutions) GetAllResoluciones(c *fiber.Ctx) error {
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
	srv := low_code_normative.NewResolution(h.db, user, h.txID)
	normative, code, err := srv.GetAllResolutionsLowCode()
	if err != nil {
		logger.Error.Printf("couldn't create Resolution, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = normative
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}
