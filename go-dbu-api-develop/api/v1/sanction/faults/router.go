package faults

import (
	"dbu-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func RouterFaults(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerFaults{db: db, txID: txID}
	v1 := app.Group("/v1")
	faults := v1.Group("/faltas")
	faults.Use(middleware.JWTProtected())
	faults.Post("/", h.CreateFault)
	faults.Put("/", h.UpdateFault)
	faults.Get("/", h.GetAllFaults)
	faults.Delete("/:id", h.DeleteFaults)
	faults.Get("/:id", h.GetFaultsByID)
	faults.Get("/:id/alumnos", h.GetAllStudentsByFault)
	faults.Patch("/:id/estado", h.UpdateEstadoFalta)
	v1.Get("/alumnos/perfil/:dni", h.GetAlumnoProfileByDNI)
	v1.Get("/alumnos/perfil-codigo/:codigo_estudiante", h.GetAlumnoProfileByCodigo)
	faults.Post("/:id/documentos", h.CrearDocumentoFalta)
	v1.Post("/faltas/upload", h.UploadDocumento)
	faults.Get("/:id/detalle", h.GetDetalleFaltaAgrupado)
	app.Static("/uploads", "./uploads")
	faults.Get("/documentos/:id", h.DescargarDocumento)
	faults.Get(":id/documentos", h.GetDocumentosPorFalta)
	v1.Get("/alumnos/validar-postulacion/:dni/:servicio_id", h.ValidarPostulacionPorDNI)
	v1.Post("/alumnos/validar-postulacion-multiple/:dni", h.ValidarPostulacionMultiple)
	// Endpoint para obtener todos los servicios
	v1.Get("/servicios", h.GetAllServicios)
}
