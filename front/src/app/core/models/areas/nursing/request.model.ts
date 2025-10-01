import {RequestConsultationRecord} from "../areas.model";
import {NursingConsultationBase} from "./nursing.model";

export interface NursingConsultationRequest extends NursingConsultationBase {
  consulta_enfermeria: RequestConsultationRecord;
}

export interface RequestUpdateAreaAssigned {
  "consulta_id": string;
  "area_asignada": string;
}
