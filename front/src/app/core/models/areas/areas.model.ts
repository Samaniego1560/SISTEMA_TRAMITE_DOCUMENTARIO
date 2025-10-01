import {Patient} from "./patient.model";

export interface RequestConsultationRecord {
  id: string;
  fecha_consulta: string;
  paciente_id: string;
}

export interface DataForm<T> {
  patient: Patient,
  data: T
}

export interface IntegralAttention {
  id: string;
  fecha: string;
  hora: string;
  numero_recibo: string;
  costo: string;
  fecha_pago: string;
  edad: string;
  motivo_consulta: string;
  tiempo_enfermedad: string;
  apetito: string;
  sed: string;
  suenio: string;
  estado_animo: string;
  orina: string;
  deposiciones: string;
  temperatura: string;
  presion_arterial: string;
  frecuencia_cardiaca: string;
  frecuencia_respiratoria: string;
  peso: string;
  talla: string;
  indice_masa_corporal: string;
  diagnostico: string;
  tratamiento: string;
  examenes_axuliares: string;
  referencia: string;
  observacion: string;
}

export interface ConsultationRecord {
  paciente:  Paciente;
  consultas: Consulta[];
}

export interface Consulta {
  id:             string;
  fecha_consulta: string;
  area_asignada:  string;
  area_origen:    string;
  servicios:      string;
}

export interface Paciente {
  codigo_sga:   string;
  dni:          string;
  nombres:      string;
  apellidos:    string;
  tipo_persona: string;
}

export interface PaymentConcept {
  id: number;
  cod_recibo: string;
  concepto_pagado: string;
  importe_pagado: string;
  precio_unit: string;
  fecha_pago: string;
}
