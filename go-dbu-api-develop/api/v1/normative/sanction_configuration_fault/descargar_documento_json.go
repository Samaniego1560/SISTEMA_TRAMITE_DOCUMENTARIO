package sanction_configuration_fault

import (
	"dbu-api/models"
	"encoding/base64"
	"net/http"
	"github.com/gofiber/fiber/v2"
)

// Descargar documento de apelación como JSON (base64)
func (h *handlerSancionesConfi) DescargarDocumentoApelacionJSON(c *fiber.Ctx) error {
	documentoID := c.Params("documento_id")
	if documentoID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "documento_id es requerido"})
	}
	var doc models.ApelacionDocumento
	err := h.db.Get(&doc, "SELECT id, apelacion_id, documento, nombre, tipo, created_at FROM apelacion_documentos WHERE id = ?", documentoID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Documento no encontrado"})
	}
	base64Doc := base64.StdEncoding.EncodeToString(doc.Documento)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id":          doc.ID,
		"apelacion_id": doc.ApelacionID,
		"nombre":      doc.Nombre,
		"tipo":        doc.Tipo,
		"created_at":  doc.CreatedAt,
		"documento":   base64Doc,
	})
}
