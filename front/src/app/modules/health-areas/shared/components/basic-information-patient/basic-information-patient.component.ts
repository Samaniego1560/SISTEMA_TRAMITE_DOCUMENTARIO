import {
  ChangeDetectionStrategy,
  Component, computed,
  inject,
  Input,
  OnChanges,
  output,
  signal,
  SimpleChanges,
} from '@angular/core';
import {PatientsService} from "../../../../../core/services/health-areas/patients.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {Patient} from "../../../../../core/models/areas/patient.model";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {ToggleSwitchComponent} from "../../../../../core/ui/toggle/toggle-switch.component";
import {RequiredVaccines} from "../../../../../core/models/areas/nursing/required-vaccines.model";
import {ViewVaccinesDetailComponent} from "./view-vaccines-detail/view-vaccines-detail.component";
import {ViewPatientDetailComponent} from "./view-patient-detail/view-patient-detail.component";

@Component({
  selector: 'app-basic-information-patient',
  standalone: true,
  imports: [
    BlockUiComponent,
    FormsModule,
    EcTableModule,
    ReactiveFormsModule,
    ToggleSwitchComponent,
    ViewVaccinesDetailComponent,
    ViewPatientDetailComponent
  ],
  templateUrl: './basic-information-patient.component.html',
  styleUrl: './basic-information-patient.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class BasicInformationPatientComponent implements OnChanges {
  @Input() infoPatientDni: string | undefined;
  @Input() showVaccinesRequired: boolean = true;
  public patient = signal<Patient | undefined>(undefined);
  patientFound = output<Patient>()

  public searchPatientID = '';

  public active = signal<boolean>(true);
  public isBlockPage = signal<boolean>(false);
  public showInfoPatient = signal<boolean>(false);
  public vaccines = signal<RequiredVaccines[] | null>(null);
  private _patientsService: PatientsService = inject(PatientsService)
  private _nursingService: NursingService = inject(NursingService)
  private _toastService: ToastService = inject(ToastService)

  hasNextVaccines = computed((): boolean => !!this.vaccines()?.length)

  hasRequiredVaccines = computed((): boolean => {
    console.log(this.vaccines())
    if (this.vaccines() === null) return false
    return this.vaccines()?.some(vaccine => vaccine.requerido) ?? false
  })

  requiredVaccines = computed((): RequiredVaccines[] => this.vaccines()?.filter(vaccine => vaccine.requerido) ?? [])

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoPatientDni'].currentValue) {
      this.searchPatientID = this.infoPatientDni!;
      this.onSearchPatient(this.searchPatientID)
    }
  }

  onSearchPatient(id: string) {
    if (!id) return

    this.isBlockPage.set(true)
    this._patientsService.getPatientByDni(id).subscribe({
      next: (res: any) => {
        if (res.error || !res.data) {
          this.isBlockPage.set(false)
          this._toastService.add({
            type: 'error',
            message: 'Paciente no encontrado'
          })
        }

        this.patient.set(res.data)
        this.patientFound.emit(res.data);
        if (this.showVaccinesRequired) {
          this.searchPatientVaccines();
          return
        }

        this.isBlockPage.set(false)
      },
      error: () => {
        this.isBlockPage.set(false)
        this._toastService.add({
          type: 'error',
          message: 'Ocurrió un error al buscar el paciente'
        })
      }
    })
  }

  searchPatientVaccines() {
    this._nursingService.getRequiredVaccines(this.patient()!.id).subscribe({
      next: (res) => {
        this.isBlockPage.set(false)
        if (res.error) {
          this._toastService.add({
            type: 'error',
            message: 'Ocurrió un error al buscar las vacunas requeridas'
          })
          return
        }
        this.vaccines.set(res.data.required)
      },
      error: () => {
        this.isBlockPage.set(false)
        this._toastService.add({
          type: 'error',
          message: 'Ocurrió un error al buscar el paciente'
        })
      }
    })
  }

  onShowDetailPatient() {
    this.showInfoPatient.set(true)
  }
}
