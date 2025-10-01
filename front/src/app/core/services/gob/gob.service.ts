import { Injectable } from '@angular/core';
import {Observable} from "rxjs";
import {IResponse} from "../../models/response";
import {IAnnouncement} from "../../models/announcement";
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {EnvServiceFactory} from "../env/env.service.provider";

@Injectable({
  providedIn: 'root'
})
export class GobService {
  private urlGob = "http://localhost:8081";

  constructor(private http: HttpClient) { }

  public getCaptchaValue(): Observable<any> {
    return this.http.get(this.urlGob+'/code');
  }

  validaDNI(dni: string, code: string): Observable<any> {
    return this.http.post(this.urlGob +'/send-data', {dni:dni, code: code});
  }
}
