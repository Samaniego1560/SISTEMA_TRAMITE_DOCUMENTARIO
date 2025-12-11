package report

import (
	repository_sqlserver "dbu-api/pkg/normative/convocatorias"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type ConvocatoriaHandler struct {
	repo *repository_sqlserver.ConvocatoriaRepository
}

func NewConvocatoriaHandler(db *sqlx.DB) *ConvocatoriaHandler {
	repo := repository_sqlserver.NewConvocatoriaRepository(db)
	return &ConvocatoriaHandler{repo: repo}
}

func (h *ConvocatoriaHandler) GetAllConvocatorias(c *fiber.Ctx) error {
	convocatorias, err := h.repo.GetAllConvocatorias()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudieron obtener las convocatorias"})
	}
	return c.JSON(fiber.Map{"convocatorias": convocatorias})
}
