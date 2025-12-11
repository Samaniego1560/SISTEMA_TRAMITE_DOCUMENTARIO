package faults_incisos

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesFaultIncisoRepository interface {
	create(m *FaultInciso) error
	delete(id string) error
	getAll() ([]*FaultInciso, error)
}

func NewStorage(db *sqlx.DB, txID string) ServicesFaultIncisoRepository {
	return newFaultIncisoSqlServerRepository(db, txID)
}
