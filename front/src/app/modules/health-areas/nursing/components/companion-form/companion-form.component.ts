import {ChangeDetectionStrategy, Component, forwardRef, inject, Input, OnInit, signal} from '@angular/core';
import {
  AbstractControl,
  ControlValueAccessor,
  FormBuilder,
  FormGroup,
  NG_VALIDATORS,
  NG_VALUE_ACCESSOR,
  ReactiveFormsModule, ValidationErrors, Validator
} from "@angular/forms";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";

@Component({
  selector: 'app-companion-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ControlErrorsDirective
  ],
  templateUrl: './companion-form.component.html',
  styleUrl: './companion-form.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => CompanionFormComponent),
      multi: true
    },
    {
      provide: NG_VALIDATORS,
      useExisting: forwardRef(() => CompanionFormComponent),
      multi: true
    }
  ],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CompanionFormComponent implements OnInit, ControlValueAccessor, Validator {
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.companionForm.markAllAsTouched();
  }

  public active = signal<boolean>(true);
  public companionForm: FormGroup;
  public _onChange: Function = (_value: any) => {
  };
  public _onTouched: Function = (_value: any) => {
  };

  private _fb = inject(FormBuilder);

  constructor() {
    this.companionForm = this._buildCompanionForm();
  }

  ngOnInit() {
    this.companionForm.valueChanges.subscribe(value => this._onChange(value));
  }

  private _buildCompanionForm() {
    return this._fb.group({
      dni: [null],
      nombres_apellidos: [null],
      edad: [null],
    })
  }

  registerOnChange(fn: any): void {
    this._onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this._onTouched = fn;
  }

  writeValue(obj: any): void {
    if (obj) {
      this.companionForm.patchValue(obj);
    }
  }

  validate(_control: AbstractControl): ValidationErrors | null {
    return this.companionForm.valid ? null : {invalidCompanion: true};
  }

  setDisabledState(isDisabled: boolean) {
    if (isDisabled) {
      this.companionForm.disable();
    } else {
      this.companionForm.enable();
    }
  }

}
