package models

type Apelacion struct {
	ID                     string  `db:"id" json:"id"`                                               // ID del apelante
	SancionFaltaAsignadaID string  `db:"sancion_falta_asignada_id" json:"sancion_falta_asignada_id"` // ID de la sanción
	Motivo                 string  `db:"motivo" json:"motivo"`                                       // Motivo de la apelación
	Estado                 string  `db:"estado" json:"estado"`                                       // PENDIENTE, APROBADA, RECHAZADA
	UsuarioApela           string  `db:"usuario_apela" json:"usuario_apela"`                         // ID del estudiante
	Observaciones          string  `db:"observaciones" json:"observaciones"`                         // Observaciones del revisor
	VeredictoObservacion   *string `db:"veredicto_observacion" json:"veredicto_observacion"`         // Observación del veredicto de la comisión
	FechaApelacion         *string `db:"fecha_apelacion" json:"fecha_apelacion"`                     // Fecha de la apelación
	FechaRevision          *string `db:"fecha_revision" json:"fecha_revision"`                       // Fecha en que se revisa
	FechaResolucion        *string `db:"fecha_resolucion" json:"fecha_resolucion"`                   // Fecha en que se resuelve
	CreatedAt              *string `db:"created_at" json:"created_at"`                               // Fecha de creación
	UpdatedAt              *string `db:"updated_at" json:"updated_at"`                               // Fecha de actualización
}
