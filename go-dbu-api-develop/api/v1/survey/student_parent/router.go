package student_parent

import (
	"dbu-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterStudentParentSurvey(app *fiber.App, db *sqlx.DB, txID string) {
	handler := handlerSurvey{db: db, txID: txID}

	group := app.Group("/api/survey/student-parent")
	group.Post("/verify", handler.Verify)
	group.Post("/", handler.Store)
	group.Get("/export", middleware.JWTProtected(), handler.Export)
	group.Get("/stats", middleware.JWTProtected(), handler.Stats)
}
