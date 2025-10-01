import {Component, OnDestroy} from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  ValidationErrors,
  ValidatorFn
} from "@angular/forms";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {catchError, of, Subject} from "rxjs";
import {Response} from "../../../../core/models/global.model";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {NgxValidators} from "../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {formatDate} from "@angular/common";
import {DentalService} from "../../../../core/services/health-areas/dental/dental.service";
import {ReportGeneralFormComponent} from "../report-general-form/report-general-form.component";
import {ControlErrorsDirective} from "../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormSubmitDirective} from "../../../../lib/validator-dynamic/directives/form-submit.directive";
import {CalendarYearComponent} from "../../../../core/ui/calendar-year/calendar-year.component";

@Component({
  selector: 'app-reports-dentistry',
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
  templateUrl: './reports-dentistry.component.html',
  styleUrl: './reports-dentistry.component.scss',
  providers: [ToastService, DentalService]
})
export class ReportsDentistryComponent implements OnDestroy {
  reportForm: FormGroup;

  public trimestres = [
    {name: '1er trimestre', value: '1,2,3'},
    {name: '2do trimestre', value: '4,5,6'},
    {name: '3er trimestre', value: '7,8,9'},
    {name: '4to trimestre', value: '10,11,12'},
  ]

  cuadros = [
    {name: '1. ATENCIONES DE SERVICIO DE ODONTOLOGÍA POR PROCEDIMIENTOS REALIZADO, SEGÚN TIPO DE PERSONAL', value: '1'},
    {
      name: '2. ATENCIONES DE SERVICIO DE ODONTOLOGÍA POR PROCEDIMIENTOS REALIZADO, SEGÚN ESCUELA PROFESIONAL',
      value: '2'
    },
  ]
  isBlockPage = false;
  activeTab: 1 | 2 = 1;

  private _subject$ = new Subject();

  constructor(
    private _fb: FormBuilder,
    private _dentalService: DentalService,
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
        numero_cuadro: [null, NgxValidators.required()],
        anio: [currentYear, [NgxValidators.required()]],
        mes_inicio: [null, [NgxValidators.required()]],
        mes_fin: [null, [NgxValidators.required()]],
      },
      {
        validators: [this.isValidRangeMonths()]
      })
  }

  isValidRangeMonths(): ValidatorFn {
    return (group: AbstractControl): ValidationErrors | null => {
      const monthStart = group.get('mes_inicio')?.value;
      const monthEnd = group.get('mes_fin')?.value;

      if (!monthStart || !monthEnd) {
        return null;
      }

      if (Number(monthEnd) >= Number(monthStart)) {
        return null;
      }

      return {
        rangeMonths: true
      }
    }
  }

  onSearchReport(): void {
    if (this.reportForm.invalid) {
      this.reportForm.markAllAsTouched();

      if (this.reportForm.errors) {
        this._toastService.add({
          message: 'El mes de inicio debe ser menor al mes de fin',
          type: 'error',
          life: 5000
        });
      }
      return
    }

    this.searchReport()
  }

  searchReport() {
    this.isBlockPage = true;
    const filterParams = this._getParamsReport();
    this._dentalService.getReportDetail(filterParams)
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
    console.log(formValues['anio'])
    return {
      ...formValues,
      anio: formValues['anio'] ? formatDate(formValues['anio'], 'yyyy', 'en') : null,
    }
  }

  clearFilters() {
    this.reportForm.reset()
  }
}
