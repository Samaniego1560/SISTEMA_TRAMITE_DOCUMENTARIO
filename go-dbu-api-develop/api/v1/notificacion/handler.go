package notificacion

import (
	"dbu-api/internal/dbx"
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"dbu-api/pkg/notificacion/smtp"
	"dbu-api/pkg/sanction/fault"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type handlerNotification struct {
	smtpService    *smtp.SMTPService
	txID           string
	faultService   fault.PortsServerFault
	sancionService smtp.SancionService
}

func (h *handlerNotification) validateNotificationData(data smtp.NotificationData) error {
	if data.Alumno.CorreoInstitucional == "" {
		return fmt.Errorf("el correo institucional del alumno es requerido")
	}
	if data.Alumno.Nombres == "" || data.Alumno.ApellidoPaterno == "" {
		return fmt.Errorf("el nombre completo del alumno es requerido")
	}
	if data.Servicio == "" {
		return fmt.Errorf("el servicio es requerido")
	}
	if data.FechaFalta == "" {
		return fmt.Errorf("la fecha de la falta es requerida")
	}
	return nil
}

func (h *handlerNotification) SendNotificationByFaltaID(c *fiber.Ctx) error {
	var faltaID string

	faltaID = c.Params("falta_id")

	if faltaID == "" {
		var req struct {
			FaltaID string `json:"falta_id"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Error al parsear datos JSON",
				"error":   err.Error(),
			})
		}

		faltaID = req.FaltaID
	}

	if faltaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "El falta_id es requerido (puede venir por URL params o en el body JSON)",
		})
	}

	agrupado, err := h.faultService.GetDetalleFaltaAgrupado(faltaID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	convocatorias, err := ObtenerConvocatorias()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error al obtener convocatorias",
			"error":   err.Error(),
		})
	}

	periodo := ""
	if p, err := dbx.GetConvocatoriaPeriodoPorFaltaID(dbx.GetConnection().DB, faltaID); err == nil {
		periodo = p
	}

	sancion, _ := smtp.GetSancionAsignadaPorFaltaInternal(agrupado.FaltaID, h.sancionService)

	notificationData := MapeaFaltaAgrupadaANotificationData(agrupado, convocatorias, sancion)
	notificationData.SemestreAcademico = periodo

	if err := h.validateNotificationData(notificationData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := h.smtpService.SendNotification(notificationData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error al enviar notificación por email",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"message":       "Notificación enviada exitosamente",
		"email_sent_to": notificationData.Alumno.CorreoInstitucional,
	})
}

func ObtenerConvocatorias() ([]models.Convocatoria, error) {
	return []models.Convocatoria{}, nil
}

func (h *handlerNotification) GetSancionPorFaltaID(c *fiber.Ctx) error {
	faltaID := c.Params("falta_id")
	if faltaID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "El falta_id es requerido en la URL",
		})
	}

	sancion, err := smtp.GetSancionAsignadaPorFaltaInternal(faltaID, h.sancionService)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "No hay sanción asignada para esa falta",
		})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"falta_id": faltaID,
		"sancion":  sancion,
	})
}

func (h *handlerNotification) SendNotificationEmail(c *fiber.Ctx) error {
	var notificationData smtp.NotificationData
	if err := c.BodyParser(&notificationData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Error al parsear datos JSON",
			"error":   err.Error(),
		})
	}

	if err := h.validateNotificationData(notificationData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := h.smtpService.SendNotification(notificationData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error al enviar notificación por email",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"message":       "Notificación enviada exitosamente",
		"email_sent_to": notificationData.Alumno.CorreoInstitucional,
	})
}

func (h *handlerNotification) SendNotificationEmailWithCC(c *fiber.Ctx) error {
	var requestData struct {
		NotificationData smtp.NotificationData `json:"notification_data"`
		CCEmails         []string              `json:"cc_emails"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Error al parsear datos JSON",
			"error":   err.Error(),
		})
	}

	if err := h.validateNotificationData(requestData.NotificationData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := h.smtpService.SendNotificationWithCC(requestData.NotificationData, requestData.CCEmails); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error al enviar notificación por email",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"message":       "Notificación enviada exitosamente con copias",
		"email_sent_to": requestData.NotificationData.Alumno.CorreoInstitucional,
		"cc_emails":     requestData.CCEmails,
	})
}

func (h *handlerNotification) TestSMTPConnection(c *fiber.Ctx) error {
	if err := h.smtpService.TestConnection(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error en conexión SMTP",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Conexión SMTP exitosa",
	})
}

// ✅ VERSIÓN MEJORADA con logging detallado del error
func (h *handlerNotification) GetNotificationInfoByFaltaID(c *fiber.Ctx) error {
	faltaID := c.Params("falta_id")
	if faltaID == "" {
		logger.Error.Printf("GetNotificationInfoByFaltaID: falta_id vacío")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "El falta_id es requerido en la URL",
		})
	}

	logger.Info.Printf("GetNotificationInfoByFaltaID: Consultando info para falta_id=%s", faltaID)

	// Intentar obtener la información existente
	info, err := dbx.GetNotificationInfo(faltaID)
	if err != nil {
		logger.Warning.Printf("GetNotificationInfoByFaltaID: No existe info para falta_id=%s, creando nueva. Error: %v", faltaID, err)

		// Si no existe, intentar crear
		numeroNotificacion, err2 := dbx.GetOrCreateNotificationNumber(faltaID)
		if err2 != nil {
			// ✅ AQUÍ ESTÁ EL CAMBIO: Retornar el error específico
			logger.Error.Printf("GetNotificationInfoByFaltaID: Error al crear número de notificación para falta_id=%s: %v", faltaID, err2)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "No se pudo generar el número de notificación",
				"error":   err2.Error(), // ✅ Agregar el error específico
				"details": fmt.Sprintf("falta_id: %s", faltaID),
			})
		}

		logger.Info.Printf("GetNotificationInfoByFaltaID: Número creado exitosamente: %s", numeroNotificacion)

		// Intentar obtener de nuevo después de crear
		info, err = dbx.GetNotificationInfo(faltaID)
		if err != nil {
			logger.Error.Printf("GetNotificationInfoByFaltaID: Error al obtener info después de crear para falta_id=%s: %v", faltaID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "No se pudo obtener la información de notificación después de crearla",
				"error":   err.Error(),
			})
		}
	}

	logger.Info.Printf("GetNotificationInfoByFaltaID: Info obtenida exitosamente para falta_id=%s: %s", faltaID, info.NumeroNotificacion)

	return c.JSON(fiber.Map{
		"success":             true,
		"falta_id":            faltaID,
		"numero_notificacion": info.NumeroNotificacion,
		"fecha":               info.Fecha,
	})
}

// ✅ NUEVO ENDPOINT: Para verificar la conexión y estado de la BD
func (h *handlerNotification) TestDatabaseConnection(c *fiber.Ctx) error {
	db := dbx.GetConnection()

	// Test ping
	if err := db.Ping(); err != nil {
		logger.Error.Printf("TestDatabaseConnection: Error en ping: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error al hacer ping a la base de datos",
			"error":   err.Error(),
		})
	}

	// Test query simple
	var result int
	if err := db.Get(&result, "SELECT 1"); err != nil {
		logger.Error.Printf("TestDatabaseConnection: Error en query simple: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error al ejecutar query simple",
			"error":   err.Error(),
		})
	}

	// Test existencia de tablas
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'notificacion'")
	if err != nil || count == 0 {
		logger.Error.Printf("TestDatabaseConnection: Tabla 'notificacion' no existe")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "La tabla 'notificacion' no existe en la base de datos",
			"error":   fmt.Sprintf("count=%d, err=%v", count, err),
		})
	}

	err = db.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'notificacion_secuencia'")
	if err != nil || count == 0 {
		logger.Error.Printf("TestDatabaseConnection: Tabla 'notificacion_secuencia' no existe")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "La tabla 'notificacion_secuencia' no existe en la base de datos",
			"error":   fmt.Sprintf("count=%d, err=%v", count, err),
		})
	}

	logger.Info.Printf("TestDatabaseConnection: Todas las verificaciones pasaron exitosamente")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Conexión a base de datos exitosa",
		"details": fiber.Map{
			"engine":                       dbx.DBEngine,
			"tabla_notificacion":           "existe",
			"tabla_notificacion_secuencia": "existe",
		},
	})
}
