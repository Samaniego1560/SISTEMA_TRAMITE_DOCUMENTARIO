import {
  ChangeDetectionStrategy,
  Component,
  inject,
  Input,
  OnChanges,
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
import {LaboratoryExam} from "../../../../../core/models/areas/nursing/nursing.model";

@Component({
  selector: 'app-laboratory-exam-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ControlErrorsDirective
  ],
  templateUrl: './laboratory-exam-form.component.html',
  styleUrl: './laboratory-exam-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LaboratoryExamFormComponent extends FormBaseComponent<LaboratoryExam> implements OnChanges {
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }
  deleteForm = output();

  public active = signal<boolean>(true);
  public form: FormGroup;
  private _name = signal('examen_laboratorio');

  private _fb = inject(FormBuilder);

  constructor() {
    super();
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if(changes['infoData'] && this.infoData) {
      this.form.patchValue(this.infoData);
      this.active.set(false);
      if(!this.canSaveForm) this.form.disable();
    }
  }

  private _buildForm() {
    return this._fb.nonNullable.group({
      serologia: ['',],
      bk: ['',],
      hemograma: ['',],
      examen_orina: ['',],
      colesterol: ['',],
      glucosa: ['',],
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
