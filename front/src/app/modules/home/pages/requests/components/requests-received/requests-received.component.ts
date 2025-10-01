import {Component, effect, Input, input, model, OnInit} from '@angular/core';
import {NgClass, NgIf} from "@angular/common";
import {HttpErrorResponse} from "@angular/common/http";
import {ManagerService} from "../../../../../../core/services/manager/manager.service";
import {Subscription} from "rxjs";
import {IAnnouncement} from "../../../../../../core/models/announcement";
import {ToastService} from "../../../../../../core/services/toast/toast.service";
import {IRequest, IUpdateService} from "../../../../../../core/models/requests";
import * as XLSX from "xlsx";
import {saveAs} from "file-saver";
import {FilterService} from "../../../../../../core/services/filter/filter.service";
import {ModalComponent} from "../../../../../../core/ui/modal/modal.component";
import {FormControl, ReactiveFormsModule, Validators} from "@angular/forms";
import {PostulationFormComponent} from "../../../../../../core/components/postulation-form/postulation-form.component";
import {BlockUiComponent} from "../../../../../../core/ui/block-ui/block-ui.component";
import {RequestEvaluationComponent} from "./components/request-evaluation/request-evaluation.component";
import {Student} from "../../../../../../core/models/student";
import {
  PostulationDynamicComponent
} from "../../../../../../core/components/postulation-dynamic/postulation-dynamic.component";

@Component({
  selector: 'app-requests-received',
  standalone: true,
  imports: [
    NgIf,
    NgClass,
    ReactiveFormsModule,
    PostulationFormComponent,
    BlockUiComponent,
    RequestEvaluationComponent,
    PostulationDynamicComponent
  ],
  templateUrl: './requests-received.component.html',
  styleUrl: './requests-received.component.scss'
})
export class RequestsReceivedComponent {
  viewMain = model.required<string>();
  role = input.required<number>();

  announcement = input.required<IAnnouncement>();

  @Input({alias: 'departments', required: true}) departmentsMap: Record<string, string> = {};
  @Input({alias: 'provinces', required: true}) provincesMap: Record<string, string> = {};
  @Input({alias: 'districts', required: true}) districtsMap: Record<string, string> = {};

  protected postulation: IAnnouncement = {} as IAnnouncement;

  protected dataServices: {comedor: number, internado: number} = {comedor: 0, internado: 0};

  protected requestsDisplay: IRequest[] = [];

  protected view: string = 'list';
  protected requests: IRequest[] = [];
  private _subscriptions: Subscription = new Subscription();

  protected isLoading: boolean = false;

  protected showModelAddPostulation: boolean = false;

  protected studentForm: Student = {} as Student;

  protected leftLimit: number = 0;

  protected currentPage: number = 1;

  protected totalPage: number = 0;


  protected rightLimit: number = 10;

  protected messageService: FormControl;

  protected student!: IRequest;

  public countServiceWomen: number = 0;
  public countServiceMen: number = 0;

  constructor(private _managerService: ManagerService,
              private _toastService: ToastService,
              private _filterService: FilterService,
              ) {
    this.messageService = new FormControl<string>('', Validators.required);
    effect(() => {
      this.getRequests();
      this.getNumberOfVacancies();
    });
  }

  public getNumberOfVacancies() {
    const id = this.announcement().id;
    this.isLoading = true;
    this._subscriptions.add(
      this._managerService
        .getNumberOfVacancies(id!)
        .subscribe({
          next: (res: any) => {
            this.isLoading = false;
            if (!res.detalle) {
              this._toastService.add({ type: 'error', message: res.msg });
              return;
            }

            this.dataServices.comedor = res.detalle.comedor;
            this.dataServices.internado = res.detalle.internado;

          },
          error: (err: HttpErrorResponse) => {
            console.error(err);
            this.isLoading = false;
            this._toastService.add({
              type: 'error',
              message:
                'No se pudo detalles del servicio',
            });
            return;
          },
        })
    );
  }


  protected exportToExcel(): void {
    this.isLoading = true;
    const id = this.announcement().id!;
    if (!id) {
      this._toastService.add({
        type: 'error',
        message: 'No tiene una convocatoria seleccionada!'
      });
      return;
    }

    this._subscriptions.add(
      this._managerService.exportDataRequestById(id).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          const detail = res.detalle;

          if (detail && detail.length > 0) {
            const originalHeaders = Object.keys(detail[0]);
            const headers = originalHeaders.map(header =>
              header
                .split('_')
                .map(word => word.charAt(0).toUpperCase() + word.slice(1))
                .join(' ')
            );
            const worksheetData = [
              headers,
              ...detail.map((item: any) =>
                originalHeaders.map(originalHeader => item[originalHeader])
              )
            ];

            const worksheet: XLSX.WorkSheet = XLSX.utils.aoa_to_sheet(worksheetData);
            const workbook: XLSX.WorkBook = { Sheets: { 'Datos': worksheet }, SheetNames: ['Datos'] };

            const excelBuffer: any = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });

            this.saveAsExcelFile(excelBuffer, this.announcement().nombre);
          }
          this._toastService.add({type: 'info', message: "Archivo exportado correctamente"});

        },
        error: (err: HttpErrorResponse) => {
          console.error(err);
          this.isLoading = false;
          this._toastService.add({
            type: 'error',
            message: 'No se pudo exportar el archivo, intente nuevamente!'
          });
        }
      })
    );
  }

  protected findRequest(event: any): void {
    const dataToFilter: IRequest[] = JSON.parse(JSON.stringify(this.requests));
    const filterValue = event.target.value;
    if (!filterValue || !filterValue.length) {
      this.requestsDisplay = dataToFilter;
      this.initPagination();
      return;
    }

    const searchFields: string[] = [
      'codigo_estudiante',
      'nombres',
      'DNI',
      'apellido_paterno',
      'apellido_materno',
      'correo_institucional',
      'correo_personal',
      'celular_estudiante',
      'celular_padre',
      'facultad',
      'escuela_profesional',
      'modalidad_ingreso',
      'lugar_procedencia',
      'lugar_nacimiento',
      'direccion',
      'fecha_nacimiento',
      'edad'
    ];
    const data = this._filterService.filter(dataToFilter, searchFields, filterValue, 'contains');
    this.initPagination(data);
  }

  protected getStudentRequest(student: IRequest): void {
    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.getStudentRequest(student.id).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({type: 'info', message: res.msg});
            return;
          }
          this.view = 'student';
          this.student = res.detalle;
        },
        error: (err: HttpErrorResponse) => {
          console.error(err);
          this.isLoading = false;
          this._toastService.add({
            type: 'error',
            message: 'No se pudo obtener la solicitud del estudiante, intente de nuevo'
          });
        }
      })
    );
  }

  protected nextPage(): void {
    this.leftLimit = this.currentPage * 10;
    this.rightLimit = this.leftLimit + 10;
    this.currentPage++;
    this.paginate();
  }

  protected previousPage(): void {
    this.currentPage--;
    this.leftLimit = this.leftLimit - 10;
    this.rightLimit = this.leftLimit + 10;
    this.paginate();
  }

  saveAsExcelFile(buffer: any, fileName: string): void {
    const EXCEL_TYPE = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8';
    const EXCEL_EXTENSION = '.xlsx';

    const data: Blob = new Blob([buffer], { type: EXCEL_TYPE });
    saveAs(data, fileName + '_export_' + new Date().getTime() + EXCEL_EXTENSION);
  }

  protected initPagination(data?: IRequest[]): void {
    this.leftLimit = 0;
    this.rightLimit = 10;
    this.currentPage = 1;
    const totalRegister = data ? data.length : this.requests.length;
    this.totalPage = Math.ceil(totalRegister / 10);
    this.paginate(data);
  }

  protected paginate(data?: IRequest[]): void {
    const requests: IRequest[] = JSON.parse(JSON.stringify(data || this.requests));
    this.requestsDisplay = requests.slice(this.leftLimit, this.rightLimit);
  }

  private getRequests(): void {
    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.getRequests(this.announcement().id || 0).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }

          const dataNew: any = [];

          this.requests = [];
          for (const re of res.detalle) {
            const serviciosSolicitados = re.servicios_solicitados;
            serviciosSolicitados.forEach((servicio: any) => {
              if (servicio.servicio_id === 2 && servicio.estado === 'aprobado' && re.alumno.sexo === 'F') this.countServiceWomen++;
              if (servicio.servicio_id === 2 && servicio.estado === 'aprobado' && re.alumno.sexo === 'M') this.countServiceMen++;
            })

            if (re.alumno.facultad === 'FACULTAD DE INGENIERIA EN INDUSTRIAS ALIMENTARIAS'
              && re.alumno.codigo_estudiante.startsWith('2023')
              && re.alumno.sexo === 'F') {
              dataNew.push(re);
            }


            const index = this.requests.findIndex(r => r.alumno.DNI === re.alumno.DNI);
            if (index === -1) {
              this.requests.push(re);
              continue;
            }

            if (re.id > this.requests[index].id) {
              this.requests[index] = re;
            }
          }
          this.initPagination();
        },
        error: (err: HttpErrorResponse) => {
          console.error(err);
          this.isLoading = false;
          this._toastService.add({
            type: 'error',
            message: 'No se pudo obtener el listado de solicitantes, intente de nuevo'
          });
        }
      })
    );
  }

  public backRequest(): void {
    this.viewMain.set('list')
  }

  public onPostulation(event: { postulation: IAnnouncement, student: Student }) {
    this.postulation = event.postulation;
    this.studentForm = event.student;
    this.view = 'postulation';
    this.showModelAddPostulation = false;
  }
}
