import {Component, inject, OnInit, signal} from '@angular/core';
import {BlockUiComponent} from "../../../core/ui/block-ui/block-ui.component";
import {EcTableModule} from "../../../core/ui/ec-table/ec-table.module";
import {FormsModule} from "@angular/forms";
import {ActivatedRoute, Router, RouterOutlet} from "@angular/router";
import {ToastService} from "../../../core/services/toast/toast.service";
import {ToastComponent} from "../../../core/ui/toast/toast.component";
import {ConsultationListComponent} from "./components/consultation-list/consultation-list.component";
import {NursingService} from "../../../core/services/health-areas/nursing/nursing.service";
import {ConsultationRecord} from "../../../core/models/areas/areas.model";
import {DentalService} from "../../../core/services/health-areas/dental/dental.service";
import {ReloadHealthAreasService} from "../../../core/services/health-areas/reload-health-areas.service";
import {HealthAreasUtil} from "../../../core/utils/health-areas/health-areas.util";

@Component({
  selector: 'app-dental-consultation',
  standalone: true,
  imports: [
    BlockUiComponent,
    EcTableModule,
    FormsModule,
    RouterOutlet,
    ToastComponent,
    ConsultationListComponent
  ],
  templateUrl: './dentistry.component.html',
  styleUrl: './dentistry.component.scss',
  providers: [ToastService, NursingService]
})
export class DentistryComponent implements OnInit {
  isBlockPage = signal(false);
  patientRecord = signal<ConsultationRecord | undefined>(undefined);

  private _router = inject(Router)
  private _activatedRoute = inject(ActivatedRoute)
  private _dentalService = inject(DentalService)
  private _toastService = inject(ToastService)
  private _reloadConsultingService = inject(ReloadHealthAreasService)
  dniPatient = ''

  ngOnInit() {
    this._reloadConsultingService.reloadDentistryConsulting$.subscribe({
      next: () => {
        if (this.dniPatient) {
          this.searchConsultation()
        }
      }
    })
  }

  onRegisterConsulting() {
    this._router.navigate(['create'], {relativeTo: this._activatedRoute}).then()
  }

  searchConsultation() {
    this.isBlockPage.set(true);
    this._dentalService.getConsultationsByDni(this.dniPatient).subscribe({
      next: (resp) => {
        this.isBlockPage.set(false);
        if (resp.error || !resp.data) {
          this._toastService.add({
            type: 'warning',
            message: 'No se encontraron registros com el DNI ingresado!',
            life: 5000
          })
          return
        }

        const consultas = HealthAreasUtil.sortConsultations(resp.data.consultas)
        this.patientRecord.set({
          ...resp.data,
          consultas
        })
      },
      error: () => {
        this.isBlockPage.set(false);
        this._toastService.add({
          type: 'error',
          message: 'Ocurrió un error al buscar el registro, intente nuevamente!',
          life: 5000
        })
      }
    })
  }

  onRefresh() {
    this.searchConsultation()
  }
}
