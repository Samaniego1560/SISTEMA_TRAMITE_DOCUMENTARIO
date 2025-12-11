package convocatorias

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Convocatorias struct {
	ID            int64      `json:"id" db:"id" valid:"-"`
	FechaInicio   *time.Time `json:"fecha_inicio" db:"fecha_inicio" valid:"optional"`
	FechaFin      *time.Time `json:"fecha_fin" db:"fecha_fin" valid:"optional"`
	Nombre        string     `json:"nombre" db:"nombre" valid:"required"`
	UserId        int64      `json:"user_id" db:"user_id" valid:"required"`
	CreditoMinimo *int       `json:"credito_minimo" db:"credito_minimo" valid:"optional"`
	NotaMinima    *int       `json:"nota_minima" db:"nota_minima" valid:"optional"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}

// ConvocatoriaServicio representa la relación entre convocatoria y servicio
type ConvocatoriaServicio struct {
	ID             int64      `json:"id" db:"id"`
	ConvocatoriaID int64      `json:"convocatoria_id" db:"convocatoria_id"`
	ServicioID     int64      `json:"servicio_id" db:"servicio_id" valid:"required"`
	Cantidad       int        `json:"cantidad" db:"cantidad" valid:"required"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
}

// Seccion representa una sección dentro de una convocatoria
type Seccion struct {
	ID             int64      `json:"id" db:"id"`
	ConvocatoriaID int64      `json:"convocatoria_id" db:"convocatoria_id"`
	Descripcion    string     `json:"descripcion" db:"descripcion" valid:"required"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
}

// Requisito representa un requisito dentro de una sección
type Requisito struct {
	ID              int64      `json:"id" db:"id"`
	SeccionID       int64      `json:"seccion_id" db:"seccion_id"`
	Nombre          string     `json:"nombre" db:"nombre" valid:"required"`
	Descripcion     *string    `json:"descripcion" db:"descripcion"`
	UrlGuia         *string    `json:"url_guia" db:"url_guia" valid:"optional,url"`
	UrlPlantilla    *string    `json:"url_plantilla" db:"url_plantilla" valid:"optional,url"`
	Opciones        *string    `json:"opciones" db:"opciones"`
	TipoRequisitoID int64      `json:"tipo_requisito_id" db:"tipo_requisito_id" valid:"required"`
	UserID          int64      `json:"user_id" db:"user_id"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
}

// CreateConvocatoriaRequest representa el request completo para crear una convocatoria
type CreateConvocatoriaRequest struct {
	FechaInicio          string                    `json:"fecha_inicio" valid:"required"`
	FechaFin             string                    `json:"fecha_fin" valid:"required"`
	Nombre               string                    `json:"nombre" valid:"required,stringlength(1|255)"`
	CreditoMinimo        *int                      `json:"credito_minimo"`
	NotaMinima           *int                      `json:"nota_minima"`
	ConvocatoriaServicio []ConvocatoriaServicioReq `json:"convocatoria_servicio" valid:"required"`
	Secciones            []SeccionReq              `json:"secciones" valid:"required"`
}

// ConvocatoriaServicioReq representa un servicio en el request
type ConvocatoriaServicioReq struct {
	ServicioID int64 `json:"servicio_id" valid:"required"`
	Cantidad   int   `json:"cantidad" valid:"required"`
}

// SeccionReq representa una sección en el request
type SeccionReq struct {
	Descripcion string         `json:"descripcion" valid:"required,stringlength(1|255)"`
	Requisitos  []RequisitoReq `json:"requisitos" valid:"required"`
}

// RequisitoReq representa un requisito en el request
type RequisitoReq struct {
	Nombre          string  `json:"nombre" valid:"required,stringlength(1|255)"`
	Descripcion     *string `json:"descripcion"`
	UrlGuia         *string `json:"url_guia" valid:"optional,url"`
	UrlPlantilla    *string `json:"url_plantilla" valid:"optional,url"`
	Opciones        *string `json:"opciones"`
	TipoRequisitoID int64   `json:"tipo_requisito_id" valid:"required"`
}

// ConvocatoriaResponse representa la respuesta completa de una convocatoria
type ConvocatoriaResponse struct {
	ID                   int64                  `json:"id"`
	FechaInicio          *time.Time             `json:"fecha_inicio"`
	FechaFin             *time.Time             `json:"fecha_fin"`
	Nombre               string                 `json:"nombre"`
	UserID               int64                  `json:"user_id"`
	CreditoMinimo        *int                   `json:"credito_minimo"`
	NotaMinima           *int                   `json:"nota_minima"`
	ConvocatoriaServicio []ConvocatoriaServicio `json:"convocatoria_servicio,omitempty"`
	Secciones            []SeccionResponse      `json:"secciones,omitempty"`
	CreatedAt            *time.Time             `json:"created_at"`
	UpdatedAt            *time.Time             `json:"updated_at"`
}

// SeccionResponse representa una sección en la respuesta
type SeccionResponse struct {
	ID             int64       `json:"id"`
	ConvocatoriaID int64       `json:"convocatoria_id"`
	Descripcion    string      `json:"descripcion"`
	Requisitos     []Requisito `json:"requisitos,omitempty"`
	CreatedAt      *time.Time  `json:"created_at"`
	UpdatedAt      *time.Time  `json:"updated_at"`
}

func NewSubmissions(id int64, fechaInicio *time.Time, fechaFin *time.Time, nombre string, userId int64, creditoMinimo *int, notaMinima *int) *Convocatorias {
	return &Convocatorias{
		ID:            id,
		FechaInicio:   fechaInicio,
		FechaFin:      fechaFin,
		Nombre:        nombre,
		UserId:        userId,
		CreditoMinimo: creditoMinimo,
		NotaMinima:    notaMinima,
	}
}

func NewCreateSubmissions(fechaInicio *time.Time, fechaFin *time.Time, nombre string, userId int64, creditoMinimo *int, notaMinima *int) *Convocatorias {
	return &Convocatorias{
		FechaInicio:   fechaInicio,
		FechaFin:      fechaFin,
		Nombre:        nombre,
		UserId:        userId,
		CreditoMinimo: creditoMinimo,
		NotaMinima:    notaMinima,
	}
}

func (m *Convocatorias) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Valid valida el request de creación
func (r *CreateConvocatoriaRequest) Valid() (bool, error) {
	return govalidator.ValidateStruct(r)
}
