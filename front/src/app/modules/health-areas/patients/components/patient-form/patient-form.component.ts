import {Component, inject, Input, OnChanges, OnInit, output, signal, SimpleChanges} from '@angular/core';
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormSubmitDirective} from "../../../../../lib/validator-dynamic/directives/form-submit.directive";
import {FormBuilder, FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {NgxValidators} from "../../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {ToastService} from "../../../../../core/services/toast/toast.service";
import {Antecedent, Patient, PatientBase} from "../../../../../core/models/areas/patient.model";
import {BLOOD_GROUPS, MARITAL_STATUS, TYPE_SEX} from "../../../../../core/contans/areas/patient.model";
import {PatientsService} from "../../../../../core/services/health-areas/patients.service";
import {ToggleSwitchComponent} from "../../../../../core/ui/toggle/toggle-switch.component";
import {EcTableModule} from "../../../../../core/ui/ec-table/ec-table.module";
import {EcTemplate} from "../../../../../core/directives/ec-template.directive";
import {BlockUiComponent} from "../../../../../core/ui/block-ui/block-ui.component";
import {v4 as uuidV4} from "uuid";
import {OnlyNumbersDirective} from "../../../../../core/directives/only-numbers.directive";
import {EcDialogComponent} from "../../../../../core/ui/ec-dialog/ec-dialog.component";

@Component({
  selector: 'app-patient-form',
  standalone: true,
  imports: [
    ControlErrorsDirective,
    FormSubmitDirective,
    ReactiveFormsModule,
    ToggleSwitchComponent,
    EcTableModule,
    EcTemplate,
    BlockUiComponent,
    OnlyNumbersDirective,
    FormsModule,
    EcDialogComponent
  ],
  templateUrl: './patient-form.component.html',
  styleUrl: './patient-form.component.scss'
})
export class PatientFormComponent implements OnChanges, OnInit {
  @Input() patient: Patient | undefined;
  save = output<PatientBase>();
  cancel = output<void>();

  private _fb = inject(FormBuilder);
  private _toastService: ToastService = inject(ToastService)
  private _patientsService: PatientsService = inject(PatientsService)

  public patientForm: FormGroup;
  public antecedentsForm: FormGroup;
  public BLOOD_GROUPS = BLOOD_GROUPS;
  public TYPE_SEX = TYPE_SEX;
  public MARITAL_STATUS = MARITAL_STATUS;
  antecedents = signal<Antecedent[]>([])
  isBlockPage = signal<boolean>(false)
  typeAntecedents = [
    'Tuberculosis',
    'IRA',
    'Inf. Transmision Sexual',
    'VHI - SIDA',
    'ITU',
    'Hepatitis',
    'Diabetes',
    'HTA',
    'Sobrepeso',
    'Dislipidemia',
    'Convulsiones',
    'Alergias',
    'Hospitalización',
    'Interv. Quirúrgica',
    'Transfusiones',
    'Accidentes',
    'Cáncer cérvix',
    'Patología prostática',
    'Dispacidad',
    'Riesgo ocupacional',
    'Gastritis',
    'Otro',
  ]
  showModalAntecedents = signal<boolean>(false)

  constructor() {
    this.patientForm = this._buildPatientForm();
    this.antecedentsForm = this._buildAntecedentsForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['patient'].currentValue) {
      this.patientForm.patchValue(changes['patient'].currentValue)
      this.validateRam(changes['patient'].currentValue.ram)
      this.antecedents.set(changes['patient'].currentValue.antecedentes || [])
      this.validateTypePerson()
    }
  }

  ngOnInit() {
    this.calculateAge()
  }

  onSearchPatient(event: Event, dni: string) {
    event.preventDefault();
    event.stopImmediatePropagation();
    this.searchPerson(dni)
  }

  private _buildPatientForm(): FormGroup {
    return this._fb.nonNullable.group({
      tipo_persona: ['', [NgxValidators.required()]],
      codigo_sga: ['', [NgxValidators.required()]], //check
      dni: ['', [NgxValidators.required()]], //check
      apellidos: ['', [NgxValidators.required()]], //check
      nombres: ['', [NgxValidators.required()]], //check
      sexo: ['', [NgxValidators.required()]], //check
      edad: new FormControl({value: '', disabled: true}), //check
      estado_civil: ['', [NgxValidators.required()]], //check
      grupo_sanguineo: ['', []], //check
      fecha_nacimiento: ['', [NgxValidators.required()]], //check
      lugar_nacimiento: ['', [NgxValidators.required()]], //check
      procedencia: ['', [NgxValidators.required()]], //check
      factor_rh: ['', []],
      escuela_profesional: ['', [NgxValidators.required()]], //check
      ocupacion: ['', []], //check
      correo_electronico: ['', [Validators.required]], //check
      numero_celular: ['', [Validators.required]], //check
      direccion: ['', [NgxValidators.required()]], //check
      ram: [false, [NgxValidators.required()]], //check
      alergias: [{value: '', disabled: true}, [NgxValidators.required()]], //check
    })
  }

  private _buildAntecedentsForm(): FormGroup {
    return this._fb.nonNullable.group({
      id: [''],
      tipo_antecedente: ['', [NgxValidators.required()]],
      nombre_antecedente: ['', [NgxValidators.required()]],
      estado_antecedente: ['', [NgxValidators.required()]], //check
    })
  }

  public searchPerson(dni: string) {
    this.isBlockPage.set(true)
    this._patientsService.getPatient(dni).subscribe({
      next: (resp) => {
        this.isBlockPage.set(false)
        if(!resp.detalle || Array.isArray(resp.detalle)) {
          this._toastService.add({
            type: 'error',
            message: 'Paciente no encontrado'
          })
          return
        }

        this.patientForm.get('apellidos')?.setValue(resp.detalle.appaterno + ' ' + resp.detalle.apmaterno)
        this.patientForm.get('nombres')?.setValue(resp.detalle.nombre)
        this.patientForm.get('edad')?.setValue(this.getAge(resp.detalle.fecnac))
        this.patientForm.get('codigo_sga')?.setValue(resp.detalle.codalumno?.toString())
        this.patientForm.get('sexo')?.setValue(resp.detalle.sexo)
        this.patientForm.get('fecha_nacimiento')?.setValue(resp.detalle.fecnac)
        this.patientForm.get('direccion')?.setValue(resp.detalle.direccion)
        this.patientForm.get('escuela_profesional')?.setValue(resp.detalle.nomesp)
        this.patientForm.get('correo_electronico')?.setValue(resp.detalle.emailinst || resp.detalle.email)
        this.patientForm.get('numero_celular')?.setValue(resp.detalle.telcelular?.toString())
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

  calculateAge(): void {
    this.patientForm.controls['fecha_nacimiento'].valueChanges.subscribe(birthday => {
      this.patientForm.controls['edad'].setValue(this.getAge(birthday), {emitEvent: false})
    })
  }

  private getAge(birthday: string) {
    const birthdayDate = new Date(birthday)
    const diff = Date.now() - birthdayDate.getTime()
    const ageDate = new Date(diff)
    return Math.abs(ageDate.getUTCFullYear() - 1970).toString()
  }

  public savePatient() {
    if (this.patientForm.invalid) {
      this.patientForm.markAllAsTouched();
      this._toastService.add({
        type: 'error',
        message: 'Formulario inválido'
      });
      return
    }

    this.save.emit({
      ...this.patientForm.getRawValue(),
      antecedentes: this.antecedents()
    });
  }

  onChangeRam(ram: boolean) {
    this.validateRam(ram)
    this.patientForm.get('alergias')?.reset()
  }

  validateRam(ram: boolean) {
    if (ram) {
      this.patientForm.get('alergias')?.enable()
      return
    }
    this.patientForm.get('alergias')?.disable()
  }

  saveFormRecord() {
    if (this.antecedentsForm.invalid) {
      this.antecedentsForm.markAllAsTouched()
      return
    }

    if (this.antecedentsForm.controls['id'].value) {
      this.updateRecord();
      this.cancelModalAntecedent()
    } else {
      this.addRecord()
    }

    this.antecedentsForm.reset();
  }

  addRecord() {
    this.antecedents.update(records => {
      return [
        ...records, {
          nombre_antecedente: this.antecedentsForm.value.nombre_antecedente,
          estado_antecedente: this.antecedentsForm.value.estado_antecedente,
          id: uuidV4()
        }]
    })
  }

  updateRecord() {
    this.antecedents.update(antecedents => {
      const procedureIndex = antecedents.findIndex((procedure) => procedure.id === this.antecedentsForm.value.id)
      antecedents[procedureIndex] = {
        ...antecedents[procedureIndex],
        nombre_antecedente: this.antecedentsForm.value.nombre_antecedente,
        estado_antecedente: this.antecedentsForm.value.estado_antecedente,
      }
      return antecedents
    })
  }

  editAntecedent(record: Antecedent) {
    this.antecedentsForm.patchValue({
      id: record.id,
      tipo_antecedente: record.nombre_antecedente,
      nombre_antecedente: record.nombre_antecedente,
      estado_antecedente: record.estado_antecedente,
    });

    const existAntecedente = this.typeAntecedents.includes(record.nombre_antecedente)
    if (!existAntecedente) {
      this.antecedentsForm.patchValue({
        tipo_antecedente: 'Otro',
      });
    }

    this.showModalAntecedent();
  }

  onDeleteAntecedent(antecedentID: string) {
    this.antecedents.update(records => {
      return records.filter((record) => record.id !== antecedentID)
    })
  }

  validateTypePerson() {
    if (this.patientForm.get('tipo_persona')?.value === 'Estudiante') {
      this.patientForm.get('codigo_sga')?.enable()
      this.patientForm.get('escuela_profesional')?.enable()
      this.patientForm.get('codigo_sga')?.setValidators([Validators.required])
      this.patientForm.get('escuela_profesional')?.setValidators([Validators.required])
      this.patientForm.get('codigo_sga')?.updateValueAndValidity()
      this.patientForm.get('escuela_profesional')?.updateValueAndValidity()
      return
    }
    this.patientForm.get('codigo_sga')?.disable()
    this.patientForm.get('escuela_profesional')?.disable()
    this.patientForm.get('codigo_sga')?.clearValidators()
    this.patientForm.get('escuela_profesional')?.clearValidators()
    this.patientForm.get('codigo_sga')?.updateValueAndValidity()
    this.patientForm.get('escuela_profesional')?.updateValueAndValidity()
  }

  selectAntecedente() {
    if (this.antecedentsForm.get('tipo_antecedente')?.value === 'Otro') {
      this.antecedentsForm.get('nombre_antecedente')?.setValue('')
      return
    }

    this.antecedentsForm.get('nombre_antecedente')?.setValue(this.antecedentsForm.get('tipo_antecedente')?.value)
  }

  cancelModalAntecedent() {
    this.showModalAntecedents.set(false)
    this.antecedentsForm.reset()
  }

  showModalAntecedent() {
    this.showModalAntecedents.set(true)
  }
}
