package chapters

import (
	"dbu-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterChapters(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerChapters{db: db, txID: txID}
	v1 := app.Group("/v1")
	chapters := v1.Group("/capitulos")
	chapters.Use(middleware.JWTProtected())
	chapters.Post("/", h.CreateChapters)
	chapters.Put("/", h.UpdateChapters)
	chapters.Get("/resolucion-id/:id", h.GetChaptersByResolutionId)
}
