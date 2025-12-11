package oauth

import (
	"dbu-api/internal/env"
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	auth_orchestrator "dbu-api/pkg/auth/oauth"
	"dbu-api/pkg/orchestrator/response_messages"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type handlerOAuth struct {
	db   *sqlx.DB
	txID string
	msg  response_messages.Message
}

// MicrosoftLogin godoc
// @Summary Iniciar sesión con Microsoft
// @Description Redirige al usuario a Microsoft para autenticarse
// @Tags OAuth
// @Accept json
// @Produce json
// @Success 307 "Redirección a Microsoft"
// @Router /api/v1/auth/oauth/microsoft [GET]
func (h *handlerOAuth) MicrosoftLogin(c *fiber.Ctx) error {
	// Crear proveedor de Microsoft
	provider := createMicrosoftProvider()

	// Crear servicio OAuth
	srv := auth_orchestrator.NewOAuthLogin(h.db, provider, h.txID)

	// Generar state token (en producción esto debería guardarse en sesión/cache)
	state := "state-token-" + h.txID

	// Obtener URL de autorización
	authURL := srv.GetAuthorizationURL(state)

	logger.Info.Printf("%s - Redirecting to Microsoft OAuth: %s", h.txID, authURL)

	// Redirigir al usuario a Microsoft
	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

// MicrosoftCallback godoc
// @Summary Callback de Microsoft OAuth
// @Description Procesa el callback de Microsoft después de la autenticación
// @Tags OAuth
// @Accept json
// @Produce json
// @Param code query string true "Código de autorización"
// @Param state query string false "State token"
// @Success 307 "Redirección al frontend con token"
// @Failure 400 {object} models.Response{error=boolean,data=interface{},code=integer,type=string,msg=string} "Error en la solicitud"
// @Router /api/v1/auth/oauth/microsoft/callback [GET]
func (h *handlerOAuth) MicrosoftCallback(c *fiber.Ctx) error {
	res := models.Response{Error: true}
	cfg := env.NewConfiguration()

	// Obtener código de autorización
	code := c.Query("code")
	if code == "" {
		logger.Error.Printf("%s - OAuth callback missing code parameter", h.txID)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(1)
		res.Msg = "Missing authorization code"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Validar state (en producción validar contra sesión/cache)
	state := c.Query("state")
	logger.Info.Printf("%s - OAuth callback received: state=%s", h.txID, state)

	// Crear proveedor de Microsoft
	provider := createMicrosoftProvider()

	// Crear servicio OAuth
	srv := auth_orchestrator.NewOAuthLogin(h.db, provider, h.txID)

	// Obtener información del cliente
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	// Procesar callback y obtener tokens
	token, code_response, err := srv.ProcessCallback(c.Context(), code, c.BaseURL(), ipAddress, userAgent)
	if err != nil {
		logger.Error.Printf("%s - OAuth callback failed: %v", h.txID, err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(code_response)
		if res.Msg == "" {
			res.Msg = "Authentication failed"
		}

		// Redirigir al frontend con error
		errorURL := fmt.Sprintf("%s//auth/callback?error_code=%d&error_message=%s", cfg.OAuth.FrontendURL, res.Code, res.Msg)
		return c.Redirect(errorURL, fiber.StatusTemporaryRedirect)
	}

	// Redirigir al frontend con el token
	redirectURL := fmt.Sprintf("%s//auth/callback?token=%s",
		cfg.OAuth.FrontendURL, token)

	logger.Info.Printf("%s - OAuth login successful, redirecting to frontend", h.txID)

	return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
}

// ValidateSession godoc
// @Summary Validar y activar sesión
// @Description Valida el token JWT y activa la sesión si no ha expirado
// @Tags OAuth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{error=boolean,data=interface{},code=integer,type=string,msg=string} "Sesión activada exitosamente"
// @Failure 400 {object} models.Response{error=boolean,data=interface{},code=integer,type=string,msg=string} "Error en la solicitud"
// @Failure 401 {object} models.Response{error=boolean,data=interface{},code=integer,type=string,msg=string} "Token expirado o inválido"
// @Router /api/v1/oauth/validate-session [POST]
func (h *handlerOAuth) ValidateSession(c *fiber.Ctx) error {
	res := models.Response{Error: true}

	// Obtener session_id del middleware JWT
	sessionID, ok := c.Locals("session_id").(string)
	if !ok || sessionID == "" {
		logger.Error.Printf("%s - Session ID not found in context", h.txID)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(101)
		if res.Msg == "" {
			res.Msg = "Session ID not found in context"
		}
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Crear servicio de alumnos
	srv := auth_orchestrator.NewOAuthLogin(h.db, nil, h.txID)

	// Activar la sesión
	if err := srv.UpdateSessionStatus(sessionID, true); err != nil {
		logger.Error.Printf("%s - Error updating session status: %v", h.txID, err)
		res.Code, res.Type, res.Msg = h.msg.GetByCode(102)
		if res.Msg == "" {
			res.Msg = "Failed to activate session"
		}
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	logger.Info.Printf("%s - Session activated successfully: %s", h.txID, sessionID)

	res.Error = false
	res.Code, res.Type, res.Msg = h.msg.GetByCode(224)
	if res.Msg == "" {
		res.Msg = "Session activated successfully"
	}
	res.Data = map[string]interface{}{
		"session_id": sessionID,
		"active":     true,
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
