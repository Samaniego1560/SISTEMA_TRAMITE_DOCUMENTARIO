import {Component, inject, OnInit, signal} from '@angular/core';
import {EcTableModule} from "../../../core/ui/ec-table/ec-table.module";
import {EcTemplate} from "../../../core/directives/ec-template.directive";
import {ModalComponent} from "../../../core/ui/modal/modal.component";
import {PatientsService} from "../../../core/services/health-areas/patients.service";
import {ActivatedRoute, Router, RouterOutlet} from "@angular/router";
import {Patient} from "../../../core/models/areas/patient.model";
import {ToastService} from "../../../core/services/toast/toast.service";
import {Subscription} from "rxjs";
import {BlockUiComponent} from "../../../core/ui/block-ui/block-ui.component";
import {ToastComponent} from "../../../core/ui/toast/toast.component";
import {ControlErrorsDirective} from "../../../lib/validator-dynamic/directives/control-error.directive";
import {OnlyNumbersDirective} from "../../../core/directives/only-numbers.directive";
import {FormBuilder, FormGroup, ReactiveFormsModule} from "@angular/forms";
import {Metadata} from "../../../core/models/global.model";

@Component({
  selector: 'app-patients',
  standalone: true,
  imports: [
    EcTableModule,
    EcTemplate,
    RouterOutlet,
    BlockUiComponent,
    ToastComponent,
    ModalComponent,
    ControlErrorsDirective,
    OnlyNumbersDirective,
    ReactiveFormsModule
  ],
  templateUrl: './patients.component.html',
  styleUrl: './patients.component.scss',
  providers: [ToastService]
})
export class PatientsComponent implements OnInit {
  public isBlockPage: boolean = false;
  private _subscription = new Subscription();

  public showConfirmDelete = signal<{
    show: boolean,
    patient: Patient | undefined,
  }>({
    show: false,
    patient: undefined,
  })

  patients = signal<Patient[]>([]);
  metadata = signal<Metadata | null>(null);
  filterForm: FormGroup;

  private _toastService: ToastService = inject(ToastService)
  private _patientsService = inject(PatientsService)
  private _router = inject(Router)
  private _activatedRoute = inject(ActivatedRoute)
  private _fb = inject(FormBuilder)

  constructor() {
    this.filterForm = this.buildForm()
  }

  ngOnInit() {
    this.searchPatients()
  }

  buildForm() {
    return this._fb.nonNullable.group({
      dni: [''],
      apellidos: [''],
      nombres: ['']
    })
  }

  searchPatients(pagination?: string) {
    this.isBlockPage = true;
    this._patientsService.getPatients(this.filterForm.value, pagination).subscribe({
      next: (resp) => {
        this.isBlockPage = false;
        if (resp.error) {
          this._toastService.add({
            type: 'error',
            message: 'No se pudo obtener los pacientes, intente nuevamente!',
            life: 5000
          })
          return
        }

        this.patients.set(resp.data ?? [])
        this.metadata.set(resp.metadata)
      }
    })
  }

  clearFilter() {
    this.filterForm.reset()
    this.searchPatients()
  }

  nextSearchPatients() {
    if(this.metadata()!.page === this.metadata()!.total_pages) return
    this.searchPatients(this.metadata()?.links.next)
  }

  backSearchPatients() {
    if(this.metadata()?.offset === 0) return
    this.searchPatients(this.metadata()?.links.prev)
  }

  public registerPatient() {
    this._router.navigate(['create'], {relativeTo: this._activatedRoute}).then()
  }

  public onEditPatient(patient: Patient) {
    this._router.navigate([patient.id, 'edit'], {relativeTo: this._activatedRoute}).then()
  }

  public deletePatient(patient: Patient) {
    this.isBlockPage = true;
    this._subscription.add(
      this._patientsService.removePatient(patient.id).subscribe({
        next: (resp: any) => {
          this.isBlockPage = false;
          if (resp.error) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo eliminar la información del paciente, intente nuevamente!',
              life: 5000
            })
            return
          }

          this._toastService.add({type: 'success', message: 'Paciente eliminado.', life: 5000})
          this.searchPatients()
          this.closeModalDelete()
        }
      })
    )
  }

  onShowModalConfirmDelete(patient: Patient): void {
    this.showConfirmDelete.set({
      show: true,
      patient: patient
    })
  }

  onConfirmDelete(): void {
    this.deletePatient(this.showConfirmDelete().patient!)
  }

  onCancelDelete(): void {
    this.closeModalDelete()
  }

  closeModalDelete(): void {
    this.showConfirmDelete.set({
      show: false,
      patient: undefined
    })
  }
}
