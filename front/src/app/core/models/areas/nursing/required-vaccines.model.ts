export interface NextVaccinesPatient {
  required: RequiredVaccines[] | null,
  notRequired?: RequiredVaccines[] | null,
}

export interface RequiredVaccines {
  fecha_requerida?: string;
  nombre: string;
  requerido: boolean;
}
