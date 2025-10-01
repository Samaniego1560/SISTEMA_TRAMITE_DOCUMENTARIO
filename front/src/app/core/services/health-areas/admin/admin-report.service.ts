import {inject, Injectable} from "@angular/core";
import {Response} from "../../../models/global.model";
import {EnvServiceFactory} from "../../env/env.service.provider";
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class AdminReportService {
  private urlBase = EnvServiceFactory().API_GO;
  private _apiUrl = this.urlBase + '/v1/area_medica';

  private _http = inject(HttpClient);

  public getReport(params: Record<string, string>){
    const queryParams = new URLSearchParams(params).toString();
    return this._http.get<Response<string>>(this._apiUrl + '/reporte/admin?' + queryParams);
  }
}
