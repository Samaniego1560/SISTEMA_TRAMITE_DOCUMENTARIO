package auth

import (
	"dbu-api/internal/middleware"
	"dbu-api/pkg/notificacion/smtp"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// RegisterRoutes registra las rutas de autenticación de estudiantes
func RegisterRoutes(app *fiber.App, db *sqlx.DB, smtpService *smtp.SMTPService, txID string) {
	h := handlerStudentAuth{
		db:          db,
		txID:        txID,
		smtpService: smtpService,
	}

	v1 := app.Group("/api/v1")
	studentAuth := v1.Group("/student/auth")

	// Rutas públicas (sin autenticación)
	studentAuth.Post("/request-otp", h.RequestOTP)
	studentAuth.Post("/verify-otp", h.VerifyOTP)

	// Rutas protegidas (requieren autenticación de estudiante)
	studentAuth.Post("/logout", middleware.StudentJWTProtected(), middleware.ValidateStudentSession(db), h.Logout)
}
