package pharmacy_medicines

import (
	"github.com/asaskevich/govalidator"
	"time"
)

// Medicine representa un medicamento en el catálogo
type Medicine struct {
	ID                 string     `json:"id" db:"id" valid:"uuid,required"`
	Codigo             string     `json:"codigo" db:"codigo" valid:"required"`
	NombreGenerico     string     `json:"nombre_generico" db:"nombre_generico" valid:"required"`
	NombreComercial    *string    `json:"nombre_comercial" db:"nombre_comercial" valid:"-"`
	FormaFarmaceutica  string     `json:"forma_farmaceutica" db:"forma_farmaceutica" valid:"required"`
	Concentracion      string     `json:"concentracion" db:"concentracion" valid:"required"`
	UnidadBase         string     `json:"unidad_base" db:"unidad_base" valid:"-"` // Siempre "UNIDAD"
	ViaAdministracion  *string    `json:"via_administracion" db:"via_administracion" valid:"-"`
	RequiereReceta     bool       `json:"requiere_receta" db:"requiere_receta" valid:"-"`
	Controlado         bool       `json:"controlado" db:"controlado" valid:"-"`
	Descripcion        *string    `json:"descripcion" db:"descripcion" valid:"-"`
	Estado             string     `json:"estado" db:"estado" valid:"required"`
	IsDeleted          bool       `json:"is_deleted" db:"is_deleted"`
	UserDeleted        *string    `json:"user_deleted" db:"user_deleted"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
	UserCreator        string     `json:"user_creator" db:"user_creator"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// MedicineWithStock representa un medicamento con información de stock
type MedicineWithStock struct {
	Medicine
	StockTotal     int `json:"stock_total" db:"stock_total"`
	LotesActivos   int `json:"lotes_activos" db:"lotes_activos"`
}

// NewMedicine crea una nueva instancia de Medicine
func NewMedicine(id, codigo, nombreGenerico string, nombreComercial *string, formaFarmaceutica, concentracion string, viaAdministracion *string, requiereReceta, controlado bool, descripcion *string, estado, userCreator string) *Medicine {
	now := time.Now()
	return &Medicine{
		ID:                id,
		Codigo:            codigo,
		NombreGenerico:    nombreGenerico,
		NombreComercial:   nombreComercial,
		FormaFarmaceutica: formaFarmaceutica,
		Concentracion:     concentracion,
		UnidadBase:        "UNIDAD", // Siempre UNIDAD
		ViaAdministracion: viaAdministracion,
		RequiereReceta:    requiereReceta,
		Controlado:        controlado,
		Descripcion:       descripcion,
		Estado:            estado,
		IsDeleted:         false,
		UserCreator:       userCreator,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// Valid valida la estructura del medicamento
func (m *Medicine) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
