package estadistica

import (
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServiceEstadisticaRepository interface {
	GetDataChartArea(inicio, fin string) ([]DailyStat, error)
	GetDataChartPie(inicio, fin string) ([]DataPie, error)
	GetDataChartBarra(inicio, fin string) ([]DataBarras, error)
}

func FactoryStorage(db *sqlx.DB) ServiceEstadisticaRepository {
	return NewEstadisticaRepository(db)
}
