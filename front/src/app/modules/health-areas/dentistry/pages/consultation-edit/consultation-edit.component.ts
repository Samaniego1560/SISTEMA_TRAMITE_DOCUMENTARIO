import {Component, inject, input, OnInit, Signal, signal} from '@angular/core';
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {ConsultationFormComponent} from "../../components/consultation-form/consultation-form.component";
import {ActivatedRoute, Router} from "@angular/router";
import {v4 as uuidV4} from "uuid";
import {DentalService} from "../../../../../core/services/health-areas/dental/dental.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {catchError, of, Subscription} from "rxjs";
import {
  DentistryConsultationForm,
  DentistryConsultationRequest
} from "../../../../../core/models/areas/dentistry.model";
import {Patient} from "../../../../../core/models/areas/patient.model";
import {ConsultationRecord} from "../../../../../core/models/areas/areas.model";
import {ReloadHealthAreasService} from "../../../../../core/services/health-areas/reload-health-areas.service";

@Component({
  selector: 'app-consultation-view',
  standalone: true,
  imports: [
    BlockUiComponent,
    EcDialogComponent,
    EcTemplate,
    ConsultationFormComponent
  ],
  templateUrl: './consultation-edit.component.html',
  styleUrl: './consultation-edit.component.scss'
})
export default class ConsultationEditComponent implements OnInit {
  infoPatientDni = input.required<string>({alias: 'dni'})
  consultationID = input.required<string>()
  type = input.required<string>()
  public isBlockPage = signal<boolean>(false)
  private _subscription = new Subscription()

  private _router: Router = inject(Router)
  private _dentalService = inject(DentalService)
  private _toastService = inject(ToastService)
  public infoConsultation = signal<DentistryConsultationRequest | undefined>(undefined)
  private _reloadConsultingService = inject(ReloadHealthAreasService)

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
    this.isBlockPage.set(true);
    this._dentalService.getConsultation(this.consultationID()!).subscribe({
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
      }
    })
  }

  public onUpdateConsultation(formValue: DentistryConsultationForm) {
    const request: DentistryConsultationRequest = {
      ...formValue.consultation,
      consulta_odontologia: {
        id: this.infoConsultation()?.consulta_odontologia.id!,
        fecha_consulta: this.infoConsultation()?.consulta_odontologia.fecha_consulta!,
        paciente_id: formValue.patient.id
      }
    }
    if (request.examen_bucal) {
      request.examen_bucal.id = this.infoConsultation()?.examen_bucal?.id ?? uuidV4()
    }
    if (request.consulta) request.consulta.id = this.infoConsultation()?.consulta?.id ?? uuidV4()
    if (request.procedimiento) request.procedimiento.id = this.infoConsultation()?.procedimiento?.id ?? uuidV4()

    this.isBlockPage.set(true);
    this._subscription.add(
      this._dentalService.updateConsulting(request)
        .subscribe({
          next: (resp) => {
            this.isBlockPage.set(false);
            if (resp.error) {
              this._toastService.add({
                type: 'error',
                message: 'No se pudo registrar la consulta, intente nuevamente!',
                life: 5000
              })
              return
            }

            this._toastService.add({type: 'success', message: 'Consulta actualizada.', life: 5000})
            this._redirectToPatients()
          },
          error: (error) => {
            console.log('error', error)
            this.isBlockPage.set(false);
            this._toastService.add({
              type: 'error',
              message: 'No se pudo registrar la consulta, intente nuevamente!',
              life: 5000
            })
          }
        }))
  }

  public onCancelCreate() {
    this._redirectToPatients()
  }

  private _redirectToPatients() {
    this._router.navigate(['/home/dental-consultation']).then(
      () => this._reloadConsultingService.reloadDentistryConsulting()
    )
  }
}
