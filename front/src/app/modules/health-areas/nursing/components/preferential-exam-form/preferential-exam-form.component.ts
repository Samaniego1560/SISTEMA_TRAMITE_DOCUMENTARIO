import {
  ChangeDetectionStrategy,
  Component,
  inject,
  Input,
  OnChanges, OnInit,
  output,
  signal,
  SimpleChanges
} from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
} from "@angular/forms";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {PreferentialExam} from "../../../../../core/models/areas/nursing/nursing.model";
import {ACTION_TYPES} from "../../../../../core/contans/areas/permissions.constant";
import {PermissionsService} from "../../../../../core/services/permissions/permissions.service";

@Component({
  selector: 'app-preferential-exam-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ControlErrorsDirective
  ],
  templateUrl: './preferential-exam-form.component.html',
  styleUrl: './preferential-exam-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class PreferentialExamFormComponent extends FormBaseComponent<PreferentialExam> implements OnChanges {
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }

  deleteForm = output();

  public active = signal<boolean>(true);
  public form!: FormGroup;
  private _name = signal<string>('examen_preferencial');

  private _fb = inject(FormBuilder);

  constructor() {
    super();
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['infoData'] && this.infoData) {
      this.form.patchValue(this.infoData);
      this.active.set(false);
      if(!this.canSaveForm) this.form.disable();
    }
  }

  private _buildForm() {
    return this._fb.group({
      aparato_respiratorio: [''],
      aparato_cardiovascular: [''],
      aparato_digestivo: [''],
      aparato_genitourinario: [''],
      papanicolau: [''],
      examen_mama: [''],
      comentarios: ['',],
    })
  }

  public get name() {
    return this._name();
  }

  public isValid(): boolean {
    return this.form.valid
  }

  public getData() {
    return this.form.value
  }
}
