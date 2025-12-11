package sanctioned

import (
	"dbu-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterSanctioned(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerSanctioned{db: db, txID: txID}
	v1 := app.Group("/v1")
	sanctioned := v1.Group("/sancionados")
	sanctioned.Use(middleware.JWTProtected())
	sanctioned.Put("/", h.UpdateSanctioned)
}
