import { Injectable } from '@angular/core';
import {EnvServiceFactory} from "../env/env.service.provider";
import {Observable} from "rxjs";
import {IResponse} from "../../models/response";
import {ThesisPracticing} from "../../models/ThesisPracticing";
import {HttpClient} from "@angular/common/http";
import {StudentFileReport} from "../../models/report";

@Injectable({
  providedIn: 'root'
})
export class ReportService {

  private urlBase = EnvServiceFactory().API_DBU;
  private urlReport = this.urlBase + '/report';

  constructor(private http: HttpClient) { }

  public getReportProfile(dni:string): Observable<IResponse<StudentFileReport>> {
    const urlReportProfile = this.urlReport + '/profile/' + dni;
    return this.http.get<IResponse<StudentFileReport>>(urlReportProfile);
  }

  public getReportProfileDownload(dni:string): Observable<any> {
    const urlReportProfile = this.urlReport + '/profile/download/' + dni;
    return this.http.get(urlReportProfile, { responseType: 'blob' });
  }
}
