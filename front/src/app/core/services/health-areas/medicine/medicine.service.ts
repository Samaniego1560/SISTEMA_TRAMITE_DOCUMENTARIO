import {inject, Injectable} from "@angular/core";
import {EnvServiceFactory} from "../../env/env.service.provider";
import {HttpClient} from "@angular/common/http";
import {Response} from "../../../models/global.model";
import {ConsultationRecord} from "../../../models/areas/areas.model";

@Injectable()
export class MedicineService {
  private urlBase = EnvServiceFactory().API_GO;
  private _api = this.urlBase + '/v1/area_medica';

  private _http = inject(HttpClient);

  public getConsultations(){
    return this._http.get<Response<ConsultationRecord[]>>(this._api + '/consultas_medicina');
  }

  public getConsultation(consultationID: string){
    return this._http.get<Response<any>>(this._api + '/consulta_medicina/' + consultationID);
  }

  public getConsultationsByDni(dni: string){
    return this._http.get<Response<ConsultationRecord>>(this._api + '/consultas_medicina/paciente/dni/' + dni);
  }

  public getReportDetail(params: Record<string, string>){
    const queryParams = new URLSearchParams(params).toString();
    return this._http.get<Response<string>>(this._api + '/reporte/medicina?' + queryParams);
  }

  public getReportGeneral(params: Record<string, string>){
    const queryParams = new URLSearchParams(params).toString();
    return this._http.get<Response<string>>(this._api + '/reporte/consultas?' + queryParams);
  }
}
