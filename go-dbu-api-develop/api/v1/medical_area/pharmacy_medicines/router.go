package pharmacy_medicines

import (
	"dbu-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// RouterPharmacyMedicines registra las rutas del módulo de medicamentos
func RouterPharmacyMedicines(app *fiber.App, db *sqlx.DB, txID string) {
	h := handlerPharmacyMedicines{db: db, txID: txID}

	v1 := app.Group("/v1")
	pharmacy := v1.Group("/pharmacy")
	pharmacy.Use(middleware.JWTProtected())

	medicines := pharmacy.Group("/medicines")
	medicines.Post("/", h.createMedicine)
	medicines.Put("/:id", h.updateMedicine)
	medicines.Delete("/:id", h.deleteMedicine)
	medicines.Get("/:id", h.getMedicineByID)
	medicines.Get("/code/:code", h.getMedicineByCode)
	medicines.Get("/", h.searchMedicines)
	medicines.Get("/all/with-stock", h.getAllMedicinesWithStock)
}
