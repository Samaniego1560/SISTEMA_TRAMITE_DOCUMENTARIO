package configuration_sanction_fault

import (
	"context"
	"dbu-api/models"
	internalmodels "dbu-api/internal/models"
)

type Repository interface {
	// ...otros métodos existentes...
	AsignarSancionFalta(ctx context.Context, sfa *internalmodels.SancionFaltaAsignada) error
	RegistrarApelacion(ctx context.Context, ap *models.Apelacion) error
}
