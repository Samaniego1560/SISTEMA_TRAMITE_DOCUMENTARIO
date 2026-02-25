import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators, FormControl, FormGroup } from '@angular/forms';
import { StudentParentSurveyService } from './student-parent-survey.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-student-parent-survey',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './student-parent-survey.component.html',
  styleUrl: './student-parent-survey.component.scss'
})
export class StudentParentSurveyComponent {
  verified = false;
  alreadyAnswered = false;
  submitting = false;
  successMessage = '';
  formMessage = '';
  errorMessage = '';
  studentInfo: any = null;
  dniHint = 'Ingrese su DNI (8 dígitos numéricos).';
  private lastVerifiedDni: string | null = null;
  private verifyingDni = false;

  showSectionIII = false;
  showSectionIV = false;
  private lastChildrenCount: number | null = null;
  private prevShowSectionIII = false;
  private prevShowSectionIV = false;

  form: FormGroup<{
    dni: FormControl<string>;
    student_code: FormControl<string>;
    school: FormControl<string | null>;
    academic_cycle: FormControl<string>;
    age: FormControl<number | null>;
    gender: FormControl<'F' | 'M' | null>;
    is_pregnant: FormControl<boolean | null>;
    has_children: FormControl<boolean | null>;
    children_count: FormControl<number | null>;
    gestation_selected: FormControl<boolean>;
    child_1_age: FormControl<number | null>;
    child_2_age: FormControl<number | null>;
    child_3_age: FormControl<number | null>;
    child_4_age: FormControl<number | null>;
    residence_type: FormControl<'familiar' | 'alquiler' | 'residencia' | 'otro' | null>;
    residence_other: FormControl<string | null>;
    lives_with_partner: FormControl<boolean | null>;
    lives_with_children: FormControl<boolean | null>;
    caregiver_type: FormControl<'ambos' | 'familia' | 'padre' | 'madre' | 'otro' | null>;
    caregiver_other: FormControl<string | null>;
    participates_in_upbringing: FormControl<boolean | null>;
    economic_support: FormControl<boolean | null>;
    monthly_amount: FormControl<number | null>;
    has_academic_support: FormControl<boolean | null>;
    support_provider: FormControl<'ambos' | 'familia' | 'padre' | 'madre' | 'otro' | null>;
    support_provider_other: FormControl<string | null>;
    interrupted_semester: FormControl<boolean | null>;
  }>;

  constructor(
    private fb: FormBuilder,
    private survey: StudentParentSurveyService,
    private router: Router
  ) {
    this.form = this.fb.group({
      dni: this.fb.control('', { validators: [Validators.required, Validators.maxLength(12)], nonNullable: true }),
      student_code: this.fb.control({ value: '', disabled: true }, { validators: [Validators.required], nonNullable: true }),
      school: this.fb.control<string | null>(null, { validators: [Validators.required] }),
      academic_cycle: this.fb.control({ value: '', disabled: true }, { validators: [Validators.required], nonNullable: true }),
      age: this.fb.control<number | null>(null, { validators: [Validators.required, Validators.min(15), Validators.max(90)] }),
      gender: this.fb.control<'F' | 'M' | null>(null, { validators: [Validators.required] }),

      is_pregnant: this.fb.control<boolean | null>(null, { validators: [Validators.required] }),
      has_children: this.fb.control<boolean | null>(null, { validators: [Validators.required] }),

      children_count: this.fb.control<number | null>(null),
      gestation_selected: this.fb.control(false, { nonNullable: true }),
      child_1_age: this.fb.control<number | null>(null),
      child_2_age: this.fb.control<number | null>(null),
      child_3_age: this.fb.control<number | null>(null),
      child_4_age: this.fb.control<number | null>(null),

      residence_type: this.fb.control<'familiar' | 'alquiler' | 'residencia' | 'otro' | null>(null),
      residence_other: this.fb.control<string | null>(null),
      lives_with_partner: this.fb.control<boolean | null>(null),
      lives_with_children: this.fb.control<boolean | null>(null),

      caregiver_type: this.fb.control<'ambos' | 'familia' | 'padre' | 'madre' | 'otro' | null>(null),
      caregiver_other: this.fb.control<string | null>(null),
      participates_in_upbringing: this.fb.control<boolean | null>(null),
      economic_support: this.fb.control<boolean | null>(null),
      monthly_amount: this.fb.control<number | null>(null),

      has_academic_support: this.fb.control<boolean | null>(null),
      support_provider: this.fb.control<'ambos' | 'familia' | 'padre' | 'madre' | 'otro' | null>(null),
      support_provider_other: this.fb.control<string | null>(null),

      interrupted_semester: this.fb.control<boolean | null>(null)
    });
    this.form.get('is_pregnant')?.valueChanges.subscribe(() => this.syncSections());
    this.form.get('has_children')?.valueChanges.subscribe(() => this.syncSections());
    this.form.get('gestation_selected')?.valueChanges.subscribe(() => this.updateChildrenCountValidation());
    this.form.get('gestation_selected')?.valueChanges.subscribe((value) => {
      if (value === true && this.isPregnantCtrl.value !== true) {
        this.isPregnantCtrl.setValue(true);
      }
    });
    this.form.get('is_pregnant')?.valueChanges.subscribe((value) => {
      if (value === false && this.gestationCtrl.value === true) {
        this.gestationCtrl.setValue(false);
      }
    });
    this.form.get('dni')?.valueChanges.subscribe((value) => this.onDniChange(value ?? ''));

    this.form.get('children_count')?.valueChanges.subscribe((value) => {
      const parsed = value === null ? null : Number(value);
      if (this.lastChildrenCount !== parsed) {
        this.resetChildrenAges();
        this.lastChildrenCount = parsed;
      }
      this.syncChildrenAges(parsed ?? 0);
    });

    this.form.get('residence_type')?.valueChanges.subscribe((value) => {
      this.toggleOtherField('residence_other', value === 'otro');
    });
    this.form.get('caregiver_type')?.valueChanges.subscribe((value) => {
      this.toggleOtherField('caregiver_other', value === 'otro');
    });
    this.form.get('support_provider')?.valueChanges.subscribe((value) => {
      this.toggleOtherField('support_provider_other', value === 'otro');
    });
    this.form.get('economic_support')?.valueChanges.subscribe((value) => {
      if (!value) {
        this.clearControl('monthly_amount');
      } else {
        this.requireControl('monthly_amount', [Validators.required, Validators.min(0), Validators.max(99999999.99)]);
      }
    });
    this.form.get('has_academic_support')?.valueChanges.subscribe((value) => {
      if (!value) {
        this.clearControl('support_provider');
        this.clearControl('support_provider_other');
      } else {
        this.requireControl('support_provider');
      }
    });
    this.form.get('lives_with_children')?.valueChanges.subscribe((value) => {
      if (value === false) {
        this.requireControl('caregiver_type');
      } else {
        this.clearControl('caregiver_type');
        this.clearControl('caregiver_other');
      }
    });
  }

  schoolOptions: string[] = [
    'AGRONOMIA',
    'ZOOTECNIA',
    'INGENIERIA EN INDUSTRIAS ALIMENTARIAS',
    'INGENIERIA FORESTAL',
    'INGENIERIA EN CONSERVACION DE SUELOS Y AGUA',
    'ADMINISTRACION',
    'CONTABILIDAD',
    'ECONOMIA',
    'INGENIERIA EN INFORMATICA Y SISTEMAS',
    'INGENIERIA AMBIENTAL',
    'INGENIERIA EN RECURSOS NATURALES RENOVABLES',
    'INGENIERIA MECANICA ELECTRICA',
    'INGENIERIA CIVIL',
    'TURISMO Y HOTELERIA',
    'INGENIERIA EN CIBERSEGURIDAD'
  ];

  get dniCtrl(): FormControl<string> {
    return this.form.controls.dni;
  }

  get studentCodeCtrl(): FormControl<string> {
    return this.form.controls.student_code;
  }

  get schoolCtrl(): FormControl<string | null> {
    return this.form.controls.school;
  }

  get academicCycleCtrl(): FormControl<string> {
    return this.form.controls.academic_cycle;
  }

  get ageCtrl(): FormControl<number | null> {
    return this.form.controls.age;
  }

  get genderCtrl(): FormControl<'F' | 'M' | null> {
    return this.form.controls.gender;
  }

  get isPregnantCtrl(): FormControl<boolean | null> {
    return this.form.controls.is_pregnant;
  }

  get hasChildrenCtrl(): FormControl<boolean | null> {
    return this.form.controls.has_children;
  }

  get childrenCountCtrl(): FormControl<number | null> {
    return this.form.controls.children_count;
  }

  get gestationCtrl(): FormControl<boolean> {
    return this.form.controls.gestation_selected;
  }

  get hasGestationSelected(): boolean {
    return this.gestationCtrl.value === true;
  }

  get hasChildSelection(): boolean {
    return this.hasGestationSelected || (this.childrenCountCtrl.value !== null && this.childrenCountCtrl.value > 0);
  }

  get canSubmit(): boolean {
    if (!this.verified || this.submitting) {
      return false;
    }
    if (this.showSectionIII && !this.hasChildSelection) {
      return false;
    }
    return this.form.valid;
  }

  get residenceTypeCtrl(): FormControl<'familiar' | 'alquiler' | 'residencia' | 'otro' | null> {
    return this.form.controls.residence_type;
  }

  get livesWithPartnerCtrl(): FormControl<boolean | null> {
    return this.form.controls.lives_with_partner;
  }

  get livesWithChildrenCtrl(): FormControl<boolean | null> {
    return this.form.controls.lives_with_children;
  }


  get caregiverTypeCtrl(): FormControl<'ambos' | 'familia' | 'padre' | 'madre' | 'otro' | null> {
    return this.form.controls.caregiver_type;
  }

  get participatesCtrl(): FormControl<boolean | null> {
    return this.form.controls.participates_in_upbringing;
  }

  get economicSupportCtrl(): FormControl<boolean | null> {
    return this.form.controls.economic_support;
  }

  get academicSupportCtrl(): FormControl<boolean | null> {
    return this.form.controls.has_academic_support;
  }

  get supportProviderCtrl(): FormControl<'ambos' | 'familia' | 'padre' | 'madre' | 'otro' | null> {
    return this.form.controls.support_provider;
  }

  get interruptedSemesterCtrl(): FormControl<boolean | null> {
    return this.form.controls.interrupted_semester;
  }

  showPlaceholder(control: FormControl<any>): boolean {
    return control.value === null || control.value === undefined;
  }

  get isPregnantValue(): boolean {
    return this.isPregnantCtrl.value === true;
  }

  get hasChildrenValue(): boolean {
    return this.hasChildrenCtrl.value === true;
  }

  get childrenCountValue(): number | null {
    return this.childrenCountCtrl.value;
  }

  get residenceTypeValue(): 'familiar' | 'alquiler' | 'residencia' | 'otro' | null {
    return this.residenceTypeCtrl.value;
  }

  get caregiverTypeValue(): 'ambos' | 'familia' | 'padre' | 'madre' | 'otro' | null {
    return this.form.controls.caregiver_type.value;
  }

  get supportProviderValue(): 'ambos' | 'familia' | 'padre' | 'madre' | 'otro' | null {
    return this.form.controls.support_provider.value;
  }

  get economicSupportValue(): boolean {
    return this.economicSupportCtrl.value === true;
  }

  get hasAcademicSupportValue(): boolean {
    return this.academicSupportCtrl.value === true;
  }

  get hasChild1(): boolean {
    const count = this.childrenCountValue;
    return count !== null && count >= 1;
  }

  get hasChild2(): boolean {
    const count = this.childrenCountValue;
    return count !== null && count >= 2;
  }

  get hasChild3(): boolean {
    const count = this.childrenCountValue;
    return count !== null && count >= 3;
  }

  get hasChild4(): boolean {
    const count = this.childrenCountValue;
    return count !== null && count >= 4;
  }

  get agesHint(): string {
    const count = this.childrenCountValue ?? 0;
    const gest = this.hasGestationSelected;
    if (!gest && count <= 0) {
      return 'Seleccione gestación o cantidad de hijos para habilitar edades.';
    }
    if (gest && count <= 0) {
      return 'Gestación seleccionada: la edad se registra automáticamente en 0.';
    }
    if (gest && count > 0) {
      return `Digite ${count} edad${count > 1 ? 'es' : ''} (y gestación en 0).`;
    }
    return `Digite ${count} edad${count > 1 ? 'es' : ''} según su selección.`;
  }

  verifyDni(): void {
    this.alreadyAnswered = false;
    this.successMessage = '';
    this.formMessage = '';
    this.errorMessage = '';
    this.studentInfo = null;
    const dni = this.dniCtrl.value;
    if (!dni) {
      this.dniCtrl.markAsTouched();
      return;
    }
    if (!/^\d{8}$/.test(dni)) {
      this.dniHint = 'El DNI debe tener 8 dígitos numéricos.';
      this.verified = false;
      return;
    }
    if (this.verifyingDni) {
      return;
    }
    this.verifyingDni = true;
    this.dniHint = 'Verificando DNI...';
    this.survey.verify(dni).subscribe({
      next: (data) => {
        this.verified = true;
        this.studentInfo = data;
        this.studentCodeCtrl.setValue(data?.student_code || '');
        this.schoolCtrl.setValue(data?.school || null);
        this.academicCycleCtrl.setValue(data?.academic_cycle || '');
        this.ageCtrl.setValue(data?.age ?? null);
        this.genderCtrl.setValue(data?.gender ?? null);
        this.lastVerifiedDni = dni;
        this.dniHint = 'DNI verificado correctamente.';
        this.verifyingDni = false;
        this.scrollToSection('section-family');
      },
      error: (err) => {
        this.verified = false;
        this.verifyingDni = false;
        if (err?.status === 409) {
          this.alreadyAnswered = true;
          this.dniHint = 'Este DNI ya registró la encuesta.';
          return;
        }
        if (err?.status === 404) {
          this.errorMessage = 'Usted no se encuentra matriculado actualmente.';
          this.dniHint = 'No se encontró matrícula para este DNI.';
          return;
        }
        this.dniHint = 'No se pudo verificar el DNI. Intente nuevamente.';
      }
    });
  }

  private onDniChange(raw: string): void {
    const dni = (raw || '').trim();
    this.alreadyAnswered = false;
    this.successMessage = '';
    this.formMessage = '';
    this.errorMessage = '';
    this.studentInfo = null;
    this.verified = false;

    if (dni.length === 0) {
      this.dniHint = 'Ingrese su DNI (8 dígitos numéricos).';
      return;
    }
    if (!/^\d*$/.test(dni)) {
      this.dniHint = 'Solo se permiten números en el DNI.';
      return;
    }
    if (dni.length < 8) {
      this.dniHint = `Faltan ${8 - dni.length} dígitos para completar el DNI.`;
      return;
    }
    if (dni.length > 8) {
      this.dniHint = 'El DNI debe tener exactamente 8 dígitos.';
      return;
    }
    if (this.lastVerifiedDni === dni) {
      this.dniHint = 'DNI verificado correctamente.';
      this.verified = true;
      return;
    }
    this.verifyDni();
  }

  submit(): void {
    this.successMessage = '';
    this.formMessage = '';
    if (!this.verified) {
      this.form.markAllAsTouched();
      this.formMessage = 'Por favor complete todas las preguntas obligatorias antes de enviar.';
      return;
    }
    if (this.showSectionIII && !this.hasChildSelection) {
      this.formMessage = 'Debe seleccionar gestación o al menos un hijo.';
      return;
    }
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      this.formMessage = 'Por favor complete todas las preguntas obligatorias antes de enviar.';
      return;
    }
    this.submitting = true;
    this.survey.submit(this.form.getRawValue()).subscribe({
      next: (data: any) => {
        this.submitting = false;
        this.successMessage = data?.partial
          ? 'Gracias por su participación.'
          : 'Encuesta terminada con éxito.';
        this.formMessage = '';
        window.scrollTo({ top: 0, behavior: 'smooth' });
        setTimeout(() => {
          this.router.navigate(['/home/services']);
        }, 2500);
      },
      error: (err) => {
        this.submitting = false;
        if (err?.status === 409) {
          this.alreadyAnswered = true;
        }
      }
    });
  }

  private syncSections(): void {
    const isPregnant = this.isPregnantValue;
    const hasChildren = this.hasChildrenValue;

    this.showSectionIII = isPregnant || hasChildren;
    this.showSectionIV = isPregnant || hasChildren;

    let targetSection: string | null = null;
    if (this.showSectionIII && !this.prevShowSectionIII) {
      targetSection = 'section-children';
    } else if (this.showSectionIV && !this.prevShowSectionIV) {
      targetSection = 'section-support';
    }
    if (targetSection) {
      this.scrollToSection(targetSection);
    }

    this.prevShowSectionIII = this.showSectionIII;
    this.prevShowSectionIV = this.showSectionIV;

    if (!this.showSectionIII) {
      this.clearSectionIII();
    } else {
      this.updateChildrenCountValidation();
    }

    if (!this.showSectionIV) {
      this.clearSectionIV();
    } else {
      this.requireSectionIV();
    }

    const isPregnantAnswered = this.isPregnantCtrl.value === false;
    const hasChildrenAnswered = this.hasChildrenCtrl.value === false;
    if (isPregnantAnswered && hasChildrenAnswered) {
      this.formMessage = 'Has seleccionado No/No. Puedes enviar la encuesta ahora.';
    } else if (this.formMessage === 'Has seleccionado No/No. Puedes enviar la encuesta ahora.') {
      this.formMessage = '';
    }

    // Guardado solo con el botón "Enviar encuesta"
  }

  private syncChildrenAges(count: number): void {
    if (!this.showSectionIII) {
      return;
    }

    if (Number.isNaN(count) || count <= 0) {
      this.disableControl('child_1_age');
      this.disableControl('child_2_age');
      this.disableControl('child_3_age');
      this.disableControl('child_4_age');
      this.resetChildrenAges();
      return;
    }

    this.enableControl('child_1_age');
    this.requireControl('child_1_age');
    if (count >= 2) {
      this.enableControl('child_2_age');
      this.requireControl('child_2_age');
    } else {
      this.disableControl('child_2_age');
      this.clearControl('child_2_age');
    }
    if (count >= 3) {
      this.enableControl('child_3_age');
      this.requireControl('child_3_age');
    } else {
      this.disableControl('child_3_age');
      this.clearControl('child_3_age');
    }
    if (count >= 4) {
      this.enableControl('child_4_age');
      this.requireControl('child_4_age');
    } else {
      this.disableControl('child_4_age');
      this.clearControl('child_4_age');
    }
  }

  private toggleOtherField(controlName: string, enabled: boolean): void {
    if (enabled) {
      this.requireControl(controlName);
    } else {
      this.clearControl(controlName);
    }
  }

  private requireSectionIV(): void {
    this.requireControl('residence_type');
    this.requireControl('lives_with_partner');
    this.requireControl('lives_with_children');
    this.requireControl('participates_in_upbringing');
    this.requireControl('economic_support');
    this.requireControl('has_academic_support');
    this.requireControl('interrupted_semester');
  }

  private clearSectionIII(): void {
    this.clearControl('children_count');
    this.gestationCtrl.setValue(false, { emitEvent: false });
    this.resetChildrenAges();
    this.disableControl('child_1_age');
    this.disableControl('child_2_age');
    this.disableControl('child_3_age');
    this.disableControl('child_4_age');
    this.lastChildrenCount = null;
  }

  private updateChildrenCountValidation(): void {
    if (!this.showSectionIII) {
      return;
    }
    if (this.hasGestationSelected) {
      this.childrenCountCtrl.clearValidators();
      this.childrenCountCtrl.updateValueAndValidity({ emitEvent: false });
      return;
    }
    this.requireControl('children_count');
  }

  private clearSectionIV(): void {
    this.clearControl('residence_type');
    this.clearControl('residence_other');
    this.clearControl('lives_with_partner');
    this.clearControl('lives_with_children');
    this.clearControl('caregiver_type');
    this.clearControl('caregiver_other');
    this.clearControl('participates_in_upbringing');
    this.clearControl('economic_support');
    this.clearControl('monthly_amount');
    this.clearControl('has_academic_support');
    this.clearControl('support_provider');
    this.clearControl('support_provider_other');
    this.clearControl('interrupted_semester');
  }

  private requireControl(name: string, validators: any[] = [Validators.required]): void {
    const control = this.form.get(name);
    if (!control) {
      return;
    }
    control.setValidators(validators);
    control.updateValueAndValidity({ emitEvent: false });
  }

  private clearControl(name: string): void {
    const control = this.form.get(name);
    if (!control) {
      return;
    }
    control.clearValidators();
    control.setValue(null, { emitEvent: false });
    control.updateValueAndValidity({ emitEvent: false });
  }

  private disableControl(name: string): void {
    const control = this.form.get(name);
    if (!control || control.disabled) {
      return;
    }
    control.disable({ emitEvent: false });
  }

  private enableControl(name: string): void {
    const control = this.form.get(name);
    if (!control || control.enabled) {
      return;
    }
    control.enable({ emitEvent: false });
  }

  private resetChildrenAges(): void {
    this.clearControl('child_1_age');
    this.clearControl('child_2_age');
    this.clearControl('child_3_age');
    this.clearControl('child_4_age');
  }

  private scrollToSection(id: string): void {
    setTimeout(() => {
      const el = document.getElementById(id);
      if (el) {
        el.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    }, 150);
  }

  onChildrenCountMouseDown(value: number, event: MouseEvent): void {
    const current = this.childrenCountCtrl.value;
    if (current === value) {
      event.preventDefault();
      this.childrenCountCtrl.setValue(null);
      this.lastChildrenCount = null;
      this.syncChildrenAges(0);
    }
  }

  toggleGestation(): void {
    this.gestationCtrl.setValue(!this.gestationCtrl.value);
  }

  selectChildrenCount(value: number): void {
    const current = this.childrenCountCtrl.value;
    if (current === value) {
      this.childrenCountCtrl.setValue(null);
      this.lastChildrenCount = null;
      this.syncChildrenAges(0);
      return;
    }
    this.childrenCountCtrl.setValue(value);
  }
}
