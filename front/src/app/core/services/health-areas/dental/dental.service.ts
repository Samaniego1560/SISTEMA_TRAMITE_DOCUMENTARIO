import {inject, Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {EnvServiceFactory} from "../../env/env.service.provider";
import {Response} from "../../../models/global.model";
import {DentistryConsultationRequest, SearchPaymentConcept} from "../../../models/areas/dentistry.model";
import {ConsultationRecord, PaymentConcept} from "../../../models/areas/areas.model";

@Injectable({
  providedIn: 'root'
})
export class DentalService {
  private urlBase = EnvServiceFactory().API_GO;
  private _apiUrl = this.urlBase + '/v1/area_medica';

  private _http = inject(HttpClient);

  public getConsultations(){
    return this._http.get<Response<ConsultationRecord[]>>(this._apiUrl + '/consultas_odontologia');
  }

  public createConsulting(consultation: DentistryConsultationRequest){
    return this._http.post<Response<DentistryConsultationRequest>>(this._apiUrl + '/consulta_odontologia', consultation);
  }

  public getConsultation(consultationID: string){
    return this._http.get<Response<DentistryConsultationRequest>>(this._apiUrl + '/consulta_odontologia/' + consultationID);
  }

  public getConsultationsByDni(dni: string){
    return this._http.get<Response<ConsultationRecord>>(this._apiUrl + '/consultas_odontologia/paciente/dni/' + dni);
  }

  public updateConsulting(request: DentistryConsultationRequest){
    return this._http.put<Response<any>>(this._apiUrl + '/consulta_odontologia', request);
  }

  public deleteConsulting(consultationID: string){
    return this._http.delete<Response>(this._apiUrl + '/consulta_odontologia/' + consultationID);
  }

  public getReportDetail(params: Record<string, string>){
    const queryParams = new URLSearchParams(params).toString();
    return this._http.get<Response<string>>(this._apiUrl + '/reporte/odontologia?' + queryParams);
  }

  public getReportGeneral(params: Record<string, string>){
    const queryParams = new URLSearchParams(params).toString();
    return this._http.get<Response<string>>(this._apiUrl + '/reporte/consultas?' + queryParams);
  }

  public getPayment(params: SearchPaymentConcept){
    return this._http.post<Response<PaymentConcept | null>>(this._apiUrl + '/odontologia/pagos/procedimiento/buscar', params);
  }
}
