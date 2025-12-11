package dbx

import (
	"fmt"
)

// GetResolucionNombreVigente obtiene el nombre de la resolución vigente para un capitulo dado
func GetResolucionNombreVigente(capituloID string) (string, error) {
	db := GetConnection()
	var resolucionID string
	err := db.Get(&resolucionID, "SELECT resolucion_id FROM capitulos WHERE id = ?", capituloID)
	if err != nil || resolucionID == "" {
		return "", fmt.Errorf("No se encontró resolucion_id para el capitulo %s", capituloID)
	}
	var nombre string
	err = db.Get(&nombre, "SELECT nombre FROM resoluciones WHERE id = ? AND estado = 1", resolucionID)
	if err != nil {
		return "", fmt.Errorf("No se encontró resolución vigente para el capitulo %s", capituloID)
	}
	return nombre, nil
}
