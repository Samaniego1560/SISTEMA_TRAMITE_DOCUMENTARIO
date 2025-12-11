package smtp

import (
	"dbu-api/internal/dbx"
	"fmt"
	"log"
	"strings"

	"gopkg.in/gomail.v2"
)

type SMTPService struct {
	config         *SMTPConfig
	dialer         *gomail.Dialer
	sancionService SancionService // ✅ Usa la interfaz local
}

// ✅ Constructor actualizado con interfaz local
func NewSMTPService(config *SMTPConfig, sancionService SancionService) *SMTPService {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	return &SMTPService{
		config:         config,
		dialer:         dialer,
		sancionService: sancionService,
	}
}

func (s *SMTPService) SendNotification(data NotificationData) error {
	return s.SendNotificationWithCC(data, nil)
}

func (s *SMTPService) SendNotificationWithCC(data NotificationData, ccEmails []string) error {
	numeroNotificacion, err := dbx.GetOrCreateNotificationNumber(data.FaltaID)
	if err != nil {
		return fmt.Errorf("error generando/consultando número de notificación: %v", err)
	}

	// ✅ Usa la llamada interna
	sancion, _ := GetSancionAsignadaPorFaltaInternal(data.FaltaID, s.sancionService)
	templateData := MapNotificationDataToTemplate(data, numeroNotificacion, sancion)

	htmlContent, err := GenerateNotificationHTML(templateData)
	if err != nil {
		return fmt.Errorf("error al generar HTML del email: %v", err)
	}

	faltasDetalle := BuildFaltasDetalleText(data)

	message := gomail.NewMessage()
	message.SetHeader("From", s.config.From)
	message.SetHeader("To", data.Alumno.CorreoInstitucional)
	if len(ccEmails) > 0 {
		message.SetHeader("Cc", ccEmails...)
	}
	message.SetHeader("Subject", fmt.Sprintf("Notificación N° %s - Dirección de Bienestar Universitario", numeroNotificacion))
	message.SetBody("text/html", htmlContent)

	textContent := s.generateTextContent(templateData, faltasDetalle)
	message.AddAlternative("text/plain", textContent)

	if err := s.dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("error al enviar email: %v", err)
	}

	log.Printf("Notificación enviada exitosamente a %s con CC a %v",
		data.Alumno.CorreoInstitucional, ccEmails)

	return nil
}

func (s *SMTPService) generateTextContent(data EmailTemplate, faltasDetalle string) string {
	resolucionNombre := ""
	if len(data.Resolucion.Capitulos) > 0 {
		capituloID := data.Resolucion.Capitulos[0].CapituloID
		nombre, err := dbx.GetResolucionNombreVigente(capituloID)
		if err == nil {
			resolucionNombre = nombre
		}
	}
	return fmt.Sprintf(`
	Universidad Nacional Agraria de la Selva - UNAS
	Tingo María
	Dirección de Bienestar Universitario

	Notificación N° %s

	Estudiante: %s
	Facultad: %s
	Fecha: %s
	Dirección: %s

	Por la presente se notifica a UD. que según el %s, usted ha incurrido en %d una FALTA %s, de acuerdo con la Directiva que regula el uso del servicio de %s de la UNAS, aprobada mediante Resolución (%s).
	La falta está tipificada en los siguientes casos:
	%s

	La sanción correspondiente se encuentra en el %s, %s, Inc. %s de la misma Directiva: que establece "%s."
	Esta sanción será registrada en el sistema de la Dirección de Bienestar Universitario (DBU) y en su expediente personal (file).

	NOTA: Caso contrario usted omita la presente notificación, se estarán tomando otras medidas que defina la comisión de disciplina.

	---
	Este es un documento oficial generado por el Sistema de la Dirección de Bienestar Universitario de la Universidad Nacional Agraria de la Selva - UNAS
	`,
		data.NumeroNotificacion,
		data.Estudiante,
		data.Facultad,
		data.Fecha,
		data.Direccion,
		data.FuenteInformacion,
		data.ContadorFaltas,
		data.GravedadFalta,
		data.Servicio,
		resolucionNombre,
		faltasDetalle,
		data.CapituloSancion,
		data.ArticuloSancion,
		data.IncisoSancion,
		data.DescripcionSancion,
	)
}

func BuildFaltasDetalleText(data NotificationData) string {
	var faltasDetalle strings.Builder

	for _, capitulo := range data.Resolucion.Capitulos {
		for _, articulo := range capitulo.Articulos {
			for _, inciso := range articulo.Incisos {
				faltasDetalle.WriteString(fmt.Sprintf(
					"%s %s, Gravedad: %s, Inc. %s: \"%s\"\n",
					capitulo.CapituloNombre,
					articulo.ArticuloDescripcion,
					articulo.ArticuloGravedad,
					inciso.IncisoNombre,
					inciso.IncisoDescripcion,
				))
			}
		}
	}
	return faltasDetalle.String()
}

func (s *SMTPService) TestConnection() error {
	message := gomail.NewMessage()
	message.SetHeader("From", s.config.From)
	message.SetHeader("To", s.config.From)
	message.SetHeader("Subject", "Test de conexión SMTP")
	message.SetBody("text/plain", "Esta es una prueba de conexión SMTP.")

	if err := s.dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("error en test de conexión SMTP: %v", err)
	}

	log.Println("Test de conexión SMTP exitoso")
	return nil
}

// SendOTPCode envía un código OTP al correo del estudiante
func (s *SMTPService) SendOTPCode(correo, nombreEstudiante, codigoOTP string, minutosExpiracion int) error {
	// Generar HTML del template
	htmlContent, err := GenerateOTPHTML("./templates/templateOTP.html", nombreEstudiante, codigoOTP, minutosExpiracion)
	if err != nil {
		return fmt.Errorf("error al generar HTML del email OTP: %v", err)
	}

	// Crear mensaje
	message := gomail.NewMessage()
	message.SetHeader("From", s.config.From)
	message.SetHeader("To", correo)
	message.SetHeader("Subject", "Código de Verificación - Portal Estudiantil DBU")
	message.SetBody("text/html", htmlContent)

	// Alternativa en texto plano
	textContent := fmt.Sprintf(`
Universidad Nacional Agraria de la Selva - UNAS
Dirección de Bienestar Universitario

Código de Verificación - Portal Estudiantil

Hola %s,

Tu código de verificación es: %s

Este código expira en %d minutos.

IMPORTANTE:
- No compartas este código con nadie
- El personal de la DBU nunca te pedirá este código
- Si no solicitaste este código, ignora este correo

Este es un correo automático. Por favor, no respondas a este mensaje.
	`, nombreEstudiante, codigoOTP, minutosExpiracion)

	message.AddAlternative("text/plain", textContent)

	// Enviar email
	if err := s.dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("error al enviar email OTP: %v", err)
	}

	log.Printf("Código OTP enviado exitosamente a %s", correo)
	return nil
}
