import {Component, inject, Input, OnInit, signal} from '@angular/core';
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {PatientsService} from "../../../../../core/services/health-areas/patients.service";
import {Patient, PatientBase, PatientRequest} from "../../../../../core/models/areas/patient.model";
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {PatientFormComponent} from "../../components/patient-form/patient-form.component";
import {ActivatedRoute, Router} from "@angular/router";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {Subscription} from "rxjs";

@Component({
  selector: 'app-edit-patient',
  standalone: true,
  imports: [
    EcDialogComponent,
    EcTemplate,
    PatientFormComponent,
    BlockUiComponent
  ],
  templateUrl: './patient-edit.component.html',
  styleUrl: './patient-edit.component.scss'
})
export default class PatientEditComponent implements OnInit {
  @Input() id: string = '';
  public isBlockPage: boolean = false;
  private _subscription = new Subscription();

  private _toastService: ToastService = inject(ToastService)
  private _patientsService: PatientsService = inject(PatientsService)
  private _router: Router = inject(Router)
  public selectedPatient = signal<Patient | undefined>(undefined);

  constructor() {
  }

  ngOnInit() {
    this._processPatient()
  }

  private _processPatient() {
    if (!this.id) {
      this._redirectToPatients()
      return
    }

    this.getPatient()
  }

  public getPatient() {
    this.isBlockPage = true;
    this._subscription.add(
      this._patientsService.getPatientByID(this.id).subscribe({
        next: (resp) => {
          this.isBlockPage = false;
          if (resp.error) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo obtener la información del paciente, intente nuevamente!',
              life: 5000
            })
            this._redirectToPatients()
            return
          }

          this.selectedPatient.set(resp.data)
        }
      })
    )
  }

  public onEditPatient(patient: PatientBase) {
    const patientRequest: PatientRequest = {
      ...patient,
      id: this.selectedPatient()!.id
    }

    this.isBlockPage = true;
    this._subscription.add(
      this._patientsService.updatePatient(patientRequest).subscribe({
        next: (resp: any) => {
          this.isBlockPage = false;
          if (resp.error) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo actualizar la información del paciente, intente nuevamente!',
              life: 5000
            })
            return
          }

          this._toastService.add({type: 'success', message: 'Paciente actualizado.', life: 5000})
          this._redirectToPatients()
        }
      })
    )
  }

  public onCancelEdit() {
    this._redirectToPatients()
  }

  private _redirectToPatients() {
    this._router.navigate(['/home/patients']).then()
  }
}
