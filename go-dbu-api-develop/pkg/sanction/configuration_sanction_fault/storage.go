package configuration_sanction_fault

import (
	"context"
	"dbu-api/models"
	internalmodels "dbu-api/internal/models"
	"github.com/jmoiron/sqlx"
)

// Constante para identificar motor, útil si alguna vez cambias de DB.
const (
	SqlServer = "sqlserver"
)

// SancionRepository expone las operaciones sobre la tabla "sancionesAFaltas"
type SancionRepository interface {
	Create(m *Sancion) error
	Update(m *Sancion) error
	Delete(id string) error
	GetByID(id string) (*Sancion, error)
	GetAll() ([]*Sancion, error)
	AsignarSancionFalta(ctx context.Context, sfa *internalmodels.SancionFaltaAsignada) error
	RegistrarApelacion(ctx context.Context, ap *models.Apelacion) error
	GetSancionesAsignadasPorFalta(faltaID string) ([]*internalmodels.SancionAsignadaDetalle, error)
}

// FactoryStorage retorna la instancia concreta del repositorio
func FactoryStorage(db *sqlx.DB, txID string) SancionRepository {
	return NewSancionSqlServerRepository(db, txID)
}
