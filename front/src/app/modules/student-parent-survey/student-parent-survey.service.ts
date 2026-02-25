import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { EnvServiceFactory } from '../../core/services/env/env.service.provider';

@Injectable({ providedIn: 'root' })
export class StudentParentSurveyService {
  private baseUrl = this.buildBaseUrl();

  constructor(private http: HttpClient) {}

  private buildBaseUrl(): string {
    const env = EnvServiceFactory();
    const apiBase = (env.API_GO || '').replace(/\/+$/, '');
    return `${apiBase}/api/survey/student-parent`;
  }

  verify(dni: string) {
    return this.http.post<any>(`${this.baseUrl}/verify`, { dni });
  }

  submit(payload: any) {
    return this.http.post<any>(this.baseUrl, payload);
  }

  exportExcel() {
    return this.http.get(`${this.baseUrl}/export`, { responseType: 'blob' });
  }

  stats() {
    return this.http.get<any>(`${this.baseUrl}/stats`);
  }
}
