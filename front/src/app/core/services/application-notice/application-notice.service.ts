import { Injectable } from '@angular/core';
import {EnvServiceFactory} from "../env/env.service.provider";
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {IResponse} from "../../models/response";
import {
  ApplicationNoticeModel
} from "../../../modules/home/pages/application-notice/application-notice-models/application-notice-model";

@Injectable({
  providedIn: 'root'
})
export class ApplicationNoticeService {
  private urlBase = EnvServiceFactory().API_DBU;
  private urlApplicationNotice = this.urlBase + '/application/notice';

  constructor(private http: HttpClient) { }
  public updateApplicationNotice(data: ApplicationNoticeModel[]): Observable<IResponse> {
    return this.http.post<IResponse>(this.urlApplicationNotice + '/update', data);
  }

  public getApplicationNotice(): Observable<IResponse> {
    return this.http.get<IResponse>(this.urlApplicationNotice + '/list');
  }
}
