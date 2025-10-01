import {IntegralAttention} from "../areas.model";
import {Patient} from "../patient.model";

export interface NursingConsultationForm {
  patient: Patient,
  attentionDate: string,
  consultation: NursingConsultationBase
}

export interface NursingConsultationBase {
  revision_rutina: RevisionRutina;
  datos_acompanante?: CompanionData;
  examenes: Exam;
}

export interface CompanionData {
  id: string;
  dni: string;
  nombres_apellidos: string;
  edad: string;
}

export interface Exam {
  vacunas?: Vacuna[];
  examen_fisico?: PhysicalExam;
  examen_laboratorio?: LaboratoryExam;
  examen_preferencial?: PreferentialExam;
  examen_sexualidad?: SexualityExam;
  examen_visual?: VisualExam;
  tratamiento_medicamentoso?: TratamientoMedicamentoso;
  procedimiento_realizado?: ProcedimientoRealizado;
  consulta_general?: IntegralAttention;
  consulta_medicina_general?: ConsultaMedicinaGeneral,
  consulta?: ConsultaBucalGeneral;
}

export interface ConsultaBucalGeneral {
  id: string;
  relato: string;//no
  diagnostico: string;
  examen_clinico: string;
  examen_auxiliar: string;
  tratamiento: string;
  indicaciones: string;
  comentarios: string;
}

export interface ConsultaMedicinaGeneral {
  id: string;
  fecha_hora: string;
  anamnesis: string;
  examen_clinico: string;
  indicaciones: string;
}

export interface PhysicalExam {
  id: string;
  talla_peso: string;
  perimetro_cintura: string;
  indice_masa_corporal_img: string;
  presion_arterial: string;
  comentarios: string;
}

export interface LaboratoryExam {
  id: string;
  serologia: string;
  bk: string;
  hemograma: string;
  examen_orina: string;
  colesterol: string;
  glucosa: string;
  comentarios: string;
}

export interface PreferentialExam {
  id: string;
  aparato_respiratorio: string;
  aparato_cardiovascular: string;
  aparato_digestivo: string;
  aparato_genitourinario: string;
  papanicolau: string;
  examen_mama: string;
  comentarios: string;
}

export interface SexualityExam {
  id: string;
  actividad_sexual: string;
  planificacion_familiar: string;
  comentarios: string;
}

export interface VisualExam {
  id: string;
  ojo_derecho: string;
  ojo_izquierdo: string;
  presion_ocular: string;
  comentarios: string;
}

export interface TratamientoMedicamentoso {
  "id": string,
  "nombre_generico_medicamento": string,
  "via_administracion": string,
  "hora_aplicacion": string,
  "responsable_atencion": string,
  "observaciones": string,
}

export interface ProcedimientoRealizado {
  "id": string,
  "procedimiento": string,
  "numero_recibo": string,
  "costo": string,
  "fecha_pago": string,
  "comentario": string,
}

export interface Vacuna {
  id: string;
  tipo_vacuna: string;
  fecha_dosis: string;
  minsa: boolean;
  tipo_atencion: string;
  indicaciones: string;
  observaciones: string;
}

export interface RevisionRutina {
  id: string;
  fiebre_ultimo_quince_dias: string;
  tos_mas_quince_dias: string;
  secrecion_lesion_genitales: string;
  fecha_ultima_regla: string;
  comentarios: string;
}


/***/
export interface VaccineListCard {
  name: string;
  vaccines: Vacuna[];
}
