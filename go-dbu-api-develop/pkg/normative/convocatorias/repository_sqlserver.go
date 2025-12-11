package repository_sqlserver

import (
	"dbu-api/internal/models"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ConvocatoriaRepository struct {
	DB *sqlx.DB
}

func NewConvocatoriaRepository(db *sqlx.DB) *ConvocatoriaRepository {
	return &ConvocatoriaRepository{DB: db}
}

// ConvocatoriaSemestre representa la estructura de respuesta para el frontend
type ConvocatoriaSemestre struct {
	ID       uint64 `json:"id"`
	Semestre string `json:"semestre"`
}

// GetAllConvocatoriasFiltradas retorna solo los semestres únicos extraídos del nombre, según reglas de negocio, junto con el id de la primera convocatoria encontrada para ese semestre
func (r *ConvocatoriaRepository) GetAllConvocatorias() ([]ConvocatoriaSemestre, error) {
	var convocatorias []*models.Convocatoria
	query := "SELECT * FROM convocatorias"
	err := r.DB.Select(&convocatorias, query)
	if err != nil {
		return nil, err
	}

	semestreMap := make(map[string]ConvocatoriaSemestre)
	re := regexp.MustCompile(`(20[0-9]{2}-[IVX]+)`)
	for _, c := range convocatorias {
		nombre := c.Nombre
		lower := strings.ToLower(nombre)
		if strings.Contains(lower, "bolsa de trabajo") || strings.Contains(lower, "segunda convocatoria") || strings.Contains(lower, "ampliacion de convocatoria") || strings.Contains(lower, "extraordinaria") || strings.Contains(lower, "extraordinario") {
			continue
		}
		match := re.FindStringSubmatch(nombre)
		if len(match) > 1 {
			semestre := "SEMESTRE " + match[1]
			if _, exists := semestreMap[semestre]; !exists {
				semestreMap[semestre] = ConvocatoriaSemestre{
					ID:       c.ID,
					Semestre: semestre,
				}
			}
		}
	}
	// Convertir el map a slice
	var result []ConvocatoriaSemestre
	for _, v := range semestreMap {
		result = append(result, v)
	}
	return result, nil
}
