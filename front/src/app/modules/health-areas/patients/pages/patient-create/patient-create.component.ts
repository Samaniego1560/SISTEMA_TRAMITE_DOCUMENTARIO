import {Component, inject} from '@angular/core';
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {PatientFormComponent} from "../../components/patient-form/patient-form.component";
import {PatientBase} from "../../../../../core/models/areas/patient.model";
import {PatientsService} from "../../../../../core/services/health-areas/patients.service";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {Subscription} from "rxjs";
import {ActivatedRoute, Router} from "@angular/router";
import { v4 as uuidV4 } from 'uuid';
import {ToastComponent} from "../../../../../core/ui/toast/toast.component";

@Component({
  selector: 'app-create-patient',
  standalone: true,
    imports: [
        EcDialogComponent,
        EcTableModule,
        EcTemplate,
        FormsModule,
        ReactiveFormsModule,
        PatientFormComponent,
        BlockUiComponent,
        ToastComponent,
    ],
  templateUrl: './patient-create.component.html',
  styleUrl: './patient-create.component.scss'
})
export default class PatientCreateComponent {
  public isBlockPage: boolean = false;
  private _subscription = new Subscription()

  private _toastService: ToastService = inject(ToastService)
  private _patientsService: PatientsService = inject(PatientsService)
  private _router: Router = inject(Router)
  private _activatedRoute: ActivatedRoute = inject(ActivatedRoute)

  public onCreatePatient(patient: PatientBase) {
    const request = {
      ...patient,
      id: uuidV4()
    }
    this.isBlockPage = true;
    this._subscription.add(
      this._patientsService.createPatient(request).subscribe({
        next: (resp: any) => {
          this.isBlockPage = false;
          if (resp.error) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo registrar el paciente, intente nuevamente!',
              life: 5000
            })
            return
          }

          this._toastService.add({type: 'success', message: 'Paciente creado.', life: 5000})
          this._redirectToPatients()
        }
      }))
  }

  public onCancelCreate() {
    this._redirectToPatients()
  }

  private _redirectToPatients() {
    this._router.navigate(['../'], {relativeTo: this._activatedRoute}).then()
  }
}
