package fault

import (
	"database/sql"
	"time"

	"github.com/asaskevich/govalidator"
)

// Fault Model struct Fault - ACTUALIZADO para manejar NULL
type Fault struct {
	ID                 string         `json:"id" valid:"uuid,required" db:"id"`
	AlumnoID           int64          `json:"alumno_id" valid:"required" db:"alumno_id"`
	ServicioId         int64          `json:"servicio_id" valid:"required" db:"servicio_id"`
	ConvocatoriaId     int64          `json:"convocatoria_id" valid:"required" db:"convocatoria_id"`
	FuenteInformacion  string         `json:"fuente_informacion" valid:"required" db:"fuente_informacion"`
	FechaFalta         time.Time      `json:"fecha_falta" valid:"required" db:"fecha_falta"`
	Estado             string         `json:"estado" valid:"required" db:"estado"`
	Apelable           bool           `json:"apelable" db:"apelable"`
	ApelacionDocumento sql.NullString `json:"apelacion_documento,omitempty" db:"apelacion_documento"` // YA USAS sql.NullString implícitamente con *string
	MotivoResolucion   sql.NullString `json:"motivo_resolucion,omitempty" db:"motivo_resolucion"`     // CAMBIO AQUÍ: de string a sql.NullString
	Observacion        string         `json:"observacion" valid:"-" db:"observacion"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
}

type FaultWithStudent struct {
	Fault
	DNI                     string `db:"dni" json:"dni"`
	Nombres                 string `db:"nombres" json:"nombres"`
	ApellidoPaterno         string `db:"apellido_paterno" json:"apellido_paterno"`
	ApellidoMaterno         string `db:"apellido_materno" json:"apellido_materno"`
	EscuelaProfesional      string `db:"escuela_profesional" json:"escuela_profesional"`
	RoomNumber              string `db:"room_number" json:"room_number"`
	ResidenceName           string `db:"residence_name" json:"residence_name"`
	AdmissionDate           string `db:"admission_date" json:"admission_date"`
	CelularEstudiante       string `db:"celular_estudiante" json:"celular_estudiante"`
	CelularPadre            string `db:"celular_padre" json:"celular_padre"`
	DepartamentoProcedencia string `db:"departamento_procedencia" json:"departamento_procedencia"`
	ProvinciaProcedencia    string `db:"provincia_procedencia" json:"provincia_procedencia"`
	DistritoProcedencia     string `db:"distrito_procedencia" json:"distrito_procedencia"`
	LugarProcedencia        string `db:"lugar_procedencia" json:"lugar_procedencia"`
	Direccion               string `db:"direccion" json:"direccion"`
	CorreoInstitucional     string `db:"correo_institucional" json:"correo_institucional"`
	CodigoEstudiante        string `db:"codigo_estudiante" json:"codigoEstudiante"`
	Edad                    string `db:"edad" json:"edad"`
	Sexo                    string `db:"sexo" json:"sexo"`
	Gravedades              string `db:"gravedades" json:"gravedades"`
	Gravedad                string `db:"gravedad" json:"gravedad"`
}

// FaultArticulo model para tabla faults_articulos
type FaultArticulo struct {
	ID         string    `json:"id" valid:"uuid,required" db:"id"`
	FaultID    string    `json:"fault_id" valid:"uuid,required" db:"falta_id"`
	ArticuloID string    `json:"articulo_id" valid:"uuid,required" db:"articulo_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// FaultInciso model para tabla faults_incisos
type FaultInciso struct {
	ID        string    `json:"id" valid:"uuid,required" db:"id"`
	FaultID   string    `json:"fault_id" valid:"uuid,required" db:"falta_id"`
	IncisoID  string    `json:"inciso_id" valid:"uuid,required" db:"inciso_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewFault - ACTUALIZADO para manejar sql.NullString
func NewFault(id string, alumnoId int64, servicioId int64, convocatoriaId int64, fuenteInformacion string, fechaFalta time.Time, estado string, apelable bool, apelacionDocumento string, motivoResolucion string, observacion string) *Fault {
	fault := &Fault{
		ID:                id,
		AlumnoID:          alumnoId,
		ServicioId:        servicioId,
		ConvocatoriaId:    convocatoriaId,
		FuenteInformacion: fuenteInformacion,
		FechaFalta:        fechaFalta,
		Estado:            estado,
		Apelable:          apelable,
		Observacion:       observacion,
	}

	// Manejo de ApelacionDocumento
	if apelacionDocumento != "" {
		fault.ApelacionDocumento = sql.NullString{String: apelacionDocumento, Valid: true}
	} else {
		fault.ApelacionDocumento = sql.NullString{Valid: false}
	}

	// Manejo de MotivoResolucion
	if motivoResolucion != "" {
		fault.MotivoResolucion = sql.NullString{String: motivoResolucion, Valid: true}
	} else {
		fault.MotivoResolucion = sql.NullString{Valid: false}
	}

	return fault
}

// Constructor para FaultArticulo
func NewFaultArticulo(id, faultID, articuloID string) *FaultArticulo {
	return &FaultArticulo{
		ID:         id,
		FaultID:    faultID,
		ArticuloID: articuloID,
	}
}

// Constructor para FaultInciso
func NewFaultInciso(id, faultID, incisoID string) *FaultInciso {
	return &FaultInciso{
		ID:       id,
		FaultID:  faultID,
		IncisoID: incisoID,
	}
}

func (m *Fault) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Métodos helper para trabajar fácilmente con sql.NullString
func (f *Fault) GetMotivoResolucion() string {
	if f.MotivoResolucion.Valid {
		return f.MotivoResolucion.String
	}
	return ""
}

func (f *Fault) SetMotivoResolucion(motivo string) {
	if motivo == "" {
		f.MotivoResolucion = sql.NullString{Valid: false}
	} else {
		f.MotivoResolucion = sql.NullString{String: motivo, Valid: true}
	}
}

func (f *Fault) GetApelacionDocumento() string {
	if f.ApelacionDocumento.Valid {
		return f.ApelacionDocumento.String
	}
	return ""
}

func (f *Fault) SetApelacionDocumento(documento string) {
	if documento == "" {
		f.ApelacionDocumento = sql.NullString{Valid: false}
	} else {
		f.ApelacionDocumento = sql.NullString{String: documento, Valid: true}
	}
}

type Alumno struct {
	ID                  int64  `db:"id" json:"id"`
	CodigoEstudiante    string `db:"codigo_estudiante" json:"codigoEstudiante"`
	DNI                 string `db:"DNI" json:"dni"`
	Nombres             string `db:"nombres" json:"nombres"`
	ApellidoPaterno     string `db:"apellido_paterno" json:"apellidoPaterno"`
	ApellidoMaterno     string `db:"apellido_materno" json:"apellidoMaterno"`
	Sexo                string `db:"sexo" json:"sexo"`
	Facultad            string `db:"facultad" json:"facultad"`
	EscuelaProfesional  string `db:"escuela_profesional" json:"escuelaProfesional"`
	Edad                int    `db:"edad" json:"edad"`
	CorreoInstitucional string `db:"correo_institucional" json:"correoInstitucional"`
	Direccion           string `db:"direccion" json:"direccion"`
	LugarProcedencia    string `db:"lugar_procedencia" json:"lugarProcedencia"`
	CelularEstudiante   string `db:"celular_estudiante" json:"celularEstudiante"`
}

type FaultDocumento struct {
	ID        string    `db:"id" json:"id"`
	FaultID   string    `db:"falta_id" json:"falta_id"`
	URL       string    `db:"url" json:"url"`
	Archivo   []byte    `db:"archivo" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewFaultDocumento(id, faultID, url string) *FaultDocumento {
	return &FaultDocumento{
		ID:        id,
		FaultID:   faultID,
		URL:       url,
		CreatedAt: time.Now(),
	}
}

type FaultDetalle struct {
	FaltaID             string    `json:"falta_id" db:"falta_id"`
	Observacion         string    `json:"observacion" db:"observacion"`
	FuenteInformacion   string    `json:"fuente_informacion" db:"fuente_informacion"`
	DNI                 string    `json:"dni" db:"dni"`
	Nombres             string    `json:"nombres" db:"nombres"`
	ApellidoPaterno     string    `db:"apellido_paterno" json:"apellidoPaterno"`
	ApellidoMaterno     string    `db:"apellido_materno" json:"apellidoMaterno"`
	Sexo                string    `db:"sexo" json:"sexo"`
	Facultad            string    `db:"facultad" json:"facultad"`
	EscuelaProfesional  string    `db:"escuela_profesional" json:"escuelaProfesional"`
	Edad                int       `db:"edad" json:"edad"`
	CorreoInstitucional string    `db:"correo_institucional" json:"correoInstitucional"`
	Direccion           string    `db:"direccion" json:"direccion"`
	LugarProcedencia    string    `db:"lugar_procedencia" json:"lugarProcedencia"`
	CelularEstudiante   string    `db:"celular_estudiante" json:"celularEstudiante"`
	ArticuloID          string    `json:"articulo_id" db:"articulo_id"`
	ArticuloDescripcion string    `json:"articulo_descripcion" db:"articulo_descripcion"`
	ArticuloGravedad    string    `json:"articulo_gravedad" db:"articulo_gravedad"`
	IncisoID            string    `json:"inciso_id" db:"inciso_id"`
	IncisoNombre        string    `json:"inciso_nombre" db:"inciso_nombre"`
	IncisoDescripcion   string    `json:"inciso_descripcion" db:"inciso_descripcion"`
	CapituloID          string    `json:"capitulo_id" db:"capitulo_id"`
	CapituloNombre      string    `json:"capitulo_nombre" db:"capitulo_nombre"`
	ResolucionID        string    `json:"resolucion_id" db:"resolucion_id"`
	ResolucionNombre    string    `json:"resolucion_nombre" db:"resolucion_nombre"`
	DocumentoURL        string    `json:"documento_url" db:"documento_url"`
	FechaFalta          time.Time `db:"fecha_falta" json:"fecha_falta"`
	ServicioID          int64     `db:"servicio_id" json:"servicio_id"`
}

type IncisoDetalle struct {
	IncisoID          string `json:"inciso_id"`
	IncisoNombre      string `json:"inciso_nombre"`
	IncisoDescripcion string `json:"inciso_descripcion"`
}

type ArticuloDetalle struct {
	ArticuloID          string          `json:"articulo_id"`
	ArticuloDescripcion string          `json:"articulo_descripcion"`
	ArticuloGravedad    string          `json:"articulo_gravedad"`
	Incisos             []IncisoDetalle `json:"incisos"`
}

type CapituloDetalle struct {
	CapituloID     string            `json:"capitulo_id"`
	CapituloNombre string            `json:"capitulo_nombre"`
	Articulos      []ArticuloDetalle `json:"articulos"`
}

type ResolucionDetalle struct {
	ResolucionID     string            `json:"resolucion_id"`
	ResolucionNombre string            `json:"resolucion_nombre"`
	Capitulos        []CapituloDetalle `json:"capitulos"`
}

type DetalleFaltaAgrupado struct {
	FaltaID           string            `json:"falta_id"`
	FuenteInformacion string            `json:"fuente_informacion"`
	Alumno            Alumno            `json:"alumno"`
	Resolucion        ResolucionDetalle `json:"resolucion"`
	Documentos        []string          `json:"documentos"`
	FechaFalta        string            `json:"fecha_falta"`
	Servicio          string            `json:"servicio"`
	ConvocatoriaId    uint64            `json:"convocatoria_id"`
}
