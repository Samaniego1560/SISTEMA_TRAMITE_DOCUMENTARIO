package nursing_consultation_vaccine

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesVaccineRepository interface {
	create(m *Vaccine) error
	update(m *Vaccine) error
	delete(id string) error
	deleteByIDConsultation(id string) error
	getByID(id string) ([]*Vaccine, error)
	getByIDPatient(id string) ([]*Vaccine, error)
	getAll() ([]*Vaccine, error)
	getAllTypesVaccines() ([]*VaccineType, error)
	getAllVaccinesByPatientDni(dni string) ([]*Vaccine, error)
	getAllRequiredVaccineTypes() ([]*VaccineType, error)
	getAllNotRequiredVaccineTypes() ([]*VaccineType, error)
	getVaccineIntervalsByVaccineType(tipoVacunaID string) ([]*VaccineInterval, error)
	getNextVaccineIntervalsByVaccineType(tipoVacunaID string, currentDosis int) (*VaccineInterval, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesVaccineRepository {
	return newVaccineSqlServerRepository(db, txID)
}
