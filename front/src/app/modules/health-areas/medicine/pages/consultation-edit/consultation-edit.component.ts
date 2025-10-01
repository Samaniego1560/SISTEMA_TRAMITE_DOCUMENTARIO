import {Component, inject, input, OnInit, signal, Signal} from '@angular/core';
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {ConsultationFormComponent} from "../../components/consultation-form/consultation-form.component";
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {Subscription} from "rxjs";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {Router} from "@angular/router";
import {MedicineService} from "../../../../../core/services/health-areas/medicine/medicine.service";
import {DataForm} from "../../../../../core/models/areas/areas.model";
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {NursingConsultationRequest} from "../../../../../core/models/areas/nursing/request.model";
import {ConsultaMedicinaGeneral} from "../../../../../core/models/areas/nursing/nursing.model";
import {v4 as uuidV4} from "uuid";
import {ReloadHealthAreasService} from "../../../../../core/services/health-areas/reload-health-areas.service";

@Component({
  selector: 'medicine-consultation-view',
  standalone: true,
  imports: [
    BlockUiComponent,
    ConsultationFormComponent,
    EcDialogComponent,
    EcTemplate
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
  private _medicineService: MedicineService = inject(MedicineService)
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

  public onEditConsultation(dataForm: DataForm<Omit<ConsultaMedicinaGeneral, "id">>) {
    const request = this.infoConsultation();
    this.isBlockPage.set(true);

    if(request?.examenes) {
      request.examenes = {
        ...request.examenes,
        consulta_medicina_general: {
          ...dataForm.data,
          id: request?.examenes['consulta_medicina_general']?.id ?? uuidV4()
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
    this._router.navigate(['home/medicine-consultation']).then(
      () => this._reloadConsultingService.reloadMedicineConsulting()
    )
  }
}
