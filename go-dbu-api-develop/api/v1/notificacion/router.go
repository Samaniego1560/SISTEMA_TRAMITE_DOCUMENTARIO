package notificacion

import (
	"dbu-api/internal/middleware"
	"dbu-api/pkg/notificacion/smtp"
	"dbu-api/pkg/sanction/fault"

	"github.com/gofiber/fiber/v2"
)

func RouterNotification(
	app *fiber.App,
	smtpService *smtp.SMTPService,
	txID string,
	faultService fault.PortsServerFault,
	sancionService smtp.SancionService,
) {
	h := handlerNotification{
		smtpService:    smtpService,
		txID:           txID,
		faultService:   faultService,
		sancionService: sancionService,
	}

	v1 := app.Group("/v1")
	notification := v1.Group("/notificacion")
	notification.Use(middleware.JWTProtected())

	notification.Post("/send", h.SendNotificationEmail)
	notification.Post("/send-with-cc", h.SendNotificationEmailWithCC)
	notification.Get("/test", h.TestSMTPConnection)
	notification.Get("/test-db", h.TestDatabaseConnection) // ✅ NUEVO
	notification.Post("/send/:falta_id", h.SendNotificationByFaltaID)
	notification.Get("/info/:falta_id", h.GetNotificationInfoByFaltaID)
	notification.Get("/sancion/:falta_id", h.GetSancionPorFaltaID)
}
