import {Component, input} from '@angular/core';
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule} from "@angular/forms";
import {NursingService} from "../../../../core/services/health-areas/nursing/nursing.service";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {formatDate} from "@angular/common";
import {NgxValidators} from "../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {catchError, of} from "rxjs";
import {Response} from "../../../../core/models/global.model";
import {MedicineService} from "../../../../core/services/health-areas/medicine/medicine.service";
import {DentalService} from "../../../../core/services/health-areas/dental/dental.service";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ControlErrorsDirective} from "../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormSubmitDirective} from "../../../../lib/validator-dynamic/directives/form-submit.directive";

@Component({
  selector: 'app-report-general-form',
  standalone: true,
  imports: [
    FormsModule,
    ReactiveFormsModule,
    BlockUiComponent,
    ControlErrorsDirective,
    FormSubmitDirective
  ],
  templateUrl: './report-general-form.component.html',
  styleUrl: './report-general-form.component.scss',
  providers: [NursingService, MedicineService, DentalService, ToastService]
})
export class ReportGeneralFormComponent {
  area = input.required<'enfermería' | 'medicina' | 'odontología'>();
  reportForm: FormGroup;
  isBlockPage = false;

  constructor(
    private _fb: FormBuilder,
    private _nursingService: NursingService,
    private _medicineService: MedicineService,
    private _dentalService: DentalService,
    private _toastService: ToastService
  ) {
    this.reportForm = this.buildFormReport();
  }

  buildFormReport() {
    const currentDate = formatDate(new Date(), 'yyyy-MM-dd', 'en');
    return this._fb.group({
      fecha_inicio: [currentDate, NgxValidators.required()],
      fecha_fin: [currentDate, NgxValidators.required()],
    })
  }

  onSearchReport(): void {
    if (this.reportForm.invalid) {
      this.reportForm.markAllAsTouched();
      return
    }

    this.searchReport()
  }

  searchReport() {
    this.isBlockPage = true;
    this.getServiceArea()
      .pipe(catchError(() => (of({error: true} as Response<string>))))
      .subscribe({
        next: (res) => {
          this.isBlockPage = false;
          if (res.error) {
            this._toastService.add({
              message: 'Error al obtener el reporte',
              type: 'error',
              life: 5000
            });
            return
          }

          const link = document.createElement('a');
          link.href = `data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,${res.data}`;
          link.download = 'Reporte';
          link.click();
        }
      })
  }

  getServiceArea() {
    const filters = this._getParamsReport();
    if (this.area() === 'enfermería') {
      return this._nursingService.getReportGeneral(filters)
    }
    if (this.area() === 'medicina') {
      return this._medicineService.getReportGeneral(filters)
    }

    return this._dentalService.getReportGeneral(filters)
  }

  private _getParamsReport() {
    const formValues = this.reportForm.value;
    return {
      fecha_inicio: formatDate(formValues['fecha_inicio'], 'yyyy-MM-dd', 'en'),
      fecha_fin: formatDate(formValues['fecha_fin'], 'yyyy-MM-dd', 'en'),
      area_medica: this.area()
    }
  }

  clearFilterReport() {
    this.reportForm.reset()
  }
}
