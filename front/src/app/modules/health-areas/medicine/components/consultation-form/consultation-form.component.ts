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
import {
  BasicInformationPatientComponent
} from "../../../shared/components/basic-information-patient/basic-information-patient.component";
import {FormSubmitDirective} from "../../../../../lib/validator-dynamic/directives/form-submit.directive";
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule
} from "@angular/forms";
import {Patient} from "../../../../../core/models/areas/patient.model";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {DataForm} from "../../../../../core/models/areas/areas.model";
import {NursingConsultationRequest} from "../../../../../core/models/areas/nursing/request.model";
import {ConsultaMedicinaGeneral} from "../../../../../core/models/areas/nursing/nursing.model";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {formatDate} from "@angular/common";

@Component({
  selector: 'medicine-consultation-form',
  standalone: true,
    imports: [
        BasicInformationPatientComponent,
        FormSubmitDirective,
        FormsModule,
        ReactiveFormsModule,
        ControlErrorsDirective,
    ],
  templateUrl: './consultation-form.component.html',
  styleUrl: './consultation-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConsultationFormComponent implements OnChanges {
  @Input() infoPatientDni: string | undefined;
  @Input() infoConsultation: NursingConsultationRequest | undefined;
  public patient = signal<Patient | undefined>(undefined);
  save = output<DataForm<Omit<ConsultaMedicinaGeneral, 'id'>>>();
  cancel = output<void>();

  private _fb = inject(FormBuilder);
  private _toastService = inject(ToastService);

  public form: FormGroup;
  showInformacionGeneral = signal<boolean>(true);
  showExamenFisico = signal<boolean>(true);
  showConsulting = signal<boolean>(true);

  attentionDate = signal<string>(formatDate(new Date(), 'yyyy-MM-dd HH:mm', 'en-US'));

  get canSaveForm() {
    return this.form.enabled;
  }

  constructor() {
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoConsultation'] && this.infoConsultation?.examenes.consulta_medicina_general) {
      this.form.patchValue(this.infoConsultation.examenes?.consulta_medicina_general);
      this.attentionDate.set(this.infoConsultation.examenes.consulta_medicina_general.fecha_hora)
      this.form.disable()
    }
  }

  private _buildForm(): FormGroup {
    return this._fb.group({
      fecha_hora: [''],
      anamnesis: [''],
      examen_clinico: [''],
      indicaciones: [''],
    })
  }

  public onSaveForm() {
    if(this.form.disabled) return

    if (this.form.invalid) {
      this._toastService.add({
        type: 'warning',
        message: 'Por favor, complete los campos requeridos'
      })
      return
    }

    if (this.attentionDate() === '') {
      this._toastService.add({
        type: 'error',
        message: 'Es obligatorio ingresar la fecha de la consulta'
      });
      return;
    }

    const attentionDate = new Date(this.attentionDate());
    if(isNaN(attentionDate.getTime())){
      this._toastService.add({
        type: 'error',
        message: 'La fecha de la consulta es inválida'
      });
      return;
    }

    const dateCurrent = new Date()
    const isInvalidDate = attentionDate.getTime() > dateCurrent.getTime()
    if (isInvalidDate) {
      this._toastService.add({
        type: 'error',
        message: 'La fecha de la consulta no puede ser mayor a la fecha actual'
      });
      return;
    }

    const data: DataForm<Omit<ConsultaMedicinaGeneral, 'id'>> = {
      patient: this.patient()!,
      data: {
        ...this.form.value,
        fecha_hora: formatDate(this.attentionDate(), 'yyyy-MM-dd HH:mm', 'en-US'),
      }
    }
    this.save.emit(data);
  }

  changeAttentionDate(event:Event){
    this.attentionDate.set((<HTMLInputElement>event.target).value)
  }
}
