package estadistica

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
)

type EstadisticaRepository struct {
	DB *sqlx.DB
}

func NewEstadisticaRepository(db *sqlx.DB) *EstadisticaRepository {
	return &EstadisticaRepository{DB: db}
}

func (r *EstadisticaRepository) GetDataChartArea(inicio, fin string) ([]DailyStat, error) {

	if inicio == "" || fin == "" {
		return nil, errors.New("fecha de inicio y fin son requeridas")
	}

	query := `
		WITH RECURSIVE fechas AS (
		SELECT DATE(?) AS fecha
		UNION ALL
		SELECT DATE_ADD(fecha, INTERVAL 1 DAY)
		FROM fechas
		WHERE fecha < DATE(?)
	)
	SELECT 
		DATE_FORMAT(f.fecha, '%Y-%m-%d') AS date,
		COUNT(CASE WHEN h.es_srq = 1 THEN 1 END) AS countSrq,
		COUNT(CASE WHEN h.es_srq = 0 THEN 1 END) AS count
	FROM fechas f
	LEFT JOIN historial h ON DATE(h.fecha_respuesta) = f.fecha
	GROUP BY f.fecha
	ORDER BY f.fecha;
	`

	var stats []DailyStat
	err := r.DB.Select(&stats, query, inicio, fin)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *EstadisticaRepository) GetDataChartPie(inicio, fin string) ([]DataPie, error) {
	if inicio == "" || fin == "" {
		return nil, errors.New("fecha de inicio y fin son requeridas")
	}

	query := `
	SELECT 
		h.estado_evaluacion as estados_evaluacion,
		COUNT(*) AS total
	FROM historial h
	WHERE h.fecha_respuesta BETWEEN ? AND ?
	GROUP BY h.estado_evaluacion
	ORDER BY total DESC;
	`

	var stats []DataPie
	err := r.DB.Select(&stats, query, inicio, fin)
	if err != nil {
		log.Printf("error consultando estados evaluacion: %v", err)
		return nil, errors.New("error al consultar estadísticas de evaluación")
	}

	return stats, nil
}
func (r *EstadisticaRepository) GetDataChartBarra(inicio, fin string) ([]DataBarras, error) {

	if inicio == "" || fin == "" {
		return nil, errors.New("fecha de inicio y fin son requeridas")
	}

	query := `
	SELECT 
		e.nombre AS escuela,
		COUNT(p.id_participante) AS total
	FROM escuelas_profesionales e
	LEFT JOIN participantes p 
		ON p.escuela = e.nombre
		AND p.created_at BETWEEN ? AND ?
	WHERE e.is_deleted = 0
	GROUP BY e.nombre
	ORDER BY e.nombre ASC;
	`

	var stats []DataBarras
	err := r.DB.Select(&stats, query, inicio, fin)
	if err != nil {
		log.Printf("error consultando participantes por escuelas activas: %v", err)
		return nil, errors.New("error al consultar estadísticas por escuelas activas")
	}

	return stats, nil
}
