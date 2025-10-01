export interface PatientBase {
  tipo_persona: string;
  codigo_sga: string;
  dni: string;
  nombres: string;
  apellidos: string;
  sexo: string;
  edad: string;
  estado_civil: string;
  grupo_sanguineo: string;
  fecha_nacimiento: string;
  lugar_nacimiento: string;
  procedencia: string;
  factor_rh: string;
  escuela_profesional: string;
  ocupacion: string;
  correo_electronico: string;
  numero_celular: string;
  direccion: string;
  ram: boolean;
  alergias: string;
  antecedentes: Antecedent[];
}


export interface Patient extends PatientBase {
  id: string;
  is_deleted: boolean;
  user_deleted?: string;
  deleted_at?: string;
  user_creator: string;
  created_at: string;
  updated_at: string;
}

export interface PatientRequest extends PatientBase {
  id: string;
}

export interface Antecedent {
  id: string;
  nombre_antecedente: string;
  estado_antecedente: string;
}

export interface RequestSearchPatient {
  dni: string;
  nombres: string;
  apellidos: string;
}
