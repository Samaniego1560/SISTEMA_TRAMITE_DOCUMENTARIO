package faults_articulos

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesFaultArticuloRepository interface {
	create(m *FaultArticulo) error
	delete(id string) error
	getAll() ([]*FaultArticulo, error)
}

func NewStorage(db *sqlx.DB, txID string) ServicesFaultArticuloRepository {
	return newFaultArticuloSqlServerRepository(db, txID)
}
