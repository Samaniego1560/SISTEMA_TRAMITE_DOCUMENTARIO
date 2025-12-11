package oauth

import (
	"dbu-api/internal/middleware"
	"dbu-api/pkg/orchestrator/response_messages"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// RouterOAuth registra las rutas de autenticación OAuth
func RouterOAuth(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerOAuth{
		db:   db,
		txID: txID,
		msg:  response_messages.NewMsg(db),
	}

	// Rutas OAuth bajo /api/v1/auth/oauth
	oauth := app.Group("/api/v1/oauth")

	// Rutas de Microsoft
	oauth.Get("/microsoft", h.MicrosoftLogin)
	oauth.Get("/microsoft/callback", h.MicrosoftCallback)

	// Ruta para validar y activar sesión (requiere autenticación JWT)
	oauth.Post("/validate-session", middleware.StudentJWTProtected(), middleware.ValidateStudentSession(db), h.ValidateSession)

	// Aquí se pueden agregar más proveedores en el futuro:
	// oauth.Get("/google", h.GoogleLogin)
	// oauth.Get("/google/callback", h.GoogleCallback)
}
