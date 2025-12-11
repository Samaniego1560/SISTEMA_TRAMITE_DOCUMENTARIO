package dbx

import (
	"database/sql"
	"regexp"
	"strings"
)

// GetConvocatoriaPeriodoPorFaltaID recupera el periodo (ej: 2025-II) de la convocatoria asociada a una falta
func GetConvocatoriaPeriodoPorFaltaID(db *sql.DB, faltaID string) (string, error) {
	var convocatoriaID int64
	err := db.QueryRow("SELECT convocatoria_id FROM faltas WHERE id = ?", faltaID).Scan(&convocatoriaID)
	if err != nil {
		return "", err
	}

	var nombre string
	err = db.QueryRow("SELECT nombre FROM convocatorias WHERE id = ?", convocatoriaID).Scan(&nombre)
	if err != nil {
		return "", err
	}

	// Buscar el patrón tipo 2025-II al final del nombre
	re := regexp.MustCompile(`([0-9]{4}-[IVX]+)$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(nombre))
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", nil // No encontrado
}
