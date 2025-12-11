package report

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// RouterReport configura el grupo de rutas para el reporte normativo
func RouterReport(app *fiber.App, db *sqlx.DB) {
	api := app.Group("/api/v1") // <-- Corrige aquí el prefijo
	normative := api.Group("/normative")
	reportGroup := normative.Group("/report")
	handler := &ReportHandler{db: db}

	// Convocatoria handler
	convocatoriaHandler := NewConvocatoriaHandler(db)

	// Cambia la ruta para usar parámetro de ruta
	reportGroup.Get("/excel/:convocatoria_id", handler.GetStudentReportExcel)

	// Nueva ruta para obtener todas las convocatorias
	reportGroup.Get("/convocatorias", convocatoriaHandler.GetAllConvocatorias)
}

// Asegúrate de que en tu main.go realmente estés llamando a RouterReport(app, db)
// y que no estés usando RegisterReportRoutes ni otro registro duplicado.
// Ejemplo correcto en main.go:

/*
import (
	// ...existing code...
	"c:\Users\Junior\Documents\DESARROLLO\go-dbu-api\api\v1\normative\report"
	"github.com/jmoiron/sqlx"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	db, err := sqlx.Connect("mysql", "usuario:password@tcp(localhost:3306)/nombre_db")
	if err != nil {
		panic(err)
	}
	report.RouterReport(app, db) // <-- Debe estar así
	app.Listen(":4021")
}
*/

// Si esto ya está así y sigue saliendo 404, revisa:
// - Que no tengas middlewares que bloqueen rutas.
// - Que no haya errores de importación o inicialización.
// - Que no estés usando otro archivo/router que registre rutas diferentes o sobrescriba el grupo.
// - Que el servidor esté corriendo en el puerto correcto y sin errores previos en consola.

// El código del router está correcto para exponer /api/v1/normative/report/excel.
