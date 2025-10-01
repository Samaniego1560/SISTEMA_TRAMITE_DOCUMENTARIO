export interface VaccinesType {
  id: string;
  nombre: string;
  dosis_minimas: number;
  dosis_maximas: number | null;
  intervalo_entre_dosis: number;
  estado: boolean;
}
