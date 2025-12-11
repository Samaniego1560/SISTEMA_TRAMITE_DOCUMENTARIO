package normative

import (
	"dbu-api/internal/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerReportStats struct {
	db *sqlx.DB
}

func NewHandlerReportStats(db *sqlx.DB) *HandlerReportStats {
	return &HandlerReportStats{db: db}
}

// GetReportStats godoc
// @Summary Obtiene estadísticas globales para reportes gráficos
// @Description Devuelve conteos de faltas graves, leves, por sexo, por escuela profesional, etc.
// @Tags Reportes
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Router /v1/normative/report-stats [get]
func (h *HandlerReportStats) GetReportStats(c *fiber.Ctx) error {
	// Llamar a processReportStats sin filtro (convocatoriaID = 0)
	return h.processReportStats(c, 0)
}

// GetReportStatsBySemestre godoc
// @Summary Obtiene estadísticas filtradas por ID de convocatoria (semestre académico)
// @Description Devuelve las mismas estadísticas pero filtradas por una convocatoria específica
// @Tags Reportes
// @Accept json
// @Produce json
// @Param convocatoria_id path int true "ID de la convocatoria"
// @Success 200 {object} models.Response
// @Router /v1/normative/report-stats/:convocatoria_id [get]
func (h *HandlerReportStats) GetReportStatsBySemestre(c *fiber.Ctx) error {
	convocatoriaIDStr := c.Params("convocatoria_id")

	if convocatoriaIDStr == "" {
		return c.Status(400).JSON(models.Response{
			Error: true,
			Data:  "Parametro convocatoria_id es requerido",
		})
	}

	convocatoriaID, err := strconv.ParseInt(convocatoriaIDStr, 10, 64)
	if err != nil || convocatoriaID <= 0 {
		return c.Status(400).JSON(models.Response{
			Error: true,
			Data:  "Parametro convocatoria_id debe ser numerico y mayor a 0",
		})
	}

	// Validar que la convocatoria existe
	var exists int
	queryExists := "SELECT COUNT(*) FROM convocatorias WHERE id = ?"
	if err := h.db.Get(&exists, queryExists, convocatoriaID); err != nil || exists == 0 {
		return c.Status(404).JSON(models.Response{
			Error: true,
			Data:  "Convocatoria no encontrada",
		})
	}

	return h.processReportStats(c, convocatoriaID)
}

// GetComparacionSemestres godoc
// @Summary Compara estadísticas entre múltiples convocatorias
// @Description Devuelve estadísticas comparativas de las convocatorias especificadas. Si no se especifican, devuelve todas.
// @Tags Reportes
// @Accept json
// @Produce json
// @Param convocatorias query string false "IDs de convocatorias separados por comas (ej: 10,15,20)"
// @Success 200 {object} models.Response
// @Router /v1/normative/report-stats/comparacion [get]
func (h *HandlerReportStats) GetComparacionSemestres(c *fiber.Ctx) error {
	convocatoriasStr := c.Query("convocatorias", "")

	var convocatoriaIDs []int64

	// Si no se especifican convocatorias, obtener todas
	if convocatoriasStr == "" {
		query := "SELECT id FROM convocatorias ORDER BY id DESC"
		if err := h.db.Select(&convocatoriaIDs, query); err != nil {
			return c.Status(500).JSON(models.Response{
				Error: true,
				Data:  "Error al obtener convocatorias",
			})
		}
	} else {
		// Parsear IDs separados por comas
		idsStr := strings.Split(convocatoriasStr, ",")
		for _, idStr := range idsStr {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil || id <= 0 {
				return c.Status(400).JSON(models.Response{
					Error: true,
					Data:  "IDs de convocatorias inválidos",
				})
			}
			convocatoriaIDs = append(convocatoriaIDs, id)
		}
	}

	// Obtener estadísticas básicas por cada convocatoria
	type ConvocatoriaStats struct {
		ConvocatoriaID     int64  `db:"convocatoria_id" json:"convocatoria_id"`
		ConvocatoriaNombre string `db:"convocatoria_nombre" json:"convocatoria_nombre"`
		TotalFaltas        int    `db:"total_faltas" json:"total_faltas"`
		FaltasGraves       int    `db:"faltas_graves" json:"faltas_graves"`
		FaltasLeves        int    `db:"faltas_leves" json:"faltas_leves"`
		TotalAlumnos       int    `db:"total_alumnos" json:"total_alumnos"`
		Reincidentes       int    `db:"reincidentes" json:"reincidentes"`
		Apeladas           int    `db:"apeladas" json:"apeladas"`
		Resueltas          int    `db:"resueltas" json:"resueltas"`
	}

	var resultados []ConvocatoriaStats

	for _, convID := range convocatoriaIDs {
		var stat ConvocatoriaStats
		stat.ConvocatoriaID = convID

		// Nombre de convocatoria
		if err := h.db.Get(&stat.ConvocatoriaNombre, "SELECT nombre FROM convocatorias WHERE id = ?", convID); err != nil {
			continue
		}

		lowerNombre := strings.ToLower(stat.ConvocatoriaNombre)
		if strings.Contains(lowerNombre, "bolsa de trabajos") ||
			strings.Contains(lowerNombre, "bolsa de trabajo") ||
			strings.Contains(lowerNombre, "bolsa de trabaj") ||
			strings.Contains(lowerNombre, "segunda convocatoria") ||
			strings.Contains(lowerNombre, "ampliacion de convocatoria") ||
			strings.Contains(lowerNombre, "extraordinaria") ||
			strings.Contains(lowerNombre, "extraordinario") {
			continue
		}

		// Total de faltas
		h.db.Get(&stat.TotalFaltas, "SELECT COUNT(*) FROM faltas WHERE convocatoria_id = ?", convID)

		// Faltas graves
		h.db.Get(&stat.FaltasGraves, `
			SELECT COUNT(*) FROM faltas f 
			WHERE f.convocatoria_id = ? AND EXISTS (
				SELECT 1 FROM articulos ar
				JOIN faltas_articulos fa ON fa.articulo_id = ar.id
				WHERE fa.falta_id = f.id AND ar.gravedad = 'grave'
			)`, convID)

		// Faltas leves
		h.db.Get(&stat.FaltasLeves, `
			SELECT COUNT(*) FROM faltas f 
			WHERE f.convocatoria_id = ? AND EXISTS (
				SELECT 1 FROM articulos ar
				JOIN faltas_articulos fa ON fa.articulo_id = ar.id
				WHERE fa.falta_id = f.id AND ar.gravedad = 'leve'
			)`, convID)

		// Total de alumnos con faltas
		h.db.Get(&stat.TotalAlumnos, "SELECT COUNT(DISTINCT alumno_id) FROM faltas WHERE convocatoria_id = ?", convID)

		// Reincidentes
		h.db.Get(&stat.Reincidentes, `
			SELECT COUNT(*) FROM (
				SELECT alumno_id FROM faltas 
				WHERE convocatoria_id = ? 
				GROUP BY alumno_id 
				HAVING COUNT(*) > 1
			) as sub`, convID)

		// Apeladas
		h.db.Get(&stat.Apeladas, `
			SELECT COUNT(*) FROM faltas 
			WHERE convocatoria_id = ? 
			AND apelacion_documento IS NOT NULL 
			AND apelacion_documento != ''`, convID)

		// Resueltas
		h.db.Get(&stat.Resueltas, "SELECT COUNT(*) FROM faltas WHERE convocatoria_id = ? AND estado IN ('resuelta','apelada')", convID)

		resultados = append(resultados, stat)
	}

	// Calcular tendencias entre semestres consecutivos
	type TendenciaComparativa struct {
		ConvocatoriaActual   int64   `json:"convocatoria_actual"`
		ConvocatoriaAnterior int64   `json:"convocatoria_anterior"`
		DiferenciaFaltas     int     `json:"diferencia_faltas"`
		DiferenciaAlumnos    int     `json:"diferencia_alumnos"`
		PorcentajeCambio     float64 `json:"porcentaje_cambio"`
	}

	var tendencias []TendenciaComparativa

	for i := 1; i < len(resultados); i++ {
		actual := resultados[i]
		anterior := resultados[i-1]

		tendencia := TendenciaComparativa{
			ConvocatoriaActual:   actual.ConvocatoriaID,
			ConvocatoriaAnterior: anterior.ConvocatoriaID,
			DiferenciaFaltas:     actual.TotalFaltas - anterior.TotalFaltas,
			DiferenciaAlumnos:    actual.TotalAlumnos - anterior.TotalAlumnos,
		}

		if anterior.TotalFaltas > 0 {
			tendencia.PorcentajeCambio = float64(actual.TotalFaltas-anterior.TotalFaltas) / float64(anterior.TotalFaltas) * 100
		}

		tendencias = append(tendencias, tendencia)
	}

	return c.JSON(models.Response{
		Error: false,
		Data: map[string]interface{}{
			"convocatorias":       resultados,
			"tendencias":          tendencias,
			"total_convocatorias": len(resultados),
		},
	})
}

// GetSemestreConAnterior godoc
// @Summary Obtiene estadísticas de un semestre comparado con el anterior
// @Description Devuelve estadísticas completas del semestre actual y comparación con el anterior
// @Tags Reportes
// @Accept json
// @Produce json
// @Param convocatoria_id path int true "ID de la convocatoria actual"
// @Success 200 {object} models.Response
// @Router /v1/normative/report-stats/:convocatoria_id/comparar [get]
func (h *HandlerReportStats) GetSemestreConAnterior(c *fiber.Ctx) error {
	convocatoriaIDStr := c.Params("convocatoria_id")

	if convocatoriaIDStr == "" {
		return c.Status(400).JSON(models.Response{
			Error: true,
			Data:  "Parametro convocatoria_id es requerido",
		})
	}

	convocatoriaID, err := strconv.ParseInt(convocatoriaIDStr, 10, 64)
	if err != nil || convocatoriaID <= 0 {
		return c.Status(400).JSON(models.Response{
			Error: true,
			Data:  "Parametro convocatoria_id debe ser numerico y mayor a 0",
		})
	}

	// Validar que la convocatoria existe
	var exists int
	queryExists := "SELECT COUNT(*) FROM convocatorias WHERE id = ?"
	if err := h.db.Get(&exists, queryExists, convocatoriaID); err != nil || exists == 0 {
		return c.Status(404).JSON(models.Response{
			Error: true,
			Data:  "Convocatoria no encontrada",
		})
	}
	// ========== OBTENER NOMBRES DE CONVOCATORIAS ==========
	var nombreActual, nombreAnterior string
	// Obtener nombre de la convocatoria actual
	queryNombreActual := "SELECT nombre FROM convocatorias WHERE id = ?"
	if err := h.db.Get(&nombreActual, queryNombreActual, convocatoriaID); err != nil {
		nombreActual = fmt.Sprintf("Semestre %d", convocatoriaID)
	}
	// Obtener el ID de la convocatoria anterior
	var convocatoriaAnteriorID int64
	queryAnteriores := `
	SELECT id, nombre 
	FROM convocatorias 
	WHERE id < ? 
	ORDER BY id DESC
`
	rows, err := h.db.Queryx(queryAnteriores, convocatoriaID)
	tieneAnterior := false
	if err == nil && rows != nil {
		defer rows.Close()

		// Buscar la primera convocatoria válida (que no sea bolsa, segunda, etc.)
		for rows.Next() {
			var tempID int64
			var tempNombre string

			if err := rows.Scan(&tempID, &tempNombre); err == nil {
				// Filtrar convocatorias no deseadas
				lowerNombre := strings.ToLower(tempNombre)
				if strings.Contains(lowerNombre, "bolsa de trabajos") ||
					strings.Contains(lowerNombre, "bolsa de trabajo") ||
					strings.Contains(lowerNombre, "segunda convocatoria") ||
					strings.Contains(lowerNombre, "ampliacion de convocatoria") ||
					strings.Contains(lowerNombre, "extraordinaria") ||
					strings.Contains(lowerNombre, "extraordinario") {
					continue // Saltar esta convocatoria
				}

				// Encontramos una convocatoria válida
				convocatoriaAnteriorID = tempID
				nombreAnterior = tempNombre
				tieneAnterior = true
				break
			}
		}
	}
	nombreActual = formatearNombreSemestre(nombreActual)
	if tieneAnterior {
		nombreAnterior = formatearNombreSemestre(nombreAnterior)
	}
	// Obtener estadísticas del semestre actual
	type StatsSimples struct {
		TotalFaltas  int `db:"total_faltas"`
		FaltasGraves int `db:"faltas_graves"`
		FaltasLeves  int `db:"faltas_leves"`
		TotalAlumnos int `db:"total_alumnos"`
		Varones      int `db:"varones"`
		Mujeres      int `db:"mujeres"`
		Reincidentes int `db:"reincidentes"`
		Apeladas     int `db:"apeladas"`
		Resueltas    int `db:"resueltas"`
	}

	var statsActual, statsAnterior StatsSimples

	// Stats del semestre actual
	h.db.Get(&statsActual.TotalFaltas, "SELECT COUNT(*) FROM faltas WHERE convocatoria_id = ?", convocatoriaID)
	h.db.Get(&statsActual.FaltasGraves, `
		SELECT COUNT(*) FROM faltas f 
		WHERE f.convocatoria_id = ? AND EXISTS (
			SELECT 1 FROM articulos ar
			JOIN faltas_articulos fa ON fa.articulo_id = ar.id
			WHERE fa.falta_id = f.id AND ar.gravedad = 'grave'
		)`, convocatoriaID)
	h.db.Get(&statsActual.FaltasLeves, `
		SELECT COUNT(*) FROM faltas f 
		WHERE f.convocatoria_id = ? AND EXISTS (
			SELECT 1 FROM articulos ar
			JOIN faltas_articulos fa ON fa.articulo_id = ar.id
			WHERE fa.falta_id = f.id AND ar.gravedad = 'leve'
		)`, convocatoriaID)
	h.db.Get(&statsActual.TotalAlumnos, "SELECT COUNT(DISTINCT alumno_id) FROM faltas WHERE convocatoria_id = ?", convocatoriaID)
	h.db.Get(&statsActual.Varones, `
		SELECT COUNT(DISTINCT a.id) FROM alumnos a 
		WHERE a.sexo = 'M' AND a.id IN (SELECT alumno_id FROM faltas WHERE convocatoria_id = ?)`, convocatoriaID)
	h.db.Get(&statsActual.Mujeres, `
		SELECT COUNT(DISTINCT a.id) FROM alumnos a 
		WHERE a.sexo = 'F' AND a.id IN (SELECT alumno_id FROM faltas WHERE convocatoria_id = ?)`, convocatoriaID)
	h.db.Get(&statsActual.Reincidentes, `
		SELECT COUNT(*) FROM (
			SELECT alumno_id FROM faltas 
			WHERE convocatoria_id = ? 
			GROUP BY alumno_id 
			HAVING COUNT(*) > 1
		) as sub`, convocatoriaID)
	h.db.Get(&statsActual.Apeladas, `
		SELECT COUNT(*) FROM faltas 
		WHERE convocatoria_id = ? 
		AND apelacion_documento IS NOT NULL 
		AND apelacion_documento != ''`, convocatoriaID)
	h.db.Get(&statsActual.Resueltas, "SELECT COUNT(*) FROM faltas WHERE convocatoria_id = ? AND estado IN ('resuelta','apelada')", convocatoriaID)

	// Stats del semestre anterior (si existe)
	if tieneAnterior {
		h.db.Get(&statsAnterior.TotalFaltas, "SELECT COUNT(*) FROM faltas WHERE convocatoria_id = ?", convocatoriaAnteriorID)
		h.db.Get(&statsAnterior.FaltasGraves, `
			SELECT COUNT(*) FROM faltas f 
			WHERE f.convocatoria_id = ? AND EXISTS (
				SELECT 1 FROM articulos ar
				JOIN faltas_articulos fa ON fa.articulo_id = ar.id
				WHERE fa.falta_id = f.id AND ar.gravedad = 'grave'
			)`, convocatoriaAnteriorID)
		h.db.Get(&statsAnterior.FaltasLeves, `
			SELECT COUNT(*) FROM faltas f 
			WHERE f.convocatoria_id = ? AND EXISTS (
				SELECT 1 FROM articulos ar
				JOIN faltas_articulos fa ON fa.articulo_id = ar.id
				WHERE fa.falta_id = f.id AND ar.gravedad = 'leve'
			)`, convocatoriaAnteriorID)
		h.db.Get(&statsAnterior.TotalAlumnos, "SELECT COUNT(DISTINCT alumno_id) FROM faltas WHERE convocatoria_id = ?", convocatoriaAnteriorID)
		h.db.Get(&statsAnterior.Varones, `
			SELECT COUNT(DISTINCT a.id) FROM alumnos a 
			WHERE a.sexo = 'M' AND a.id IN (SELECT alumno_id FROM faltas WHERE convocatoria_id = ?)`, convocatoriaAnteriorID)
		h.db.Get(&statsAnterior.Mujeres, `
			SELECT COUNT(DISTINCT a.id) FROM alumnos a 
			WHERE a.sexo = 'F' AND a.id IN (SELECT alumno_id FROM faltas WHERE convocatoria_id = ?)`, convocatoriaAnteriorID)
		h.db.Get(&statsAnterior.Reincidentes, `
			SELECT COUNT(*) FROM (
				SELECT alumno_id FROM faltas 
				WHERE convocatoria_id = ? 
				GROUP BY alumno_id 
				HAVING COUNT(*) > 1
			) as sub`, convocatoriaAnteriorID)
		h.db.Get(&statsAnterior.Apeladas, `
			SELECT COUNT(*) FROM faltas 
			WHERE convocatoria_id = ? 
			AND apelacion_documento IS NOT NULL 
			AND apelacion_documento != ''`, convocatoriaAnteriorID)
		h.db.Get(&statsAnterior.Resueltas, "SELECT COUNT(*) FROM faltas WHERE convocatoria_id = ? AND estado IN ('resuelta','apelada')", convocatoriaAnteriorID)
	}

	// Calcular diferencias y tendencias
	calcularTendencia := func(actual, anterior int) map[string]interface{} {
		diferencia := actual - anterior
		porcentaje := 0.0

		// Si anterior es 0 pero hay valores actuales, considerarlo como 100% de incremento
		if anterior == 0 && actual > 0 {
			porcentaje = 100.0
		} else if anterior > 0 {
			// Cálculo normal del porcentaje
			porcentaje = (float64(diferencia) / float64(anterior)) * 100.0
		}
		// Si ambos son 0, porcentaje queda en 0
		return map[string]interface{}{
			"valor_actual":      actual,
			"valor_anterior":    anterior,
			"diferencia":        diferencia,
			"porcentaje_cambio": porcentaje,
			"tendencia":         getTendenciaTexto(diferencia),
		}
	}

	response := map[string]interface{}{
		"convocatoria_actual_id":     convocatoriaID,
		"convocatoria_actual_nombre": nombreActual, // NUEVO
		"tiene_anterior":             tieneAnterior,
	}

	if tieneAnterior {
		response["convocatoria_anterior_id"] = convocatoriaAnteriorID
		response["convocatoria_anterior_nombre"] = nombreAnterior
		response["comparacion"] = map[string]interface{}{
			"total_faltas":  calcularTendencia(statsActual.TotalFaltas, statsAnterior.TotalFaltas),
			"faltas_graves": calcularTendencia(statsActual.FaltasGraves, statsAnterior.FaltasGraves),
			"faltas_leves":  calcularTendencia(statsActual.FaltasLeves, statsAnterior.FaltasLeves),
			"total_alumnos": calcularTendencia(statsActual.TotalAlumnos, statsAnterior.TotalAlumnos),
			"varones":       calcularTendencia(statsActual.Varones, statsAnterior.Varones),
			"mujeres":       calcularTendencia(statsActual.Mujeres, statsAnterior.Mujeres),
			"reincidentes":  calcularTendencia(statsActual.Reincidentes, statsAnterior.Reincidentes),
			"apeladas":      calcularTendencia(statsActual.Apeladas, statsAnterior.Apeladas),
			"resueltas":     calcularTendencia(statsActual.Resueltas, statsAnterior.Resueltas),
		}
	} else {
		response["stats_actuales"] = statsActual
		response["mensaje"] = "No hay semestre anterior para comparar"
	}

	return c.JSON(models.Response{
		Error: false,
		Data:  response,
	})
}
func formatearNombreSemestre(nombre string) string {
	// Si ya tiene "SEMESTRE", retornar tal cual
	upperNombre := strings.ToUpper(nombre)
	if strings.Contains(upperNombre, "SEMESTRE") {
		return nombre
	}

	// Extraer año y periodo (ej: 2025-II)
	re := regexp.MustCompile(`(\d{4}[-\s]+(I{1,3}|II?))`)
	match := re.FindString(nombre)
	if match != "" {
		return fmt.Sprintf("SEMESTRE %s", strings.ToUpper(match))
	}

	return nombre
}
func getTendenciaTexto(diferencia int) string {
	if diferencia > 0 {
		return "incremento"
	} else if diferencia < 0 {
		return "reduccion"
	}
	return "sin_cambios"
}

// processReportStats procesa las estadísticas con o sin filtro de convocatoria
// Si convocatoriaID es 0, se obtienen estadísticas globales
// Si convocatoriaID > 0, se filtran por esa convocatoria específica
func (h *HandlerReportStats) processReportStats(c *fiber.Ctx, convocatoriaID int64) error {
	// Construir la cláusula WHERE y parámetros según si hay filtro o no
	whereClause := ""
	if convocatoriaID > 0 {
		whereClause = fmt.Sprintf(" AND f.convocatoria_id = %d", convocatoriaID)
	}

	// Subquery de alumnos con faltas (con o sin filtro)
	subqueryAlumnos := fmt.Sprintf("(SELECT DISTINCT alumno_id FROM faltas f WHERE 1=1%s)", whereClause)

	res := models.Response{Error: true}

	// ==================== ESTADÍSTICAS DE FALTAS ====================

	// 1. Conteo de faltas graves y leves
	var faltasGraves, faltasLeves int

	queryGraves := fmt.Sprintf(`
		SELECT COUNT(*) FROM faltas f 
		WHERE EXISTS (
			SELECT 1 FROM articulos ar
			JOIN faltas_articulos fa ON fa.articulo_id = ar.id
			WHERE fa.falta_id = f.id AND ar.gravedad = 'grave'
		)%s`, whereClause)
	h.db.Get(&faltasGraves, queryGraves)

	queryLeves := fmt.Sprintf(`
		SELECT COUNT(*) FROM faltas f 
		WHERE EXISTS (
			SELECT 1 FROM articulos ar
			JOIN faltas_articulos fa ON fa.articulo_id = ar.id
			WHERE fa.falta_id = f.id AND ar.gravedad = 'leve'
		)%s`, whereClause)
	h.db.Get(&faltasLeves, queryLeves)

	// 2. Porcentaje de sanciones apeladas vs. no apeladas
	var totalSanciones, sancionesApeladas int
	h.db.Get(&totalSanciones, fmt.Sprintf("SELECT COUNT(*) FROM faltas f WHERE 1=1%s", whereClause))
	h.db.Get(&sancionesApeladas, fmt.Sprintf(
		"SELECT COUNT(*) FROM faltas f WHERE (apelacion_documento IS NOT NULL AND apelacion_documento != '')%s",
		whereClause))

	porcentajeApeladas := 0.0
	if totalSanciones > 0 {
		porcentajeApeladas = float64(sancionesApeladas) / float64(totalSanciones) * 100
	}

	// 3. Ranking de reincidencia (top 10 alumnos con más sanciones)
	rankingReincidencia := []map[string]interface{}{}
	queryRanking := fmt.Sprintf(`
		SELECT alumno_id, COUNT(*) as total 
		FROM faltas f 
		WHERE 1=1%s 
		GROUP BY alumno_id 
		HAVING total > 1 
		ORDER BY total DESC 
		LIMIT 10`, whereClause)

	rowsRanking, err := h.db.Queryx(queryRanking)
	if err == nil && rowsRanking != nil {
		defer rowsRanking.Close()
		for rowsRanking.Next() {
			var alumnoID int64
			var total int
			if err := rowsRanking.Scan(&alumnoID, &total); err == nil {
				rankingReincidencia = append(rankingReincidencia, map[string]interface{}{
					"alumno_id": alumnoID,
					"total":     total,
				})
			}
		}
	}

	// 4. Conteo de sanciones por periodo académico
	sancionesPorPeriodo := map[string]int{}
	queryPeriodo := fmt.Sprintf(`
		SELECT 
			SUBSTRING_INDEX(SUBSTRING_INDEX(c.nombre, ' ', -1), ' ', 1) AS periodo_convocatoria,
			COUNT(f.id) as total
		FROM faltas f
		INNER JOIN convocatorias c ON f.convocatoria_id = c.id
		WHERE 1=1%s
		GROUP BY periodo_convocatoria`, whereClause)

	rowsPeriodo, err := h.db.Queryx(queryPeriodo)
	if err == nil && rowsPeriodo != nil {
		defer rowsPeriodo.Close()
		for rowsPeriodo.Next() {
			var periodo string
			var total int
			if err := rowsPeriodo.Scan(&periodo, &total); err == nil {
				sancionesPorPeriodo[periodo] = total
			}
		}
	}

	// 5. Conteo de faltas por servicio
	faltasPorServicio := map[string]int{}
	queryServicio := fmt.Sprintf(`
		SELECT s.nombre, COUNT(f.id) as total 
		FROM faltas f 
		INNER JOIN servicios s ON f.servicio_id = s.id 
		WHERE 1=1%s 
		GROUP BY s.nombre`, whereClause)

	rowsServicio, err := h.db.Queryx(queryServicio)
	if err == nil && rowsServicio != nil {
		defer rowsServicio.Close()
		for rowsServicio.Next() {
			var nombre string
			var total int
			if err := rowsServicio.Scan(&nombre, &total); err == nil {
				faltasPorServicio[nombre] = total
			}
		}
	}

	// 6. Tipo de sanción aplicada
	tiposSancion := map[string]int{}
	querySancion := fmt.Sprintf(`
		SELECT sad.detalle_sancion, COUNT(*) as total 
		FROM sancion_asignada_detalle sad
		INNER JOIN sancion_asignada sa ON sad.sancion_asignada_id = sa.id
		INNER JOIN faltas f ON sa.falta_id = f.id
		WHERE 1=1%s
		GROUP BY sad.detalle_sancion`, whereClause)

	rowsSancion, err := h.db.Queryx(querySancion)
	if err == nil && rowsSancion != nil {
		defer rowsSancion.Close()
		for rowsSancion.Next() {
			var tipo string
			var total int
			if err := rowsSancion.Scan(&tipo, &total); err == nil {
				tiposSancion[tipo] = total
			}
		}
	}

	// 7. Estados de apelación
	apelacionEstados := map[string]int{}
	queryApel := fmt.Sprintf(`
		SELECT 
			CASE 
				WHEN apelacion_documento IS NOT NULL AND apelacion_documento != '' THEN 'apelada'
				ELSE 'no_apelada' 
			END as estado_apelacion,
			CASE 
				WHEN estado IN ('resuelta','apelada') THEN 'resuelta'
				WHEN estado = 'pendiente' THEN 'pendiente'
				ELSE 'otro' 
			END as estado_resolucion,
			COUNT(*) as total
		FROM faltas f
		WHERE 1=1%s
		GROUP BY estado_apelacion, estado_resolucion`, whereClause)

	rowsApel, err := h.db.Queryx(queryApel)
	if err == nil && rowsApel != nil {
		defer rowsApel.Close()
		for rowsApel.Next() {
			var estadoApelacion, estadoResolucion string
			var total int
			if err := rowsApel.Scan(&estadoApelacion, &estadoResolucion, &total); err == nil {
				key := estadoApelacion + ":" + estadoResolucion
				apelacionEstados[key] = total
			}
		}
	}

	// ==================== ESTADÍSTICAS DEMOGRÁFICAS ====================

	// 8. Conteo por sexo
	var varones, mujeres int
	h.db.Get(&varones, fmt.Sprintf("SELECT COUNT(*) FROM alumnos WHERE sexo = 'M' AND id IN %s", subqueryAlumnos))
	h.db.Get(&mujeres, fmt.Sprintf("SELECT COUNT(*) FROM alumnos WHERE sexo = 'F' AND id IN %s", subqueryAlumnos))

	// 9. Conteo por escuela profesional
	escuelas := map[string]int{}
	rows, err := h.db.Queryx(fmt.Sprintf(`
		SELECT escuela_profesional, COUNT(*) as total 
		FROM alumnos 
		WHERE id IN %s 
		GROUP BY escuela_profesional`, subqueryAlumnos))
	if err == nil && rows != nil {
		defer rows.Close()
		for rows.Next() {
			var nombre string
			var total int
			if err := rows.Scan(&nombre, &total); err == nil {
				escuelas[nombre] = total
			}
		}
	}

	// 10. Conteo por facultad
	facultades := map[string]int{}
	rows2, err := h.db.Queryx(fmt.Sprintf(`
		SELECT facultad, COUNT(*) as total 
		FROM alumnos 
		WHERE id IN %s 
		GROUP BY facultad`, subqueryAlumnos))
	if err == nil && rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var nombre string
			var total int
			if err := rows2.Scan(&nombre, &total); err == nil {
				facultades[nombre] = total
			}
		}
	}

	// 11. Conteo por modalidad de ingreso
	modalidades := map[string]int{}
	rows3, err := h.db.Queryx(fmt.Sprintf(`
		SELECT modalidad_ingreso, COUNT(*) as total 
		FROM alumnos 
		WHERE id IN %s 
		GROUP BY modalidad_ingreso`, subqueryAlumnos))
	if err == nil && rows3 != nil {
		defer rows3.Close()
		for rows3.Next() {
			var nombre string
			var total int
			if err := rows3.Scan(&nombre, &total); err == nil {
				modalidades[nombre] = total
			}
		}
	}

	// 12. Conteo por estado de matrícula
	estadosMatricula := map[string]int{}
	rows4, err := h.db.Queryx(fmt.Sprintf(`
		SELECT estado_matricula, COUNT(*) as total 
		FROM alumnos 
		WHERE id IN %s 
		GROUP BY estado_matricula`, subqueryAlumnos))
	if err == nil && rows4 != nil {
		defer rows4.Close()
		for rows4.Next() {
			var nombre string
			var total int
			if err := rows4.Scan(&nombre, &total); err == nil {
				estadosMatricula[nombre] = total
			}
		}
	}

	// 13. Conteo por semestre académico
	semestres := map[string]int{}
	querySemestres := fmt.Sprintf(`
		SELECT 
			SUBSTRING_INDEX(SUBSTRING_INDEX(c.nombre, ' ', -1), ' ', 1) AS periodo_convocatoria,
			COUNT(*) as total
		FROM alumnos a
		INNER JOIN faltas f ON f.alumno_id = a.id
		INNER JOIN convocatorias c ON f.convocatoria_id = c.id
		WHERE 1=1%s
		GROUP BY periodo_convocatoria`, whereClause)

	rows5, err := h.db.Queryx(querySemestres)
	if err == nil && rows5 != nil {
		defer rows5.Close()
		for rows5.Next() {
			var periodo string
			var total int
			if err := rows5.Scan(&periodo, &total); err == nil {
				semestres[periodo] = total
			}
		}
	}

	// 14. Distribución por rango de edad
	edades := map[string]int{}
	rows7, err := h.db.Queryx(fmt.Sprintf(`
		SELECT 
			CASE 
				WHEN edad BETWEEN 16 AND 18 THEN '16-18'
				WHEN edad BETWEEN 19 AND 21 THEN '19-21'
				WHEN edad BETWEEN 22 AND 25 THEN '22-25'
				ELSE '26+' 
			END as rango, 
			COUNT(*) as total
		FROM alumnos 
		WHERE id IN %s 
		GROUP BY rango`, subqueryAlumnos))
	if err == nil && rows7 != nil {
		defer rows7.Close()
		for rows7.Next() {
			var rango string
			var total int
			if err := rows7.Scan(&rango, &total); err == nil {
				edades[rango] = total
			}
		}
	}

	// 15. Conteo de reincidencias
	var reincidentes int
	h.db.Get(&reincidentes, fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM (
			SELECT alumno_id 
			FROM faltas f 
			WHERE 1=1%s 
			GROUP BY alumno_id 
			HAVING COUNT(*) > 1
		) as sub`, whereClause))

	// 16. Tendencia temporal (faltas por mes del año actual)
	tendencia := map[string]int{}
	queryTendencia := fmt.Sprintf(`
		SELECT DATE_FORMAT(fecha_falta, '%%Y-%%m') as mes, COUNT(*) as total 
		FROM faltas f 
		WHERE YEAR(fecha_falta) = YEAR(CURDATE())%s 
		GROUP BY mes`, whereClause)

	rows8, err := h.db.Queryx(queryTendencia)
	if err == nil && rows8 != nil {
		defer rows8.Close()
		for rows8.Next() {
			var mes string
			var total int
			if err := rows8.Scan(&mes, &total); err == nil {
				tendencia[mes] = total
			}
		}
	}

	// 17. Porcentaje de faltas resueltas
	var totalFaltas, faltasResueltas int
	h.db.Get(&totalFaltas, fmt.Sprintf("SELECT COUNT(*) FROM faltas f WHERE 1=1%s", whereClause))
	h.db.Get(&faltasResueltas, fmt.Sprintf("SELECT COUNT(*) FROM faltas f WHERE estado IN ('resuelta','apelada')%s", whereClause))

	porcentajeResueltas := 0.0
	if totalFaltas > 0 {
		porcentajeResueltas = float64(faltasResueltas) / float64(totalFaltas) * 100
	}

	// ==================== ESTADÍSTICAS GEOGRÁFICAS ====================

	// 18. Procedencia geográfica concatenada
	procedencias := map[string]int{}
	queryProcedencia := fmt.Sprintf(`
		SELECT
        CONCAT(
            COALESCE((
                SELECT d.name
                FROM solicitudes sol
                INNER JOIN (
                    SELECT alumno_id, MAX(created_at) as max_created
                    FROM solicitudes
                    WHERE alumno_id = a.id
                    GROUP BY alumno_id
                ) ultima ON sol.alumno_id = ultima.alumno_id AND sol.created_at = ultima.max_created
                JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
                JOIN requisitos req ON req.id = ds.requisito_id
                JOIN departaments d ON d.id = ds.opcion_seleccion
                WHERE sol.alumno_id = a.id AND req.nombre = 'Departamento de procedencia'
                LIMIT 1
            ), 'Sin especificar'),
            '/',
            COALESCE((
                SELECT p.name
                FROM solicitudes sol
                INNER JOIN (
                    SELECT alumno_id, MAX(created_at) as max_created
                    FROM solicitudes
                    WHERE alumno_id = a.id
                    GROUP BY alumno_id
                ) ultima ON sol.alumno_id = ultima.alumno_id AND sol.created_at = ultima.max_created
                JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
                JOIN requisitos req ON req.id = ds.requisito_id
                JOIN provinces p ON p.id = ds.opcion_seleccion
                WHERE sol.alumno_id = a.id AND req.nombre = 'Provincia de procedencia'
                LIMIT 1
            ), 'Sin especificar'),
            '/',
            COALESCE((
                SELECT dist.name
                FROM solicitudes sol
                INNER JOIN (
                    SELECT alumno_id, MAX(created_at) as max_created
                    FROM solicitudes
                    WHERE alumno_id = a.id
                    GROUP BY alumno_id
                ) ultima ON sol.alumno_id = ultima.alumno_id AND sol.created_at = ultima.max_created
                JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
                JOIN requisitos req ON req.id = ds.requisito_id
                JOIN districts dist ON dist.id = ds.opcion_seleccion
                WHERE sol.alumno_id = a.id AND req.nombre = 'Distrito de procedencia'
                LIMIT 1
            ), 'Sin especificar')
        ) AS lugar_procedencia,
        COUNT(*) as total
    FROM alumnos a
    WHERE id IN %s
    GROUP BY lugar_procedencia
    HAVING total > 0`, subqueryAlumnos)

	rows9, err := h.db.Queryx(queryProcedencia)
	if err == nil && rows9 != nil {
		defer rows9.Close()
		for rows9.Next() {
			var lugar string
			var total int
			if err := rows9.Scan(&lugar, &total); err == nil {
				procedencias[lugar] = total
			}
		}
	}

	// 19. Por departamento
	porDepartamento := map[string]int{}
	rowsDep, err := h.db.Queryx(fmt.Sprintf(`
		SELECT COALESCE(d.name, 'Sin especificar') as nombre, COUNT(DISTINCT a.id) as total
    FROM alumnos a
    INNER JOIN solicitudes sol ON sol.alumno_id = a.id
    INNER JOIN (
        SELECT alumno_id, MAX(created_at) as max_created
        FROM solicitudes
        WHERE alumno_id IN %s
        GROUP BY alumno_id
    ) ultima_sol ON sol.alumno_id = ultima_sol.alumno_id AND sol.created_at = ultima_sol.max_created
    INNER JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
    INNER JOIN requisitos req ON req.id = ds.requisito_id
    LEFT JOIN departaments d ON d.id = ds.opcion_seleccion
    WHERE a.id IN %s AND req.nombre = 'Departamento de procedencia'
    GROUP BY d.name
    HAVING total > 0`, subqueryAlumnos, subqueryAlumnos))
	if err == nil && rowsDep != nil {
		defer rowsDep.Close()
		for rowsDep.Next() {
			var nombre string
			var total int
			if err := rowsDep.Scan(&nombre, &total); err == nil {
				porDepartamento[nombre] = total
			}
		}
	}

	// 20. Por provincia
	porProvincia := map[string]int{}
	rowsProv, err := h.db.Queryx(fmt.Sprintf(`
		SELECT COALESCE(p.name, 'Sin especificar') as nombre, COUNT(DISTINCT a.id) as total
    FROM alumnos a
    INNER JOIN solicitudes sol ON sol.alumno_id = a.id
    INNER JOIN (
        SELECT alumno_id, MAX(created_at) as max_created
        FROM solicitudes
        WHERE alumno_id IN %s
        GROUP BY alumno_id
    ) ultima_sol ON sol.alumno_id = ultima_sol.alumno_id AND sol.created_at = ultima_sol.max_created
    INNER JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
    INNER JOIN requisitos req ON req.id = ds.requisito_id
    LEFT JOIN provinces p ON p.id = ds.opcion_seleccion
    WHERE a.id IN %s AND req.nombre = 'Provincia de procedencia'
    GROUP BY p.name
    HAVING total > 0`, subqueryAlumnos, subqueryAlumnos))
	if err == nil && rowsProv != nil {
		defer rowsProv.Close()
		for rowsProv.Next() {
			var nombre string
			var total int
			if err := rowsProv.Scan(&nombre, &total); err == nil {
				porProvincia[nombre] = total
			}
		}
	}

	// 21. Por distrito
	porDistrito := map[string]int{}
	rowsDist, err := h.db.Queryx(fmt.Sprintf(`
		SELECT COALESCE(dist.name, 'Sin especificar') as nombre, COUNT(DISTINCT a.id) as total
    FROM alumnos a
    INNER JOIN solicitudes sol ON sol.alumno_id = a.id
    INNER JOIN (
        SELECT alumno_id, MAX(created_at) as max_created
        FROM solicitudes
        WHERE alumno_id IN %s
        GROUP BY alumno_id
    ) ultima_sol ON sol.alumno_id = ultima_sol.alumno_id AND sol.created_at = ultima_sol.max_created
    INNER JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
    INNER JOIN requisitos req ON req.id = ds.requisito_id
    LEFT JOIN districts dist ON dist.id = ds.opcion_seleccion
    WHERE a.id IN %s AND req.nombre = 'Distrito de procedencia'
    GROUP BY dist.name
    HAVING total > 0`, subqueryAlumnos, subqueryAlumnos))
	if err == nil && rowsDist != nil {
		defer rowsDist.Close()
		for rowsDist.Next() {
			var nombre string
			var total int
			if err := rowsDist.Scan(&nombre, &total); err == nil {
				porDistrito[nombre] = total
			}
		}
	}

	// 22. Reporte detallado con CTE optimizada
	var reporteProcedencia []struct {
		Departamento string `db:"departamento" json:"departamento"`
		Provincia    string `db:"provincia" json:"provincia"`
		Distrito     string `db:"distrito" json:"distrito"`
		TotalAlumnos int    `db:"total_alumnos" json:"total_alumnos"`
	}

	cteQuery := fmt.Sprintf(`
     WITH UltimaUbicacionAlumno AS (
        SELECT
            sol.alumno_id,
            MAX(CASE WHEN req.nombre = 'Departamento de procedencia' THEN ds.opcion_seleccion END) AS departamento_id,
            MAX(CASE WHEN req.nombre = 'Provincia de procedencia' THEN ds.opcion_seleccion END) AS provincia_id,
            MAX(CASE WHEN req.nombre = 'Distrito de procedencia' THEN ds.opcion_seleccion END) AS distrito_id
        FROM
            solicitudes sol
        INNER JOIN (
            SELECT alumno_id, MAX(created_at) AS max_created_at
            FROM solicitudes
            WHERE alumno_id IN %s
            GROUP BY alumno_id
        ) AS ultima_sol ON sol.alumno_id = ultima_sol.alumno_id AND sol.created_at = ultima_sol.max_created_at
        JOIN detalle_solicitudes ds ON ds.solicitud_id = sol.id
        JOIN requisitos req ON req.id = ds.requisito_id
        WHERE
            req.nombre IN ('Departamento de procedencia', 'Provincia de procedencia', 'Distrito de procedencia')
            AND sol.alumno_id IN %s
        GROUP BY
            sol.alumno_id
    )
    SELECT
        COALESCE(d.name, 'Sin especificar') AS departamento,
        COALESCE(p.name, 'Sin especificar') AS provincia,
        COALESCE(dist.name, 'Sin especificar') AS distrito,
        COUNT(a.id) AS total_alumnos
    FROM
        alumnos a
    JOIN
        UltimaUbicacionAlumno uba ON a.id = uba.alumno_id
    LEFT JOIN
        departaments d ON d.id = uba.departamento_id
    LEFT JOIN
        provinces p ON p.id = uba.provincia_id
    LEFT JOIN
        districts dist ON dist.id = uba.distrito_id
    WHERE
        a.id IN %s
    GROUP BY
        d.name,
        p.name,
        dist.name
    HAVING
        COUNT(a.id) > 0
    ORDER BY
        total_alumnos DESC,
        departamento,
        provincia,
        distrito`, subqueryAlumnos, subqueryAlumnos, subqueryAlumnos)

	// No usar el if err aquí, dejar que el slice vacío se retorne si falla
	h.db.Select(&reporteProcedencia, cteQuery)

	// ==================== ESTADÍSTICAS ACADÉMICAS ====================

	// 23. Promedio de créditos matriculados
	var promedioCreditos float64
	h.db.Get(&promedioCreditos, fmt.Sprintf(
		"SELECT AVG(CAST(creditos_matriculados AS UNSIGNED)) FROM alumnos WHERE id IN %s",
		subqueryAlumnos))

	// 24. Promedio de PPA
	var promedioPPA float64
	h.db.Get(&promedioPPA, fmt.Sprintf(
		"SELECT AVG(CAST(ppa AS DECIMAL(5,2))) FROM alumnos WHERE id IN %s",
		subqueryAlumnos))

	// 25. Promedio de TCA
	var promedioTCA float64
	h.db.Get(&promedioTCA, fmt.Sprintf(
		"SELECT AVG(CAST(tca AS DECIMAL(5,2))) FROM alumnos WHERE id IN %s",
		subqueryAlumnos))

	// ==================== RESPUESTA FINAL ====================

	res.Data = map[string]interface{}{
		// Metadata
		"convocatoria_filtrada": convocatoriaID,

		// Estadísticas de faltas
		"faltas_graves":           faltasGraves,
		"faltas_leves":            faltasLeves,
		"porcentaje_apeladas":     porcentajeApeladas,
		"ranking_reincidencia":    rankingReincidencia,
		"sanciones_por_periodo":   sancionesPorPeriodo,
		"faltas_por_servicio":     faltasPorServicio,
		"por_tipo_sancion":        tiposSancion,
		"sanciones_por_apelacion": apelacionEstados,
		"reincidentes":            reincidentes,
		"tendencia_mensual":       tendencia,
		"porcentaje_resueltas":    porcentajeResueltas,

		// Estadísticas demográficas
		"varones":                 varones,
		"mujeres":                 mujeres,
		"por_escuela_profesional": escuelas,
		"por_facultad":            facultades,
		"por_modalidad_ingreso":   modalidades,
		"por_estado_matricula":    estadosMatricula,
		"por_semestre":            semestres,
		"por_rango_edad":          edades,

		// Estadísticas geográficas
		"por_procedencia":     procedencias,
		"por_departamento":    porDepartamento,
		"por_provincia":       porProvincia,
		"por_distrito":        porDistrito,
		"reporte_procedencia": reporteProcedencia,

		// Estadísticas académicas
		"promedio_creditos": promedioCreditos,
		"promedio_ppa":      promedioPPA,
		"promedio_tca":      promedioTCA,
	}

	res.Error = false
	return c.JSON(res)
}
