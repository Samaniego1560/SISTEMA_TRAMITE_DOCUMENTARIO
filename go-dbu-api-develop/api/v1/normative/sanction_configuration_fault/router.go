package sanction_configuration_fault

import (
	"dbu-api/pkg/sanction/configuration_sanction_fault" // Importa SOLO si necesitas instanciar el service

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Recuerda usar tu handler constructor correctamente:
func RouterSancionesConfiguration(app *fiber.App, db *sqlx.DB, txID string) {
	v1 := app.Group("/v1")
	sancionesConfiguration := v1.Group("/configuration")
	repo := configuration_sanction_fault.FactoryStorage(db, txID)
	service := configuration_sanction_fault.NewSanctionFaultService(repo, nil, txID)
	h := NewHandlerSancionesConfi(db, txID, service)

	sancionesConfiguration.Get("/:sancion_id", h.GetSancionByID)
	sancionesConfiguration.Post("/", h.CreateSancion)
	sancionesConfiguration.Put("/", h.UpdateSancion)
	sancionesConfiguration.Delete("/:id", h.DeleteSancion)
	sancionesConfiguration.Get("/", h.GetAllSanciones)
	sancionesConfiguration.Post("/asignar-sancion-falta/:falta_id", h.AsignarSancionFalta)
	sancionesConfiguration.Get("/asignadas-por-falta/:falta_id", h.GetSancionesAsignadasPorFalta)
	sancionesConfiguration.Post("/registrar-apelacion/:sancion_falta_asignada_id", h.RegistrarApelacion)
	sancionesConfiguration.Get("/apelaciones/:sancion_falta_asignada_id", h.GetApelacionesPorSancionFaltaAsignada)
	sancionesConfiguration.Get("/detalle-apelacion/:apelacion_id", h.GetDetalleApelacion)
	sancionesConfiguration.Get("/documento-apelacion/:documento_id", h.DescargarDocumentoApelacion)
	sancionesConfiguration.Put("/resolver-apelacion/:apelacion_id", h.ResolverApelacion)
	sancionesConfiguration.Get("/por-falta/:falta_id", h.GetSancionesByFalta)
	sancionesConfiguration.Get("/sanciones-sugeridas/:falta_id", h.GetSancionesSugeridas)
}
