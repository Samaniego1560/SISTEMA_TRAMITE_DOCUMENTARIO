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
import {ConsultaBucalGeneral, ConsultaMedicinaGeneral} from "../../../../../core/models/areas/nursing/nursing.model";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {OnlyNumbersDirective} from "../../../../../core/directives/only-numbers.directive";
import {ToggleSwitchComponent} from "../../../../../core/ui/toggle/toggle-switch.component";

@Component({
  selector: 'dentistry-nursing-consultation-form',
  standalone: true,
    imports: [
        BasicInformationPatientComponent,
        FormSubmitDirective,
        FormsModule,
        ReactiveFormsModule,
        ControlErrorsDirective,
    ],
  templateUrl: './consultation-nursing-form.component.html',
  styleUrl: './consultation-nursing-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConsultationFormComponent implements OnChanges {
  @Input() infoPatientDni: string | undefined;
  @Input() infoConsultation: NursingConsultationRequest | undefined;
  public patient = signal<Patient | undefined>(undefined);
  save = output<DataForm<Omit<ConsultaBucalGeneral, 'id'>>>();
  cancel = output<void>();

  private _fb = inject(FormBuilder);
  private _toastService = inject(ToastService);

  public form: FormGroup<any>;
  showInformacionGeneral = signal<boolean>(true);
  showExamenFisico = signal<boolean>(true);

  constructor() {
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoConsultation'] && this.infoConsultation?.examenes.consulta) {
      this.form.patchValue(this.infoConsultation.examenes?.consulta);
      this.form.disable()
    }
  }

  private _buildForm(): FormGroup {
    return this._fb.nonNullable.group({
      relato: [''],
      diagnostico: [''],
      examen_clinico: [''],
      examen_auxiliar: [''],
      tratamiento: [''],
      indicaciones: [''],
      comentarios: [''],
    })
  }

  public onSaveForm() {
    if (this.form.invalid) {
      this._toastService.add({
        type: 'warning',
        message: 'Por favor, complete los campos requeridos'
      })
      return
    }

    const data: DataForm<Omit<ConsultaBucalGeneral, 'id'>> = {
      patient: this.patient()!,
      data: {
        ...this.form.value
      }
    }
    this.save.emit(data);
  }
}
