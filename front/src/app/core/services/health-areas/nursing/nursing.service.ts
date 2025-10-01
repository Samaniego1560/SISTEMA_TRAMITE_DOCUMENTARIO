import {inject, Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {EnvServiceFactory} from "../../env/env.service.provider";
import {Response} from "../../../models/global.model";
import {ConsultationRecord} from "../../../models/areas/areas.model";
import {NextVaccinesPatient, RequiredVaccines} from "../../../models/areas/nursing/required-vaccines.model";
import {NursingConsultationRequest, RequestUpdateAreaAssigned} from "../../../models/areas/nursing/request.model";
import {VaccinesType} from "../../../models/areas/nursing/vaccines.model";
import {Vacuna} from "../../../models/areas/nursing/nursing.model";

@Injectable()
export class NursingService {
  private urlBase = EnvServiceFactory().API_GO;
  private apiPatients = this.urlBase + '/v1/area_medica';

  private _http = inject(HttpClient);

  public getConsultations(){
    return this._http.get<Response<ConsultationRecord[]>>(this.apiPatients + '/consultas_enfermeria');
  }

  public getConsultation(consultationID: string){
    return this._http.get<Response<NursingConsultationRequest>>(this.apiPatients + '/consulta_enfermeria/' + consultationID);
  }

  public getConsultationsByDni(dni: string){
    return this._http.get<Response<ConsultationRecord>>(this.apiPatients + '/consultas_enfermeria/paciente/dni/' + dni);
  }

  public createConsulting(request: NursingConsultationRequest){
    return this._http.post<Response<any>>(this.apiPatients + '/consulta_enfermeria', request);
  }

  public updateConsulting(request: NursingConsultationRequest){
    return this._http.put<Response<any>>(this.apiPatients + '/consulta_enfermeria', request);
  }

  public deleteConsultation(consultationID: string){
    return this._http.delete<Response>(this.apiPatients + '/consulta_enfermeria' + '/' + consultationID);
  }

  public updateAssignedArea(request: RequestUpdateAreaAssigned){
    return this._http.put<Response>(this.apiPatients + '/consulta_enfermeria_asignacion', request);
  }

  public getRequiredVaccines(patientID: string){
    return this._http.get<Response<NextVaccinesPatient>>(this.apiPatients + '/consultas_enfermeria/vacunas_requeridas/paciente/' + patientID);
  }

  public getVaccines(){
    return this._http.get<Response<VaccinesType[]>>(this.apiPatients + '/consultas_enfermeria/vacunas');
  }

  public getReportDetail(params: Record<string, string>){
    const queryParams = new URLSearchParams(params).toString();
    return this._http.get<Response<string>>(this.apiPatients + '/reporte/enfermeria?' + queryParams);
  }

  public getReportGeneral(params: Record<string, string>){
    const queryParams = new URLSearchParams(params).toString();
    return this._http.get<Response<string>>(this.apiPatients + '/reporte/enfermeria/rua?' + queryParams);
  }

  public getVaccinesPatientByDni(dni: string){
    return this._http.get<Response<Vacuna[]>>(this.apiPatients + '/consultas_enfermeria/vacunas/paciente/' + dni);
  }
}
