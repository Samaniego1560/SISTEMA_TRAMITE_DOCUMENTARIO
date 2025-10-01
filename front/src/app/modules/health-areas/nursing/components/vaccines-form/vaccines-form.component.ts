import {
  ChangeDetectionStrategy,
  Component, computed, effect,
  inject, input,
  Input,
  OnChanges, OnInit,
  output,
  signal, SimpleChanges,
  WritableSignal
} from '@angular/core';
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule} from "@angular/forms";
import {NgxValidators} from "../../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {
  Vacuna, VaccineListCard
} from "../../../../../core/models/areas/nursing/nursing.model";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {v4 as uuidV4} from "uuid";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormSubmitDirective} from "../../../../../lib/validator-dynamic/directives/form-submit.directive";
import {NursingService} from "../../../../../core/services/health-areas/nursing/nursing.service";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {VaccinesType} from "../../../../../core/models/areas/nursing/vaccines.model";
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {DatePipe, formatDate} from "@angular/common";
import {
  VaccinesListComponent
} from "../../../shared/components/vaccines-list/vaccines-list.component";
import {ToggleSwitchComponent} from "../../../../../core/ui/toggle/toggle-switch.component";

@Component({
  selector: 'app-vaccines-form',
  standalone: true,
  imports: [
    FormsModule,
    EcTableModule,
    ControlErrorsDirective,
    FormSubmitDirective,
    ReactiveFormsModule,
    EcDialogComponent,
    BlockUiComponent,
    VaccinesListComponent,
    ToggleSwitchComponent
  ],
  templateUrl: './vaccines-form.component.html',
  styleUrl: './vaccines-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class VaccinesFormComponent extends FormBaseComponent<Vacuna[]> implements OnInit, OnChanges {
  infoPatientDni = input.required<string>();
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }
  deleteForm = output();
  public typesVaccines = signal<VaccinesType[]>([])
  public typesVaccinesAvailable = signal<VaccinesType[]>([])
  vaccines = signal<Vacuna[]>([])

  private _fb = inject(FormBuilder);
  private _nursingService = inject(NursingService);
  private _toastService = inject(ToastService);

  public form: FormGroup;
  newVaccinesPatient: WritableSignal<Vacuna[]> = signal<Vacuna[]>([]);
  isEdit = signal<boolean>(false);
  private _name = signal('vacunas');
  public active = signal<boolean>(true);

  showModalVaccine = signal<boolean>(false);
  isBlockPage = signal<boolean>(false);
  vaccinesListCard = computed<VaccineListCard[]>(() =>{
    const names = new Set(this.newVaccinesPatient().map(record => record.tipo_vacuna));
    const nameVaccines = Array.from(names);
    nameVaccines.sort((a, b) => a.localeCompare(b))
    return this.sortVaccinesCard(nameVaccines);
  });

  constructor() {
    super();
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if(changes['infoData'] && this.infoData) {
      this.newVaccinesPatient.set(this.infoData);
      this.active.set(false);
    }
  }

  ngOnInit() {
    this.loadTypeVaccines();
  }

  sortVaccinesCard(nameVaccines: string[]): VaccineListCard[] {
    return nameVaccines.map(name => {
      const vaccines = this.newVaccinesPatient().filter(record => record.tipo_vacuna === name);
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

  loadTypeVaccines() {
    this.isBlockPage.set(true);
    this._nursingService.getVaccines().subscribe({
      next: (resp) => {
        this.isBlockPage.set(false);
        if (resp.error) {
          this._toastService.add({
            type: 'error',
            message: 'No se pudo obtener las vacunas, intente nuevamente!',
            life: 5000
          })
          return
        }

        this.typesVaccines.set(resp.data);
        this.loadVaccines()
      },
      error: () => {
        this.isBlockPage.set(false);
        this._toastService.add({
          type: 'error',
          message: 'No se pudo obtener las vacunas, intente nuevamente!',
          life: 5000
        })
      }
    })
  }

  loadVaccines() {
    this.isBlockPage.set(true)
    this._nursingService.getVaccinesPatientByDni(this.infoPatientDni()).subscribe({
      next: (resp) => {
        this.isBlockPage.set(false)
        if (resp.error) {
          this._toastService.add({
            type: 'error',
            message: 'Ocurrió un error al cargar las vacunas, intente nuevamente!',
          })
          return
        }

        this.vaccines.set(resp.data ?? [])
        this.validateAvailableVaccines()
      }
    })
  }

  validateAvailableVaccines() {
    this.typesVaccinesAvailable.set([])
    for (const typeVaccine of this.typesVaccines()) {
      if(typeVaccine.dosis_maximas === null) {
        this.typesVaccinesAvailable.update(types => [...types, typeVaccine])
        continue
      }

      const vaccines = this.vaccines().filter(vaccine => vaccine.tipo_vacuna === typeVaccine.nombre)
      const newVaccinesPatient = this.newVaccinesPatient().filter(record => record.tipo_vacuna === typeVaccine.nombre)
      const totalVaccines = newVaccinesPatient.length + vaccines.length

      if(totalVaccines === 0) {
        this.typesVaccinesAvailable.update(types => [...types, typeVaccine])
        continue
      }

      if(totalVaccines === typeVaccine.dosis_maximas) continue

      this.typesVaccinesAvailable.update(types => [...types, typeVaccine])
    }
  }

  private _buildForm() {
    return this._fb.nonNullable.group({
      id: [''],
      tipo_vacuna: ['', [NgxValidators.required()]],
      fecha_dosis: ['', [NgxValidators.required()]],
      minsa: [false],
      tipo_atencion: ['', []],
      indicaciones: ['', []],
      observaciones: ['', []],
    })
  }

  public get name() {
    return this._name();
  }

  isValid(): boolean {
    return !!this.newVaccinesPatient().length;
  }

  getData() {
    return this.newVaccinesPatient();
  }

  saveFormRecord() {
    if (this.form.invalid) {
      this.form.markAllAsTouched()
      return
    }

    if (this.isEdit()) {
      this.updateRecord();
      this.validateAvailableVaccines();
      this.isEdit.set(false)
    } else {
      this.addRecord();
      this.validateAvailableVaccines();
    }

    this.form.reset();
  }

  addRecord() {
    this.newVaccinesPatient.update(records => {
      return [
        ...records, {
          ...this.form.value,
          fecha_dosis: formatDate(this.form.value.fecha_dosis, 'yyyy-MM-dd', 'en-US'),
          id: uuidV4()
        }]
    })
  }

  updateRecord() {
    this.newVaccinesPatient.update(procedures => {
      const procedureIndex = procedures.findIndex((procedure) => procedure.id === this.form.value.id)
      procedures[procedureIndex] = {
        ...this.form.value,
        fecha_dosis: formatDate(this.form.value.fecha_dosis, 'yyyy-MM-dd', 'en-US'),
      }
      return structuredClone(procedures)
    })
  }

  onEditRecord(record: Vacuna) {
    this.form.patchValue({
      ...record,
      fecha_dosis: formatDate(record.fecha_dosis, 'yyyy-MM-dd', 'en-US')
    });
    this.isEdit.set(true);
  }

  onDeleteRecord(recordID: string) {
    this.newVaccinesPatient.update(records => {
      return records.filter((record) => record.id !== recordID)
    })
  }

  clearForm(): void {
    this.form.reset();
    this.isEdit.set(false);
  }
}
