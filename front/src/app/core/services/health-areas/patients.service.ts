import {inject, Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {EnvServiceFactory} from "../env/env.service.provider";
import {Patient, PatientBase, PatientRequest, RequestSearchPatient} from "../../models/areas/patient.model";
import {Response} from "../../models/global.model";
import {InfoPerson} from "../../models/areas/person.model";

@Injectable({
  providedIn: 'root'
})
export class PatientsService {
  private urlBase = EnvServiceFactory().API_GO;
  private urlBase2 = EnvServiceFactory().API_DBU;
  private apiPatients = this.urlBase + '/v1/area_medica';

  private _http = inject(HttpClient);

  public getPatients(request: RequestSearchPatient, pagination?: string) {
    const baseUrl = this.apiPatients + '/pacientes/get';
    const url = baseUrl + (pagination ? pagination : "?limit=20&offset=0");

    return this._http.post<Response<Patient[]>>(url, request);
  }

  public getPatient(dni: string) {
    return this._http.get<InfoPerson>(this.urlBase2 + '/DatosAlumnoAcademico/show/' + dni);
  }

  public createPatient(patient: PatientRequest) {
    return this._http.post<Response<Patient>>(this.apiPatients + '/paciente', patient);
  }

  public updatePatient(patient: PatientRequest) {
    return this._http.put<Response<Patient>>(this.apiPatients + '/paciente', patient);
  }

  public getPatientByID(id: string) {
    return this._http.get<Response<Patient>>(this.apiPatients + '/paciente/' + id);
  }

  public getPatientByDni(id: string) {
    return this._http.get<Response<Patient>>(this.apiPatients + '/paciente/dni/' + id);
  }

  public removePatient(patientID: string) {
    return this._http.delete<Response>(this.apiPatients + '/paciente' + '/' + patientID);
  }
}
