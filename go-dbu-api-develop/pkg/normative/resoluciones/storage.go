package resolutions

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServiceResolutionRepository interface {
	create(m *Resolution) error
	update(m *Resolution) error
	delete(id string) error
	getByID(id string) (*Resolution, error)
	getAll() ([]*Resolution, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServiceResolutionRepository {
	return newResolutionSqlServerRepository(db, txID)
}
