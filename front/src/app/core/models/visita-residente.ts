
export interface VisitaResidente {
    id:                  number;
    alumno_id:           number;
    estado:              string;
    comentario:          string;
    imagen_url?:          string; // de momento no tomar en cuenta es campo
    id_usuario:          number;
    fecha_creacion:      Date;
    fecha_actualizacion: Date;
    alumno_nombre:       string;
    alumno_codigo:       string;
    escuela_profesional: string;
    lugar_procedencia:   string;
    usuario_nombre:      string;
}
export interface CreateVisitaResidenteRequest {
    alumno_id:           number;
    estado:              string;
    comentario:          string;
    imagen_url:          string;
}
export interface UpdateVisitaResidenteRequest {
    alumno_id:           number;
    estado:              string;
    comentario:          string;
    imagen_url:          string;
}

export interface ResidentesPorVisitar {
    alumno_id:           number;
    codigo:              string;
    nombre:              string;
    dni:                 string;
    celular:             string;
    direccion:           string;
    escuela_profesional: string;
    lugar_procedencia:   string;
    solicitud_id:        number;
    convocatoria_id:     number;
    convocatoria_nombre: string;
}


export interface ReporteVisitaResidente {
    total_visitas:      number;
    pendientes:         number;
    verificadas:        number;
    observadas:         number;
    visitas_del_mes:    number;
    alumnos_sin_visita: number;
}
export interface ReportePorEscuelaProfesional {
    escuela_profesional: string;
    total_visitados: number;
}
export interface LugarProcedencia {
    departamento: string;
    total_visitados: number;
}

export interface VisitasPorDepartamento {
    alumno_id:           number;
    codigo:              string;
    nombre:              string;
    escuela_profesional: string;
    departamento:        string;
    provincia:           string;
    distrito:            string;
    direccion:           string;
    celular:             string;
    convocatoria_nombre: string;
}



export interface ExportDataExcel  {
    alumno_id:           number;
    codigo:              string;
    nombre:              string;
    dni:                 string;
    celular:             string;
    celular_padre:       string;
    direccion:           string;
    escuela_profesional: string;
    departamento:        string;
    provincia:           string;
    distrito:            string;
    solicitud_id:        number;
    convocatoria_id:     number;
    convocatoria_nombre: string;
    estado_visita:       string;
}
