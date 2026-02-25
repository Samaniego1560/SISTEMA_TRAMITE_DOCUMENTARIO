package api

import (
	"dbu-api/api/health"
	authentication "dbu-api/api/v1/auth"
	"dbu-api/api/v1/automation/room_assignment"
	"dbu-api/api/v1/general_visit/visita_general"
	"dbu-api/api/v1/general_visit/visita_residente"
	"dbu-api/api/v1/medical_area/announcement_signatures"
	"dbu-api/api/v1/medical_area/dentistry_consultation"
	"dbu-api/api/v1/medical_area/exam_toxicologico"
	"dbu-api/api/v1/medical_area/medical_consultation"
	"dbu-api/api/v1/medical_area/nursing_consultation"
	"dbu-api/api/v1/medical_area/nursing_consultation_assignment"
	"dbu-api/api/v1/medical_area/patients"
	"dbu-api/api/v1/medical_area/payments"
	"dbu-api/api/v1/medical_area/pharmacy_medicines"
	"dbu-api/api/v1/medical_area/reports"
	"dbu-api/api/v1/normative"
	"dbu-api/api/v1/normative/articles"
	"dbu-api/api/v1/normative/chapters"
	incisos "dbu-api/api/v1/normative/inciso"
	"dbu-api/api/v1/normative/report"
	"dbu-api/api/v1/normative/resolutions"
	sanctionConfig "dbu-api/api/v1/normative/sanction_configuration_fault"
	"dbu-api/api/v1/notificacion"
	"dbu-api/api/v1/oauth"
	"dbu-api/api/v1/psicopedagogia"
	"dbu-api/api/v1/residence/residences"
	"dbu-api/api/v1/residence/rooms"
	"dbu-api/api/v1/sanction/faults"
	"dbu-api/api/v1/survey/student_parent"
	"dbu-api/api/v1/student/assignment"
	studentAuth "dbu-api/api/v1/student/auth"
	"dbu-api/api/v1/submission/submissions"

	"dbu-api/internal/env"
	"dbu-api/pkg/notificacion/smtp"
	"dbu-api/pkg/sanction/configuration_sanction_fault" // ✅ Agregado
	"dbu-api/pkg/sanction/fault"
	"log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func routes(loggerHttp bool, allowedOrigins string, db *sqlx.DB) *fiber.App {
	app := fiber.New()

	prometheus := fiberprometheus.New("DBU API REST")
	prometheus.RegisterAt(app, "/metrics")

	app.Get("/doc/*", swagger.New(swagger.Config{
		URL:         "/doc/doc.json",
		DeepLinking: false,
	}))

	app.Use(recover.New())
	app.Use(prometheus.Middleware)
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, X-Requested-With, Content-Type, Accept, Authorization, signature, submission_id",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
	}))

	if loggerHttp {
		app.Use(logger.New())
	}

	TxID := uuid.New().String()

	// ✅ Crear el servicio de sanciones PRIMERO
	sancionRepo := configuration_sanction_fault.FactoryStorage(db, TxID)
	sancionService := configuration_sanction_fault.NewSanctionFaultService(sancionRepo, nil, TxID)

	// ✅ Inicializar el servicio SMTP usando el singleton de configuración
	var smtpService *smtp.SMTPService
	cfg := env.NewConfiguration()

	// Verificar si hay configuración SMTP en config.json
	if cfg.SMTP.Host != "" && cfg.SMTP.Port > 0 {
		smtpConfig := smtp.NewSMTPConfigFromJSON(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.Username,
			cfg.SMTP.Password,
			cfg.SMTP.From,
		)
		smtpService = smtp.NewSMTPService(smtpConfig, sancionService)
		log.Println("SMTP service inicializado desde config.json")
	}

	loadRoutes(app, TxID, db, smtpService, sancionService) // ✅ Pasar sancionService

	return app
}

// ✅ Agregar sancionService como parámetro
func loadRoutes(app *fiber.App, TxID string, db *sqlx.DB, smtpService *smtp.SMTPService, sancionService smtp.SancionService) {
	// Core modules
	submissions.RouterSubmissions(app, db, TxID)
	residences.RouterResidences(app, db, TxID)
	rooms.RouterRooms(app, db, TxID)
	room_assignment.RouterRoomAssignment(app, db, TxID)
	health.RouterHealth(app, db, TxID)
	authentication.RouterAuthentication(app, db, TxID)
	oauth.RouterOAuth(app, db, TxID) // OAuth authentication (Microsoft, etc.)

	// Student Portal - Portal de Estudiantes
	studentAuth.RegisterRoutes(app, db, smtpService, TxID) // Autenticación con OTP
	assignment.RegisterRoutes(app, db)                     // Consulta de asignaciones

	// 🏥 MEDICAL AREA - Área Médica
	patients.RouterMedicalArea(app, db, TxID)                        // Gestión de pacientes
	nursing_consultation.RouterMedicalArea(app, db, TxID)            // Consultas de enfermería
	dentistry_consultation.RouterMedicalArea(app, db, TxID)          // Consultas odontológicas
	medical_consultation.RouterMedicalArea(app, db, TxID)            // Consultas médicas generales
	nursing_consultation_assignment.RouterMedicalArea(app, db, TxID) // Asignación de consultas
	announcement_signatures.RouterMedicalArea(app, db, TxID)         // Firmas de anuncios
	reports.RouterMedicalArea(app, db, TxID)                         // Reportes médicos
	exam_toxicologico.RouterMedicalArea(app, db, TxID)               // Exámenes toxicológicos
	payments.Router(app, db, TxID)                                   // Pagos
	pharmacy_medicines.RouterPharmacyMedicines(app, db, TxID)        // Farmacia - Medicamentos
	// Normative
	resolutions.RouterResoluciones(app, db, TxID)
	articles.RouterArticles(app, db, TxID)
	chapters.RouterChapters(app, db, TxID)
	incisos.RouterIncisos(app, db, TxID)
	sanctionConfig.RouterSancionesConfiguration(app, db, TxID)
	report.RouterReport(app, db)
	normative.RouterReportStats(app, db)

	// Sanctions
	faults.RouterFaults(app, db, TxID)
	faultRepo := fault.FactoryStorage(db, TxID)
	faultService := fault.NewFaultService(faultRepo, nil, TxID)

	// ✅ Notifications - Pasar sancionService
	notificacion.RouterNotification(app, smtpService, TxID, faultService, sancionService)

	// Psicopedagogía y visitas
	psicopedagogia.RouterPsicopedagogia(app, db, TxID)
	visita_general.RouterVisitaGeneral(app, db, TxID)
	visita_residente.RouterVisitaResidente(app, db, TxID)

	// Encuesta estudiantes (padres/madres)
	student_parent.RouterStudentParentSurvey(app, db, TxID)
}
