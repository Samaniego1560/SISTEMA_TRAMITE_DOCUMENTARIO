import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable,of} from "rxjs";
import {EnvServiceFactory} from "../../env/env.service.provider";
import { Residence, ResidenceRequest, ResidenceResponse } from '../../../models/residence';
import { Base64 } from '../../../utils/statics/base64';
import { StudentDebtInfo } from '../../../models/residence';

@Injectable({
  providedIn: 'root'
})
export class ResidencesService {

  private urlBaseGo = EnvServiceFactory().API_GO;
  private urlCaja = EnvServiceFactory().API_CAJA;
  private urlResidences = this.urlBaseGo + '/v1/residencias';
  private urlAssignation = this.urlBaseGo + '/v1/automatizacion';
  private urlDebts = this.urlCaja + '/api/cajapaymentsunasdbu';

  private base_64: string =  Base64;
  constructor(
    private http: HttpClient
  ) {
  }

  public getListResidences(): Observable<ResidenceResponse<Residence[]>>{
    return this.http.get<ResidenceResponse<Residence[]>>(this.urlResidences );
  }
  public getResidence(id: string): Observable<ResidenceResponse<Residence>> {
    return this.http.get<ResidenceResponse<Residence>>(this.urlResidences + '/' + id);
  }

  public deleteResidence(id: string): Observable<ResidenceResponse> {
    return this.http.delete<ResidenceResponse>(this.urlResidences + '/' + id);
  }

  public createResidence(data: ResidenceRequest): Observable<ResidenceResponse> {
    return this.http.post<ResidenceResponse>(this.urlResidences, data);
  }

  public updateResidence(data: ResidenceRequest, id: string): Observable<ResidenceResponse> {
    data.id = id;
    delete data.floors;
    return this.http.put<ResidenceResponse>(this.urlResidences, data);
  }

  public getStudentsResidence(idResidence: number|any, filters: any): Observable<ResidenceResponse> {
    return this.http.get<ResidenceResponse>(this.urlResidences + '/' + idResidence+ '/alumnos', {params: filters}, );
  }

  public getRoomsResidence(idResidence: number|any, filters: any): Observable<ResidenceResponse> {

    return this.http.get<ResidenceResponse>(this.urlResidences + '/' + idResidence+ '/cuartos', {params: filters}, );
  }

  public updateRulesResidence(idResidence: number|any, data: any): Observable<ResidenceResponse> {
    return this.http.put<ResidenceResponse>(this.urlResidences + '/' + idResidence + '/configuracion', data);
  }

  public sortResidenceAssignation(idResidence: any): Observable<ResidenceResponse> {
    return this.http.post<ResidenceResponse>(this.urlAssignation + '/asignacion-cuartos', { residence_id: idResidence });
  }

  public validateStudentDebts(dni: string): Observable<StudentDebtInfo> {
      return this.http.get<StudentDebtInfo>(this.urlDebts+ '/' + dni);
  }
}
