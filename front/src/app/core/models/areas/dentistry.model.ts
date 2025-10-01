import {RequestConsultationRecord} from "./areas.model";
import {Patient} from "./patient.model";

export interface DentistryConsultationBase {
  procedimiento: ProcedureExam;
  examen_bucal: OdontogramExam;
  consulta: ConsultaOdontologia;
}

export interface DentistryConsultationRequest extends DentistryConsultationBase {
  consulta_odontologia: RequestConsultationRecord;
}

export interface ConsultaOdontologia {
  "id": string,
  "relato": string,
  "diagnostico": string,
  "examen_auxiliar": string,
  "examen_clinico": string,
  "tratamiento": string,
  "indicaciones": string,
  "comentarios": string,
}

export interface ProcedureExam {
  id: string;
  tipo_procedimiento: string;
  recibo: string;
  costo: string;
  fecha_pago: string;
  pieza_dental: string;
  comentarios: string;
}

export interface OdontogramExam {
  id: string;
  cpod: string;
  ihos: string;
  observacion: string;
  comentarios: string;
  odontograma_img: string;
}

export interface DentistryConsultationForm {
  patient: Patient,
  attentionDate: string,
  consultation: DentistryConsultationBase
}

export interface SearchPaymentConcept {
  tipo_servicio: string,
  nombre_servicio: string,
  dni: string,
  recibo: string,
}
