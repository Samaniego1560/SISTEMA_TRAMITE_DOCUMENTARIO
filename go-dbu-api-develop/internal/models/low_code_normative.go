package models

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

var (
	ErrInvalidName        = errors.New("invalid name format")
	ErrInvalidDescription = errors.New("invalid description format")
	ErrInvalidServiceID   = errors.New("invalid service ID")
	ErrInvalidFilePath    = errors.New("invalid file path")
	ErrInvalidGravity     = errors.New("invalid gravity value")
)

type Normative struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	ServicioId  string `json:"servicio_id"`
	RutaArchivo string `json:"ruta_archivo"`
	Estado      int    `json:"estado"`
}

type Chapter struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	NormativeId string `json:"resolucion_id"`
}
type DetalleFalta struct {
	FaltaID     string `json:"falta_id"`
	Observacion string `json:"observacion"`
	DNI         string `json:"dni"`
	Nombres     string `json:"nombres"`

	ArticuloID          string `json:"articulo_id"`
	ArticuloDescripcion string `json:"articulo_descripcion"`
	ArticuloGravedad    string `json:"articulo_gravedad"`

	IncisoID          string `json:"inciso_id"`
	IncisoNombre      string `json:"inciso_nombre"`
	IncisoDescripcion string `json:"inciso_descripcion"`

	CapituloID     string `json:"capitulo_id"`
	CapituloNombre string `json:"capitulo_nombre"`

	NormativaID     string `json:"normativa_id"`
	NormativaNombre string `json:"normativa_nombre"`

	DocumentoURL string `json:"documento_url"`
}

type Article struct {
	ID          string `json:"id" valid:"-"`
	Descripcion string `json:"descripcion" valid:"required"`
	Gravedad    string `json:"gravedad" valid:"required"`
	CapituloId  string `json:"capitulo_id" valid:"required"`
}

type Inciso struct {
	ID          string `json:"id" valid:"-"`
	Nombre      string `json:"nombre" valid:"required"`
	Descripcion string `json:"descripcion" valid:"required"`
	ArticuloId  string `json:"articulo_id" valid:"required"`
}

// Valid valida la estructura de Normative
func (m *Normative) Valid() (bool, error) {
	if m == nil {
		return false, errors.New("normative cannot be nil")
	}

	if ok, err := govalidator.ValidateStruct(m); !ok {
		return false, err
	}

	return true, nil
}

// ValidChapter valida la estructura de Chapter
func (m *Chapter) ValidChapter() (bool, error) {
	if m == nil {
		return false, errors.New("chapter cannot be nil")
	}

	if ok, err := govalidator.ValidateStruct(m); !ok {
		return false, err
	}

	return true, nil
}
func (m *Article) ValidArticle() (bool, error) {
	if m == nil {
		return false, errors.New("article cannot be nil")
	}

	if ok, err := govalidator.ValidateStruct(m); !ok {
		return false, err
	}

	// Validación específica para el campo Gravedad
	switch m.Gravedad {
	case "leve", "grave":
		return true, nil
	default:
		return false, ErrInvalidGravity
	}
}

func (m *Inciso) ValidInciso() (bool, error) {
	if m == nil {
		return false, errors.New("inciso cannot be nil")
	}

	if ok, err := govalidator.ValidateStruct(m); !ok {
		return false, err
	}
	return true, nil
}
