package incisos

import (
	"dbu-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterIncisos(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerIncisos{db: db, txID: txID}
	v1 := app.Group("/v1")
	incisos := v1.Group("/incisos")
	incisos.Use(middleware.JWTProtected())
	incisos.Post("/", h.CreateIncisos)
	incisos.Put("/", h.UpdateIncisos)
	incisos.Get("", h.GetIncisos)
	incisos.Get("/article/:id", h.GetIncisosByArticleId)
	incisos.Delete("/:id", h.DeleteInciso)
}
