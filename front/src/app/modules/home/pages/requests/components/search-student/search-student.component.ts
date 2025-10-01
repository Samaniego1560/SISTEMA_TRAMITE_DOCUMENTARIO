import {Component, ElementRef, EventEmitter, Input, Output, ViewChild} from '@angular/core';
import { ManagerService } from "../../../../../../core/services/manager/manager.service";
import { IAnnouncement } from "../../../../../../core/models/announcement";
import { take } from "rxjs";
import { IResponse } from "../../../../../../core/models/response";
import { DatePipe, NgForOf, NgIf } from "@angular/common";
import { jsPDF } from 'jspdf';
import html2canvas from 'html2canvas';
import {ReportService} from "../../../../../../core/services/report/report.service";
import {RequirementReport, StudentFileReport} from "../../../../../../core/models/report";
import {FormsModule} from "@angular/forms";
import {ToastService} from "../../../../../../core/services/toast/toast.service";
import {ToastComponent} from "../../../../../../core/ui/toast/toast.component";
import {BlockUiComponent} from "../../../../../../core/ui/block-ui/block-ui.component";

@Component({
  selector: 'app-search-student',
  standalone: true,
  imports: [
    FormsModule,
    ToastComponent,
    BlockUiComponent
  ],
  templateUrl: './search-student.component.html',
  styleUrls: ['./search-student.component.scss']
})
export class SearchStudentComponent {
  public studentFileReport: StudentFileReport = {} as StudentFileReport;
  public dni: string = '';
  protected isLoading: boolean = false;
  @Output() onBack: EventEmitter<string> = new EventEmitter();

    @Input({alias: 'departments', required: true}) departmentsMap: Record<string, string> = {};
    @Input({alias: 'provinces', required: true}) provincesMap: Record<string, string> = {};
    @Input({alias: 'districts', required: true}) districtsMap: Record<string, string> = {};
  

  public urlProfile: string = '';

  constructor(private _reportService: ReportService,
              private _toastService: ToastService,) {}

  protected async searchStudent(dni: string) {
    this.isLoading = true;
    const response = await this.reportProfilePromise(dni);
    if (response.error) {
      this.isLoading = false;
      this._toastService.add({
        type: 'error',
        message: 'No se pudo extraer información'
      });
      return;
    }
    this.isLoading = false;
    this.studentFileReport = response.data!;
    this.searchPhotoProfile()
  }

  private searchPhotoProfile() {
    const linkProfile = this.studentFileReport.link_profile;
    if (linkProfile) this.urlProfile = linkProfile
  }

  back() {
    this.onBack.emit('list')
  }

  private reportProfilePromise(dni: string): Promise<{ data: StudentFileReport | null, error: boolean, msg: string }> {
    return new Promise((resolve) => {
      this._reportService.getReportProfile(dni).pipe(take(1)).subscribe({
        next: (res: IResponse<any>) => resolve({ data: res.detalle, error: false, msg: '' }),
        error: () => resolve({ data: null, error: true, msg: 'Error al obtener datos del estudiante.' })
      });
    });
  }

  public async downloadAsPDF() {
    this.isLoading = true;
    const response = await this.reportProfileDownloadPromise(this.dni);
    if (response.error) {
      this.isLoading = false;
      this._toastService.add({
        type: 'error',
        message: 'No se descargar pdf'
      });
      return;
    }
    this.isLoading = false;
    const url = window.URL.createObjectURL(response.data);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'document_' + this.dni + '.pdf';
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  }

  private reportProfileDownloadPromise(dni: string): Promise<{ data: any | null, error: boolean, msg: string }> {
    return new Promise((resolve) => {
      this._reportService.getReportProfileDownload(dni).pipe(take(1)).subscribe({
        next: (res: any) => resolve({ data: res, error: false, msg: '' }),
        error: () => resolve({ data: null, error: true, msg: 'Error al obtener datos del estudiante.' })
      });
    });
  }

  public getValues(requirement: RequirementReport): string {
    if (requirement.type === 4) {
      return this.departmentsMap[requirement.value] || this.provincesMap[requirement.value] || this.districtsMap[requirement.value] || requirement.value  ;
    }
    return requirement.value;
  }
}
