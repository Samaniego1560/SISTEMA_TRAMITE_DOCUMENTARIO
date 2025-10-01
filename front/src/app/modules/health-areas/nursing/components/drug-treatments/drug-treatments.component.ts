import {
  ChangeDetectionStrategy,
  Component,
  inject,
  Input,
  OnChanges, OnInit,
  output,
  signal,
  SimpleChanges
} from '@angular/core';
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {
  TratamientoMedicamentoso
} from "../../../../../core/models/areas/nursing/nursing.model";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {formatDate, TitleCasePipe} from "@angular/common";
import {FormSubmitDirective} from "../../../../../lib/validator-dynamic/directives/form-submit.directive";
import {NgxValidators} from "../../../../../lib/validator-dynamic/utils/ngx-validators.util";

@Component({
  selector: 'app-drug-treatments',
  standalone: true,
  imports: [
    FormSubmitDirective,
    ControlErrorsDirective,
    FormsModule,
    ReactiveFormsModule,
    EcTableModule,
    TitleCasePipe,
  ],
  templateUrl: './drug-treatments.component.html',
  styleUrl: './drug-treatments.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DrugTreatmentsComponent extends FormBaseComponent<TratamientoMedicamentoso> implements OnChanges, OnInit {
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }

  deleteForm = output();

  public active = signal<boolean>(true);
  public form: FormGroup;
  private _name = signal('tratamiento_medicamentoso');

  areas = [
    'odontología',
    'medicina',
    "psicología",
    "otros"
  ]

  private _fb = inject(FormBuilder);

  constructor() {
    super();
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if(changes['infoData'] && this.infoData) {
      this.form.patchValue(this.infoData);
      this.active.set(false);
      if(!this.canSaveForm) this.form.disable();
    }
  }

  ngOnInit() {
    this.setDefaultDate();
  }

  private _buildForm() {
    return this._fb.nonNullable.group({
      fecha: ['',],
      hora_aplicacion: ['',],
      nombre_generico_medicamento: ['', NgxValidators.required()],
      via_administracion: ['', NgxValidators.required()],
      responsable_atencion: ['',],
      area_solicitante: [''],
      especialista_solicitante: [''],
      observaciones: ['',],
    })
  }

  setDefaultDate() {
    this.form.patchValue({
      fecha: this._getFormatDate(),
      hora_aplicacion: this._getFormatHour()
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

  isValid(): boolean {
    return this.form.valid;
  }

  getData() {
    return this.form.value;
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
