export interface OdontogramRequest {
  consulta_odontologia: ConsultaOdontologia;
  examen_bucal:         OralExam;
  revision_odontograma: RevisionOdontograma;
}

export interface ConsultaOdontologia {
  id:          string;
  paciente_id: string;
}

export interface OralExam {
  id:                        string;
  capacidad_masticatoria:    string;
  encias:                    string;
  carles_dentales:           string;
  edentulismo_parcial_total: string;
  portador_protesis_dental:  string;
  estado_higiene_bucal:      string;
  urgencia_tratamiento:      string;
  fluorizacion:              string;
  destartraje:               string;
  comentarios:               string;
}

export interface RevisionOdontograma {
  id:                   string;
  caries:               string;
  erupcionado:          string;
  perdido:              string;
  costo:                string;
  fecha_pago:           string;
  cpod:                 string;
  urgencia_tratamiento: string;
  mes:                  string;
  comentarios:          string;
}
