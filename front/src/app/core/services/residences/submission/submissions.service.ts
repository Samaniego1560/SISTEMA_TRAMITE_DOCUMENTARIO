import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable,of} from "rxjs";
import {EnvServiceFactory} from "../../env/env.service.provider";
import { Residence, ResidenceRequest, ResidenceResponse } from '../../../models/residence';
import { Base64 } from '../../../utils/statics/base64';

@Injectable({
  providedIn: 'root'
})
export class SubmissionService {

  private urlBase = EnvServiceFactory().API_GO;
  private urlBaseAnnouncement = this.urlBase + '/v1/convocatorias';

  private base_64: string =  Base64;
  constructor(
    private http: HttpClient
  ) {
  }

  public getStudentsByAnnouncement(idAnnouncement: number, filters: any): Observable<ResidenceResponse> {
    return this.http.get<ResidenceResponse>(this.urlBaseAnnouncement + '/' + idAnnouncement +'/alumnos-aceptados', {params: filters});
  }

  public exportReportResidence(idAnnouncement: number): Observable<any> {
    let base = this.base_64;
    return this.http.get(this.urlBaseAnnouncement + '/' + idAnnouncement +'/reporte-residencias');
  }
}
