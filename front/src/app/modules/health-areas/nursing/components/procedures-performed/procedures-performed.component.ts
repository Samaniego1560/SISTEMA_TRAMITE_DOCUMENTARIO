import {Component, inject, Input, OnChanges, output, signal, SimpleChanges} from '@angular/core';
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {ProcedimientoRealizado} from "../../../../../core/models/areas/nursing/nursing.model";
import {NgxValidators} from "../../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {FormSubmitDirective} from "../../../../../lib/validator-dynamic/directives/form-submit.directive";
import {InformationPayment} from "../../../../../core/models/trasury.model";
import {formatDate, TitleCasePipe} from "@angular/common";
import {ManagerService} from "../../../../../core/services/manager/manager.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {OnlyNumbersDirective} from "../../../../../core/directives/only-numbers.directive";

@Component({
  selector: 'app-procedures-performed',
  standalone: true,
  imports: [
    FormSubmitDirective,
    ControlErrorsDirective,
    FormsModule,
    ReactiveFormsModule,
    EcTableModule,
    BlockUiComponent,
    OnlyNumbersDirective,
    TitleCasePipe,
  ],
  templateUrl: './procedures-performed.component.html',
  styleUrl: './procedures-performed.component.scss'
})
export class ProceduresPerformedComponent extends FormBaseComponent<ProcedimientoRealizado> implements OnChanges {
  @Input({required: true}) dniPatient: string | undefined;

  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }

  deleteForm = output();

  public active = signal<boolean>(true);
  public form: FormGroup;
  private _name = signal('procedimiento_realizado');
  public typesProcedures = [
    'INYECTABLES',
    'CONTROL DE FUNCIONES VITALES',
    'LAVADO DE OIDO',
    'BAÑO RAYOS INFRAROJOS',
    'EXTRACCION DE UÑERO',
    'CIRUGIA MENOR',
    'VENOCLISIS',
    'OTROS',
  ]

  areas = [
    'odontología',
    'medicina',
    "psicología",
    "otros"
  ]

  isBlockPage = signal<boolean>(false)

  private _fb = inject(FormBuilder);
  private _managerService = inject(ManagerService);
  private _toastService = inject(ToastService);

  constructor() {
    super();
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoData'] && this.infoData) {
      this.form.patchValue(this.infoData);
      this.active.set(false);
      if(!this.canSaveForm) this.form.disable();
    }
  }

  private _buildForm() {
    return this._fb.nonNullable.group({
      id: [''],
      procedimiento: ['', NgxValidators.required()],
      comentarios: ['', NgxValidators.required()],
      area_solicitante: [''],
      especialista_solicitante: [''],
      numero_recibo: ['',],
      costo: [{value: '', disabled: false}],
      fecha_pago: [{value: '', disabled: false}],
    })
  }

  get name() {
    return this._name();
  }

  isValid(): boolean {
    return this.form.valid;
  }

  getData() {
    return this.form.value;
  }

  onValidateRecibo(nroRecibo: string) {
    if (!nroRecibo) return

    this.isBlockPage.set(true)
    this._managerService.validateStudentPayment(this.dniPatient!)
      .subscribe({
        next: (res: InformationPayment) => {
          this._validatePayment(res, nroRecibo);
          this.isBlockPage.set(false)
        },
        error: () => {
          this._toastService.add({
            type: 'warning',
            message: 'No se encontró el recibo',
            life: 5000
          })
          this.isBlockPage.set(false)
        }
      })
  }

  private _validatePayment(info: InformationPayment, nroRecibo: string) {
    for (const pago of info.pagos_realizados) {
      if (pago.cod_recibo === nroRecibo) {
        this.form.get('costo')?.setValue(pago.importe_pagado);
        this.form.get('fecha_pago')?.setValue(formatDate(pago.fecha_pago, 'yyyy-MM-dd HH:mm', 'en-US'));
        this._toastService.add({
          type: 'info',
          message: `Concepto pagado: ${pago.concepto_pagado}`,
          life: 6000
        })
        return
      }
    }
  }

  onSelectArea(event: Event) {
    if((<HTMLSelectElement>event.target).value === '') {
      this.form.get('especialista_solicitante')?.clearValidators()
    }else {
      this.form.get('especialista_solicitante')?.setValidators(Validators.required)
    }
    this.form.get('especialista_solicitante')?.updateValueAndValidity()
  }
}
