package assignment

import (
	"dbu-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// RegisterRoutes registra las rutas de asignaciones de estudiantes
func RegisterRoutes(app *fiber.App, db *sqlx.DB) {

	h := NewStudentAssignment(db)

	studentAssignment := app.Group("/api/v1/student")

	// Todas las rutas requieren autenticación de estudiante
	studentAssignment.Use(middleware.StudentJWTProtected())
	studentAssignment.Use(middleware.ValidateStudentSession(db))

	// Rutas de consulta de asignaciones
	studentAssignment.Get("/submissions", h.GetSubmissions)
	studentAssignment.Get("/assignment/:convocatoria_id", h.GetAssignmentDetail)

	// Ruta de perfil del estudiante
	studentAssignment.Get("/profile", h.GetProfile)
}
