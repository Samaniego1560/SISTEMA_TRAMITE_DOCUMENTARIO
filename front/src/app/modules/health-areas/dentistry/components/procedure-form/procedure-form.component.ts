import {
  ChangeDetectionStrategy,
  Component,
  inject,
  Input,
  OnChanges,
  output,
  signal,
  SimpleChanges
} from '@angular/core';
import {FormBuilder, FormGroup, ReactiveFormsModule,} from "@angular/forms";
import {NgxValidators} from "../../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {formatDate, UpperCasePipe} from "@angular/common";
import {ProcedureExam} from "../../../../../core/models/areas/dentistry.model";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {OnlyNumbersDirective} from "../../../../../core/directives/only-numbers.directive";
import {DentalService} from "../../../../../core/services/health-areas/dental/dental.service";
import {PaymentConcept} from "../../../../../core/models/areas/areas.model";
import {Response} from "../../../../../core/models/global.model";
import {HttpErrorResponse} from "@angular/common/http";

@Component({
  selector: 'app-procedure-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ControlErrorsDirective,
    EcTableModule,
    BlockUiComponent,
    OnlyNumbersDirective,
    UpperCasePipe
  ],
  templateUrl: './procedure-form.component.html',
  styleUrl: './procedure-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ProcedureFormComponent extends FormBaseComponent<ProcedureExam> implements OnChanges {
  @Input({required: true}) dniPatient: string | undefined;
  public isBlockPage = signal<boolean>(false)

  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }

  deleteForm = output();

  public form: FormGroup;
  isPaymentProcedure = signal(false)
  public active = signal<boolean>(true);
  public successPayment = signal<boolean>(true);
  private _name = signal('procedimiento');

  private _fb = inject(FormBuilder);
  private _dentalService = inject(DentalService);
  private _toastService = inject(ToastService);

  public typeProcedure: { name: string, isPaid: boolean }[] = [
    {name: 'profilaxis dental', isPaid: true},
    {name: 'destartraje dental', isPaid: true},
    {name: 'aplicación de flúor gel', isPaid: true},
    {name: 'exodoncia simple', isPaid: true},
    {name: 'exodoncia compleja', isPaid: true},
    {name: 'curetaje alveolar', isPaid: true},
    {name: 'resina simple', isPaid: true},
    {name: 'resina compuesta', isPaid: true},
    {name: 'apertura cameral', isPaid: true},
    {name: 'endodoncia anterior', isPaid: true},
    {name: 'endodoncia posterior', isPaid: true},
    {name: 'radiografía', isPaid: true},
    {name: 'cementado de corona', isPaid: true},
    {name: 'prevención colutorio', isPaid: false},
    {name: 'prevención IHO - TCA cepillado', isPaid: false},
    {name: 'otros', isPaid: true},
  ];

  constructor() {
    super();
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoData']?.currentValue) {
      this.form.patchValue(this.infoData!);
      this.validatePayment();
      this.form.disable();
      this.active.set(false);
    }
  }

  private _buildForm() {
    return this._fb.nonNullable.group({
      tipo_procedimiento: ['', [NgxValidators.required()]],
      recibo: [{value: '', disabled: true}, [NgxValidators.required()]],
      costo: [{value: '', disabled: true}, [NgxValidators.required()]],
      fecha_pago: [{value: '', disabled: true}, [NgxValidators.required()]],
      pieza_dental: ['', [NgxValidators.required()]],
      comentarios: ['',],
    })
  }

  public get name() {
    return this._name();
  }

  public isValid(): boolean {
    return this.form.valid && this.successPayment();
  }

  public getData() {
    return this.form.getRawValue();
  }

  public onValidateRecibo(nroRecibo: string) {
    this.isBlockPage.set(true)
    this._dentalService.getPayment({
      dni: this.dniPatient!,
      recibo: nroRecibo,
      tipo_servicio: 'procedimiento',
      nombre_servicio: this.form.get('tipo_procedimiento')?.value,
    })
      .subscribe({
        next: (res: Response<PaymentConcept | null>) => {
          this.isBlockPage.set(false)
          this._toastService.add({
            type: this.getTypeMessage(res.type),
            message: res.msg ?? 'Recibo validado correctamente',
            life: 6000
          })
          if(res.data) this._validatePayment(res.data);
        },
        error: (error) => {
          this._toastService.add({
            type: 'warning',
            message: error?.error?.msg ?? 'Ocurrió un error al validar el recibo, intente nuevamente',
            life: 5000
          })
          this.isBlockPage.set(false)
          this.invalidPayment()
        }
      })
  }

  private _validatePayment(payment: PaymentConcept | null) {
    if (payment !== null) {
      this.form.get('costo')?.setValue(payment.precio_unit);
      this.form.get('fecha_pago')?.setValue(formatDate(payment.fecha_pago, 'yyyy-MM-dd', 'en-US'));
      this.successPayment.set(true);
      return
    }

    this.invalidPayment()
    this._toastService.add({
      type: 'warning',
      message: 'No se encontró el recibo',
      life: 5000
    })
  }

  invalidPayment(): void {
    this.form.get('costo')?.setValue('');
    this.form.get('fecha_pago')?.setValue('');
    this.successPayment.set(false);
  }

  validatePaymentProcedure() {
    const typeProcedure = this.form.get('tipo_procedimiento')?.value;
    const isPaid = this.typeProcedure.find(procedure => procedure.name === typeProcedure)?.isPaid;
    if (!isPaid) {
      this.form.get('recibo')?.setValue('');
      this.form.get('costo')?.setValue('');
      this.form.get('fecha_pago')?.setValue('');
      this.successPayment.set(true);
      this.isPaymentProcedure.set(false);
      this.form.get('recibo')?.disable();
      return
    }

    this.form.get('recibo')?.setValue('');
    this.form.get('costo')?.setValue('');
    this.form.get('fecha_pago')?.setValue('');
    this.successPayment.set(false);
    this.isPaymentProcedure.set(true);
    this.form.get('recibo')?.enable();
  }

  validatePayment() {
    const typeProcedure = this.form.get('tipo_procedimiento')?.value;
    const isPaid = this.typeProcedure.find(procedure => procedure.name === typeProcedure)?.isPaid;
    if (!isPaid) {
      this.isPaymentProcedure.set(false);
      return
    }

    this.isPaymentProcedure.set(true);
  }

  getTypeMessage(type: string) {
    if(type === 'ERROR') return 'error';
    if(type === 'WARNING') return 'warning';
    if(type === 'SUCCESS') return 'success';
    return 'info';
  }

}
