package report

import "time"

// safeDate función utilitaria para formatear fechas
func safeDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// El archivo y la función RegisterReportRoutes están correctos.
// Si recibes 404, el problema está en el registro de rutas en tu archivo principal (main.go).
// Asegúrate de tener en main.go algo como esto:

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
	api := app.Group("/api/v1")
	normative := api.Group("/normative")
	reportGroup := normative.Group("/report")
	report.RegisterReportRoutes(reportGroup, db)
	app.Listen(":4021")
}
*/

// Verifica que el método sea GET y que el servidor esté corriendo en el puerto 4021.
// Si todo esto está así, la ruta /api/v1/normative/report/excel debe funcionar.
