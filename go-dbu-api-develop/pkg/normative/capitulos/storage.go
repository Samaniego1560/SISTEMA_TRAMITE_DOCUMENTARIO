package chapters

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServiceChapterRepository interface {
	create(m *Chapter) error
	update(m *Chapter) error
	delete(id string) error
	getByID(id string) (*Chapter, error)
	getAll() ([]*Chapter, error)
	GetByResolutionID(ResolucionId string) ([]*Chapter, error)
	updateOnlyCharacteristics(m *Chapter) error
}

func FactoryStorage(db *sqlx.DB, txID string) ServiceChapterRepository {
	return newChapterSqlServerRepository(db, txID)
}
