import {ChangeDetectionStrategy, Component, inject, input, output, signal} from '@angular/core';
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {ActivatedRoute, Router} from "@angular/router";
import {Subscription} from "rxjs";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {Consulta, ConsultationRecord} from "../../../../../core/models/areas/areas.model";
import {ModalComponent} from "../../../../../core/ui/modal/modal.component";
import {ReactiveFormsModule} from "@angular/forms";
import {TitleCasePipe} from "@angular/common";
import {PermissionsService} from "../../../../../core/services/permissions/permissions.service";
import {ACTION_TYPES} from "../../../../../core/contans/areas/permissions.constant";

@Component({
  selector: 'app-consultation-list',
  standalone: true,
  imports: [
    EcTableModule,
    EcTemplate,
    BlockUiComponent,
    ModalComponent,
    ReactiveFormsModule,
    TitleCasePipe
  ],
  templateUrl: './consultation-list.component.html',
  styleUrl: './consultation-list.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConsultationListComponent {
  patientRecord = input.required<ConsultationRecord>()
  refresh = output()

  public isBlockPage = signal<boolean>(false)
  public showListAreas = signal<boolean>(false)
  public selectedConsultation = signal<Consulta | undefined>(undefined)
  public _subscription = new Subscription()
  public showConfirmDelete = signal<{
    show: boolean,
    consultationID: string | undefined,
  }>({
    show: false,
    consultationID: undefined,
  })

  private _nursingService = inject(NursingService)
  private _toastService = inject(ToastService)
  private _router = inject(Router)
  private _activatedRoute = inject(ActivatedRoute)
  private _permissionsService = inject(PermissionsService)

  public onEditConsultation(consultationID: string) {
    this._router.navigate([
      'view/paciente',
      this.patientRecord().paciente.dni,
      'consulta',
      consultationID
    ], {relativeTo: this._activatedRoute}).then(err => console.log(err))
  }

  public onDeleteConsultation(consultationID: string) {
    this.isBlockPage.set(true);
    this._subscription.add(
      this._nursingService.deleteConsultation(consultationID).subscribe({
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
          this.patientRecord().consultas = this.patientRecord().consultas.filter(consulta => consulta.id !== consultationID)
          this.closeModalDelete()
        },
        error: () => {
          this.isBlockPage.set(false);
          this._toastService.add({
            type: 'error',
            message: 'No se pudo eliminar la consulta',
            life: 5000
          })
        }
      })
    )
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
          this.updateListAssignedArea(area);
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

  onCancelAssignmentArea() {
    this.cancelAssignmentArea()
  }

  cancelAssignmentArea() {
    this.showListAreas.set(false)
    this.selectedConsultation.set(undefined)
  }

  onAsignarArea(consulta: Consulta) {
    const servicios = consulta.servicios.split(', ')
    if (!servicios.includes('Atención Integral')) {
      this._toastService.add({
        type: 'error',
        message: 'No se puede asignar área sin tener una consulta',
        life: 5000
      })
      return
    }

    this.showListAreas.set(true);
    this.selectedConsultation.set(consulta)
  }

  onShowModalConfirmDelete(consultingID: string): void {
    this.showConfirmDelete.set({
      show: true,
      consultationID: consultingID
    })
  }

  onConfirmDelete(): void {
    this.onDeleteConsultation(this.showConfirmDelete().consultationID!)
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

  canEdit(): boolean {
    return this._permissionsService.hasPermission(
      'enfermería',
      'atenciones',
      ACTION_TYPES.UPDATE
    )
  }

  canView(): boolean {
    return this._permissionsService.hasPermission(
      'enfermería',
      'atenciones',
      ACTION_TYPES.READ
    )
  }
}
