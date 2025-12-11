package articles

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServiceArticleRepository interface {
	create(m *Article) error
	update(m *Article) error
	delete(id string) error
	getByID(id string) (*Article, error)
	getAll() ([]*Article, error)
	GetByChapterID(CapituloId int64) ([]*Article, error)
	updateOnlyCharacteristics(m *Article) error
}

func FactoryStorage(db *sqlx.DB, txID string) ServiceArticleRepository {
	return newArticleSqlServerRepository(db, txID)
}
