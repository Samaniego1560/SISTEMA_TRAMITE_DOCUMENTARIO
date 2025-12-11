
package models

// ApelacionDocumento representa un documento asociado a una apelación
type ApelacionDocumento struct {
	ID           string `db:"id" json:"id"`
	ApelacionID  string `db:"apelacion_id" json:"apelacion_id"`
	Documento    []byte `db:"documento" json:"documento"`
	Nombre       string `db:"nombre" json:"nombre"`
	Tipo         string `db:"tipo" json:"tipo"`
	CreatedAt    string `db:"created_at" json:"created_at"`
}
