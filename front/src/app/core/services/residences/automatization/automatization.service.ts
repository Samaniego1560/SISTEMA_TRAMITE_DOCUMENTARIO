import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable,of} from "rxjs";
import {EnvServiceFactory} from "../../env/env.service.provider";
import { Residence, ResidenceRequest, ResidenceResponse } from '../../../models/residence';
import { Base64 } from '../../../utils/statics/base64';

@Injectable({
  providedIn: 'root'
})
export class AutomatizationService {

  private urlBaseGo = EnvServiceFactory().API_GO;
  private urlResidences = this.urlBaseGo + '/v1/residencias';
  private urlBaseAnnouncementGo = this.urlBaseGo + '/v1/convocatorias';

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
    return this.http.get<ResidenceResponse>(this.urlResidences + '/' + idResidence+ '/alumnos', {params: filters});
  }

  public getRoomsResidence(idResidence: number|any, filters: any): Observable<ResidenceResponse> {
    let submission_id = filters.submission_id;
    delete filters.submission_id;
    return this.http.get<ResidenceResponse>(this.urlResidences + '/' + idResidence+ '/cuartos', {params: filters, headers:{submission_id: submission_id}}, );
  }

  public sortResidenceAssignation(idResidence: any): Observable<ResidenceResponse> {
    return this.http.get<ResidenceResponse>(this.urlResidences + '/' + idResidence + '/asignaciones');
  }
  public assignRoomStudent(idRoom: any, idStudent: any): Observable<ResidenceResponse> {
    return this.http.post<ResidenceResponse>(this.urlResidences + '/asignar', {idRoom, idStudent}); 
  }
  public exportReportResidence(): Observable<any> {
    let base = this.base_64;
    //return this.http.get(this.urlResidences + '/report');
    return of({ fileBase64: base, success: true });
  }

  public getStudentsByAnnouncement(idAnnouncement: number, filters: any): Observable<ResidenceResponse> {
    return this.http.get<ResidenceResponse>(this.urlBaseAnnouncementGo + '/' + idAnnouncement +'/alumnos-aceptados', {params: filters});
  }

  public getAnnouncementResidence(): Observable<any> {
    return this.http.get<any>(this.urlBaseGo + '/api/v1/convocatorias');
  }
}
