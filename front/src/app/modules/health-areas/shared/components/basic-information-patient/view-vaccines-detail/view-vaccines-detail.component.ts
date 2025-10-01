import {Component, computed, inject, input, OnInit, output, signal} from '@angular/core';
import {EcDialogComponent} from "../../../../../../core/ui/ec-dialog/ec-dialog.component";
import {VaccinesListComponent} from "../../vaccines-list/vaccines-list.component";
import {BlockUiComponent} from "../../../../../../core/ui/block-ui/block-ui.component";
import {ToastService} from "../../../../../../core/services/toast/toast.service";
import {VaccineListCard, Vacuna} from "../../../../../../core/models/areas/nursing/nursing.model";
import {NursingService} from "../../../../../../core/services/health-areas/nursing/nursing.service";

@Component({
  selector: 'app-view-vaccines-detail',
  standalone: true,
  imports: [
    EcDialogComponent,
    VaccinesListComponent,
    BlockUiComponent
  ],
  templateUrl: './view-vaccines-detail.component.html',
  styleUrl: './view-vaccines-detail.component.scss',
  providers: [ToastService, NursingService]
})
export class ViewVaccinesDetailComponent implements OnInit {
  dni = input.required<string>()

  vaccines = signal<Vacuna[]>([])
  isBlockPage = signal<boolean>(false)
  isShowModal = signal<boolean>(false)

  hasVaccines = computed<boolean>(() => this.vaccines().length > 0);

  vaccinesListCard = computed<VaccineListCard[]>(() =>{
    const names = new Set(this.vaccines().map(record => record.tipo_vacuna));
    const nameVaccines = Array.from(names);
    nameVaccines.sort((a, b) => a.localeCompare(b))
    return this.sortVaccinesCard(nameVaccines);
  });

  private _toastService = inject(ToastService)
  private _nursingService = inject(NursingService)

  ngOnInit() {
    this.loadVaccines()
  }

  loadVaccines() {
    this.isBlockPage.set(true)
    this._nursingService.getVaccinesPatientByDni(this.dni()).subscribe({
      next: (resp) => {
        this.isBlockPage.set(false)
        if (resp.error) {
          this._toastService.add({
            type: 'error',
            message: 'No se encontraron vacunas'
          })
          return
        }

        this.vaccines.set(resp.data ?? [])
      }
    })
  }

  sortVaccinesCard(nameVaccines: string[]): VaccineListCard[] {
    return nameVaccines.map(name => {
      const vaccines = this.vaccines().filter(record => record.tipo_vacuna === name);
      vaccines.sort((a, b) => {
        const aDate = new Date(a.fecha_dosis).getTime();
        const bDate = new Date(b.fecha_dosis).getTime();
        return aDate > bDate ? 1 : -1;
      });
      return {
        name,
        vaccines: vaccines
      }
    });
  }
}
