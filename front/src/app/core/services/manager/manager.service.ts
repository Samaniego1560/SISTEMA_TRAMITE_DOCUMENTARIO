import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {IResponse} from "../../models/response";
import {IAnnouncement} from "../../models/announcement";
import {
  Estudiante,
  IBodyRequest,
  IFileRequest,
  IRequest,
  IResponseFile,
  IStatusRequest,
  IUpdateService
} from "../../models/requests";
import {IDebtsStudent, IValidationUser} from "../../models/user";
import {EnvServiceFactory} from "../env/env.service.provider";
import {IStatistics} from "../../models/statistics";
import { ThesisPracticing } from '../../models/ThesisPracticing';
import { Departamento, Distrito, IVisitaGeneral, Provincia } from '../../models/visita-general';
import { ExamenToxicologico, CreateExamenToxicologicoRequest, UpdateExamenToxicologicoRequest } from '../../models/exam_toxicologico';
import { CreateVisitaResidenteRequest, ExportDataExcel, LugarProcedencia, ReportePorEscuelaProfesional, ReporteVisitaResidente, UpdateVisitaResidenteRequest, VisitaResidente, VisitasPorDepartamento } from '../../models/visita-residente';
import { Residence, ResidenceRequest, ResidenceResponse } from '../../models/residence';
import {InformationPayment} from "../../models/trasury.model";
import {Response} from "../../models/global.model";
import { Base64 } from '../../utils/statics/base64';
import {forkJoin, Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ManagerService {

  private urlBase = EnvServiceFactory().API_DBU;
  private urlCaja = EnvServiceFactory().API_CAJA;
  private urlPayments = EnvServiceFactory().API_CAJA;
  private urlBaseGo = EnvServiceFactory().API_GO;
  private urlBaseAnnouncement = this.urlBase + '/convocatoria';
  private urlRequests = this.urlBase + '/solicitudes';
  private urlRequest = this.urlBase + '/solicitud';
  private urlCreate = this.urlBaseAnnouncement + '/create';
  private urlThesisPracticing = this.urlBaseAnnouncement + '/thesis-practicings';
  private urlPsicodagogia = this.urlBaseGo + '/v1/psicopedagogia';
  private urlVisitas = this.urlBaseGo + '/v1/visita-general';
  private urlVisitasResidente = this.urlBaseGo + '/v1/visita-residente';
  private urlExamenToxicologico = this.urlBaseGo + '/v1/area_medica/examen_toxicologico';
  private urlResidences = this.urlBaseGo + '/v1/residencias';
  private urlBaseAnnouncementGo = this.urlBaseGo + '/v1/convocatorias';
  private urlBaseAreaMedicaGo = this.urlBaseGo + '/v1/area_medica';

  private base_64: string =  Base64;
  constructor(
    private http: HttpClient
  ) {
  }

  public getAnnouncement(): Observable<IResponse<IAnnouncement[]>> {
    return this.http.get<IResponse<IAnnouncement[]>>(this.urlBaseAnnouncement);
  }

  public createAnnouncement(announcement: IAnnouncement): Observable<IResponse> {
    return this.http.post<IResponse>(this.urlCreate, announcement);
  }

  public getRequests(id: number): Observable<IResponse<IRequest[]>> {
    return this.http.get<IResponse<IRequest[]>>(this.urlRequests + '/' + id);
  }

  public getCurrentAnnouncement(): Observable<IResponse<IAnnouncement>> {
    return this.http.get<IResponse<IAnnouncement>>(this.urlBaseAnnouncement + '/vigente-convocatoria');
  }

  public getDataStudent(code: string): Observable<IResponse<IAnnouncement>> {
    return this.http.get<IResponse<IAnnouncement>>(this.urlRequest + '/alumno/' + code);
  }

  public getDataAnnouncementAndStudent(announcementId: number, code: string): Observable<IResponse<IAnnouncement>> {
    return this.http.get<IResponse<IAnnouncement>>(this.urlRequest + '/convocatoria/alumno/' + announcementId + '/' + code);
  }

  public uploadRequestFile(file: IFileRequest): Observable<IResponse<IResponseFile>> {
    return this.http.post<IResponse<IResponseFile>>(this.urlRequest + '/uploadDocument', file);
  }

  public createRequest(data: IBodyRequest): Observable<IResponse> {
    return this.http.post<IResponse>(this.urlRequest + '/create', data);
  }

  public validateStudent(data: IValidationUser): Observable<IResponse> {
    return this.http.post<IResponse>(this.urlRequest + '/validacion', data);
  }

  public validateStudentThesisOrPracticingOrStudent(data: IValidationUser): Observable<IResponse> {
    return this.http.post<IResponse>(this.urlRequest + '/validacion/thesis_practicing_student', data);
  }

  public validateStudentPermission(data: IValidationUser): Observable<IResponse> {
    return this.http.post<IResponse>(this.urlRequest + '/validacion/permision', data);
  }

  public validateSignAreas(dni: string): Observable<IResponse> {
    return this.http.get<IResponse>(this.urlBaseAreaMedicaGo + '/firmas/paciente/dni/' + dni);
  }

  public getStudentRequest(code: number): Observable<IResponse<IRequest>> {
    return this.http.get<IResponse<IRequest>>(this.urlRequest + '/show/' + code);
  }

  public updateStatusService(data: IUpdateService): Observable<IResponse> {
    return this.http.put<IResponse>(this.urlRequest + '/servicio', data);
  }

  public getNumberOfVacancies(id: number): Observable<IResponse> {
    return this.http.get<IResponse>(this.urlRequest + '/servicio/cantidad/vacantes/' + id);
  }

  public validateStudentDebts(code: string): Observable<IDebtsStudent[]> {
    return this.http.get<IDebtsStudent[]>(this.urlCaja + '/api/report/hasdebt/' + code);
  }

  public validateStudentPayment(dni: string): Observable<InformationPayment> {
    return this.http.get<InformationPayment>(this.urlPayments + '/api/cajapaymentsunasdbu-normal-payments/' + dni);
  }

  public exportRequest(): Observable<Blob> {
    return this.http.get(this.urlRequest + '/export', {responseType: 'blob'});
  }

  public exportRequestById(id: number): Observable<Blob> {
    return this.http.get(this.urlRequest + '/export/convocatoria/' + id, {responseType: 'blob'});
  }

  public exportDataRequestById(id: number): Observable<any> {
    return this.http.get(this.urlRequest + '/export/data/convocatoria/' + id);
  }

  public getRequestStatus(dni: string): Observable<IResponse<Estudiante>> {
    return this.http.get<IResponse<Estudiante>>(this.urlRequest + '/servicioSolicitado/' + dni);
  }

  public getStatistics(code: number): Observable<IResponse<IStatistics>> {
    return this.http.get<IResponse<IStatistics>>(this.urlBaseAnnouncement + '/reporte/' + code);
  }

  public getStatisticsAndService(code: number, typeService: number): Observable<IResponse<IStatistics>> {
    return this.http.get<IResponse<IStatistics>>(this.urlBaseAnnouncement + '/reporte/' + code+ '/' + typeService);
  }

  // register thesis o praticacing

  public createThesisOrPracticing(data: ThesisPracticing): Observable<IResponse> {
    return this.http.post<IResponse>(this.urlThesisPracticing + '/register', data);
  }

  public updateThesisOrPracticing(data: ThesisPracticing, id: number): Observable<IResponse> {
    return this.http.put<IResponse>(this.urlThesisPracticing + '/update/' + id, data);
  }

  public deleteThesisOrPracticing(id: number): Observable<IResponse> {
    return this.http.delete<IResponse>(this.urlThesisPracticing + '/destroy/' + id);
  }

  public getThesisPracticing(): Observable<IResponse<ThesisPracticing[]>> {
    return this.http.get<IResponse<ThesisPracticing[]>>(this.urlThesisPracticing );
  }


  //LuisReynaga: route psicopedagogia
  public getStudentQuestions(numDocumento: string, typeParticipante: string): Observable<[IResponse<any>, IResponse<any>]> {
    const urlStudent = `${this.urlPsicodagogia}/estudiantes/${numDocumento}?tipoParticipante=${typeParticipante}`;
    const urlQuestions = `${this.urlPsicodagogia}/preguntas?page=1&pageSize=100`;

    return forkJoin([
      this.http.get<IResponse<any>>(urlStudent),
      this.http.get<IResponse<any>>(urlQuestions)
    ]);
  }

  public getDatosAlumnoAcademico(dni: string): Observable<IResponse<any>> {
    const url = `${this.urlBase}/DatosAlumnoAcademico/show/${dni}`;
    return this.http.get<IResponse<any>>(url)
  }

  public saveAnswers(data: any) {
    this.urlRequest = `${this.urlPsicodagogia}/historial`;
    return this.http.post<IResponse>(this.urlRequest, data);
  }

  public getResultQuestionnaire(page: number): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/participantes?page=${page}`;
    return this.http.get<IResponse<any>>(url)
  }

  public hasActiveSRQ(): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/encuestas/active-srq`;
    return this.http.get<IResponse<any>>(url)
  }

  public getResultAnswers(patientId: number, numeroAtencion: number, page: number = 1, pageSize: number = 60): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/respuestas/participante`;
    const body = {
      idParticipante: patientId,
      numeroAtencion: numeroAtencion,
      page: page,
      pageSize: pageSize
    };

    return this.http.post<IResponse<any>>(url, body);
  }


  public updateStatusEvaluationParticipante(dataUpdate: any, patientId: number): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/participantes/${patientId}`;
    const body = { estado_evaluacion: dataUpdate.statusEvaluation, notes: dataUpdate.notes };
    return this.http.put<IResponse<any>>(url, body);
  }

  public getHaveSurvey(numDocumento: string): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/questions/have-survey/${numDocumento}`;
    return this.http.get<IResponse<any>>(url)
  }

  public searchQuestionnaires(filters: { [key: string]: string }): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/participantes/search-result`;
    return this.http.post<IResponse<any>>(url, filters);
  }

  public getParticipantByDocument(documento: string, tipoParticipante: string): Observable<IResponse<any>> {
    if (!documento) {
      throw new Error('El número de documento es requerido');
    }
    const url = `${this.urlPsicodagogia}/estudiantes/${documento}?tipoParticipante=${tipoParticipante}`;
    return this.http.get<IResponse<any>>(url);
  }

  public getHistorial(body: any): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/historial/filtered`;
    return this.http.post<IResponse<any>>(url, body);
  }

  public keyUrlExists(keyId: string): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/historial/key-url-exists?key_url=${keyId}`;
    return this.http.get<IResponse<any>>(url);
  }

  public hasStudentResponded(dni: string, idEncuesta: number): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/historial/has-student-responded`;
    const body = { dni: dni, id_encuesta: idEncuesta };

    return this.http.post<IResponse<any>>(url, body);
  }

  public getSurveys(): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/encuestas`;
    return this.http.get<IResponse<any>>(url);
  }
  public updateSurvey(data: any, idEncuesta: number): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/encuestas/${idEncuesta}`;
    return this.http.put<IResponse<any>>(url, data);
  }

  public createSurvey(data: any): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/encuestas`;
    return this.http.post<IResponse<any>>(url, data);
  }

  public getAllDiagnostico(): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/diagnostico`;
    return this.http.get<IResponse<any>>(url);
  }

  public createDiagnostico(data: any): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/diagnostico`;
    return this.http.post<IResponse<any>>(url, data);
  }

  public saveHistoryNoSrq(data: any): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/historial/save`;
    return this.http.post<IResponse<any>>(url, data);
  }

  public getPdfDowloader(studentId: number): Observable<IResponse<any>> {
    const url = `${this.urlPsicodagogia}/pdfs/srq/${studentId}`;
    return this.http.get<IResponse<any>>(url);
  }
  //end

  // Ario Masgo
  //  modulo Visita General

  public getAllVisitaGeneral(): Observable<IResponse<IVisitaGeneral[]>> {
    return this.http.get<IResponse<IVisitaGeneral[]>>(this.urlVisitas);
  }


  public getVisitaGeneralById(id: string): Observable<IResponse<IVisitaGeneral>> {
    return this.http.get<IResponse<IVisitaGeneral>>(`${this.urlVisitas}/${id}`);
  }

  public createVisitaGeneral(data: Partial<IVisitaGeneral>): Observable<IResponse<IVisitaGeneral>> {
    return this.http.post<IResponse<IVisitaGeneral>>(this.urlVisitas, data);
  }

  public updateVisitaGeneral(data: Partial<IVisitaGeneral>, id: string): Observable<IResponse<IVisitaGeneral>> {
    const requestData = {id, ...data};
    return this.http.put<IResponse<IVisitaGeneral>>(this.urlVisitas, requestData);
  }
  public geDepartamentos(): Observable<IResponse<Departamento[]>> {
    return this.http.get<IResponse<Departamento[]>>(`${this.urlVisitas}/departamentos`);
  }
  public getProvinciasByDepartamento(departamentoId: string): Observable<IResponse<Provincia[]>> {
    return this.http.get<IResponse<Provincia[]>>(`${this.urlVisitas}/departamentos/${departamentoId}/provincias`);
  }
  public getDistritosByProvincia(provinciaId: string): Observable<IResponse<Distrito[]>> {
    return this.http.get<IResponse<Distrito[]>>(`${this.urlVisitas}/provincias/${provinciaId}/distritos`);
  }
  public createCita(data: any): Observable<IResponse<any>>  {
    const url = `${this.urlPsicodagogia}/citas`;
    return this.http.post<IResponse<any>>(url, data);
  }

  public getDatachartArea(params: Record<string, string>) {
    const queryParams = new URLSearchParams(params).toString();
    return this.http.get<Response<string>>(this.urlPsicodagogia + '/estadistica/chart-area?' + queryParams);
  }

  public getDataChartPie(params: Record<string, string>) {
    const queryParams = new URLSearchParams(params).toString();
    return this.http.get<Response<string>>(this.urlPsicodagogia + '/estadistica/chart-pie?' + queryParams);
  }

  public getDataChartBarras(params: Record<string, string>) {
    const queryParams = new URLSearchParams(params).toString();
    return this.http.get<Response<string>>(this.urlPsicodagogia + '/estadistica/chart-barras?' + queryParams);
  }

  public getReportAttentionsPsicology(params: Record<string, string>){
    const queryParams = new URLSearchParams(params).toString();
    return this.http.get<Response<string>>(this.urlPsicodagogia + '/reporte/atenciones?' + queryParams);
  }

  //end

  public deleteVisitaGeneral(id: string): Observable<IResponse<IVisitaGeneral>> {
    return this.http.delete<IResponse<IVisitaGeneral>>(`${this.urlVisitas}/${id}`);
  }
  // Examen Toxicologico
  public getAllExamenToxicologico(): Observable<IResponse<ExamenToxicologico[]>> {
    return this.http.get<IResponse<ExamenToxicologico[]>>(`${this.urlExamenToxicologico}`);
  }

  public getExamenToxicologicoByConvocatoria(convocatoriaId: string): Observable<IResponse<ExamenToxicologico[]>> {
    return this.http.get<IResponse<ExamenToxicologico[]>>(`${this.urlExamenToxicologico}/convocatoria/${convocatoriaId}`);
  }
  public getExamenToxicologicoById(id: string): Observable<IResponse<ExamenToxicologico>> {
    return this.http.get<IResponse<ExamenToxicologico>>(`${this.urlExamenToxicologico}/${id}`);
  }

  public createOrUpdateExamenToxicologico(data: CreateExamenToxicologicoRequest): Observable<IResponse<ExamenToxicologico>> {
    return this.http.post<IResponse<ExamenToxicologico>>(this.urlExamenToxicologico, data);
  }

  public updateExamenToxicologico(id: string, data: UpdateExamenToxicologicoRequest): Observable<IResponse<ExamenToxicologico>> {
    return this.http.put<IResponse<ExamenToxicologico>>(`${this.urlExamenToxicologico}/${id}`, data);
  }
  public deleteExamenToxicologico(id: string): Observable<IResponse<ExamenToxicologico>> {
    return this.http.delete<IResponse<ExamenToxicologico>>(`${this.urlExamenToxicologico}/${id}`);
  }

  //end
  // visita residente
  public getAllVisitaResidente():Observable<IResponse<VisitaResidente[]>>
  {
    return this.http.get<IResponse<VisitaResidente[]>>(this.urlVisitasResidente);
  }
  public getVisitaResidenteById(id: string): Observable<IResponse<VisitaResidente>> {
    return this.http.get<IResponse<VisitaResidente>>(`${this.urlVisitasResidente}/${id}`);
  }
  public getAlumnosPendientesVisitaResidente(convocatoriaId: number): Observable<IResponse<any[]>> {
    return this.http.get<IResponse<any[]>>(`${this.urlVisitasResidente}/alumnos-pendientes?convocatoria_id=${convocatoriaId}`);
  }
  public createVisitaResidente(data: CreateVisitaResidenteRequest): Observable<IResponse<VisitaResidente>> {
    return this.http.post<IResponse<VisitaResidente>>(this.urlVisitasResidente, data);
  }
  public updateVisitaResidente(id: string, data: UpdateVisitaResidenteRequest): Observable<IResponse<VisitaResidente>> {
    return this.http.put<IResponse<VisitaResidente>>(`${this.urlVisitasResidente}/${id}`, data);
  }
  public deleteVisitaResidente(id: string): Observable<IResponse<VisitaResidente>> {
    return this.http.delete<IResponse<VisitaResidente>>(`${this.urlVisitasResidente}/${id}`);
  }
  public getReporteVisitaResidente(): Observable<IResponse<ReporteVisitaResidente>> {
    return this.http.get<IResponse<ReporteVisitaResidente>>(`${this.urlVisitasResidente}/estadisticas`);
  }
  public getReporteVisitaResidentePorEscuelaProfesional(): Observable<IResponse<ReportePorEscuelaProfesional[]>> {
    return this.http.get<IResponse<ReportePorEscuelaProfesional[]>>(`${this.urlVisitasResidente}/estadisticas/escuela-profesional`);
  }
  public getLugarProcedencia(): Observable<IResponse<LugarProcedencia[]>> {
    return this.http.get<IResponse<LugarProcedencia[]>>(`${this.urlVisitasResidente}/estadisticas/lugar-procedencia`);
  }
  public getAllVisitasPorDepartamento(convocatoriaId: number, departamento:string): Observable<IResponse<VisitasPorDepartamento[]>> {
    return this.http.get<IResponse<VisitasPorDepartamento[]>>(`${this.urlVisitasResidente}/alumnos-pendientes/departamento?convocatoria_id=${convocatoriaId}&departamento=${departamento}`);
  }
  public getExportDataExcel(convocatoriaId: number): Observable<IResponse<ExportDataExcel>> {
    return this.http.get<IResponse<ExportDataExcel>>(`${this.urlVisitasResidente}/alumnos-completos?convocatoria_id=${convocatoriaId}`);
  }

}
