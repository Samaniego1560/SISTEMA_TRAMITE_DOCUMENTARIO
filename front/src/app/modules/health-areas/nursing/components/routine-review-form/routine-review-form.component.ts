import {Component, forwardRef, inject, Input, OnInit} from '@angular/core';
import {
  AbstractControl, ControlValueAccessor,
  FormBuilder,
  FormGroup, NG_VALIDATORS,
  NG_VALUE_ACCESSOR,
  ReactiveFormsModule,
  ValidationErrors, Validator
} from "@angular/forms";
import {NgxValidators} from "../../../../../lib/validator-dynamic/utils/ngx-validators.util";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {ToggleSwitchComponent} from "../../../../../core/ui/toggle/toggle-switch.component";

@Component({
  selector: 'app-routine-review-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ControlErrorsDirective,
    ToggleSwitchComponent
  ],
  templateUrl: './routine-review-form.component.html',
  styleUrl: './routine-review-form.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => RoutineReviewFormComponent),
      multi: true
    },
    {
      provide: NG_VALIDATORS,
      useExisting: forwardRef(() => RoutineReviewFormComponent),
      multi: true
    }
  ],
})
export class RoutineReviewFormComponent implements OnInit, ControlValueAccessor, Validator {
  @Input() genero: string = 'M';

  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.routineReviewForm.markAllAsTouched();
  }

  public routineReviewForm: FormGroup;
  public _onChange: Function = (_value: any) => {
  };
  public _onTouched: Function = (_value: any) => {
  };

  private _fb = inject(FormBuilder);

  constructor() {
    this.routineReviewForm = this._buildRoutineForm();
  }

  ngOnInit() {
    this.routineReviewForm.valueChanges.subscribe(value => this._onChange(value));
  }

  private _buildRoutineForm() {
    return this._fb.group({
      fiebre_ultimo_quince_dias: ['no'],
      tos_mas_quince_dias: ['no'],
      secrecion_lesion_genitales: ['no'],
      fecha_ultima_regla: ['no',],
      comentarios: ['',],
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
      this.routineReviewForm.patchValue(obj);
    }
  }

  validate(_control: AbstractControl): ValidationErrors | null {
    return this.routineReviewForm.valid ? null : {invalidRoutineReview: true};
  }

  setDisabledState(isDisabled: boolean) {
    if (isDisabled) {
      this.routineReviewForm.disable();
    } else {
      this.routineReviewForm.enable();
    }
  }
}
