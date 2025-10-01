import {
  ChangeDetectionStrategy,
  Component,
  inject,
  input,
  OnDestroy,
  OnInit,
  output,
  signal,
  Signal
} from '@angular/core';
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {ActivatedRoute, Router} from "@angular/router";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {DentalService} from "../../../../../core/services/health-areas/dental/dental.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {Consulta, ConsultationRecord} from "../../../../../core/models/areas/areas.model";
import {TitleCasePipe} from "@angular/common";
import {Subscription} from "rxjs";
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {ModalComponent} from "../../../../../core/ui/modal/modal.component";
import {ReactiveFormsModule} from "@angular/forms";

@Component({
  selector: 'app-consultation-list',
  standalone: true,
  imports: [
    EcTableModule,
    EcTemplate,
    BlockUiComponent,
    TitleCasePipe,
    ModalComponent,
    ReactiveFormsModule
  ],
  templateUrl: './consultation-list.component.html',
  styleUrl: './consultation-list.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConsultationListComponent implements OnDestroy {
  patientRecord = input.required<ConsultationRecord>()
  refresh = output()

  public isBlockPage = signal<boolean>(false)
  public showConfirmDelete = signal<{
    show: boolean,
    consultationID: string | undefined,
  }>({
    show: false,
    consultationID: undefined,
  })

  private _router = inject(Router)
  private _activatedRoute = inject(ActivatedRoute)
  private _dentalService = inject(DentalService)
  private _toastService = inject(ToastService)
  private _nursingService = inject(NursingService)
  _subscription = new Subscription()

  public showListAreas = signal<boolean>(false)
  public selectedConsultation = signal<Consulta | undefined>(undefined)

  ngOnDestroy() {
    this._subscription.unsubscribe()
  }

  public onEditConsulting(consulta: Consulta) {
    if(consulta.area_origen === 'odontología'){
      this._router.navigate([
        'view/paciente',
        this.patientRecord().paciente.dni,
        'consulta',
        consulta.id
      ], {relativeTo: this._activatedRoute}).then()
      return
    }

    this._router.navigate([
      'edit/nursing/paciente',
      this.patientRecord().paciente.dni,
      'consulta',
      consulta.id
    ], {relativeTo: this._activatedRoute}).then(err => console.log(err))
  }

  onDeleteConsulting(consultationID: string) {
    this.isBlockPage.set(true);
    this._dentalService.deleteConsulting(consultationID).subscribe({
      next: (resp) => {
        this.isBlockPage.set(false);
        if (resp.error) {
          this._toastService.add({
            type: 'error',
            message: 'No se pudo eliminar la consulta, intente nuevamente!',
            life: 5000
          })
          return
        }

        this._toastService.add({type: 'success', message: 'Consulta eliminada.', life: 5000})
        this.patientRecord().consultas = this.patientRecord().consultas.filter(consulta => consulta.id !== consultationID);
        this.closeModalDelete()
      },
      error: () => {
        this.isBlockPage.set(false);
        this._toastService.add({
          type: 'error',
          message: 'No se pudo eliminar la consulta, intente nuevamente!',
          life: 5000
        })
      }
    })
  }
  onShowAssignModal(item: Consulta): void {
    if(item.area_origen === 'odontología') return;

    this.showListAreas.set(true);
    this.selectedConsultation.set(item)
  }

  public onUpdateAssignedArea(area: string) {
    if(!area) return

    this.isBlockPage.set(true);
    this._subscription.add(
      this._nursingService.updateAssignedArea({
        consulta_id: this.selectedConsultation()?.id!,
        area_asignada: area
      }).subscribe({
        next: (resp) => {
          this.isBlockPage.set(false);
          if (resp.error) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo actualizar asignación, intente nuevamente!',
              life: 5000
            })
            return
          }

          this._toastService.add({type: 'success', message: 'Asignación actualizada', life: 5000})
          this.updateListAssignedArea(area)
          this.cancelAssignmentArea();
        },
        error: () => {
          this.isBlockPage.set(false);
          this._toastService.add({
            type: 'error',
            message: 'Ocurrío un error al actualizar la asignación',
            life: 5000
          })
        }
      })
    )
  }

  updateListAssignedArea(area: string) {
    const records = structuredClone(this.patientRecord().consultas)
    const recordConsultaIdx = records.findIndex(consulta => consulta.id === this.selectedConsultation()!.id)
    if(recordConsultaIdx === -1) return

    records[recordConsultaIdx].area_asignada = area
    this.patientRecord().consultas = records;
  }

  cancelAssignmentArea() {
    this.showListAreas.set(false)
    this.selectedConsultation.set(undefined)
  }

  onCancelAssignmentArea() {
    this.cancelAssignmentArea()
  }

  onShowModalConfirmDelete(consultingID: string): void {
    this.showConfirmDelete.set({
      show: true,
      consultationID: consultingID
    })
  }

  onConfirmDelete(): void {
    this.onDeleteConsulting(this.showConfirmDelete().consultationID!)
  }

  onCancelDelete(): void {
    this.closeModalDelete()
  }

  closeModalDelete(): void {
    this.showConfirmDelete.set({
      show: false,
      consultationID: undefined
    })
  }
}
