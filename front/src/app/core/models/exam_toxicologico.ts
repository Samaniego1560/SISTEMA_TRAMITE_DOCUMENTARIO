export interface ExamenToxicologico {
    id?:               number;
    alumno_id:         number;
    convocatoria_id:   number;
    codigo_estudiante: string;
    nombres:           string;
    apellido_paterno:  string;
    apellido_materno:  string;
    escuela_profesional:string;
    estado:            Estado;
    comentario:        string;
    fecha_examen:      Date | null;
    usuario_nombre:    string;
    fecha_creacion?:   string;
    fecha_actualizacion?: string;
    id_usuario?:       number;
}
export interface CreateExamenToxicologicoRequest {
    alumno_id: number;
    convocatoria_id: number;
    estado: Estado;
    comentario: string;
}
export interface UpdateExamenToxicologicoRequest {
    estado: Estado;
    comentario: string;
}
export enum Estado {
    Observado = "observado",
    Verificado = "verificado",
    Pendiente = "pendiente",
}