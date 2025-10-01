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
import {FormBuilder, FormGroup, ReactiveFormsModule,} from "@angular/forms";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {PhysicalExam} from "../../../../../core/models/areas/nursing/nursing.model";

@Component({
  selector: 'app-physical-exam-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    ControlErrorsDirective
  ],
  templateUrl: './physical-exam-form.component.html',
  styleUrl: './physical-exam-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class PhysicalExamFormComponent extends FormBaseComponent<PhysicalExam> implements OnChanges {
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }
  deleteForm = output();

  public active = signal<boolean>(true);
  public form: FormGroup;
  private _name = signal('examen_fisico');

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
      talla_peso: [''],
      perimetro_cintura: [''],
      indice_masa_corporal_img: [''],
      presion_arterial: [''],
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
