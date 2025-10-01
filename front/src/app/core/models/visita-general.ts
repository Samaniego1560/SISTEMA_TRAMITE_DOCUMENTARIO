export interface IVisitaGeneral {
    id:                 string;
    tipo_usuario:       TipoUsuario;
    codigo_estudiante:  null | string;
    dni:                null | string;
    nombre_completo:    string;
    genero:             Genero;
    edad:               null;
    escuela:            null | string;
    area:               null;
    motivo_atencion:    string;
    descripcion_motivo: string;
    url_imagen:         null;
    departamento:       null;
    provincia:          null;
    distrito:           null;
    lugar_atencion:     string;
    created_by:         number;
    created_at:         Date;
    updated_by:         number;
    updated_at:         Date;
}

export enum Genero {
    F = "F",
    M = "M",
}

export enum TipoUsuario {
    Administrativo = "administrativo",
    Alumno = "alumno",
    Docente = "docente",
}

export interface Departamento {
    id:         string;
    name:       string;
    created_at: null;
    updated_at: null;
}
export interface Provincia {
    id:            string;
    name:          string;
    department_id: string;
    created_at:    null;
    updated_at:    null;
}
export interface Distrito {
    id:          string;
    name:        string;
    province_id: string;
    created_at:  null;
    updated_at:  null;
}



