package pharmacy_medicines

import (
	"dbu-api/pkg/medical_area/pharmacy_medicines"
)

// Request para crear medicamento
type requestCreateMedicine struct {
	ID                string  `json:"id"`
	Codigo            string  `json:"codigo" validate:"required"`
	NombreGenerico    string  `json:"nombre_generico" validate:"required"`
	NombreComercial   *string `json:"nombre_comercial"`
	FormaFarmaceutica string  `json:"forma_farmaceutica" validate:"required"`
	Concentracion     string  `json:"concentracion" validate:"required"`
	ViaAdministracion *string `json:"via_administracion"`
	RequiereReceta    bool    `json:"requiere_receta"`
	Controlado        bool    `json:"controlado"`
	Descripcion       *string `json:"descripcion"`
	Estado            string  `json:"estado" validate:"required"`
}

// Request para actualizar medicamento
type requestUpdateMedicine struct {
	Codigo            string  `json:"codigo" validate:"required"`
	NombreGenerico    string  `json:"nombre_generico" validate:"required"`
	NombreComercial   *string `json:"nombre_comercial"`
	FormaFarmaceutica string  `json:"forma_farmaceutica" validate:"required"`
	Concentracion     string  `json:"concentracion" validate:"required"`
	ViaAdministracion *string `json:"via_administracion"`
	RequiereReceta    bool    `json:"requiere_receta"`
	Controlado        bool    `json:"controlado"`
	Descripcion       *string `json:"descripcion"`
	Estado            string  `json:"estado" validate:"required"`
}

// Response para crear medicamento
type responseCreateMedicine struct {
	Error   bool                         `json:"error"`
	Code    int                          `json:"code"`
	Message string                       `json:"message"`
	Data    *pharmacy_medicines.Medicine `json:"data,omitempty"`
}

// Response para actualizar medicamento
type responseUpdateMedicine struct {
	Error   bool                         `json:"error"`
	Code    int                          `json:"code"`
	Message string                       `json:"message"`
	Data    *pharmacy_medicines.Medicine `json:"data,omitempty"`
}

// Response para eliminar medicamento
type responseDeleteMedicine struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response para obtener medicamento
type responseGetMedicine struct {
	Error   bool                         `json:"error"`
	Code    int                          `json:"code"`
	Message string                       `json:"message"`
	Data    *pharmacy_medicines.Medicine `json:"data,omitempty"`
}

// Response para buscar medicamentos
type responseSearchMedicines struct {
	Error   bool                                    `json:"error"`
	Code    int                                     `json:"code"`
	Message string                                  `json:"message"`
	Data    []*pharmacy_medicines.MedicineWithStock `json:"data,omitempty"`
	Total   int64                                   `json:"total"`
}

// Response para obtener todos los medicamentos con stock
type responseGetAllMedicinesWithStock struct {
	Error   bool                                    `json:"error"`
	Code    int                                     `json:"code"`
	Message string                                  `json:"message"`
	Data    []*pharmacy_medicines.MedicineWithStock `json:"data,omitempty"`
}
