import {Component, computed, effect, EventEmitter, Input, input, model, OnInit, Output, signal} from '@angular/core';
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {DecimalPipe, JsonPipe, NgIf} from "@angular/common";
import {IRequirement} from "../../../../../../core/models/announcement";
import {ModalComponent} from "../../../../../../core/ui/modal/modal.component";
import {ToastService} from "../../../../../../core/services/toast/toast.service";

@Component({
  selector: 'app-form-requirement',
  standalone: true,
  imports: [
    FormsModule,
    NgIf,
    ReactiveFormsModule,
    ModalComponent,
    JsonPipe,
    DecimalPipe
  ],
  templateUrl: './form-requirement.component.html',
  styleUrl: './form-requirement.component.scss'
})
export class FormRequirementComponent {
  typeProcess = input.required<string>();
  showModal = model<boolean>(false);
  requirement = input.required<IRequirement>();

  requirementNames = input.required<string[]>();

  protected formRequirement: FormGroup;
  @Output() onSaveRequirement: EventEmitter<IRequirement> = new EventEmitter<IRequirement>();

  protected titleModal: string = 'Creación de requisito';

  constructor(private fb: FormBuilder,
              private _toastService: ToastService,) {
    this.formRequirement = this.fb.group({
      name: ['', Validators.required],
      description: [''],
      guide: [''],
      format: [0, Validators.required],
      active: [false, Validators.required],
      url_template: [''],
      options: [''],
      text_color: [''],
      text_type: [''],
      text_size: [''],
      type_input: [''],
      is_recoverable: [''],
      is_dependent: [''],
      field_dependent: [''],
      value_dependent: [''],
      show_dependent: [''],
    });

    effect(() => {
      if (this.typeProcess() === 'edit') {
        this.titleModal = 'Edición de requisito';
        this.loadFormRequirement(this.requirement());
        this.formRequirement.enable();
      }
      if (this.typeProcess() === 'view') {
        this.titleModal = 'Vista del requisito';
        this.loadFormRequirement(this.requirement());
        this.formRequirement.disable()
      }

      if (this.typeProcess() === 'create') {
        this.titleModal = 'Creación de requisito';
        this.resetFormRequirement()
        this.formRequirement.enable();
      }
    });
  }

  public formatOptions = [
    { value: 1, label: 'Documento' },
    { value: 2, label: 'Imagen' },
    { value: 3, label: 'Formulario' },
    { value: 4, label: 'Selección' },
    { value: 5, label: 'Radio' },
    { value: 6, label: 'Checkbox' }
  ];

  protected processFile(event: any): void {
    const file: File = event.target.files[0]
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = (res: any) => {
      this.formRequirement.get('guide')?.setValue(res.target?.result);
    };
  }

  protected changeFormat(): void {
    if (this.formRequirement.value.format === '1' || this.formRequirement.value.format === '2') {
      this.formRequirement.get('guide')?.setValidators([Validators.required]);
      this.formRequirement.get('guide')?.updateValueAndValidity();
      return;
    }

    this.formRequirement.get('guide')?.clearValidators();
    this.formRequirement.get('guide')?.updateValueAndValidity();
    this.formRequirement.get('guide')?.setValue('');
    return;
  }

  private loadFormRequirement(req: IRequirement) {
    this.formRequirement.patchValue({
      name: req.nombre,
      description: req.descripcion,
      guide: req.url_guia,
      format: req.tipo_requisito_id.toString(),
      options: req.opciones,
      text_color: req.text_color,
      text_type: req.text_type,
      text_size: req.text_size,
      type_input: req.type_input,
      active: req.activo,
      is_recoverable: req.is_recoverable,
      is_dependent: req.is_dependent,
      field_dependent: req.field_dependent,
      value_dependent: req.value_dependent,
      show_dependent: req.show_dependent,
    })
  }

  protected createOrUpdateRequirement(): void {
    if (this.formRequirement.invalid) {
      this._toastService.add({type: 'error', message: 'Complete los campos correctamente'});
      this.formRequirement.markAllAsTouched();
      return;
    }

    const requirement: IRequirement = {
      activo: this.formRequirement.value.active,
      url_guia: this.formRequirement.value.guide,
      descripcion: this.formRequirement.value.description,
      nombre: this.formRequirement.value.name,
      url_plantilla: this.formRequirement.value.url_template,
      tipo_requisito_id: this.formRequirement.value.format,
      text_color: this.formRequirement.value.text_color,
      text_type: this.formRequirement.value.text_type,
      text_size: this.formRequirement.value.text_size,
      type_input: this.formRequirement.value.type_input,
      is_recoverable: this.formRequirement.value.is_recoverable,
      is_dependent: this.formRequirement.value.is_dependent,
      field_dependent: this.formRequirement.value.field_dependent,
      value_dependent: this.formRequirement.value.value_dependent,
      show_dependent: this.formRequirement.value.show_dependent,
      opciones: this.formRequirement.value.options,
    }
    this.onSaveRequirement.emit(requirement);
    this.resetFormRequirement();
    this.showModal.set(false);
  }

  private resetFormRequirement() {
    this.formRequirement.reset();
  }
}
