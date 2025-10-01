import {ChangeDetectionStrategy, Component, inject, input, OnInit, output, signal} from '@angular/core';
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {Subscription} from "rxjs";
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {ActivatedRoute, Router} from "@angular/router";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {Consulta, ConsultationRecord} from "../../../../../core/models/areas/areas.model";
import {TitleCasePipe} from "@angular/common";
import {ModalComponent} from "../../../../../core/ui/modal/modal.component";
import {ReactiveFormsModule} from "@angular/forms";

@Component({
  selector: 'medicine-consultation-list',
  standalone: true,
  imports: [
    BlockUiComponent,
    EcTableModule,
    EcTemplate,
    TitleCasePipe,
    ModalComponent,
    ReactiveFormsModule
  ],
  templateUrl: './consultation-list.component.html',
  styleUrl: './consultation-list.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ConsultationListComponent implements OnInit {
  patientRecord = input.required<ConsultationRecord>()
  refresh = output()

  public isBlockPage = signal<boolean>(false)
  public _subscription = new Subscription()

  private _nursingService = inject(NursingService)
  private _toastService = inject(ToastService)
  private _router = inject(Router)
  private _activatedRoute = inject(ActivatedRoute)


  public showListAreas = signal<boolean>(false)
  public selectedConsultation = signal<Consulta | undefined>(undefined)

  ngOnInit() {
  }

  public onEditConsultation(consultationID: string) {
    this._router.navigate([
      'edit/paciente',
      this.patientRecord().paciente.dni,
      'consulta',
      consultationID
    ], {relativeTo: this._activatedRoute}).then(err => console.log(err))
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

  onCancelAssignmentArea() {
    this.cancelAssignmentArea()
  }

  cancelAssignmentArea() {
    this.showListAreas.set(false)
    this.selectedConsultation.set(undefined)
  }
}
