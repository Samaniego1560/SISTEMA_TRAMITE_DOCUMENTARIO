package nursing_consultation_vaccine

import (
	"github.com/asaskevich/govalidator"
	"time"
)

type Vaccine struct {
	ID                   string     `json:"id" db:"id"`
	CreatedAt            *time.Time `json:"created_at" db:"created_at" valid:"required"`
	UpdatedAt            *time.Time `json:"updated_at" db:"updated_at" valid:"required"`
	IsDeleted            bool       `json:"is_deleted" db:"is_deleted"`
	UserDeleted          *string    `json:"user_deleted" db:"user_deleted"`
	DeletedAt            *time.Time `json:"deleted_at" db:"deleted_at"`
	UserCreator          *string    `json:"user_creator" db:"user_creator"`
	IDConsultaEnfermeria string     `json:"consulta_enfermeria_id" db:"consulta_enfermeria_id"`
	TipoVacuna           string     `json:"tipo_vacuna" db:"tipo_vacuna" valid:"-"`
	FechaDosis           string     `json:"fecha_dosis" db:"fecha_dosis" valid:"-"`
	Minsa                bool       `json:"minsa" db:"minsa" valid:"-"`
	TipoAtencion         string     `json:"tipo_atencion" db:"tipo_atencion" valid:"-"`
	Indicaciones         string     `json:"indicaciones" db:"indicaciones" valid:"-"`
	Observaciones        string     `json:"observaciones" db:"observaciones" valid:"-"`
}

type TypesVaccines struct {
	ID            string     `json:"id" db:"id"`
	Nombre        string     `json:"nombre" db:"nombre"`
	Estado        bool       `json:"estado" db:"estado"`
	DuracionMeses string     `json:"duracion_meses" db:"duracion_meses"`
	IsDeleted     bool       `json:"is_deleted" db:"is_deleted"`
	UserDeleted   *string    `json:"user_deleted" db:"user_deleted"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
	UserCreator   *string    `json:"user_creator" db:"user_creator"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at" valid:"required"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at" valid:"required"`
}

func NewVaccine(id string, consulta_enfermeria_id string, tipo_vacuna string, fecha_dosis string, minsa bool, tipoAtencion, indicaciones, observaciones string) *Vaccine {
	now := time.Now()
	return &Vaccine{
		ID:                   id,
		IDConsultaEnfermeria: consulta_enfermeria_id,
		TipoVacuna:           tipo_vacuna,
		FechaDosis:           fecha_dosis,
		Minsa:                minsa,
		TipoAtencion:         tipoAtencion,
		Observaciones:        observaciones,
		Indicaciones:         indicaciones,
		IsDeleted:            false,
		CreatedAt:            &now,
		UpdatedAt:            &now,
	}
}

func (m *Vaccine) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

type VaccineType struct {
	ID                  string    `json:"id" db:"id"`
	Nombre              string    `json:"nombre" db:"nombre"`
	DosisMinimas        int       `json:"dosis_minimas" db:"dosis_minimas"`
	DosisMaximas        *int      `json:"dosis_maximas" db:"dosis_maximas"`
	IntervaloEntreDosis *int      `json:"intervalo_entre_dosis" db:"intervalo_entre_dosis"`
	Estado              bool      `json:"estado" db:"estado"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type VaccineInterval struct {
	ID           string    `json:"id" db:"id"`
	TipoVacunaID string    `json:"tipo_vacuna_id" db:"tipo_vacuna_id"`
	NumeroDosis  int       `json:"numero_dosis" db:"numero_dosis"`
	MesesEspera  int       `json:"meses_espera" db:"meses_espera"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
