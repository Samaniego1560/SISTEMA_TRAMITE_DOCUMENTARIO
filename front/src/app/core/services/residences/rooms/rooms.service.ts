import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable,of} from "rxjs";
import {EnvServiceFactory} from "../../env/env.service.provider";
import { ResidenceResponse } from '../../../models/residence';
import { Base64 } from '../../../utils/statics/base64';

@Injectable({
  providedIn: 'root'
})
export class RoomsService {

  private urlBaseGo = EnvServiceFactory().API_GO;
  private urlRooms = this.urlBaseGo + '/v1/residencias/cuartos';

  constructor(
    private http: HttpClient
  ) {
  }

  public updateRoom(data : any): Observable<ResidenceResponse> {
    return this.http.put<ResidenceResponse>(this.urlRooms , data); 
  }

  public assignRoomStudent(data : any): Observable<ResidenceResponse> {
    return this.http.post<ResidenceResponse>(this.urlRooms + '/'+ data.room_id + '/asignar', {"submission_id":data.submission_id, "student_id":data.student_id}); 
  }

  public deleteStudentAssignation(data : any): Observable<ResidenceResponse> {
    let room_id = data.room_id;
    delete data.room_id;
    return this.http.delete<ResidenceResponse>(this.urlRooms + '/' + room_id+'/eliminar-asignacion', {body:data}); 
  }
}
