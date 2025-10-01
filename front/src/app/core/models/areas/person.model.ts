export interface InfoPerson {
  msg:     string;
  detalle: DetailPerson;
}

export interface DetailPerson {
  codsem:         string;
  codalumno:      number;
  tdocumento:     string;
  appaterno:      string;
  apmaterno:      string;
  nombre:         string;
  sexo:           string;
  direccion:      string;
  fecnac:         string;
  mod_ingreso:    string;
  nombrecolegio:  string;
  ubigeo:         string;
  nomesp:         string;
  nomfac:         string;
  telcelular:     string;
  tel_ref:        string;
  email:          string;
  emailinst:      string;
  nume_sem_cur:   number;
  est_mat_act:    string;
  credmat:        number;
  pps:            string;
  ppa:            string;
  artincurso:     string;
  artpermanencia: string;
  tca:            number;
}
