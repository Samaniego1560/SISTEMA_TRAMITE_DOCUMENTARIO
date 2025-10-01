export interface InformationPayment {
  nombre_completo: string;
  nombres: string;
  apellido_paterno: string;
  apellido_materno: string;
  dni: string;
  pagos_realizados: PagosRealizados[];
}

export interface PagosRealizados {
  id: number;
  cod_recibo: string;
  concepto_pagado: string;
  importe_pagado: string;
  precio_unit: string;
  fecha_pago: string;
}
