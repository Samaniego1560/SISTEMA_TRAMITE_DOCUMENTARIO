export interface IAnnouncement {
  id?: number;
  user_id?: number;
  fecha_inicio: string;
  fecha_fin: string;
  nombre: string;
  convocatoria_servicio: IService[];
  secciones: ISection[];
  activo: boolean;
  credito_minimo: number;
  nota_minima: number;
  recoverable_detail_requests?: IRecoverableDetailRequests[];
}

export interface IService {
  servicio_id: number;
  cantidad: number;
  id?: number;
}

export interface ISection {
  id?: number;
  descripcion: string;
  type: string;
  requisitos: IRequirement[];
}

export interface IRequirement {
  nombre: string;
  descripcion: string;
  url_guia: string;
  tipo_requisito_id: number;
  opciones?: string;
  default?: string;
  activo: boolean;
  id?: number;
  url_plantilla?: string;
  type_input?: string;
  text_color?: string;
  text_type?: string;
  text_size?: string;
  is_recoverable?: boolean;
  is_dependent?: boolean,
  field_dependent?: string,
  value_dependent?: string,
  show_dependent?: boolean,
}

export interface IRecoverableDetailRequests {
  respuesta_formulario: string;
  requisito_id: number;
  solicitud_id: number;
  order: number;
}

//1= documento  2=Imagen 3=Formulario
