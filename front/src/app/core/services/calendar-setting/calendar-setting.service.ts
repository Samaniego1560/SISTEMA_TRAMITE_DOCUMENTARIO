import { Injectable } from '@angular/core';
import {EnvServiceFactory} from "../env/env.service.provider";
import {IBodyRequest} from "../../models/requests";
import {Observable} from "rxjs";
import {IResponse} from "../../models/response";
import {HttpClient} from "@angular/common/http";
import {ICalendarSettings} from "../../models/calendar-setting";

@Injectable({
  providedIn: 'root'
})
export class CalendarSettingService {
  private urlBase = EnvServiceFactory().API_DBU;
  private urlCalendarSetting = this.urlBase + '/calendar/setting';

  constructor(private http: HttpClient) { }
  public updateCalendarSetting(data: ICalendarSettings[]): Observable<IResponse> {
    return this.http.post<IResponse>(this.urlCalendarSetting + '/update', data);
  }

  public getCalendarSetting(): Observable<IResponse> {
    return this.http.get<IResponse>(this.urlCalendarSetting + '/list');
  }
}
