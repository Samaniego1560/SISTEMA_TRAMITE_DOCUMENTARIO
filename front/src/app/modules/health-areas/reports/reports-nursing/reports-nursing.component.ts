import {Component, OnDestroy, OnInit} from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  ValidationErrors, ValidatorFn
} from "@angular/forms";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {NursingService} from "../../../../core/services/health-areas/nursing/nursing.service";
import {catchError, of, Subject} from "rxjs";
import {Response} from "../../../../core/models/global.model";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {NgxValidators} from "../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {formatDate, NgIf} from "@angular/common";
import {ReportGeneralFormComponent} from "../report-general-form/report-general-form.component";
import {ControlErrorsDirective} from "../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormSubmitDirective} from "../../../../lib/validator-dynamic/directives/form-submit.directive";
import {CalendarYearComponent} from "../../../../core/ui/calendar-year/calendar-year.component";

@Component({
  selector: 'app-reports-nursing',
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
  templateUrl: './reports-nursing.component.html',
  styleUrl: './reports-nursing.component.scss',
  providers: [ToastService, NursingService]
})
export class ReportsNursingComponent implements OnInit, OnDestroy {
  reportForm: FormGroup;

  public trimestres = [
    {name: '1er trimestre', value: '1,2,3'},
    {name: '2do trimestre', value: '4,5,6'},
    {name: '3er trimestre', value: '7,8,9'},
    {name: '4to trimestre', value: '10,11,12'},
  ]

  cuadros = [
    {
      name: '1. ATENCIONES DE CONSULTAS DE ENFERMERIA A ESTUDIANTES POR TRIMESTRE, SEGÚN ESCUELA PROFESIONAL',
      value: '1'
    },
    {name: '2. PROCEDIMIENTOS DE ENFERMERIA REALIZADOS A ESTUDIANTES, SEGÚN ESCUELA PROFESIONAL Y MESES', value: '2'},
    {name: '3. PROCEDIMIENTOS DE ENFERMERIA REALIZADOS A ESTUDIANTES, SEGÚN ESCUELA PROFESIONAL', value: '3'},
    {name: '4. ATENCIONES DE CONSULTAS DE ENFERMERIA POR MES Y SEXO, SEGÚN TIPO DE PERSONAL', value: '4'},
    {name: '5. ATENCIONES DE ENFERMERIA POR PROCEDIMIENTOS REALIZADOS SEGÚN PERSONAL ATENDIDO Y MES', value: '5'},
    {name: '6. ATENCIONES DE ENFERMERÍA A ESTUDIANTES POR TRIMESTRE Y SEXO, SEGÚN ESCUELA PROFESIONAL', value: '6'},
    {name: '7. ATENCIONES DE ENFERMERÍA POR TRIMESTRE Y SEXO, SEGÚN TIPO DE PERSONAL', value: '7'},
    {name: '8. ATENCIONES DE ENFERMERIA A ESTUDIANTES POR MES, SEGÚN PROCEDIMIENTOS REALIZADOS', value: '8'},
    {
      name: '9. EVOLUCIÓN DE ATENCIONES DE SERVICIO DE TÓPICO POR PROCEDIMIENTO REALIZADO, SEGÚN ESCUELA PROFESIONAL',
      value: '9'
    },
  ]
  isBlockPage = false;
  activeTab: 1 | 2 = 1;

  private _subject$ = new Subject();

  constructor(
    private _fb: FormBuilder,
    private _nursingService: NursingService,
    private _toastService: ToastService
  ) {
    this.reportForm = this.buildFormReport();
  }

  ngOnInit() {
    this.initFilters()
  }

  ngOnDestroy() {
    this._subject$.next(true);
    this._subject$.complete();
  }

  buildFormReport() {
    const currentYear = new Date().getFullYear().toString();
    return this._fb.group({
        anio: [currentYear, [NgxValidators.required()]],
        anio_inicial: [currentYear, [NgxValidators.required()]],
        anio_final: [currentYear, [NgxValidators.required()],],
        meses: [null, {validators: [NgxValidators.required()]}],
        numero_cuadro: [null, [NgxValidators.required()]],
        mes_inicio: [null, [NgxValidators.required()]],
        mes_fin: [null, [NgxValidators.required()]],
      }, {
        validators: [this.isValidRangeMonths()]
      }
    )
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

  initFilters(): void {
    this._enableFieldsForReportDefault();
  }

  onSelectBox(nroBox: string): void {
    if (nroBox === '9') {
      this._enableFieldsForReport9();
      return
    }

    if (nroBox === '6' || nroBox === '7' || nroBox === '8') {
      this._enableFieldsForReport6();
      return
    }

    if (nroBox === '3') {
      this._enableFieldsForReport3();
      return
    }

    this._enableFieldsForReportDefault();
  }

  private _enableFieldsForReportDefault(): void {
    this.reportForm.get('anio')?.enable();
    this.reportForm.get('meses')?.enable();

    this.reportForm.get('mes_inicio')?.disable();
    this.reportForm.get('mes_fin')?.disable();
    this.reportForm.get('anio_inicial')?.disable();
    this.reportForm.get('anio_final')?.disable();

    this.reportForm.get('mes_inicio')?.reset();
    this.reportForm.get('mes_fin')?.reset();
    this.reportForm.get('anio_inicial')?.reset();
    this.reportForm.get('anio_final')?.reset();
  }

  private _enableFieldsForReport9(): void {
    this.reportForm.get('anio_inicial')?.enable();
    this.reportForm.get('anio_final')?.enable();

    this.reportForm.get('meses')?.disable();
    this.reportForm.get('anio')?.disable();
    this.reportForm.get('mes_inicio')?.disable();
    this.reportForm.get('mes_fin')?.disable();

    this.reportForm.get('meses')?.reset();
    this.reportForm.get('anio')?.reset();
    this.reportForm.get('mes_inicio')?.reset();
    this.reportForm.get('mes_fin')?.reset();
  }

  private _enableFieldsForReport6(): void {
    this.reportForm.get('anio')?.enable();

    this.reportForm.get('meses')?.disable();
    this.reportForm.get('anio_inicial')?.disable();
    this.reportForm.get('anio_final')?.disable();
    this.reportForm.get('mes_inicio')?.disable();
    this.reportForm.get('mes_fin')?.disable();

    this.reportForm.get('meses')?.reset('');
    this.reportForm.get('anio_inicial')?.reset('');
    this.reportForm.get('anio_final')?.reset();
    this.reportForm.get('mes_inicio')?.reset();
    this.reportForm.get('mes_fin')?.reset();
  }

  private _enableFieldsForReport3(): void {
    this.reportForm.get('anio')?.enable();
    this.reportForm.get('mes_inicio')?.enable();
    this.reportForm.get('mes_fin')?.enable();

    this.reportForm.get('anio_inicial')?.disable();
    this.reportForm.get('anio_final')?.disable();
    this.reportForm.get('meses')?.disable();

    this.reportForm.get('mes_fin')?.reset();
    this.reportForm.get('anio_inicial')?.reset();
    this.reportForm.get('anio_final')?.reset();
    this.reportForm.get('mes_inicio')?.reset();
  }

  onSearchReportDetail(): void {
    if (this.reportForm.invalid) {
      this.reportForm.markAllAsTouched();

      if(this.reportForm.errors) {
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
    const filterParams = this._getParamsReportDetail();
    this._nursingService.getReportDetail(filterParams)
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

  private _getParamsReportDetail() {
    const formValues = this.reportForm.value;
    return {
      ...formValues,
      anio: formValues['anio'] ? formatDate(formValues['anio'], 'yyyy', 'en') : null,
      anio_inicial: formValues['anio_inicial'] ? formatDate(formValues['anio_inicial'], 'yyyy', 'en') : null,
      anio_final: formValues['anio_final'] ? formatDate(formValues['anio_final'], 'yyyy', 'en') : null,
    }
  }

  clearFilterReport() {
    this.reportForm.reset()
  }
}
