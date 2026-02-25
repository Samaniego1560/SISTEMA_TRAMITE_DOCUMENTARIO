package pharmacy_medicines

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

// ServicesMedicineRepository define la interfaz del repositorio
type ServicesMedicineRepository interface {
	create(m *Medicine) error
	update(m *Medicine) error
	delete(id string) error
	getByID(id string) (*Medicine, error)
	getByCode(code string) (*Medicine, error)
	getAll() ([]*Medicine, error)
	getAllWithStock() ([]*MedicineWithStock, error)
	searchMedicines(search, estado string, limit, offset int64) ([]*MedicineWithStock, error)
	countMedicines(search, estado string) (int64, error)
	existsByCode(code string) (bool, error)
}

// FactoryStorage crea una instancia del repositorio
func FactoryStorage(db *sqlx.DB, txID string) ServicesMedicineRepository {
	return newMedicineSqlServerRepository(db, txID)
}
