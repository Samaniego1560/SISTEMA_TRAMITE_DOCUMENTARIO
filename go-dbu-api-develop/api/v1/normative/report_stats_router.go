package normative

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterReportStats(app *fiber.App, db *sqlx.DB) {
	api := app.Group("/api/v1")
	normative := api.Group("/normative")

	// Crear una instancia del handler
	handler := NewHandlerReportStats(db)

	// Rutas
	normative.Get("/report-stats", handler.GetReportStats)
	normative.Get("/report-stats/comparacion", handler.GetComparacionSemestres)
	normative.Get("/report-stats/:convocatoria_id", handler.GetReportStatsBySemestre)
	normative.Get("/report-stats/:convocatoria_id/comparar", handler.GetSemestreConAnterior)
}
