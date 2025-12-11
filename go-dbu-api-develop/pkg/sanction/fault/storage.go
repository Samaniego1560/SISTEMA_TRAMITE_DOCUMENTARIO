package fault

import (
	"github.com/jmoiron/sqlx"
)

const (
	SqlServer = "sqlserver"
)

// ServicesFaultRepository maneja las operaciones sobre la tabla "faltas"
type ServicesFaultRepository interface {
	create(m *Fault) error
	update(m *Fault) error
	delete(id string) error
	getByID(id string) (*Fault, error)

	getAll() ([]*FaultWithStudent, error)
	GetDetalleFalta(faltaID string) ([]*FaultDetalle, error)
	// Nuevos métodos para manejar relaciones con artículos e incisos
	createFaultArticulo(m *FaultArticulo) error
	createFaultInciso(m *FaultInciso) error
	GetServicioNombreByID(servicioID int64) (string, error)
	CreateFaultDocumento(doc *FaultDocumento) error
	GetFaultDocumentoByID(id string) (*FaultDocumento, error)
	// Nuevo método para obtener todos los incisos cometidos por un alumno
	GetAllIncisosByAlumnoID(alumnoID int64) ([]*FaultIncisoDetalle, error)
}

// Struct para devolver el detalle de cada inciso cometido por un alumno
type FaultIncisoDetalle struct {
	FaultID    string `db:"fault_id" json:"fault_id"`
	IncisoID   string `db:"inciso_id" json:"inciso_id"`
	Gravedad   string `db:"gravedad" json:"gravedad"`
	FechaFalta string `db:"fecha_falta" json:"fecha_falta"`
}

// FactoryStorage devuelve una nueva instancia del repositorio
func FactoryStorage(db *sqlx.DB, txID string) ServicesFaultRepository {
	return newFaultSqlServerRepository(db, txID)
}
