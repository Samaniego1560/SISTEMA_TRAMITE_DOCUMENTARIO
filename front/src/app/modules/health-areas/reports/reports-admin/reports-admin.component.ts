import {Component, OnDestroy} from '@angular/core';
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {ControlErrorsDirective} from "../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormSubmitDirective} from "../../../../lib/validator-dynamic/directives/form-submit.directive";
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  ValidationErrors,
  ValidatorFn
} from "@angular/forms";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {catchError, of, Subject} from "rxjs";
import {DentalService} from "../../../../core/services/health-areas/dental/dental.service";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {NgxValidators} from "../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {Response} from "../../../../core/models/global.model";
import {formatDate} from "@angular/common";
import {AdminReportService} from "../../../../core/services/health-areas/admin/admin-report.service";

@Component({
  selector: 'app-reports-admin',
  standalone: true,
  imports: [
    BlockUiComponent,
    ControlErrorsDirective,
    FormSubmitDirective,
    FormsModule,
    ReactiveFormsModule,
    ToastComponent
  ],
  templateUrl: './reports-admin.component.html',
  styleUrl: './reports-admin.component.scss'
})
export class ReportsAdminComponent implements OnDestroy {
  reportForm: FormGroup;

  cuadros = [
    {name: 'ATENCIONES DE LAS AREAS', value: '1'},
    {name: 'TOTALES ATENCIONES DE LAS AREAS', value: '2'},
  ]
  isBlockPage = false;

  private _subject$ = new Subject();

  constructor(
    private _fb: FormBuilder,
    private _adminReportService: AdminReportService,
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
        fecha_inicio: [null, NgxValidators.required()],
        fecha_fin: [currentYear, [NgxValidators.required()]],
      },
      {
        validators: [this.isValidRangeMonths()]
      })
  }

  isValidRangeMonths(): ValidatorFn {
    return (group: AbstractControl): ValidationErrors | null => {
      const dateStart = group.get('fecha_inicio')?.value;
      const dateEnd = group.get('fecha_fin')?.value;

      if (!dateStart || !dateEnd) {
        return null;
      }

      if (new Date(dateEnd).getTime() >= new Date(dateStart).getTime()) {
        return null;
      }

      return {
        rangeDates: true
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
    debugger
    this.isBlockPage = true;
    const filterParams = this.reportForm.value;
    this._adminReportService.getReport(filterParams)
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

  clearFilters() {
    this.reportForm.reset()
  }
}
