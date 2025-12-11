package convocatorias

import (
	"time"

	"github.com/jmoiron/sqlx"

	"dbu-api/internal/models"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServicesConvocatoriasRepository interface {
	create(m *Convocatorias) error
	update(m *Convocatorias) error
	delete(id int64) error
	getByID(id int64) (*Convocatorias, error)
	getAll() ([]*Convocatorias, error)
	getAllByService(id int64) ([]*Convocatorias, error)
	getActive() (*Convocatorias, error)
	getLast() (*Convocatorias, error)
	// Métodos para relaciones
	createConvocatoriaServicio(cs *ConvocatoriaServicio) error
	createSeccion(sec *Seccion) error
	createRequisito(req *Requisito) error
	getConvocatoriaServiciosByConvocatoriaID(convocatoriaID int64) ([]ConvocatoriaServicio, error)
	getSeccionesByConvocatoriaID(convocatoriaID int64) ([]Seccion, error)
	getRequisitosBySeccionID(seccionID int64) ([]Requisito, error)
	deleteConvocatoriaServicios(convocatoriaID int64) error
	deleteRequisitosBySeccionID(seccionID int64) error
	deleteSecciones(convocatoriaID int64) error
	checkOverlappingConvocatorias(fechaInicio time.Time, excludeID *int64) (bool, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesConvocatoriasRepository {
	return newConvocatoriasSqlServerRepository(db, user, txID)
}
