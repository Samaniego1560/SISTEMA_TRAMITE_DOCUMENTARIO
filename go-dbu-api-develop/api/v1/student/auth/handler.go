package auth

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"dbu-api/pkg/auth/otp"
	"dbu-api/pkg/notificacion/smtp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type handlerStudentAuth struct {
	db          *sqlx.DB
	txID        string
	smtpService *smtp.SMTPService
}

// RequestOTP solicita un código OTP para autenticación
// @Summary Solicitar código OTP
// @Description Genera y envía un código OTP al correo institucional del estudiante
// @Tags Student Auth
// @Accept json
// @Produce json
// @Param body body RequestOTPRequest true "DNI del estudiante"
// @Success 200 {object} models.Response{data=RequestOTPResponse}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/student/auth/request-otp [post]
func (h *handlerStudentAuth) RequestOTP(c *fiber.Ctx) error {
	txID := uuid.New().String()

	// Parsear request
	var req RequestOTPRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("%s - error parsing request: %v", txID, err)
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Error: true,
			Msg:   "Datos inválidos",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Validar DNI
	if req.DNI == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Error: true,
			Msg:   "El DNI es requerido",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Obtener IP del cliente
	ipAddress := c.IP()

	// Crear servicio OTP
	otpRepo := otp.FactoryStorage(h.db, txID)
	otpService := otp.NewOTPService(otpRepo, h.smtpService, txID)

	// Solicitar OTP (el servicio maneja todo: generación, guardado y envío de email)
	otpToken, err := otpService.RequestOTP(req.DNI, ipAddress)
	if err != nil {
		logger.Error.Printf("%s - error requesting OTP: %v", txID, err)
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Error: true,
			Msg:   err.Error(),
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	response := RequestOTPResponse{
		CorreoParcial: otpToken.CorreoDestino,
		ExpiraEnSeg:   otp.OTPExpirationMinutes,
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Error: false,
		Msg:   "Código OTP enviado a tu correo institucional",
		Data:  response,
		Code:  200,
		Type:  "success",
	})
}

// VerifyOTP verifica el código OTP y genera un token JWT
// @Summary Verificar código OTP
// @Description Valida el código OTP y retorna un token JWT para el estudiante
// @Tags Student Auth
// @Accept json
// @Produce json
// @Param body body VerifyOTPRequest true "DNI y código OTP"
// @Success 200 {object} models.Response{data=VerifyOTPResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /api/v1/student/auth/verify-otp [post]
func (h *handlerStudentAuth) VerifyOTP(c *fiber.Ctx) error {
	txID := uuid.New().String()

	// Parsear request
	var req VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("%s - error parsing request: %v", txID, err)
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Error: true,
			Msg:   "Datos inválidos",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Validar campos
	if req.DNI == "" || req.CodigoOTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Error: true,
			Msg:   "DNI y código OTP son requeridos",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	if len(req.CodigoOTP) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Error: true,
			Msg:   "El código OTP debe tener 6 dígitos",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Obtener IP y User-Agent
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	// Crear servicio OTP
	otpRepo := otp.FactoryStorage(h.db, txID)
	otpService := otp.NewOTPService(otpRepo, nil, txID)

	// Verificar OTP y generar JWT
	token, _, err := otpService.VerifyOTP(req.DNI, req.CodigoOTP, ipAddress, userAgent)
	if err != nil {
		logger.Error.Printf("%s - error verifying OTP: %v", txID, err)
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Error: true,
			Msg:   err.Error(),
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	response := VerifyOTPResponse{
		Token:     token,
		ExpiresIn: otp.SessionExpirationHours * 3600, // en segundos
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Error: false,
		Msg:   "success",
		Data:  response,
		Code:  200,
		Type:  "success",
	})
}

// Logout cierra la sesión del estudiante
// @Summary Cerrar sesión
// @Description Desactiva la sesión actual del estudiante
// @Tags Student Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /api/v1/student/auth/logout [post]
func (h *handlerStudentAuth) Logout(c *fiber.Ctx) error {
	txID := uuid.New().String()

	// Obtener session_id del contexto (viene del middleware)
	sessionID, ok := c.Locals("session_id").(string)
	if !ok || sessionID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Error: true,
			Msg:   "Sesión inválida",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	// Crear servicio OTP
	otpRepo := otp.FactoryStorage(h.db, txID)
	otpService := otp.NewOTPService(otpRepo, nil, txID)

	// Cerrar sesión
	err := otpService.Logout(sessionID)
	if err != nil {
		logger.Error.Printf("%s - error logging out: %v", txID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Error: true,
			Msg:   "Error al cerrar sesión",
			Data:  nil,
			Code:  0,
			Type:  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Error: false,
		Msg:   "Sesión cerrada exitosamente",
		Data:  nil,
		Code:  200,
		Type:  "success",
	})
}
