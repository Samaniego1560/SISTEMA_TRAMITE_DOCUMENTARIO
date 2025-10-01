import {Component, inject, input, OnInit, signal} from '@angular/core';
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {Subscription} from "rxjs";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {Router} from "@angular/router";
import {
  ConsultationsFormComponent
} from "../../components/consultation-form/consultations-form.component";
import {
  Exam,
  NursingConsultationForm,
} from "../../../../../core/models/areas/nursing/nursing.model";
import {v4 as uuidV4} from 'uuid';
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {NursingConsultationRequest} from "../../../../../core/models/areas/nursing/request.model";

@Component({
  selector: 'consultation-view',
  standalone: true,
  imports: [
    BlockUiComponent,
    EcDialogComponent,
    EcTemplate,
    ConsultationsFormComponent
  ],
  templateUrl: './consultation-view.component.html',
  styleUrl: './consultation-view.component.scss'
})
export default class ConsultationViewComponent implements OnInit {
  infoPatientDni = input.required<string>({alias: 'dni'})
  consultationID = input.required<string>()
  public isBlockPage = signal(true);
  private _subscription = new Subscription();

  private _toastService: ToastService = inject(ToastService)
  private _nursingConsultingService: NursingService = inject(NursingService)
  private _router: Router = inject(Router)
  public infoConsultation = signal<NursingConsultationRequest | undefined>(undefined)

  ngOnInit() {
    this._processPatient()
  }

  private _processPatient() {
    if (!this.consultationID() || !this.infoPatientDni()) {
      this._redirectToPatients()
      return
    }

    this._getConsultation()
  }

  private _getConsultation() {
    this._nursingConsultingService.getConsultation(this.consultationID()).subscribe({
      next: (resp: any) => {

        this.isBlockPage.set(false);
        if (resp.error) {
          this._toastService.add({
            type: 'error',
            message: 'No se pudo obtener la consulta, intente nuevamente!',
            life: 5000
          })
          this._redirectToPatients()
          return
        }

        this.infoConsultation.set(resp.data)
      },
      error: () => {
        this.isBlockPage.set(false);
        this._toastService.add({
          type: 'error',
          message: 'No se pudo obtener la consulta',
          life: 5000
        })
        this._redirectToPatients()
      }
    })
  }

  public onEditConsultation(formValue: NursingConsultationForm) {

    const consultingRequest: NursingConsultationRequest = this._requestUpdateConsulting(formValue)

    this.isBlockPage.set(true);

    this._subscription.add(
      this._nursingConsultingService.updateConsulting(consultingRequest).subscribe({
        next: (resp: any) => {
          this.isBlockPage.set(false);
          if (resp.error) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo actualizar la consulta, intente nuevamente!',
              life: 5000
            })
            return
          }

          this._toastService.add({type: 'success', message: 'Consulta actualizada.', life: 5000})
          this._redirectToPatients()
        }
      })
    )
  }

  private _requestUpdateConsulting(formValue: NursingConsultationForm) {
    const request: NursingConsultationRequest = {
      ...formValue.consultation,
      consulta_enfermeria: {
        id: this.infoConsultation()?.consulta_enfermeria.id!,
        fecha_consulta: this.infoConsultation()?.consulta_enfermeria.fecha_consulta!,
        paciente_id: this.infoConsultation()?.consulta_enfermeria.paciente_id!,
      }
    }
    request.revision_rutina.id = this.infoConsultation()?.revision_rutina.id ?? uuidV4()
    if(request.datos_acompanante) {
      request.datos_acompanante.id = this.infoConsultation()?.datos_acompanante?.id ?? uuidV4()
    }
    request.examenes = this._mapExamsCreate(formValue.consultation.examenes)

    return request
  }

  private _mapExamsCreate(exams: Exam) {
    if (exams.procedimiento_realizado) {
      exams.procedimiento_realizado.id = this.infoConsultation()?.examenes?.procedimiento_realizado?.id ?? uuidV4()
    }

    if (exams.tratamiento_medicamentoso) {
      exams.tratamiento_medicamentoso.id = this.infoConsultation()?.examenes?.tratamiento_medicamentoso?.id ?? uuidV4()
    }

    if (exams.examen_fisico) {
      exams.examen_fisico.id = this.infoConsultation()?.examenes?.examen_fisico?.id ?? uuidV4()
    }

    if (exams.examen_laboratorio) {
      exams.examen_laboratorio.id = this.infoConsultation()?.examenes?.examen_laboratorio?.id ?? uuidV4()
    }

    if (exams.examen_preferencial) {
      exams.examen_preferencial.id = this.infoConsultation()?.examenes?.examen_preferencial?.id ?? uuidV4()
    }

    if (exams.examen_sexualidad) {
      exams.examen_sexualidad.id = this.infoConsultation()?.examenes?.examen_sexualidad?.id ?? uuidV4()
    }

    if (exams.examen_visual) {
      exams.examen_visual.id = this.infoConsultation()?.examenes?.examen_visual?.id ?? uuidV4()
    }

    if (exams.consulta_general) {
      exams.consulta_general.id = this.infoConsultation()?.examenes?.consulta_general?.id ?? uuidV4()
    }
    return exams
  }

  public onCancelEdit() {
    this._redirectToPatients()
  }

  private _redirectToPatients() {
    this._router.navigate(['home/nursing-consultation']).then()
  }
}
