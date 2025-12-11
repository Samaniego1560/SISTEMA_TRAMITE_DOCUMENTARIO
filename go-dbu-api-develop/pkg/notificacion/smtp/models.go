package smtp

// NotificationData estructura que mapea con tu respuesta JSON
type NotificationData struct {
	FaltaID           string     `json:"falta_id"`
	FuenteInformacion string     `json:"fuente_informacion"`
	Alumno            Alumno     `json:"alumno"`
	Resolucion        Resolucion `json:"resolucion"`
	Documentos        []string   `json:"documentos"`
	FechaFalta        string     `json:"fecha_falta"`
	Servicio          string     `json:"servicio"`
	SemestreAcademico string `json:"semestre_academico"`
	ObsLabel          string `json:"obs_label"`
	FechaHtml         string `json:"fecha_html"`
}

type Alumno struct {
	ID                  int    `json:"id"`
	CodigoEstudiante    string `json:"codigoEstudiante"`
	DNI                 string `json:"dni"`
	Nombres             string `json:"nombres"`
	ApellidoPaterno     string `json:"apellidoPaterno"`
	ApellidoMaterno     string `json:"apellidoMaterno"`
	Sexo                string `json:"sexo"`
	Facultad            string `json:"facultad"`
	EscuelaProfesional  string `json:"escuelaProfesional"`
	Edad                int    `json:"edad"`
	CorreoInstitucional string `json:"correoInstitucional"`
	Direccion           string `json:"direccion"`
	LugarProcedencia    string `json:"lugarProcedencia"`
	CelularEstudiante   string `json:"celularEstudiante"`
}

type Resolucion struct {
	ResolucionID     string     `json:"resolucion_id"`
	ResolucionNombre string     `json:"resolucion_nombre"`
	Capitulos        []Capitulo `json:"capitulos"`
}

type Capitulo struct {
	CapituloID     string     `json:"capitulo_id"`
	CapituloNombre string     `json:"capitulo_nombre"`
	Articulos      []Articulo `json:"articulos"`
}

type Articulo struct {
	ArticuloID          string   `json:"articulo_id"`
	ArticuloDescripcion string   `json:"articulo_descripcion"`
	ArticuloGravedad    string   `json:"articulo_gravedad"`
	Incisos             []Inciso `json:"incisos"`
}

type Inciso struct {
	IncisoID          string `json:"inciso_id"`
	IncisoNombre      string `json:"inciso_nombre"`
	IncisoDescripcion string `json:"inciso_descripcion"`
}

// EmailTemplate estructura para el template del email
type EmailTemplate struct {
	NumeroNotificacion string
	Estudiante         string
	Facultad           string
	Fecha              string
	Direccion          string
	ContadorFaltas     int
	GravedadFalta      string
	Servicio           string
	NombreNotificador  string
	NombreNotificado   string
	FaltasDetalle      []string
	SancionDetalle     string
	Resolucion         Resolucion
	FuenteInformacion  string
	CapituloSancion    string
	ArticuloSancion    string
	IncisoSancion      string
	DescripcionSancion string
	SemestreAcademico   string
	ObsLabel           string
	FechaHtml          string
}
