package models

type HeaderMedicalAreaExcel struct {
	Frame string `json:"frame"`
	Title string `json:"title"`
	Area  string `json:"area"`
}

type ConsultationPatientsMedicalAreaExcel struct {
	ID                 string `json:"id" db:"id" valid:"uuid,required"`
	FechaConsulta      string `json:"fecha_consulta" db:"fecha_consulta"`
	DNI                string `json:"dni" db:"dni"`
	CodigoSga          string `json:"codigo_sga" db:"codigo_sga"`
	NombreCompleto     string `json:"nombre_completo" db:"nombre_completo"`
	Sexo               string `json:"sexo" db:"sexo"`
	FechaNacimiento    string `json:"fecha_nacimiento" db:"fecha_nacimiento"`
	Edad               string `json:"edad" db:"edad"`
	Domicilio          string `json:"domicilio" db:"domicilio"`
	NumeroCelular      string `json:"numero_celular" db:"numero_celular"`
	TipoPersona        string `json:"tipo_persona" db:"tipo_persona"`
	EscuelaProfesional string `json:"escuela_profesional" db:"escuela_profesional"`
	Ocupacion          string `json:"ocupacion" db:"ocupacion"`
	Servicios          string `json:"servicios" db:"servicios"`
	TipoProcedimiento  string `json:"tipo_procedimiento" db:"tipo_procedimiento"`
	Recibo             string `json:"recibo" db:"recibo"`
	Costo              string `json:"costo" db:"costo"`
	FechaPago          string `json:"fecha_pago" db:"fecha_pago"`
	PiezaDental        string `json:"pieza_dental" db:"pieza_dental"`
}

type PerformedProceduresExcel struct {
	ID                 string `json:"id" db:"id" valid:"uuid,required"`
	FechaConsulta      string `json:"fecha_consulta" db:"fecha_consulta"`
	TipoPersona        string `json:"tipo_persona" db:"tipo_persona"`
	EscuelaProfesional string `json:"escuela_profesional" db:"escuela_profesional"`
	Sexo               string `json:"sexo" db:"sexo"`
	TipoProcedimiento  string `json:"tipo_procedimiento" db:"tipo_procedimiento"`
}

type ConsultationIntegralAttentionExcel struct {
	ID                 string `json:"id" db:"id" valid:"uuid,required"`
	FechaConsulta      string `json:"fecha_consulta" db:"fecha_consulta"`
	TipoPersona        string `json:"tipo_persona" db:"tipo_persona"`
	EscuelaProfesional string `json:"escuela_profesional" db:"escuela_profesional"`
	Sexo               string `json:"sexo" db:"sexo"`
}

type ReportNursingRua struct {
	Fecha            string `db:"fecha"`
	Codigo           string `db:"codigo"`
	DNI              string `db:"dni"`
	ApellidosNombres string `db:"apellidos_nombres"`
	FechaNacimiento  string `db:"fecha_nacimiento"`
	Sexo             string `db:"sexo"`
	Edad             string `db:"edad"`
	Domicilio        string `db:"domicilio"`
	Procedencia      string `db:"procedencia"`
	CondicionSalud   string `db:"condicion_salud"`
	DosisVacuna      string `db:"dosis_vacuna"`
	Vacuna           string `db:"vacuna"`
	Minsa            string `db:"minsa"`
	Diagnostico      string `db:"diagnostico"`
	Tratamiento      string `db:"tratamiento"`
	Servicio         string `db:"servicio"`
	Responsable      string `db:"responsable"`
	TipoPaciente     string `db:"tipo_paciente"`
	Escuela          string `db:"escuela"`
	Celular          string `db:"celular"`
	Observaciones    string `db:"observaciones"`
	Recibo           string `db:"recibo"`
	Monto            string `db:"monto"`
}

type Report1Admin struct {
	Fecha              string `db:"fecha"`
	Codigo             string `db:"codigo"`
	TipoPersona        string `db:"tipo_persona"`
	DNI                string `db:"dni"`
	ApellidosNombres   string `db:"apellidos_nombres"`
	FechaNacimiento    string `db:"fecha_nacimiento"`
	Sexo               string `db:"sexo"`
	Edad               string `db:"edad"`
	EscuelaProfesional string `db:"escuela"`
	Ocupacion          string `db:"ocupacion"`
	AreaMedica         string `db:"area_medica"`
}

type Report2Admin struct {
	TotalEnfermeria  string `db:"total_enfermeria"`
	TotalMedicina    string `db:"total_medicina"`
	TotalOdontologia string `db:"total_odontologia"`
}
