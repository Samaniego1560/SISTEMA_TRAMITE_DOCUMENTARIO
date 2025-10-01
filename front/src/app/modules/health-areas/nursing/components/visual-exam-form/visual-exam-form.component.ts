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
import {VisualExam, PhysicalExam} from "../../../../../core/models/areas/nursing/nursing.model";

@Component({
  selector: 'app-visual-exam-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ControlErrorsDirective
  ],
  templateUrl: './visual-exam-form.component.html',
  styleUrl: './visual-exam-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class VisualExamFormComponent extends FormBaseComponent<VisualExam> implements OnChanges {
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }

  deleteForm = output();

  public active = signal<boolean>(true);
  public form: FormGroup;
  private _name = signal('examen_visual');

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
    return this._fb.group({
      ojo_derecho: [''],
      ojo_izquierdo: [''],
      presion_ocular: [''],
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
