package resolutions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterResoluciones(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerResolutions{db: db, txID: txID}
	v1 := app.Group("/v1")
	Resoluciones := v1.Group("/resoluciones")
	//Resoluciones.Use(middleware.JWTProtected())
	Resoluciones.Post("/", h.CreateResoluciones)
	Resoluciones.Put("/:id", h.UpdateResoluciones)
	Resoluciones.Get("/", h.GetAllResoluciones)
	Resoluciones.Delete("/:id", h.DeleteResoluciones)
}
