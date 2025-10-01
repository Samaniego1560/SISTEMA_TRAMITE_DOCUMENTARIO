import {Component, inject, Input, OnChanges, output, signal, SimpleChanges} from '@angular/core';
import {ControlErrorsDirective} from "../../../../../lib/validator-dynamic/directives/control-error.directive";
import {FormBuilder, FormGroup, ReactiveFormsModule} from "@angular/forms";
import {FormBaseComponent} from "../../../../../core/models/areas/form-base.model";
import {ConsultaOdontologia} from "../../../../../core/models/areas/dentistry.model";

@Component({
  selector: 'app-integral-attention-form',
  standalone: true,
    imports: [
        ControlErrorsDirective,
        ReactiveFormsModule
    ],
  templateUrl: './integral-attention-form.component.html',
  styleUrl: './integral-attention-form.component.scss'
})
export class IntegralAttentionFormComponent extends FormBaseComponent<ConsultaOdontologia> implements OnChanges {
  @Input() set touched(isTouched: boolean) {
    if (isTouched) this.form.markAllAsTouched();
  }

  deleteForm = output();

  public active = signal<boolean>(true);
  public form: FormGroup;
  private _name = signal('consulta');

  private _fb = inject(FormBuilder);

  constructor() {
    super();
    this.form = this._buildForm();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['infoData']?.currentValue && this.infoData) {
      this.form.patchValue(this.infoData);
      this.active.set(false);
      this.form.disable()
    }
  }

  private _buildForm() {
    return this._fb.group({
      relato: ['',], //listo
      diagnostico: ['',], //listo
      examen_auxiliar: ['',], //listo
      examen_clinico: ['',], //listo
      tratamiento: ['',], //listo
      indicaciones: ['',], //listo
      comentarios: ['',], //listo
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
