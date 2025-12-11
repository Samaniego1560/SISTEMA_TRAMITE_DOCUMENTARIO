package sanctions

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesSanctionRepository interface {
	create(m *Sanction) error
	update(m *Sanction) error
	delete(id string) error
	getByID(id string) (*Sanction, error)
	getAll() ([]*Sanction, error)
}

func NewStorage(db *sqlx.DB, txID string) ServicesSanctionRepository {
	return newSanctionSqlServerRepository(db, txID)
}
