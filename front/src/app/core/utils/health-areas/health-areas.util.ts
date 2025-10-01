import {Consulta} from "../../models/areas/areas.model";

export class HealthAreasUtil {
  static sortConsultations(consultations: Consulta[]) {
    consultations = [...consultations]
    return consultations.sort((consultaA, consultaB) => {
      const timeB = new Date(consultaB.fecha_consulta).getTime()
      const timeA = new Date(consultaA.fecha_consulta).getTime()
      return timeB - timeA
    })
  }
}
