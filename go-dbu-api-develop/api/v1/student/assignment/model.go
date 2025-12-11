package assignment

import "time"

// SubmissionResponse representa una convocatoria donde el estudiante tiene asignación
type SubmissionResponse struct {
	ConvocatoriaID     int64      `json:"convocatoria_id"`
	ConvocatoriaNombre string     `json:"convocatoria_nombre"`
	FechaInicio        *time.Time `json:"fecha_inicio"`
	FechaFin           *time.Time `json:"fecha_fin"`
}

// SubmissionWithDetailResponse representa una convocatoria con toda la información del assignment
type SubmissionWithDetailResponse struct {
	Convocatoria ConvocatoriaInfo `json:"convocatoria"`
	Residencia   ResidenciaInfo   `json:"residencia"`
	Cuarto       CuartoInfo       `json:"cuarto"`
	Companeros   []CompaneroInfo  `json:"companeros"`
	Objetos      []ObjetoInfo     `json:"objetos"`
}

// AssignmentDetailResponse representa el detalle completo de la asignación del estudiante
type AssignmentDetailResponse struct {
	Convocatoria ConvocatoriaInfo `json:"convocatoria"`
	Residencia   ResidenciaInfo   `json:"residencia"`
	Cuarto       CuartoInfo       `json:"cuarto"`
	Companeros   []CompaneroInfo  `json:"companeros"`
	Objetos      []ObjetoInfo     `json:"objetos"`
}

// ConvocatoriaInfo información de la convocatoria
type ConvocatoriaInfo struct {
	Nombre string `json:"nombre"`
}

// ResidenciaInfo información de la residencia
type ResidenciaInfo struct {
	Nombre      string `json:"nombre"`
	Direccion   string `json:"direccion"`
	Descripcion string `json:"descripcion"`
}

// CuartoInfo información del cuarto asignado
type CuartoInfo struct {
	Numero       int    `json:"numero"`
	Piso         int    `json:"piso"`
	Capacidad    int    `json:"capacidad"`
	FechaIngreso string `json:"fecha_ingreso"`
}

// CompaneroInfo información de un compañero de cuarto
type CompaneroInfo struct {
	CodigoEstudiante    string `json:"codigo_estudiante"`
	NombreCompleto      string `json:"nombre_completo"`
	Carrera             string `json:"carrera"`
	CorreoInstitucional string `json:"correo_institucional"`
}

// ObjetoInfo información de objetos asignados
type ObjetoInfo struct {
	Nombre          string `json:"nombre"`
	Descripcion     string `json:"descripcion"`
	Categoria       string `json:"categoria"`
	Estado          string `json:"estado"`
	FechaAsignacion string `json:"fecha_asignacion"`
}

// StudentProfileResponse representa los datos personales básicos del estudiante
type StudentProfileResponse struct {
	NombreCompleto   string `json:"nombre_completo"`
	Genero           string `json:"genero"`
	Edad             int32  `json:"edad"`
	DNI              string `json:"dni"`
	CodigoEstudiante string `json:"codigo_estudiante"`
}
