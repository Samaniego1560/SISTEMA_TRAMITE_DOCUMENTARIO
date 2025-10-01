import {Component, inject, OnInit, signal} from '@angular/core';
import {BlockUiComponent} from "../../../core/ui/block-ui/block-ui.component";
import {EcTableModule} from "../../../core/ui/ec-table/ec-table.module";
import {ActivatedRoute, Router, RouterOutlet} from "@angular/router";
import {FormsModule} from "@angular/forms";
import {ToastService} from "../../../core/services/toast/toast.service";
import {ConsultationListComponent} from "./components/consultation-list/consultation-list.component";
import {MedicineService} from "../../../core/services/health-areas/medicine/medicine.service";
import {NursingService} from "../../../core/services/health-areas/nursing/nursing.service";
import {ConsultationRecord} from "../../../core/models/areas/areas.model";
import {ReloadHealthAreasService} from "../../../core/services/health-areas/reload-health-areas.service";
import {ToastComponent} from "../../../core/ui/toast/toast.component";
import {HealthAreasUtil} from "../../../core/utils/health-areas/health-areas.util";

@Component({
  selector: 'app-medicine-consultation',
  standalone: true,
  imports: [
    BlockUiComponent,
    EcTableModule,
    RouterOutlet,
    FormsModule,
    ConsultationListComponent,
    ToastComponent
  ],
  templateUrl: './medicine.component.html',
  styleUrl: './medicine.component.scss',
  providers: [ToastService, MedicineService, NursingService]
})
export class MedicineComponent implements OnInit {
  isBlockPage = signal(false);
  patientRecord = signal<ConsultationRecord | undefined>(undefined);

  private _medicineService = inject(MedicineService)
  private _toastService = inject(ToastService)
  private _reloadConsultingService = inject(ReloadHealthAreasService)

  dniPatient = ''

  ngOnInit() {
    this._reloadConsultingService.reloadMedicineConsulting$.subscribe({
      next: () => {
        if(this.dniPatient){
          this.searchConsultation()
        }
      }
    })
  }

  searchConsultation() {
    this.isBlockPage.set(true);
    this._medicineService.getConsultationsByDni(this.dniPatient).subscribe({
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
          message: 'No se pudo obtener las consultas, intente nuevamente!',
          life: 5000
        })
      }
    })
  }

  onRefresh() {
    this.searchConsultation()
  }
}
