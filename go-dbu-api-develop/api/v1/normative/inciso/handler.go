package incisos

import (
	"dbu-api/internal/authorization"
	"dbu-api/internal/logger"
	"dbu-api/internal/middleware"
	"dbu-api/internal/models"
	"dbu-api/pkg/normative"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type handlerIncisos struct {
	db   *sqlx.DB
	txID string
}

// UpdateIncisos godoc
// @Summary Actualiza una instancia de cuartos
// @Description Método que permite Actualiza una instancia del objeto cuartos en la base de datos
// @tags Cuartos
// @Accept json
// @Produce json
// @Param RoomRequest body RoomRequest true "Datos para actualizar cuartos"
// @securityDefinitions.basic  BasicAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 202 {object} models.Response

func (h *handlerIncisos) CreateIncisos(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	req := models.Inciso{}

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

	isValid, err := req.ValidInciso()
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
	article, code, err := srv.SrvInciso.CreateInciso(req.ID, req.Nombre, req.Descripcion, req.ArticuloId)
	if err != nil {
		logger.Error.Printf("couldn't update Incisos, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = article
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

// @Router /api/v1/cuartos [PUT]
func (h *handlerIncisos) UpdateIncisos(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	req := models.Article{}

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

	isValid, err := req.ValidArticle()
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
	_, code, err := srv.SrvArticle.UpdateArticle(req.ID, req.Descripcion, req.Gravedad, req.CapituloId)
	if err != nil {
		logger.Error.Printf("couldn't update Incisos, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = req
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerIncisos) GetIncisos(c *fiber.Ctx) error {
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

	srv := normative.NewServerResolution(h.db, user, h.txID)
	Incisos, err := srv.SrvArticle.GetAllArticle()
	if err != nil {
		logger.Error.Printf("couldn't GetAllArticle, error: %v", err)
		res.Code, res.Type, res.Msg = 99, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = Incisos
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerIncisos) GetIncisosByArticleId(c *fiber.Ctx) error {
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

	id := c.Params("id")
	srv := normative.NewServerResolution(h.db, user, h.txID)
	Incisos, _, err := srv.SrvInciso.GetIncisosByArticleId(id)
	if err != nil {
		logger.Error.Printf("couldn't GetAllArticle, error: %v", err)
		res.Code, res.Type, res.Msg = 99, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = Incisos
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}
func (h *handlerIncisos) DeleteInciso(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	// Obtener el ID del inciso desde la URL
	incisoId := c.Params("id")
	if incisoId == "" {
		res.Code, res.Type, res.Msg = 1, "error", "Inciso ID is required"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	// Verificar el token de autorización
	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Verificar los permisos del usuario
	err = authorization.ValidPermissions(user, h.db, c)
	if err != nil {
		logger.Error.Printf("User does not have permission to call the api, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Llamar a la función que elimina el inciso
	srv := normative.NewServerResolution(h.db, user, h.txID)
	code, err := srv.SrvInciso.DeleteInciso(incisoId)
	if err != nil {
		logger.Error.Printf("No se pudo eliminar el inciso, error: %v", err)
		res.Code, res.Type, res.Msg = code, "error", "no se pudo eliminar el inciso"
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	res.Error = false
	res.Code, res.Type, res.Msg = 0, "success", "Inciso deleted successfully"
	return c.Status(http.StatusOK).JSON(res)
}
