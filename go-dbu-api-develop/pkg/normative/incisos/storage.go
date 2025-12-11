package incisos

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServiceIncisoRepository interface {
	create(m *Inciso) error
	update(m *Inciso) error
	delete(id string) error
	getByID(id string) (*Inciso, error)
	getAll() ([]*Inciso, error)
	GetByarticleID(articuloId string) ([]*Inciso, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServiceIncisoRepository {
	return newIncisoSqlServerRepository(db, txID)
}
