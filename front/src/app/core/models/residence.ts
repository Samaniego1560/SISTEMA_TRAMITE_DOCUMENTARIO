export interface ResidenceRequest {
    id: string,
    name: string,
    gender: string,
    description: string
    address: string,
    status:string,
    floors?: Floors[]
}

export interface Residence {
    id: string,
    name: string,
    gender: string,
    description: string
    address: string,
    status:string,
    floors: Floors[],
    configuration?: ConfigurationRules,
    submissions: Submission[]
}

export interface ResidenceResponse <T = any>{
    msg: string,
    data: any,
    type: string,
    code: number,
    error: boolean
}
export interface Floors {
    floor: number,
    space_rooms?: number,
    q_rooms?: number,
    rooms: Room[]
}
export interface Room {
    id: string,
    number: number,
    capacity?: number,
    status:string,
    floor?: number|any,
}

export interface Dataset {
    name: string,
    value: string
}

export interface DatasetResidence {
    name: string,
    value: number
}

export interface ConfigurationRules {
    id?: string,
    percentage_fcea: number,
    percentage_engineering: number,
    minimum_grade_fcea: number,
    minimum_grade_engineering: number
}

export interface Submission {
    id: number,
    name: string,
    start: string,
    state: boolean
}

export interface Student {
    admission_date: string,
    code: string,
    full_name: string,
    id: number,
    professional_school: string,
    residence: string,
    room: RoomStudent,
    number_identification: string,
    sex: string
}

export interface RoomStudent {
    id: string,
    number: number
}

// export interface StudentInfo {
//     student: Student,
//     sanctions: Sanctions[],
//     assigned_goods: AssigedThings[],
//     room_mates: RoomMates[]
// }

export interface Sanctions {
    description: string,
    date: string,
    category: string|any
}

export interface AssigedThings {
    code: string
}

export interface RoomMates {
    code: string,
    full_name: string
}

export interface StudentInfo {
    student: Student,
    sanctions: Sanctions[],
    assigned_goods: AssigedThings[],
    room_mates: RoomMates[]
}

export interface Submission {
    nombre: string,
    fecha_inicio: string,
    fecha_fin: string,
    activo: boolean,
    secciones: any[],
    convocatoria_servicio: any[]
}

export interface StudentDebtInfo {
   nombre_completo: string,
   nombres: string,
   apellido_paterno: string,
   apellido_materno: string,
   dni: string,
   historial_pagos: PaymentsHistory[]
}

export interface PaymentsHistory {
    id : string,
    cod_deuda : string,
    monto_deuda : string,
    amortizado : string,
    monto_restante_deuda : string,
    estado_deuda : string,
    desc_deuda : string,
    concept_id : string,
    client_id : string,
    created_at : string
    updated_at : string,
    debt_payments: DebtsPayments[]
}

export interface DebtsPayments {
    id : string,
    monto_pagado_deuda : string,
    debt_id : string,
    receipt_detail_id : string
    created_at : string,
    updated_at : string,
    receipt_detail: ReceiptDetail
}

export interface ReceiptDetail {
    id : string,
    concepto_pagado : string,
    importe_pagado : string,
    receipt_id : string,
    receipt: Receipt
}

export interface Receipt {
    id : string,
    cod_recibo : string,
    tipo_cliente : string,
    obs_recibo : string,
    created_at : string
}



