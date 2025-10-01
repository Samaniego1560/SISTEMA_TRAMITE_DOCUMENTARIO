import {
  ChangeDetectionStrategy,
  Component, ElementRef,
  inject, Input,
  OnChanges,
  output,
  signal, SimpleChanges, viewChild, viewChildren,
} from '@angular/core';
import {FormSubmitDirective} from "../../../../../lib/validator-dynamic/directives/form-submit.directive";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule} from "@angular/forms";
import {Patient} from "../../../../../core/models/areas/patient.model";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {formatDate, TitleCasePipe} from "@angular/common";
import {CompanionFormComponent} from "../companion-form/companion-form.component";
import {RoutineReviewFormComponent} from "../routine-review-form/routine-review-form.component";
import {PhysicalExamFormComponent} from "../physical-exam-form/physical-exam-form.component";
import {VisualExamFormComponent} from "../visual-exam-form/visual-exam-form.component";
import {PreferentialExamFormComponent} from "../preferential-exam-form/preferential-exam-form.component";
import {LaboratoryExamFormComponent} from "../laboratory-exam-form/laboratory-exam-form.component";
import {SexualityExamFormComponent} from "../sexuality-exam-form/sexuality-exam-form.component";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {
  Exam, NursingConsultationForm,
} from "../../../../../core/models/areas/nursing/nursing.model";
import {
  BasicInformationPatientComponent
} from "../../../shared/components/basic-information-patient/basic-information-patient.component";
import {VaccinesFormComponent} from "../vaccines-form/vaccines-form.component";
import {DrugTreatmentsComponent} from "../drug-treatments/drug-treatments.component";
import {ProceduresPerformedComponent} from "../procedures-performed/procedures-performed.component";
import {IntegralAttentionFormComponent} from "../integral-attention-form/integral-attention-form.component";
import {NursingConsultationRequest} from "../../../../../core/models/areas/nursing/request.model";
import { timer } from 'rxjs';
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {PermissionsService} from "../../../../../core/services/permissions/permissions.service";
import {ACTION_TYPES} from "../../../../../core/contans/areas/permissions.constant";

@Component({
  selector: 'app-nursing-consultations-form',
  standalone: true,
  imports: [
    FormSubmitDirective,
    FormsModule,
    ReactiveFormsModule,
    CompanionFormComponent,
    RoutineReviewFormComponent,
    PhysicalExamFormComponent,
    TitleCasePipe,
    VisualExamFormComponent,
    PreferentialExamFormComponent,
    LaboratoryExamFormComponent,
    SexualityExamFormComponent,
    BasicInformationPatientComponent,
    VaccinesFormComponent,
    DrugTreatmentsComponent,
    ProceduresPerformedComponent,
    IntegralAttentionFormComponent,
    ControlErrorsDirective,
  ],
  templateUrl: './consultations-form.component.html',
  styleUrl: './consultations-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConsultationsFormComponent implements OnChanges {
  @Input() infoPatientDni: string | undefined;
  @Input() infoConsultation: NursingConsultationRequest | undefined;
  save = output<NursingConsultationForm>();
  cancel = output<void>();

  _formExams = viewChildren<FormBaseComponent<unknown>>('formExams')

  public consultingForm: FormGroup;
  public typesServices = [
    'consulta',
    'vacunas',
    'fisico',
    'visual',
    'preferencial',
    'laboratorio',
    'sexualidad',
    'procedimientos',
    'tratamientos',
  ]
  public listExams: Set<string> = new Set<string>();
  public patient = signal<Patient | undefined>(undefined);
  attentionDate = signal<string>(formatDate(new Date(), 'yyyy-MM-dd HH:mm', 'en-US'));

  private _permissionsService: PermissionsService = inject(PermissionsService)
  private _fb = inject(FormBuilder);
  private _toastService: ToastService = inject(ToastService)
  public selectExam = viewChild<ElementRef<HTMLSelectElement>>('selectExam');
  public examSelected = signal<string>('')

  get isModeCreation() {
    return !this.infoConsultation;
  }

  get canSaveForm() {
    return this.canCreate || this.canUpdate;
  }

  get canViewForm() {
    return this.canRead;
  }

  get canCreate() {
    return this.isModeCreation && this._permissionsService.hasPermission(
      'enfermería',
      'atenciones',
      ACTION_TYPES.CREATE
    )
  }

  get canUpdate() {
    return !this.isModeCreation && this._permissionsService.hasPermission(
      'enfermería',
      'atenciones',
      ACTION_TYPES.UPDATE
    )
  }

  get canRead() {
    return this._permissionsService.hasPermission(
      'enfermería',
      'atenciones',
      ACTION_TYPES.READ
    )
  }

  constructor() {
    this.consultingForm = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoConsultation']?.currentValue && this.infoConsultation) {
      this.attentionDate.set(this.infoConsultation.consulta_enfermeria.fecha_consulta)
      this._patchForms();

      const exams = this.infoConsultation.examenes;
      if (exams?.consulta_general) {
        this.listExams.add('consulta')
      }
      if (exams?.vacunas) {
        this.listExams.add('vacunas')
      }
      if (exams?.examen_fisico) {
        this.listExams.add('fisico')
      }
      if (exams?.examen_visual) {
        this.listExams.add('visual')
      }
      if (exams?.examen_preferencial) {
        this.listExams.add('preferencial')
      }
      if (exams?.examen_laboratorio) {
        this.listExams.add('laboratorio')
      }
      if (exams?.examen_sexualidad) {
        this.listExams.add('sexualidad')
      }
      if (this.infoConsultation?.examenes?.procedimiento_realizado) {
        this.listExams.add('procedimientos')
      }
      if (this.infoConsultation?.examenes?.tratamiento_medicamentoso) {
        this.listExams.add('tratamientos')
      }
    }
  }

  private _buildForm(): FormGroup {
    return this._fb.group({
      datos_acompanante: [null],
      revision_rutina: [{
        fiebre_ultimo_quince_dias: 'no',
        tos_mas_quince_dias: 'no',
        secrecion_lesion_genitales: 'no',
        fecha_ultima_regla: 'no',
        comentarios: '',
      }],
    })
  }

  private _patchForms() {
    if (this.infoConsultation?.datos_acompanante) {
      this.consultingForm.get('datos_acompanante')?.patchValue(this.infoConsultation.datos_acompanante);
    }
    if (this.infoConsultation?.revision_rutina) {
      this.consultingForm.get('datos_acompanante')?.patchValue(this.infoConsultation.revision_rutina);
    }

    this.consultingForm.disable();
  }

  public saveConsulting() {
    if (this.consultingForm.invalid) {
      this.consultingForm.markAllAsTouched();
      this._toastService.add({
        type: 'error',
        message: 'Formulario inválido'
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

    const validAll = this._validateFormExams()
    if (!validAll) {
      this._toastService.add({
        type: 'error',
        message: 'Hay servicios pendientes a completar, se deben completar todos los servicios seleccionados'
      });
      return;
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

    const exams = this._getExams();
    const consulting: NursingConsultationForm = {
      patient: this.patient()!,
      attentionDate: formatDate(this.attentionDate(), 'yyyy-MM-dd HH:mm', 'en-US'),
      consultation: {
        datos_acompanante: this.consultingForm.get('datos_acompanante')?.value,
        revision_rutina: this.consultingForm.get('revision_rutina')?.value,
        examenes: exams
      },
    }

    this.save.emit(consulting);
  }

  private _validateFormExams() {
    return this._formExams().every(dd => dd.isValid())
  }

  private _getExams(): Exam {
    return this._formExams().reduce((acc, dd) => {
        return {
          ...acc,
          [dd.name]: dd.getData()
        }
      },
      <Record<string, Exam[keyof Exam]>>{})
  }

  public onAddExam(typeExam: string) {
    this.listExams.add(typeExam);
    this.examSelected.set('');
    if (this.selectExam()?.nativeElement?.value) {
      this.selectExam()!.nativeElement!.value = '';
    }
  }

  public onDeleteExam(typeExam: string) {
    this.listExams.delete(typeExam);
  }

  changeAttentionDate(event:Event){
    this.attentionDate.set((<HTMLInputElement>event.target).value)
  }
}
