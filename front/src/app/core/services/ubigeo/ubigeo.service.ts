import { Injectable } from '@angular/core';
import { EnvServiceFactory } from '../env/env.service.provider';
import {Observable} from "rxjs";
import { IResponse } from '../../models/response';
import { IAnnouncement } from '../../models/announcement';
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class UbigeoService {

  private urlBase = EnvServiceFactory().API_DBU;
  private urlBaseUbigeo = this.urlBase + '/ubigeo';
  constructor(private http: HttpClient) { }

  public getUbigeosDepartaments(): Observable<any> {

    return this.http.get<any>('https://bienestar.unas.edu.pe/backend/ubigeo/departaments');
  }

  public getUbigeosProvinces(): Observable<IResponse<IAnnouncement>> {
    return this.http.get<IResponse<IAnnouncement>>(this.urlBaseUbigeo + '/provinces');
  }

  public getUbigeosDistricts(): Observable<IResponse<IAnnouncement>> {
    return this.http.get<IResponse<IAnnouncement>>(this.urlBaseUbigeo + '/districts');
  }
}
