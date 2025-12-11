package articles

import (
	"dbu-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterArticles(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerArticles{db: db, txID: txID}
	v1 := app.Group("/v1")
	articles := v1.Group("/articulos")
	articles.Use(middleware.JWTProtected())
	articles.Post("/", h.CreateArticles)
	articles.Put("/", h.UpdateArticles)
	articles.Get("", h.GetArticles)
}
