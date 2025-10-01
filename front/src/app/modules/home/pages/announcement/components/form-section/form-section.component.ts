import {Component, effect, EventEmitter, input, model, Output} from '@angular/core';
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {ModalComponent} from "../../../../../../core/ui/modal/modal.component";
import {NgIf} from "@angular/common";
import { ISection} from "../../../../../../core/models/announcement";
import {ToastService} from "../../../../../../core/services/toast/toast.service";

@Component({
  selector: 'app-form-section',
  standalone: true,
  imports: [
    FormsModule,
    ModalComponent,
    NgIf,
    ReactiveFormsModule
  ],
  templateUrl: './form-section.component.html',
  styleUrl: './form-section.component.scss'
})
export class FormSectionComponent {
  typeProcess = input.required<string>();
  showModal = model<boolean>(false);
  section = input.required<ISection>();
  protected formSection: FormGroup;
  @Output() onSaveSection: EventEmitter<ISection> = new EventEmitter<ISection>();
  protected titleModal: string = 'Creación de sección';

  constructor(private fb: FormBuilder,
              private _toastService: ToastService,) {
    this.formSection = this.fb.group({
      description: ['', Validators.required],
      type: ['form', Validators.required],
    });

    effect(() => {
      if (this.typeProcess() === 'edit') {
        this.titleModal = 'Edición de sección';
        this.loadFormSection(this.section());
        this.formSection.enable();
      }
      if (this.typeProcess() === 'view') {
        this.titleModal = 'Vista del sección';
        this.loadFormSection(this.section());
        this.formSection.disable()
      }

      if (this.typeProcess() === 'create') {
        this.titleModal = 'Creación de sección';
        this.resetFormSection()
        this.formSection.enable();
      }
    });
  }

  private loadFormSection(req: ISection) {
    this.formSection.patchValue({
      description: req.descripcion,
      type: req.type,
    })
  }

  private resetFormSection() {
    this.formSection.reset();
  }

  protected createOrUpdateSection() {
    if (this.formSection.invalid) {
      this._toastService.add({type: 'error', message: 'Complete los campos correctamente'});
      this.formSection.markAllAsTouched();
      return;
    }

    const section: ISection = {
      descripcion: this.formSection.value.description,
      type: this.formSection.value.type,
      requisitos: []
    }
    this.onSaveSection.emit(section);
    this.resetFormSection();
    this.showModal.set(false);
  }
}
