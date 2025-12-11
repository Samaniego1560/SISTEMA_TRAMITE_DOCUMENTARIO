package sanction_configuration_fault

import (
	"dbu-api/internal/authorization"
	"dbu-api/internal/logger"
	"dbu-api/internal/middleware"
	internalmodels "dbu-api/internal/models"
	"dbu-api/models"
	"dbu-api/pkg/sanction/configuration_sanction_fault"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type handlerSancionesConfi struct {
	db      *sqlx.DB
	txID    string
	service configuration_sanction_fault.PortsSanctionFault
}

// GetSancionByID devuelve el detalle de una sanción por su ID, incluyendo requiereFechas
func (h *handlerSancionesConfi) GetSancionByID(c *fiber.Ctx) error {
	sancionFaltaID := c.Params("sancion_id")
	logger.Info.Printf("[GetSancionByID] ID recibido: %s", sancionFaltaID)
	if sancionFaltaID == "" {
		logger.Warning.Printf("[GetSancionByID] sancion_id vacío")
		return c.Status(400).JSON(fiber.Map{"error": true, "msg": "El parámetro sancion_id es requerido"})
	}
	var s internalmodels.SancionFaltaAsignada
	err := h.db.Get(&s, "SELECT * FROM sanciones_faltas_normativa WHERE id = ?", sancionFaltaID)
	if err != nil {
		logger.Error.Printf("[GetSancionByID] Error SQL: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": true, "msg": "Sanción asignada no encontrada"})
	}
	logger.Info.Printf("[GetSancionByID] detalle_sancion obtenido: %s", s.DetalleSancion)
	// Lógica requiereFechas basada en detalle_sancion de sancionesAFaltas
	detalle := strings.ToLower(s.DetalleSancion)
	detalle = strings.ReplaceAll(detalle, "á", "a")
	detalle = strings.ReplaceAll(detalle, "é", "e")
	detalle = strings.ReplaceAll(detalle, "í", "i")
	detalle = strings.ReplaceAll(detalle, "ó", "o")
	detalle = strings.ReplaceAll(detalle, "ú", "u")
	detalle = strings.ReplaceAll(detalle, "  ", " ")
	requiere := false
	if strings.Contains(detalle, "separacion temporal") ||
		strings.Contains(detalle, "suspension temporal") ||
		strings.Contains(detalle, "suspencion temporal") ||
		strings.Contains(detalle, "definitiva") ||
		strings.Contains(detalle, "separacion definitiva") ||
		strings.Contains(detalle, "expulsion") {
		requiere = true
	}
	logger.Info.Printf("[GetSancionByID] requiereFechas: %v", requiere)
	// Respuesta
	return c.Status(200).JSON(fiber.Map{
		"error": false,
		"data": fiber.Map{
			"id":              s.ID,
			"detalle_sancion": s.DetalleSancion,
			"requiereFechas":  requiere,
			// ...agrega aquí otros campos relevantes de la sanción asignada si lo deseas
		},
		"msg":  "Detalle de sanción asignada",
		"code": 29,
		"type": "success",
	})
}

// Obtener sanciones asignadas a una falta
func (h *handlerSancionesConfi) GetSancionesAsignadasPorFalta(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	faltaID := c.Params("falta_id")
	if faltaID == "" {
		res.Code, res.Type, res.Msg = 1, "error", "falta_id es requerido"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	list, err := h.service.GetSancionesAsignadasPorFalta(faltaID)
	if err != nil {
		res.Code, res.Type, res.Msg = 2, "error", "Error consultando sanciones asignadas"
		return c.Status(http.StatusInternalServerError).JSON(res)
	}
	sancionAsignada := false
	if len(list) > 0 {
		sancionAsignada = true
	}
	// Mapear la respuesta para incluir el campo 'revisada' en cada sanción
	sancionesOut := make([]map[string]interface{}, len(list))
	for i, s := range list {
		// Consultar si existe apelación para esta sanción asignada y su estado
		var apelada bool = false
		var revisada bool = false
		var estadoApelacion string
		// Buscar si hay apelación y su estado
		err := h.db.Get(&estadoApelacion, "SELECT estado FROM apelaciones WHERE sancion_falta_asignada_id = ? ORDER BY created_at DESC LIMIT 1", s.ID)
		if err == nil {
			apelada = true
			if estadoApelacion == "APROBADA" || estadoApelacion == "RECHAZADA" {
				revisada = true
			}
		}
		sancionesOut[i] = map[string]interface{}{
			"id":               s.ID,
			"falta_id":         s.FaltaID,
			"resolucion_id":    s.ResolucionID,
			"sancion_id":       s.SancionID,
			"fecha_asignacion": s.FechaAsignacion,
			"fecha_inicio":     s.FechaInicio,
			"fecha_fin":        s.FechaFin,
			"observaciones":    s.Observaciones,
			"created_at":       s.CreatedAt,
			"updated_at":       s.UpdatedAt,
			"capitulo_sancion": s.CapituloSancion,
			"articulo_sancion": s.ArticuloSancion,
			"inciso_sancion":   s.IncisoSancion,
			"detalle_sancion":  s.DetalleSancion,
			"apelada":          apelada,
			"revisada":         revisada,
		}
	}
	res.Data = map[string]interface{}{
		"sancion_asignada": sancionAsignada,
		"sanciones":        sancionesOut,
	}
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Sanciones asignadas encontradas"
	return c.Status(http.StatusOK).JSON(res)
}

// Constructor para el handler
func NewHandlerSancionesConfi(db *sqlx.DB, txID string, service configuration_sanction_fault.PortsSanctionFault) *handlerSancionesConfi {
	return &handlerSancionesConfi{
		db:      db,
		txID:    txID,
		service: service,
	}
}

// DescargarDocumentoApelacion permite descargar un documento de apelación por su ID
func (h *handlerSancionesConfi) DescargarDocumentoApelacion(c *fiber.Ctx) error {
	documentoID := c.Params("documento_id")
	if documentoID == "" {
		return c.Status(400).JSON(fiber.Map{"error": true, "msg": "El parámetro documento_id es requerido"})
	}
	var doc models.ApelacionDocumento
	err := h.db.Get(&doc, "SELECT * FROM apelacion_documentos WHERE id = ?", documentoID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": true, "msg": "Documento no encontrado"})
	}
	c.Set("Content-Disposition", "attachment; filename="+doc.Nombre)
	c.Set("Content-Type", "application/octet-stream")
	return c.Send(doc.Documento)
}
func (h *handlerSancionesConfi) CreateSancion(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	req := internalmodels.SancionFaltaAsignada{}

	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	if err := authorization.ValidPermissions(user, h.db, c); err != nil {
		logger.Error.Printf("User does not have permission, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	// DEBUG: Mostrar el JSON recibido y los valores de fecha
	fmt.Println("[DEBUG] Handler AsignarSancionFalta ejecutado")
	var raw map[string]interface{}
	if err := c.BodyParser(&raw); err != nil {
		fmt.Println("[DEBUG] Error BodyParser:", err)
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "", "Invalid body"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	// DEBUG: Mostrar todos los campos relevantes recibidos desde el frontend
	fmt.Println("[DEBUG] JSON recibido:", raw)
	campos := []string{"falta_id", "id", "resolucion_id", "sancion_id", "fecha_asignacion", "fecha_inicio", "fecha_fin", "observaciones", "created_at", "updated_at"}
	for _, campo := range campos {
		if v, ok := raw[campo]; ok {
			fmt.Printf("[DEBUG] %s recibido: %v\n", campo, v)
		}
	}
	logger.Info.Printf("[DEBUG] JSON recibido: %v", raw)
	if v, ok := raw["fecha_inicio"]; ok {
		logger.Info.Printf("[DEBUG] fecha_inicio recibido: %v", v)
	}
	if v, ok := raw["fecha_fin"]; ok {
		logger.Info.Printf("[DEBUG] fecha_fin recibido: %v", v)
	}
	// Asignar los campos normales
	if v, ok := raw["falta_id"].(string); ok {
		req.FaltaID = v
	}
	if v, ok := raw["id"].(string); ok {
		req.ID = v
	}
	if v, ok := raw["resolucion_id"].(string); ok {
		req.ResolucionID = v
	}
	if v, ok := raw["sancion_id"].(string); ok {
		req.SancionID = v
	}
	if v, ok := raw["observaciones"].(string); ok {
		req.Observaciones = v
	}
	if v, ok := raw["capitulo_sancion"].(string); ok {
		req.CapituloSancion = v
	}
	if v, ok := raw["articulo_sancion"].(string); ok {
		req.ArticuloSancion = v
	}
	if v, ok := raw["inciso_sancion"].(string); ok {
		req.IncisoSancion = v
	}
	if v, ok := raw["detalle_sancion"].(string); ok {
		req.DetalleSancion = v
	}
	if v, ok := raw["articulo_id"].(string); ok {
		req.ArticuloID = v
	}
	// Parsear fechas
	if v, ok := raw["fecha_inicio"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			req.FechaInicio = &t
		}
	}
	if v, ok := raw["fecha_fin"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			req.FechaFin = &t
		}
	}
	if v, ok := raw["fecha_asignacion"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			req.FechaAsignacion = t
		}
	}
	if v, ok := raw["created_at"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			req.CreatedAt = t
		}
	}
	if v, ok := raw["updated_at"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			req.UpdatedAt = t
		}
	}
	// Si la sanción requiere fechas, respeta las enviadas en el JSON
	// Si no se envía nada, asigna la fecha actual
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	// Validar solo el campo ID para sanción normativa
	if req.ID == "" {
		logger.Error.Printf("Invalid request data: ID es obligatorio")
		res.Code, res.Type, res.Msg = 1, "", "Invalid data: ID es obligatorio"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	sancion, code, err := h.service.CreateSancion(
		req.ID, req.ResolucionID, req.ArticuloID, req.CapituloSancion, req.ArticuloSancion, req.IncisoSancion, req.DetalleSancion,
	)
	if err != nil {
		logger.Error.Printf("Error creating sanction: %v", err)
		res.Code, res.Type, res.Msg = code, "", err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = sancion
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", "Sancion creada correctamente"
	return c.Status(http.StatusCreated).JSON(res)
}
func (h *handlerSancionesConfi) UpdateSancion(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	req := configuration_sanction_fault.Sancion{}

	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	if err := authorization.ValidPermissions(user, h.db, c); err != nil {
		logger.Error.Printf("User does not have permission, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "", "Invalid body"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	if ok, err := req.Valid(); !ok {
		logger.Error.Printf("Invalid request data: %v", err)
		res.Code, res.Type, res.Msg = 1, "", "Invalid data"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	sancion, code, err := h.service.UpdateSancion(
		req.ID, req.ResolucionID, req.ArticuloID, req.CapituloSancion, req.ArticuloSancion, req.IncisoSancion, req.DetalleSancion,
	)
	if err != nil {
		logger.Error.Printf("Error updating sanction: %v", err)
		res.Code, res.Type, res.Msg = code, "", err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = sancion
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", "Sancion actualizada correctamente"
	return c.Status(http.StatusOK).JSON(res)
}
func (h *handlerSancionesConfi) DeleteSancion(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}

	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	if err := authorization.ValidPermissions(user, h.db, c); err != nil {
		logger.Error.Printf("User does not have permission, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	idStr := c.Params("id")
	if idStr == "" {
		logger.Error.Println("ID inválido")
		res.Code, res.Type, res.Msg = 1, "", "Invalid ID"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	code, err := h.service.DeleteSancion(idStr)
	if err != nil {
		logger.Error.Printf("Error deleting sanction: %v", err)
		res.Code, res.Type, res.Msg = code, "", err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", "Sancion eliminada correctamente"
	return c.Status(http.StatusOK).JSON(res)
}
func (h *handlerSancionesConfi) GetAllSanciones(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}

	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	if err := authorization.ValidPermissions(user, h.db, c); err != nil {
		logger.Error.Printf("User does not have permission, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	sanciones, err := h.service.GetAllSanciones()
	if err != nil {
		logger.Error.Printf("Error retrieving sanctions: %v", err)
		res.Code, res.Type, res.Msg = 1, "", err.Error()
		return c.Status(http.StatusInternalServerError).JSON(res)
	}
	// Lógica para asignar requiereFechas dinámicamente
	for i := range sanciones {
		detalle := strings.ToLower(sanciones[i].DetalleSancion)
		// Normalizar tildes y espacios
		detalle = strings.ReplaceAll(detalle, "á", "a")
		detalle = strings.ReplaceAll(detalle, "é", "e")
		detalle = strings.ReplaceAll(detalle, "í", "i")
		detalle = strings.ReplaceAll(detalle, "ó", "o")
		detalle = strings.ReplaceAll(detalle, "ú", "u")
		detalle = strings.ReplaceAll(detalle, "  ", " ")
		var requiere bool
		if strings.Contains(detalle, "separacion temporal") ||
			strings.Contains(detalle, "suspension temporal") ||
			strings.Contains(detalle, "suspencion temporal") ||
			strings.Contains(detalle, "definitiva") ||
			strings.Contains(detalle, "separacion definitiva") ||
			strings.Contains(detalle, "expulsion") {
			requiere = true
		} else {
			requiere = false
		}
		sanciones[i].RequiereFechas = requiere
	}
	// Serializar explícitamente para asegurar que requiereFechas aparezca
	sancionesOut := make([]map[string]interface{}, len(sanciones))
	for i, s := range sanciones {
		sancionesOut[i] = map[string]interface{}{
			"id":                   s.ID,
			"resolucion_id":        s.ResolucionID,
			"resolucion_nombre":    s.ResolucionNombre,
			"articulo_id":          s.ArticuloID,
			"articulo_descripcion": s.ArticuloDescripcion,
			"gravedad":             s.Gravedad,
			"capitulo_sancion":     s.CapituloSancion,
			"capitulo_nombre":      s.CapituloNombre,
			"articulo_sancion":     s.ArticuloSancion,
			"inciso_sancion":       s.IncisoSancion,
			"detalle_sancion":      s.DetalleSancion,
			"created_at":           s.CreatedAt,
			"updated_at":           s.UpdatedAt,
			"requiereFechas":       s.RequiereFechas,
		}
	}
	res.Error = false
	res.Data = sancionesOut
	res.Code, res.Type, res.Msg = 29, "", "Sanciones obtenidas correctamente"
	return c.Status(http.StatusOK).JSON(res)
}

// Asignar una sanción a una falta concreta
func (h *handlerSancionesConfi) AsignarSancionFalta(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	req := internalmodels.SancionFaltaAsignada{}
	var estadoReq struct {
		Estado string `json:"estado"`
	}

	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	if err := authorization.ValidPermissions(user, h.db, c); err != nil {
		logger.Error.Printf("User does not have permission, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	var raw map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "", "Invalid body"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	// Refuerzo: parsear fechas explícitamente desde el JSON recibido
	if err := c.BodyParser(&raw); err == nil {
		layout := time.RFC3339
		if fechaInicioStr, ok := raw["fecha_inicio"].(string); ok && fechaInicioStr != "" {
			if t, err := time.Parse(layout, fechaInicioStr); err == nil {
				req.FechaInicio = &t
			} else {
				logger.Error.Printf("Error parseando fecha_inicio: %v", err)
			}
		}
		if fechaFinStr, ok := raw["fecha_fin"].(string); ok && fechaFinStr != "" {
			if t, err := time.Parse(layout, fechaFinStr); err == nil {
				req.FechaFin = &t
			} else {
				logger.Error.Printf("Error parseando fecha_fin: %v", err)
			}
		}
	}
	// Parsear el estado si viene en el body
	_ = c.BodyParser(&estadoReq)
	// Asegurar que FaltaID se obtenga del body o de la URL
	if req.FaltaID == "" {
		faltaID := c.Params("falta_id")
		if faltaID != "" {
			req.FaltaID = faltaID
		}
	}
	if estadoReq.Estado != "" && req.FaltaID != "" {
		query := "UPDATE faltas SET estado = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
		_, err := h.db.Exec(query, estadoReq.Estado, req.FaltaID)
		if err != nil {
			logger.Error.Printf("No se pudo actualizar el estado de la falta: %v", err)
			// No bloquea la asignación de sanción, solo loguea
		}
	}
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	// Permitir FaltaID por body o por parámetro de ruta (si aplica)
	if req.FaltaID == "" {
		faltaID := c.Params("falta_id")
		if faltaID != "" {
			req.FaltaID = faltaID
		}
	}
	if req.FaltaID == "" {
		res.Code, res.Type, res.Msg = 5, "error", "Debe especificar el id de la falta"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	// Obtener información de la sanción para saber si requiere fechas
	sancion, _, err := h.service.GetSancionByID(req.SancionID)
	if err != nil || sancion == nil {
		logger.Error.Printf("No se encontró la sanción para validar fechas: %v", err)
		res.Code, res.Type, res.Msg = 2, "error", "No se encontró la sanción"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	// Validar fechas si vienen en el body
	if req.FechaInicio == nil || req.FechaFin == nil {
		res.Code, res.Type, res.Msg = 3, "error", "Debe especificar fecha de inicio y fin"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	if req.FechaFin.Before(*req.FechaInicio) || req.FechaFin.Equal(*req.FechaInicio) {
		res.Code, res.Type, res.Msg = 4, "error", "La fecha de fin debe ser mayor a la fecha de inicio"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	code, err := h.service.AsignarSancionFalta(&req)
	if err != nil {
		logger.Error.Printf("Error asignando sanción a falta: %v", err)
		res.Code, res.Type, res.Msg = code, "", err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}
	res.Data = req
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", "Sanción asignada correctamente"
	return c.Status(http.StatusCreated).JSON(res)
}

// Registrar una apelación sobre una sanción asignada
func (h *handlerSancionesConfi) RegistrarApelacion(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	req := models.Apelacion{}

	// Obtener el ID de la sanción asignada desde la URL
	sancionFaltaAsignadaID := c.Params("sancion_falta_asignada_id")
	if sancionFaltaAsignadaID == "" {
		res.Code, res.Type, res.Msg = 1, "error", "El parámetro sancion_falta_asignada_id es requerido en la URL"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	bearer := c.Get("Authorization")
	user, err := middleware.GetUser(bearer, h.db)
	if err != nil {
		logger.Error.Printf("Unauthenticated user, error: %v", err)
		res.Code, res.Type, res.Msg = 9, "error", "unauthenticated"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}
	if err := authorization.ValidPermissions(user, h.db, c); err != nil {
		logger.Error.Printf("User does not have permission, error: %v", err)
		res.Code, res.Type, res.Msg = 10, "error", "rejected for route permits"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Tomar el ID del alumno autenticado como usuario_apela
	req.UsuarioApela = fmt.Sprintf("%v", user.ID)

	// Procesar múltiples archivos adjuntos (documentos)
	form, err := c.MultipartForm()
	var documentosGuardados []struct {
		Nombre string
		Tipo   string
		Datos  []byte
	}
	if err == nil && form != nil {
		files := form.File["documentos"]
		for _, file := range files {
			f, err := file.Open()
			if err != nil {
				logger.Error.Printf("Error abriendo archivo de apelación: %v", err)
				res.Code, res.Type, res.Msg = 11, "error", "No se pudo abrir el documento adjunto"
				return c.Status(http.StatusInternalServerError).JSON(res)
			}
			defer f.Close()
			fileBytes := make([]byte, file.Size)
			_, err = f.Read(fileBytes)
			if err != nil {
				logger.Error.Printf("Error leyendo archivo de apelación: %v", err)
				res.Code, res.Type, res.Msg = 11, "error", "No se pudo leer el documento adjunto"
				return c.Status(http.StatusInternalServerError).JSON(res)
			}
			documentosGuardados = append(documentosGuardados, struct {
				Nombre string
				Tipo   string
				Datos  []byte
			}{Nombre: file.Filename, Tipo: file.Header.Get("Content-Type"), Datos: fileBytes})
		}
	}

	// Procesar el resto del body (motivo, estado, observaciones, etc.)
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "", "Invalid body"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	req.SancionFaltaAsignadaID = sancionFaltaAsignadaID
	code, err := h.service.RegistrarApelacion(&req)
	if err != nil {
		logger.Error.Printf("Error registrando apelación: %v", err)
		res.Code, res.Type, res.Msg = code, "", err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	// Guardar los documentos en la tabla apelacion_documentos
	for _, doc := range documentosGuardados {
		idDoc := uuid.New().String()
		_, err := h.db.Exec(`INSERT INTO apelacion_documentos (id, apelacion_id, documento, nombre, tipo, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
			idDoc, req.ID, doc.Datos, doc.Nombre, doc.Tipo, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			logger.Error.Printf("Error guardando documento de apelación en BD: %v", err)
			// No detenemos el proceso, pero puedes agregar lógica para manejar errores si lo deseas
		}
	}

	res.Data = req
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "", "Apelación registrada correctamente"
	return c.Status(http.StatusCreated).JSON(res)
}

// Listar sanciones normativas asociadas a la resolución y servicio de una falta
func (h *handlerSancionesConfi) GetSancionesByFalta(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	faltaID := c.Params("falta_id")
	if faltaID == "" {
		res.Code, res.Type, res.Msg = 1, "error", "falta_id es requerido"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	// Buscar la falta para obtener resolucion_id y servicio_id
	var falta struct {
		ResolucionID string `db:"resolucion_id"`
		ServicioID   int64  `db:"servicio_id"`
	}
	db := h.db
	err := db.Get(&falta, "SELECT resolucion_id, servicio_id FROM faltas WHERE id = ?", faltaID)
	if err != nil {
		logger.Error.Printf("No se encontró la falta: %v", err)
		res.Code, res.Type, res.Msg = 2, "error", "Falta no encontrada"
		return c.Status(http.StatusNotFound).JSON(res)
	}
	// Consultar sanciones normativas asociadas a la resolución
	var sanciones []configuration_sanction_fault.Sancion
	err = db.Select(&sanciones, `SELECT * FROM sanciones_faltas_normativa WHERE resolucion_id = ?`, falta.ResolucionID)
	if err != nil {
		logger.Error.Printf("Error consultando sanciones: %v", err)
		res.Code, res.Type, res.Msg = 3, "error", "Error consultando sanciones"
		return c.Status(http.StatusInternalServerError).JSON(res)
	}
	res.Data = sanciones
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Sanciones encontradas"
	return c.Status(http.StatusOK).JSON(res)
}

// Sugerir sanciones normativas aplicables a una falta/inciso
func (h *handlerSancionesConfi) GetSancionesSugeridas(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	faltaID := c.Params("falta_id")
	if faltaID == "" {
		res.Code, res.Type, res.Msg = 1, "error", "falta_id es requerido"
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	// 1. Obtener todos los incisos asociados a la falta
	var incisos []struct {
		IncisoID string `db:"inciso_id"`
	}
	err := h.db.Select(&incisos, "SELECT inciso_id FROM faltas_incisos WHERE falta_id = ?", faltaID)
	if err != nil || len(incisos) == 0 {
		logger.Error.Printf("No se encontraron incisos para la falta: %v", err)
		res.Code, res.Type, res.Msg = 2, "error", "No se encontraron incisos para la falta"
		return c.Status(http.StatusNotFound).JSON(res)
	}

	sancionesTotales := []configuration_sanction_fault.Sancion{}

	for _, inc := range incisos {
		// 2. Obtener el articulo_id del inciso
		var inciso struct {
			ArticuloID string `db:"articulo_id"`
		}
		err = h.db.Get(&inciso, "SELECT articulo_id FROM incisos WHERE id = ?", inc.IncisoID)
		if err != nil || inciso.ArticuloID == "" {
			logger.Error.Printf("No se encontró articulo para el inciso %s: %v", inc.IncisoID, err)
			continue
		}
		// 3. Obtener el capitulo_id y gravedad del artículo
		var articulo struct {
			CapituloID string `db:"capitulo_id"`
			Gravedad   string `db:"gravedad"`
		}
		err = h.db.Get(&articulo, "SELECT capitulo_id, gravedad FROM articulos WHERE id = ?", inciso.ArticuloID)
		if err != nil || articulo.CapituloID == "" {
			logger.Error.Printf("No se encontró capitulo para el artículo %s: %v", inciso.ArticuloID, err)
			continue
		}
		// 4. Obtener el resolucion_id del capítulo
		var capitulo struct {
			ResolucionID string `db:"resolucion_id"`
		}
		err = h.db.Get(&capitulo, "SELECT resolucion_id FROM capitulos WHERE id = ?", articulo.CapituloID)
		if err != nil || capitulo.ResolucionID == "" {
			logger.Error.Printf("No se encontró resolucion para el capítulo %s: %v", articulo.CapituloID, err)
			continue
		}
		// 5. Buscar sanciones posibles
		var sanciones []configuration_sanction_fault.Sancion
		err = h.db.Select(&sanciones, "SELECT * FROM sanciones_faltas_normativa WHERE resolucion_id = ? AND articulo_id = ?", capitulo.ResolucionID, inciso.ArticuloID)
		if err != nil {
			logger.Error.Printf("Error consultando sanciones para resolucion %s y articulo %s: %v", capitulo.ResolucionID, inciso.ArticuloID, err)
			continue
		}
		// Consultar gravedad para el articulo y asignarla a cada sanción
		var gravedad string
		err = h.db.Get(&gravedad, "SELECT gravedad FROM articulos WHERE id = ?", inciso.ArticuloID)
		if err != nil {
			gravedad = ""
		}
		for i := range sanciones {
			sanciones[i].Gravedad = gravedad
		}
		sancionesTotales = append(sancionesTotales, sanciones...)
	}

	res.Data = sancionesTotales
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Sanciones sugeridas encontradas"
	return c.Status(http.StatusOK).JSON(res)
}

// Consultar apelaciones asociadas a una sanción/falta
func (h *handlerSancionesConfi) GetApelacionesPorSancionFaltaAsignada(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	sancionFaltaAsignadaID := c.Params("sancion_falta_asignada_id")
	if sancionFaltaAsignadaID == "" {
		res.Code, res.Type, res.Msg = 1, "error", "El parámetro sancion_falta_asignada_id es requerido"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	var apelaciones []models.Apelacion
	err := h.db.Select(&apelaciones, "SELECT * FROM apelaciones WHERE sancion_falta_asignada_id = ?", sancionFaltaAsignadaID)
	if err != nil {
		logger.Error.Printf("Error consultando apelaciones: %v", err)
		res.Code, res.Type, res.Msg = 2, "error", "Error consultando apelaciones"
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	type DocumentoMeta struct {
		ID          string `json:"id"`
		ApelacionID string `json:"apelacion_id"`
		CreatedAt   string `json:"created_at"`
		Nombre      string `json:"nombre,omitempty"`
		Tipo        string `json:"tipo,omitempty"`
	}
	type VeredictoDocumento struct {
		ID          string `json:"id"`
		ApelacionID string `json:"apelacion_id"`
		CreatedAt   string `json:"created_at"`
		Nombre      string `json:"nombre"`
		Tipo        string `json:"tipo"`
		Documento   string `json:"documento_base64"`
	}
	var resultado []map[string]interface{}
	for _, apelacion := range apelaciones {
		var docs []DocumentoMeta
		err := h.db.Select(&docs, "SELECT id, apelacion_id, created_at, nombre, tipo FROM apelacion_documentos WHERE apelacion_id = ? AND tipo != 'veredicto'", apelacion.ID)
		if err != nil {
			docs = []DocumentoMeta{}
		}
		// Buscar documento de veredicto (si existe)
		var veredictoDoc VeredictoDocumento
		veredictoFound := false
		row := h.db.QueryRowx("SELECT id, apelacion_id, created_at, nombre, tipo, documento FROM apelacion_documentos WHERE apelacion_id = ? AND tipo = 'veredicto' LIMIT 1", apelacion.ID)
		var documentoBytes []byte
		err = row.Scan(&veredictoDoc.ID, &veredictoDoc.ApelacionID, &veredictoDoc.CreatedAt, &veredictoDoc.Nombre, &veredictoDoc.Tipo, &documentoBytes)
		if err == nil {
			veredictoDoc.Documento = base64.StdEncoding.EncodeToString(documentoBytes)
			veredictoFound = true
		}
		apelacionMap := map[string]interface{}{
			"id":                        apelacion.ID,
			"sancion_falta_asignada_id": apelacion.SancionFaltaAsignadaID,
			"motivo":                    apelacion.Motivo,
			"estado":                    apelacion.Estado,
			"usuario_apela":             apelacion.UsuarioApela,
			"observaciones":             apelacion.Observaciones,
			"fecha_apelacion":           apelacion.FechaApelacion,
			"fecha_resolucion":          apelacion.FechaResolucion,
			"created_at":                apelacion.CreatedAt,
			"updated_at":                apelacion.UpdatedAt,
			"documentos":                docs,
		}
		if veredictoFound {
			apelacionMap["veredicto_documento"] = veredictoDoc
		}
		resultado = append(resultado, apelacionMap)
	}
	res.Data = resultado
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Apelaciones encontradas"
	return c.Status(http.StatusOK).JSON(res)
}

// Revisar y resolver apelación
func (h *handlerSancionesConfi) ResolverApelacion(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	apelacionID := c.Params("apelacion_id")
	if apelacionID == "" {
		res.Code, res.Type, res.Msg = 1, "error", "El parámetro apelacion_id es requerido"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	var req models.Apelacion
	// Procesar posible archivo adjunto (veredicto)
	form, formErr := c.MultipartForm()
	var veredictoGuardado *struct {
		Nombre string
		Tipo   string
		Datos  []byte
	}
	if formErr == nil && form != nil {
		files := form.File["veredicto"]
		if len(files) > 0 {
			file := files[0]
			f, err := file.Open()
			if err != nil {
				logger.Error.Printf("Error abriendo archivo de veredicto: %v", err)
				res.Code, res.Type, res.Msg = 11, "error", "No se pudo abrir el documento de veredicto"
				return c.Status(http.StatusInternalServerError).JSON(res)
			}
			defer f.Close()
			fileBytes := make([]byte, file.Size)
			_, err = f.Read(fileBytes)
			if err != nil {
				logger.Error.Printf("Error leyendo archivo de veredicto: %v", err)
				res.Code, res.Type, res.Msg = 11, "error", "No se pudo leer el documento de veredicto"
				return c.Status(http.StatusInternalServerError).JSON(res)
			}
			veredictoGuardado = &struct {
				Nombre string
				Tipo   string
				Datos  []byte
			}{Nombre: file.Filename, Tipo: file.Header.Get("Content-Type"), Datos: fileBytes}
		}
	}
	if err := c.BodyParser(&req); err != nil {
		logger.Error.Printf("couldn't parse body request, error: %v", err)
		res.Code, res.Type, res.Msg = 1, "error", "Invalid body"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	// Actualizar estado, observaciones y fecha_resolucion
	now := time.Now().Format("2006-01-02 15:04:05")
	fechaResolucion := ""
	if req.Estado == "APROBADA" || req.Estado == "RECHAZADA" {
		fechaResolucion = now
	}
	logger.Info.Printf("[ResolverApelacion] apelacionID: %s, estado: %s, observaciones: %s, fechaResolucion: %s", apelacionID, req.Estado, req.Observaciones, fechaResolucion)
	query := `UPDATE apelaciones SET estado = ?, observaciones = ?, fecha_resolucion = ? WHERE id = ?`
	result, err := h.db.Exec(query, req.Estado, req.Observaciones, fechaResolucion, apelacionID)
	if err != nil {
		logger.Error.Printf("Error actualizando apelación: %v", err)
		res.Code, res.Type, res.Msg = 2, "error", "Error actualizando apelación"
		return c.Status(http.StatusInternalServerError).JSON(res)
	}
	rowsAffected, _ := result.RowsAffected()
	logger.Info.Printf("[ResolverApelacion] Filas afectadas en apelaciones: %d", rowsAffected)

	// Guardar el documento de veredicto si fue adjuntado
	if veredictoGuardado != nil {
		idDoc := uuid.New().String()
		_, err := h.db.Exec(`INSERT INTO apelacion_documentos (id, apelacion_id, documento, nombre, tipo, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
			idDoc, apelacionID, veredictoGuardado.Datos, veredictoGuardado.Nombre, "veredicto", time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			logger.Error.Printf("Error guardando documento de veredicto en BD: %v", err)
		}
	}

	// Si la apelación fue aprobada, actualizar el estado de la sanción asociada a 'APELADA'
	if req.Estado == "APROBADA" {
		// Obtener el id de la sanción asignada desde la apelación
		var sancionFaltaAsignadaID string
		err := h.db.Get(&sancionFaltaAsignadaID, "SELECT sancion_falta_asignada_id FROM apelaciones WHERE id = ?", apelacionID)
		if err == nil && sancionFaltaAsignadaID != "" {
			// Actualizar sancion_falta_asignada
			_, err := h.db.Exec("UPDATE sancion_falta_asignada SET estado = ? WHERE id = ?", "APELADA", sancionFaltaAsignadaID)
			if err != nil {
				logger.Error.Printf("Error actualizando estado de sanción asignada: %v", err)
			}
			// Obtener el id de la falta asociada
			var faltaID string
			err = h.db.Get(&faltaID, "SELECT falta_id FROM sancion_falta_asignada WHERE id = ?", sancionFaltaAsignadaID)
			if err == nil && faltaID != "" {
				// Actualizar faltas
				_, err := h.db.Exec("UPDATE faltas SET estado = ? WHERE id = ?", "APELADA", faltaID)
				if err != nil {
					logger.Error.Printf("Error actualizando estado de falta: %v", err)
				}
			}
		}
	}

	res.Data = req
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Apelación revisada/resuelta"
	return c.Status(http.StatusOK).JSON(res)
}

// Obtener todos los detalles de una apelación, incluyendo documentos en base64
func (h *handlerSancionesConfi) GetDetalleApelacion(c *fiber.Ctx) error {
	res := internalmodels.Response{Error: true}
	apelacionID := c.Params("apelacion_id")
	if apelacionID == "" {
		res.Code, res.Type, res.Msg = 1, "error", "El parámetro apelacion_id es requerido"
		return c.Status(http.StatusBadRequest).JSON(res)
	}
	var apelacion models.Apelacion
	err := h.db.Get(&apelacion, "SELECT * FROM apelaciones WHERE id = ?", apelacionID)
	if err != nil {
		res.Code, res.Type, res.Msg = 2, "error", "Apelación no encontrada"
		return c.Status(http.StatusNotFound).JSON(res)
	}
	// Obtener documentos asociados
	var docs []models.ApelacionDocumento
	err = h.db.Select(&docs, "SELECT * FROM apelacion_documentos WHERE apelacion_id = ?", apelacionID)
	if err != nil {
		docs = []models.ApelacionDocumento{}
	}
	// Codificar documentos en base64
	type DocumentoDetalle struct {
		ID          string `json:"id"`
		ApelacionID string `json:"apelacion_id"`
		Nombre      string `json:"nombre"`
		Tipo        string `json:"tipo"`
		CreatedAt   string `json:"created_at"`
		Documento   string `json:"documento_base64"`
	}
	var documentosDetalle []DocumentoDetalle
	for _, d := range docs {
		documentosDetalle = append(documentosDetalle, DocumentoDetalle{
			ID:          d.ID,
			ApelacionID: d.ApelacionID,
			Nombre:      d.Nombre,
			Tipo:        d.Tipo,
			CreatedAt:   d.CreatedAt,
			Documento:   base64.StdEncoding.EncodeToString(d.Documento),
		})
	}
	// Respuesta completa
	res.Data = map[string]interface{}{
		"id":                        apelacion.ID,
		"sancion_falta_asignada_id": apelacion.SancionFaltaAsignadaID,
		"motivo":                    apelacion.Motivo,
		"estado":                    apelacion.Estado,
		"usuario_apela":             apelacion.UsuarioApela,
		"observaciones":             apelacion.Observaciones,
		"fecha_apelacion":           apelacion.FechaApelacion,
		// "fecha_revision":            apelacion.FechaRevision, // Eliminado, no existe en la tabla
		"fecha_resolucion": apelacion.FechaResolucion,
		"created_at":       apelacion.CreatedAt,
		"updated_at":       apelacion.UpdatedAt,
		"documentos":       documentosDetalle,
	}
	res.Error = false
	res.Code, res.Type, res.Msg = 29, "success", "Detalle de apelación encontrado"
	return c.Status(http.StatusOK).JSON(res)
}
