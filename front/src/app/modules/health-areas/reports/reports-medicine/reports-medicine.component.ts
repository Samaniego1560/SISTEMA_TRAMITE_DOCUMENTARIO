import {Component, OnDestroy} from '@angular/core';
import {FormBuilder, FormGroup, ReactiveFormsModule} from "@angular/forms";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {catchError, of, Subject} from "rxjs";
import {Response} from "../../../../core/models/global.model";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {NgxValidators} from "../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {formatDate} from "@angular/common";
import {MedicineService} from "../../../../core/services/health-areas/medicine/medicine.service";
import {ReportGeneralFormComponent} from "../report-general-form/report-general-form.component";
import {ControlErrorsDirective} from "../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormSubmitDirective} from "../../../../lib/validator-dynamic/directives/form-submit.directive";
import {CalendarYearComponent} from "../../../../core/ui/calendar-year/calendar-year.component";

@Component({
  selector: 'app-reports-medicine',
  standalone: true,
    imports: [
        ReactiveFormsModule,
        ToastComponent,
        BlockUiComponent,
        ReportGeneralFormComponent,
        ControlErrorsDirective,
        FormSubmitDirective,
        CalendarYearComponent
    ],
  templateUrl: './reports-medicine.component.html',
  styleUrl: './reports-medicine.component.scss',
  providers: [ToastService, MedicineService]
})
export class ReportsMedicineComponent implements OnDestroy {
  reportForm: FormGroup;

  public trimestres = [
    {name: '1er trimestre', value: '1,2,3'},
    {name: '2do trimestre', value: '4,5,6'},
    {name: '3er trimestre', value: '7,8,9'},
    {name: '4to trimestre', value: '10,11,12'},
  ]
  activeTab: 1 | 2 = 1;

  cuadros = [
    {name: '1. CONSULTAS A ESTUDIANTES POR MES Y SEXO SEGÚN ESCUELA PROFESIONAL', value: '1'},
    {name: '2. CONSULTAS POR MES Y SEXO SEGÚN TIPO DE PERSONAL', value: '3'},
    {name: '3. ATENCIONES DE CONSULTAS MÉDICAS A ESTUDIANTES DEL I TRIMESTRE, SEGÚN ESCUELA PROFESIONAL', value: '8'},
    {name: '4. ATENCIONES DE CONSULTAS MÉDICAS A ESTUDIANTES DEL I TRIMESTRE, SEGÚN TIPO DE PERSONAL', value: '12'},
  ]
  isBlockPage = false;

  private _subject$ = new Subject();

  constructor(
    private _fb: FormBuilder,
    private _medicineService: MedicineService,
    private _toastService: ToastService
  ) {
    this.reportForm = this.buildForm();
  }

  ngOnDestroy() {
    this._subject$.next(true);
    this._subject$.complete();
  }

  buildForm() {
    const currentYear = new Date().getFullYear().toString();
    return this._fb.group({
      anio: [currentYear, [NgxValidators.required()]],
      meses: [null, NgxValidators.required()],
      numero_cuadro: [null, NgxValidators.required()],
    })
  }

  onSearchReport(): void {
    if(this.reportForm.invalid) {
      this.reportForm.markAllAsTouched();
      return
    }

    this.searchReport()
  }

  searchReport() {
    this.isBlockPage = true;
    const filterParams = this._getParamsReport();
    this._medicineService.getReportDetail(filterParams)
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

  private _getParamsReport() {
    const formValues = this.reportForm.value;
    return {
      ...formValues,
      anio: formValues['anio'] ? formatDate(formValues['anio'], 'yyyy', 'en') : null,
    }
  }

  clearFilters() {
    this.reportForm.reset()
  }
}
