import {Component, inject, input, OnInit, Signal, signal} from '@angular/core';
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {ActivatedRoute, Router} from "@angular/router";
import {v4 as uuidV4} from "uuid";
import {DentalService} from "../../../../../core/services/health-areas/dental/dental.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {catchError, of, Subscription} from "rxjs";
import {ConsultationRecord, DataForm} from "../../../../../core/models/areas/areas.model";
import {NursingConsultationRequest} from "../../../../../core/models/areas/nursing/request.model";
import {MedicineService} from "../../../../../core/services/health-areas/medicine/medicine.service";
import {ReloadHealthAreasService} from "../../../../../core/services/health-areas/reload-health-areas.service";
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {ConsultaBucalGeneral, ConsultaMedicinaGeneral} from "../../../../../core/models/areas/nursing/nursing.model";
import {
  ConsultationFormComponent
} from "../../components/consultation-nursing-form/consultation-nursing-form.component";

@Component({
  selector: 'app-consultation-view',
  standalone: true,
  imports: [
    BlockUiComponent,
    EcDialogComponent,
    EcTemplate,
    ConsultationFormComponent,
  ],
  templateUrl: './consultation-edit.component.html',
  styleUrl: './consultation-edit.component.scss'
})
export default class ConsultationEditComponent implements OnInit {
  infoPatientDni = input.required<string>({alias: 'dni'})
  consultationID = input.required<string>()
  public isBlockPage = signal<boolean>(true);
  private _subscription = new Subscription();

  private _toastService: ToastService = inject(ToastService)
  private _nursingService: NursingService = inject(NursingService)
  private _router: Router = inject(Router)
  public infoConsultation = signal<NursingConsultationRequest | undefined>(undefined)
  private _reloadConsultingService = inject(ReloadHealthAreasService)

  constructor() {
  }

  ngOnInit() {
    this._processPatient()
  }

  private _processPatient() {
    if (!this.consultationID() || !this.infoPatientDni()) {
      this._redirect()
      return
    }

    this._getConsultation()
  }

  private _getConsultation() {
    this._nursingService.getConsultation(this.consultationID()!).subscribe({
      next: (resp: any) => {
        this.isBlockPage.set(false);
        if (resp.error || !resp.data) {
          this._toastService.add({
            type: 'error',
            message: 'No se pudo obtener la consulta, intente nuevamente!',
            life: 5000
          })
          this._redirect()
          return
        }

        this.infoConsultation.set(resp.data)
      },
      error: () => {
        this.isBlockPage.set(false);
        this._toastService.add({
          type: 'error',
          message: 'No se pudo obtener la consulta, intente nuevamente!',
          life: 5000
        })
        this._redirect()
      }
    })
  }

  public onEditConsultation(dataForm: DataForm<Omit<ConsultaBucalGeneral, "id">>) {
    const request = this.infoConsultation();
    this.isBlockPage.set(true);

    if(request?.examenes) {
      request.examenes = {
        ...request.examenes,
        consulta: {
          ...dataForm.data,
          id: request?.examenes['consulta']?.id ?? uuidV4()
        }
      }
    }

    this._subscription.add(
      this._nursingService.updateConsulting(request!).subscribe({
        next: (resp: any) => {
          this.isBlockPage.set(false);
          if (resp.error) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo actualizar la información del paciente, intente nuevamente!',
              life: 5000
            })
            return
          }

          this._toastService.add({type: 'success', message: 'Paciente actualizado.', life: 5000})
          this._redirect()
        },
        error:() => {
          this.isBlockPage.set(false);
          this._toastService.add({
            type: 'error',
            message: 'No se pudo actualizar la información de la consulta, intente nuevamente!',
            life: 5000
          })
        }
      })
    )
  }

  public onCancelEdit() {
    this._redirect()
  }

  private _redirect() {
    this._router.navigate(['home/dental-consultation']).then(
      () => this._reloadConsultingService.reloadNursingConsulting()
    )
  }
}
