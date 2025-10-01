import {
  ChangeDetectionStrategy,
  Component,
  inject,
  Input,
  OnChanges,
  output,
  signal, SimpleChanges,
  viewChildren
} from '@angular/core';
import {
  BasicInformationPatientComponent
} from "../../../shared/components/basic-information-patient/basic-information-patient.component";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {Patient} from "../../../../../core/models/areas/patient.model";
import {OralExamFormComponent} from "../oral-exam-form/oral-exam-form.component";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {formatDate, TitleCasePipe} from "@angular/common";
import {ProcedureFormComponent} from "../procedure-form/procedure-form.component";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {IntegralAttentionFormComponent} from "../integral-attention-form/integral-attention-form.component";
import {
  DentistryConsultationBase, DentistryConsultationForm,
  DentistryConsultationRequest
} from "../../../../../core/models/areas/dentistry.model";

@Component({
  selector: 'app-consultation-form',
  standalone: true,
  imports: [
    BasicInformationPatientComponent,
    EcTableModule,
    FormsModule,
    ReactiveFormsModule,
    OralExamFormComponent,
    TitleCasePipe,
    ProcedureFormComponent,
    IntegralAttentionFormComponent,
  ],
  templateUrl: './consultation-form.component.html',
  styleUrl: './consultation-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConsultationFormComponent implements OnChanges {
  @Input() infoPatientDni: string | undefined;
  @Input() infoConsultation: DentistryConsultationRequest | undefined;
  save = output<DentistryConsultationForm>();
  cancel = output<void>();

  public listExams: Set<string> = new Set<string>();
  public typesServices = [
    'consulta',
    'examen',
    'procedimiento',
  ]
  attentionDate = signal<string>(formatDate(new Date(), 'yyyy-MM-dd HH:mm', 'en-US'));
  public examSelected = signal<string>('');
  private _formExams = viewChildren<FormBaseComponent<unknown>>('formExams')
  public patient = signal<Patient | undefined>(undefined);
  public touchedForms = signal<boolean>(false);

  private _toastService = inject(ToastService);

  get isModeCreation() {
    return !this.infoConsultation;
  }

  constructor() {
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoConsultation']?.currentValue && this.infoConsultation) {
      this.attentionDate.set(this.infoConsultation.consulta_odontologia.fecha_consulta)

      if(this.infoConsultation?.consulta) {
        this.listExams.add('consulta')
      }
      if(this.infoConsultation?.examen_bucal) {
        this.listExams.add('examen')
      }
      if(this.infoConsultation?.procedimiento) {
        this.listExams.add('procedimiento')
      }
    }
  }

  public onAddExam(typeExam: string) {
    this.listExams.add(typeExam);
    this.examSelected.set('')
  }

  public onDeleteExam(typeExam: string) {
    this.listExams.delete(typeExam);
  }

  private _validateFormExams() {
    return this._formExams().every(dd => dd.isValid())
  }

  public onSaveForm() {
    if(!this._formExams().length) {
      this._toastService.add({
        type: 'warning',
        message: 'Es necesario agregar almenos un examen'
      });
      return
    }

    if (!this._formExams().length) {
      this._toastService.add({
        type: 'error',
        message: 'Debe de agregar al menos un servicio'
      });
      return;
    }

    const isValidAll = this._validateFormExams()
    if (!isValidAll) {
      this.touchedForms.set(true)
      this._toastService.add({
        type: 'warning',
        message: 'Complete correctamente todos los campos requeridos'
      });
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

    const consulta_enfermeria: DentistryConsultationBase = {
      ...this._getExams()
    }

    this.save.emit({
      patient: this.patient()!,
      attentionDate: formatDate(this.attentionDate(), 'yyyy-MM-dd HH:mm', 'en-US'),
      consultation: consulta_enfermeria
    })
  }

  private _getExams(): Record<keyof DentistryConsultationBase, any> {
    return this._formExams().reduce((acc, dd) => {
        return {
          ...acc,
          [dd.name]: dd.getData()
        }
      },
      <Record<keyof DentistryConsultationRequest, any>>{})
  }

  changeAttentionDate(event:Event){
    this.attentionDate.set((<HTMLInputElement>event.target).value)
  }
}
