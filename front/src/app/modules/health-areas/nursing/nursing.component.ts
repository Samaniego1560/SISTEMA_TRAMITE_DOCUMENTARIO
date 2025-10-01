import {Component, inject, Input, OnChanges, OnInit, signal, SimpleChanges} from '@angular/core';
import {EcTableModule} from "../../../core/ui/ec-table/ec-table.module";
import {BlockUiComponent} from "../../../core/ui/block-ui/block-ui.component";
import {ActivatedRoute, Router, RouterOutlet} from "@angular/router";
import {ToastComponent} from "../../../core/ui/toast/toast.component";
import {ToastService} from "../../../core/services/toast/toast.service";
import {ConsultationListComponent} from "./components/consultation-list/consultation-list.component";
import {NursingService} from "../../../core/services/health-areas/nursing/nursing.service";
import {ConsultationRecord} from "../../../core/models/areas/areas.model";
import {FormsModule} from "@angular/forms";
import {ReloadHealthAreasService} from "../../../core/services/health-areas/reload-health-areas.service";
import {HealthAreasUtil} from "../../../core/utils/health-areas/health-areas.util";
import {Response} from "../../../core/models/global.model";

@Component({
  selector: 'app-consulting',
  standalone: true,
  imports: [
    EcTableModule,
    BlockUiComponent,
    RouterOutlet,
    ToastComponent,
    ConsultationListComponent,
    FormsModule
  ],
  templateUrl: './nursing.component.html',
  styleUrl: './nursing.component.scss',
  providers: [ToastService, NursingService]
})
export class NursingComponent implements OnInit {
  isBlockPage = signal(false);
  patientRecord = signal<ConsultationRecord | undefined>(undefined);

  private _router = inject(Router)
  private _activatedRoute = inject(ActivatedRoute)
  private _toastService = inject(ToastService)
  private _nursingService = inject(NursingService)
  private _reloadConsultingService = inject(ReloadHealthAreasService)

  dniPatient = ''

  ngOnInit() {
    this._reloadConsultingService.reloadNursingConsulting$.subscribe({
      next: () => {
        if(this.dniPatient){
          this.searchConsultation()
        }
      }
    })
  }

  onRegisterConsulting() {
    this._router.navigate(['create'], {relativeTo: this._activatedRoute}).then()
  }

  searchConsultation() {
    if(this.dniPatient === ''){
      this._toastService.add({
        type: 'warning',
        message: 'Debe ingresar un DNI!',
        life: 5000
      })
      return
    }

    this.isBlockPage.set(true);
    this._nursingService.getConsultationsByDni(this.dniPatient.trim()).subscribe({
      next: (resp) => this.handleSearchResponse(resp),
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

  handleSearchResponse(resp: Response<ConsultationRecord>): void {
    this.isBlockPage.set(false);
    if (resp.error || !resp.data) {
      this._toastService.add({
        type: 'warning',
        message: 'No se encontraron registros com el DNI ingresado!',
        life: 5000
      })
      this.patientRecord.set(undefined)
      return
    }

    const consultas = HealthAreasUtil.sortConsultations(resp.data.consultas)
    this.patientRecord.set({
      ...resp.data,
      consultas
    })
  }

  onRefresh() {
    this.searchConsultation()
  }
}
