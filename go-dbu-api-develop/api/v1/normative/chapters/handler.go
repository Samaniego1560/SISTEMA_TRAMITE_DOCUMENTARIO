package chapters

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

type handlerChapters struct {
	db   *sqlx.DB
	txID string
}

// UpdateChapters godoc
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
// @Router /api/v1/cuartos [PUT]

func (h *handlerChapters) CreateChapters(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	req := models.Chapter{}

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

	isValid, err := req.ValidChapter()
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
	_, code, err := srv.SrvChapter.CreateChapter(req.ID, req.Nombre, req.Descripcion, req.NormativeId)
	if err != nil {
		logger.Error.Printf("couldn't update Chapters, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = req
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerChapters) UpdateChapters(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	req := ChapterRequest{}

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
	_, code, err := srv.SrvChapter.UpdateOnlyCharacteristicsChapter(req.ID, req.name, req.Description)
	if err != nil {
		logger.Error.Printf("couldn't update Chapters, error: %v", err)
		res.Code, res.Type, res.Msg = code, "", ""
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = req
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", ""
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerChapters) GetChaptersByResolutionId(c *fiber.Ctx) error {
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

	resolucionId := c.Params("id")
	if resolucionId == "" {
		logger.Error.Printf("Resolución Id not found, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "Resolución Id not found"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	srv := normative.NewServerResolution(h.db, user, h.txID)
	normative, code, err := srv.SrvChapter.GetChaptersByResolution(resolucionId)
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
