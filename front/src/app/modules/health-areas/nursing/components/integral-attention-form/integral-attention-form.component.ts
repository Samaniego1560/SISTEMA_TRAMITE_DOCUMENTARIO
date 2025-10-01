import {Component, inject, Input, OnChanges, OnInit, output, signal, SimpleChanges} from '@angular/core';
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormBuilder, FormGroup, ReactiveFormsModule} from "@angular/forms";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {
  ConsultaBucalGeneral,
  ConsultaMedicinaGeneral,
  Vacuna
} from "../../../../../core/models/areas/nursing/nursing.model";
import {IntegralAttention} from "../../../../../core/models/areas/areas.model";
import {ToggleSwitchComponent} from "../../../../../core/ui/toggle/toggle-switch.component";
import {OnlyNumbersDirective} from "../../../../../core/directives/only-numbers.directive";
import {InformationPayment} from "../../../../../core/models/trasury.model";
import {ManagerService} from "../../../../../core/services/manager/manager.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {formatDate, NgTemplateOutlet} from "@angular/common";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {FormSubmitDirective} from "../../../../../lib/validator-dynamic/directives/form-submit.directive";

@Component({
  selector: 'app-integral-attention-form',
  standalone: true,
  imports: [
    ControlErrorsDirective,
    ReactiveFormsModule,
    ToggleSwitchComponent,
    OnlyNumbersDirective,
    BlockUiComponent,
    NgTemplateOutlet,
    FormSubmitDirective
  ],
  templateUrl: './integral-attention-form.component.html',
  styleUrl: './integral-attention-form.component.scss'
})
export class IntegralAttentionFormComponent extends FormBaseComponent<IntegralAttention> implements OnChanges, OnInit {
  @Input({required: true}) dniPatient: string | undefined;
  @Input() infoMedicine: ConsultaMedicinaGeneral | undefined;
  @Input() infoDentistry: ConsultaBucalGeneral | undefined;
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }

  deleteForm = output();

  public active = signal<boolean>(true);
  public form: FormGroup;
  private _name = signal('consulta_general');
  isBlockPage = signal<boolean>(false);

  showNursing = signal(true)
  showMedicine = signal(false)
  showDentistry = signal(false)

  private _fb = inject(FormBuilder);
  private _managerService = inject(ManagerService);
  private _toastService = inject(ToastService);

  constructor() {
    super();
    this.form = this._buildForm();
  }

  private _buildForm() {
    return this._fb.nonNullable.group({
      //fecha: ['', ],
      //hora: ['',],
      edad: ['',],
      numero_recibo: ['',],
      costo: ['',],
      fecha_pago: ['',],
      motivo_consulta: ['',],
      tiempo_enfermedad: ['',],
      apetito: ['no',],
      sed: ['no',],
      suenio: ['no',],
      estado_animo: ['',],
      orina: ['',],
      deposiciones: ['',],
      temperatura: ['',],
      presion_arterial: ['',],
      frecuencia_cardiaca: ['',],
      frecuencia_respiratoria: ['',],
      peso: ['',],
      talla: ['',],
      indice_masa_corporal: ['',],
      diagnostico: ['',],
      tratamiento: ['',],
      examenes_axuliares: ['',],
      referencia: ['',],
      observacion: ['', ],
    })
  }

  ngOnChanges(changes: SimpleChanges) {
    if(changes['infoData'] && this.infoData) {
      this.form.patchValue(this.infoData);
      this.active.set(false);
      if(!this.canSaveForm) this.form.disable();
    }
  }

  ngOnInit() {
    if(!this.infoData) {
      this.setDefaultDate();
    }
  }

  setDefaultDate() {
    this.form.patchValue({
      fecha: this._getFormatDate(),
      hora: this._getFormatHour()
    })
  }

  private _getFormatDate(): string {
    return formatDate(new Date(), 'yyyy-MM-dd', 'en-US')
  }

  private _getFormatHour(): string {
    return formatDate(new Date(), 'HH:mm', 'en-US')
  }

  public get name() {
    return this._name();
  }

  public isValid(): boolean {
    return this.form.valid
  }

  public getData() {
    return this.form.value
  }

  public onValidateRecibo(nroRecibo: string) {
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
    for (const pago of info?.pagos_realizados) {
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

    this.form.get('costo')?.setValue('');
    this.form.get('fecha_pago')?.setValue('');
    this._toastService.add({
      type: 'warning',
      message: 'No se encontró el recibo',
      life: 5000
    })
  }

  openNursing(): void {
    this.showNursing.set(true);
    this.showMedicine.set(false);
    this.showDentistry.set(false);
    this.form.disable()
  }

  openMedicine(): void {
    this.showNursing.set(false);
    this.showMedicine.set(true);
    this.showDentistry.set(false);
    this.form.disable()
  }

  openDentistry(): void {
    this.showMedicine.set(false);
    this.showNursing.set(false);
    this.showDentistry.set(true);
    this.form.disable()
  }
}
