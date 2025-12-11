package smtp

import (
	"dbu-api/internal/models"
	"strings"
)

// ✅ Interfaz con el nombre correcto del método
type SancionService interface {
	GetSancionesAsignadasPorFalta(faltaID string) ([]*models.SancionAsignadaDetalle, error)
}

// ✅ Actualiza el nombre del método aquí también
func GetSancionAsignadaPorFaltaInternal(faltaID string, service SancionService) (*models.SancionAsignadaDetalle, error) {
	sanciones, err := service.GetSancionesAsignadasPorFalta(faltaID) // ✅ Nombre correcto
	if err != nil {
		return nil, err
	}

	if len(sanciones) == 0 {
		return nil, nil
	}

	return sanciones[0], nil
}

func BuildSancionDetalle(sancion *models.SancionAsignadaDetalle) string {
	cap := strings.TrimSpace(sancion.CapituloSancion)
	art := strings.TrimSpace(sancion.ArticuloSancion)
	inc := strings.TrimSpace(sancion.IncisoSancion)
	detalle := strings.TrimSpace(sancion.DetalleSancion)

	partes := []string{}
	if cap != "" {
		partes = append(partes, cap)
	}
	if art != "" {
		partes = append(partes, art)
	}
	if inc != "" {
		partes = append(partes, "Inc. "+inc)
	}
	if detalle != "" {
		partes = append(partes, detalle)
	}
	if len(partes) > 0 {
		return strings.Join(partes, " ")
	}
	if sancion.SancionID != "" {
		return "Sanción ID: " + sancion.SancionID
	}
	return "Sanción no asignada"
}

func BuildEmailTemplateFromSancion(sancion *models.SancionAsignadaDetalle, base EmailTemplate) EmailTemplate {
	base.SancionDetalle = BuildSancionDetalle(sancion)
	return base
}
