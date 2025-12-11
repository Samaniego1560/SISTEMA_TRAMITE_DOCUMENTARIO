package models

type ConsultationAttentionExcel struct {
	ID                 string `json:"id" db:"id" valid:"uuid,required"`
	FechaConsulta      string `json:"fecha_consulta" db:"fecha_consulta"`
	TipoPersona        string `json:"tipo_persona" db:"tipo_persona"`
	EscuelaProfesional string `json:"escuela_profesional" db:"escuela_profesional"`
	Sexo               string `json:"sexo" db:"sexo"`
}

type SRQMonthlySummary struct {
	TipoParticipante string `db:"tipo_participante"` // "ESTUDIANTE", "DOCENTE", etc.
	Sexo             string `db:"sexo"`              // "M", "F", "C.P"
	Mes              string `db:"mes"`               // "2025-01", "2025-02", ...
	Total            int    `db:"total"`             // cantidad
}

type PatientReportExcel struct {
	DNI           string `db:"dni" json:"dni"`
	NumAttentions int    `db:"numero_atencion" json:"num_attentions"`
	PatientType   string `db:"tipo_participante" json:"patient_type"`
	PhoneNumber   string `db:"num_telefono" json:"phone_number"`
	FullName      string `db:"full_name" json:"full_name"`
	School        string `db:"escuela" json:"school"`
	Diagnosis     string `db:"diagnostico" json:"diagnosis"`
	Status        string `db:"estado_evaluacion" json:"status"`
	FechaRegistro string `db:"created_date" json:"fecha_registro"`
}
