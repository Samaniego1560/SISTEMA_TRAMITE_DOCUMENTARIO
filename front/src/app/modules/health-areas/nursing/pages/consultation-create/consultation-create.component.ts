import {Component, inject, OnInit, signal} from '@angular/core';
import {Subscription} from "rxjs";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {Router} from "@angular/router";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {
  ConsultationsFormComponent
} from "../../components/consultation-form/consultations-form.component";
import {
  Exam,
  NursingConsultationForm
} from "../../../../../core/models/areas/nursing/nursing.model";
import {v4 as uuidV4} from "uuid";
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {NursingConsultationRequest} from "../../../../../core/models/areas/nursing/request.model";
import {ReloadHealthAreasService} from "../../../../../core/services/health-areas/reload-health-areas.service";

@Component({
  selector: 'nursing-consultation-create',
  standalone: true,
  imports: [
    BlockUiComponent,
    EcDialogComponent,
    EcTemplate,
    ConsultationsFormComponent
  ],
  templateUrl: './consultation-create.component.html',
  styleUrl: './consultation-create.component.scss'
})
export default class ConsultationCreateComponent implements OnInit {
  public isBlockPage = signal(false);
  private _subscription = new Subscription()

  private _toastService: ToastService = inject(ToastService)
  private _router: Router = inject(Router)
  private _nursingConsultingService: NursingService = inject(NursingService)
  private _reloadConsultingService = inject(ReloadHealthAreasService)

  ngOnInit() {
  }

  public onCreatePatient(formValue: NursingConsultationForm) {
    const consultingRequest: NursingConsultationRequest = this._requestCreateConsulting(formValue)

    this.isBlockPage.set(true);

    this._subscription.add(
      this._nursingConsultingService.createConsulting(consultingRequest).subscribe({
        next: (resp: any) => {
          this.isBlockPage.set(false);
          if (resp.error) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo crear la consulta, intente nuevamente!',
              life: 5000
            })
            return
          }

          this._toastService.add({type: 'success', message: 'Consulta creada.', life: 5000})
          this._redirectToPatients()
        },
        error: () => {
          this.isBlockPage.set(false);
          this._toastService.add({
            type: 'error',
            message: 'No se pudo crear la consulta, intente nuevamente!',
            life: 5000
          })
        }
      })
    )
  }

  private _requestCreateConsulting(formValue: NursingConsultationForm) {
    const request: NursingConsultationRequest = {
      ...formValue.consultation,
      consulta_enfermeria: {
        id: uuidV4(),
        fecha_consulta: formValue.attentionDate,
        paciente_id: formValue.patient.id
      }
    }
    request.consulta_enfermeria.id = uuidV4()
    request.revision_rutina.id = uuidV4()
    if (request.datos_acompanante) {
      request.datos_acompanante.id = uuidV4()
    }

    request.examenes = this._mapExamsCreate(formValue.consultation.examenes)

    return request
  }

  private _mapExamsCreate(exams: Exam) {
    if (exams.procedimiento_realizado) {
      exams.procedimiento_realizado.id = uuidV4()
    }

    if (exams.tratamiento_medicamentoso) {
      exams.tratamiento_medicamentoso.id = uuidV4()
    }

    if (exams.examen_fisico) {
      exams.examen_fisico.id = uuidV4()
    }

    if (exams.examen_laboratorio) {
      exams.examen_laboratorio.id = uuidV4()
    }

    if (exams.examen_preferencial) {
      exams.examen_preferencial.id = uuidV4()
    }

    if (exams.examen_sexualidad) {
      exams.examen_sexualidad.id = uuidV4()
    }

    if (exams.examen_visual) {
      exams.examen_visual.id = uuidV4()
    }

    if (exams.consulta_general) {
      exams.consulta_general.id = uuidV4()
    }

    return exams
  }

  public onCancelCreate() {
    this._redirectToPatients()
  }

  private _redirectToPatients() {
    this._router.navigate(['/home/nursing-consultation']).then(
      () => this._reloadConsultingService.reloadNursingConsulting()
    )
  }
}
