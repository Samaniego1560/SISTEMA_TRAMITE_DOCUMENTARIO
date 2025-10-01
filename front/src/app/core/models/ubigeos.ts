export interface Department {
    id: string;
    name: string,
}

export interface Province {
    id: string;
    departament_id: string,
    name: string,
}

export interface District {
    id: string;
    province_id: string,
    name: string,
}
