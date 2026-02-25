package pharmacy_medicines

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"dbu-api/pkg/medical_area"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type handlerPharmacyMedicines struct {
	db   *sqlx.DB
	txID string
}

// createMedicine crea un nuevo medicamento
// @Summary Crear medicamento
// @Description Crea un nuevo medicamento en el catálogo
// @Tags Pharmacy - Medicines
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param medicine body requestCreateMedicine true "Datos del medicamento"
// @Success 201 {object} responseCreateMedicine
// @Router /api/v1/pharmacy/medicines [post]
func (h *handlerPharmacyMedicines) createMedicine(c *fiber.Ctx) error {
	res := responseCreateMedicine{Error: true}

	// Obtener usuario del contexto
	user := c.Locals("user").(*models.User)

	// Parsear request
	var req requestCreateMedicine
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't bind model: %v", err)
		res.Code, res.Message = 1, "invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Generar UUID si no viene
	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	// Crear servicio
	srv := medical_area.NewServerMedicalArea(h.db, user, h.txID)

	// Crear medicamento
	medicine, code, err := srv.SrvPharmacyMedicines.CreateMedicine(
		req.ID,
		req.Codigo,
		req.NombreGenerico,
		req.NombreComercial,
		req.FormaFarmaceutica,
		req.Concentracion,
		req.ViaAdministracion,
		req.RequiereReceta,
		req.Controlado,
		req.Descripcion,
		req.Estado,
	)

	if err != nil {
		logger.Error.Printf("couldn't create medicine: %v", err)
		res.Code, res.Message = code, err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	res.Data = medicine
	res.Error = false
	res.Code, res.Message = 29, "medicine created successfully"
	return c.Status(fiber.StatusCreated).JSON(res)
}

// updateMedicine actualiza un medicamento existente
// @Summary Actualizar medicamento
// @Description Actualiza un medicamento existente
// @Tags Pharmacy - Medicines
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param id path string true "ID del medicamento"
// @Param medicine body requestUpdateMedicine true "Datos del medicamento"
// @Success 200 {object} responseUpdateMedicine
// @Router /api/v1/pharmacy/medicines/{id} [put]
func (h *handlerPharmacyMedicines) updateMedicine(c *fiber.Ctx) error {
	res := responseUpdateMedicine{Error: true}

	// Obtener usuario del contexto
	user := c.Locals("user").(*models.User)

	// Obtener ID del path
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		res.Code, res.Message = 15, "invalid UUID"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Parsear request
	var req requestUpdateMedicine
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't bind model: %v", err)
		res.Code, res.Message = 1, "invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Crear servicio
	srv := medical_area.NewServerMedicalArea(h.db, user, h.txID)

	// Obtener medicamento actual
	medicine, code, err := srv.SrvPharmacyMedicines.GetMedicineByID(id)
	if err != nil {
		logger.Error.Printf("couldn't get medicine: %v", err)
		res.Code, res.Message = code, err.Error()
		return c.Status(fiber.StatusNotFound).JSON(res)
	}

	// Actualizar campos
	medicine.Codigo = req.Codigo
	medicine.NombreGenerico = req.NombreGenerico
	medicine.NombreComercial = req.NombreComercial
	medicine.FormaFarmaceutica = req.FormaFarmaceutica
	medicine.Concentracion = req.Concentracion
	medicine.ViaAdministracion = req.ViaAdministracion
	medicine.RequiereReceta = req.RequiereReceta
	medicine.Controlado = req.Controlado
	medicine.Descripcion = req.Descripcion
	medicine.Estado = req.Estado

	// Actualizar
	code, err = srv.SrvPharmacyMedicines.UpdateMedicine(medicine)
	if err != nil {
		logger.Error.Printf("couldn't update medicine: %v", err)
		res.Code, res.Message = code, err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	res.Data = medicine
	res.Error = false
	res.Code, res.Message = 29, "medicine updated successfully"
	return c.Status(fiber.StatusOK).JSON(res)
}

// deleteMedicine elimina un medicamento (soft delete)
// @Summary Eliminar medicamento
// @Description Elimina un medicamento del catálogo (soft delete)
// @Tags Pharmacy - Medicines
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param id path string true "ID del medicamento"
// @Success 200 {object} responseDeleteMedicine
// @Router /api/v1/pharmacy/medicines/{id} [delete]
func (h *handlerPharmacyMedicines) deleteMedicine(c *fiber.Ctx) error {
	res := responseDeleteMedicine{Error: true}

	// Obtener usuario del contexto
	user := c.Locals("user").(*models.User)

	// Obtener ID del path
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		res.Code, res.Message = 15, "invalid UUID"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Crear servicio
	srv := medical_area.NewServerMedicalArea(h.db, user, h.txID)

	// Eliminar
	code, err := srv.SrvPharmacyMedicines.DeleteMedicine(id)
	if err != nil {
		logger.Error.Printf("couldn't delete medicine: %v", err)
		res.Code, res.Message = code, err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	res.Error = false
	res.Code, res.Message = 28, "medicine deleted successfully"
	return c.Status(fiber.StatusOK).JSON(res)
}

// getMedicineByID obtiene un medicamento por ID
// @Summary Obtener medicamento por ID
// @Description Obtiene un medicamento por su ID
// @Tags Pharmacy - Medicines
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param id path string true "ID del medicamento"
// @Success 200 {object} responseGetMedicine
// @Router /api/v1/pharmacy/medicines/{id} [get]
func (h *handlerPharmacyMedicines) getMedicineByID(c *fiber.Ctx) error {
	res := responseGetMedicine{Error: true}

	// Obtener usuario del contexto
	user := c.Locals("user").(*models.User)

	// Obtener ID del path
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		res.Code, res.Message = 15, "invalid UUID"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Crear servicio
	srv := medical_area.NewServerMedicalArea(h.db, user, h.txID)

	// Obtener medicamento
	medicine, code, err := srv.SrvPharmacyMedicines.GetMedicineByID(id)
	if err != nil {
		logger.Error.Printf("couldn't get medicine: %v", err)
		res.Code, res.Message = code, err.Error()
		return c.Status(fiber.StatusNotFound).JSON(res)
	}

	res.Data = medicine
	res.Error = false
	res.Code, res.Message = 29, "success"
	return c.Status(fiber.StatusOK).JSON(res)
}

// getMedicineByCode obtiene un medicamento por código
// @Summary Obtener medicamento por código
// @Description Obtiene un medicamento por su código
// @Tags Pharmacy - Medicines
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param code path string true "Código del medicamento"
// @Success 200 {object} responseGetMedicine
// @Router /api/v1/pharmacy/medicines/code/{code} [get]
func (h *handlerPharmacyMedicines) getMedicineByCode(c *fiber.Ctx) error {
	res := responseGetMedicine{Error: true}

	// Obtener usuario del contexto
	user := c.Locals("user").(*models.User)

	// Obtener código del path
	code := c.Params("code")
	if code == "" {
		res.Code, res.Message = 15, "code is required"
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	// Crear servicio
	srv := medical_area.NewServerMedicalArea(h.db, user, h.txID)

	// Obtener medicamento
	medicine, resCode, err := srv.SrvPharmacyMedicines.GetMedicineByCode(code)
	if err != nil {
		logger.Error.Printf("couldn't get medicine: %v", err)
		res.Code, res.Message = resCode, err.Error()
		return c.Status(fiber.StatusNotFound).JSON(res)
	}

	res.Data = medicine
	res.Error = false
	res.Code, res.Message = 29, "success"
	return c.Status(fiber.StatusOK).JSON(res)
}

// searchMedicines busca medicamentos con filtros
// @Summary Buscar medicamentos
// @Description Busca medicamentos con filtros y paginación
// @Tags Pharmacy - Medicines
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param search query string false "Búsqueda por código, nombre genérico o comercial"
// @Param estado query string false "Filtrar por estado (ACTIVO/INACTIVO)"
// @Param limit query int false "Límite de resultados" default(10)
// @Param offset query int false "Offset para paginación" default(0)
// @Success 200 {object} responseSearchMedicines
// @Router /api/v1/pharmacy/medicines [get]
func (h *handlerPharmacyMedicines) searchMedicines(c *fiber.Ctx) error {
	res := responseSearchMedicines{Error: true}

	// Obtener usuario del contexto
	user := c.Locals("user").(*models.User)

	// Obtener parámetros de query
	search := c.Query("search", "")
	estado := c.Query("estado", "")
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.Query("offset", "0"), 10, 64)

	// Crear servicio
	srv := medical_area.NewServerMedicalArea(h.db, user, h.txID)

	// Buscar medicamentos
	medicines, total, err := srv.SrvPharmacyMedicines.SearchMedicines(search, estado, limit, offset)
	if err != nil {
		logger.Error.Printf("couldn't search medicines: %v", err)
		res.Code, res.Message = 3, err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res.Data = medicines
	res.Total = total
	res.Error = false
	res.Code, res.Message = 29, "success"
	return c.Status(fiber.StatusOK).JSON(res)
}

// getAllMedicinesWithStock obtiene todos los medicamentos con stock
// @Summary Obtener todos los medicamentos con stock
// @Description Obtiene todos los medicamentos con información de stock
// @Tags Pharmacy - Medicines
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} responseGetAllMedicinesWithStock
// @Router /api/v1/pharmacy/medicines/all/with-stock [get]
func (h *handlerPharmacyMedicines) getAllMedicinesWithStock(c *fiber.Ctx) error {
	res := responseGetAllMedicinesWithStock{Error: true}

	// Obtener usuario del contexto
	user := c.Locals("user").(*models.User)

	// Crear servicio
	srv := medical_area.NewServerMedicalArea(h.db, user, h.txID)

	// Obtener medicamentos
	medicines, err := srv.SrvPharmacyMedicines.GetAllMedicinesWithStock()
	if err != nil {
		logger.Error.Printf("couldn't get medicines with stock: %v", err)
		res.Code, res.Message = 3, err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res.Data = medicines
	res.Error = false
	res.Code, res.Message = 29, "success"
	return c.Status(fiber.StatusOK).JSON(res)
}
