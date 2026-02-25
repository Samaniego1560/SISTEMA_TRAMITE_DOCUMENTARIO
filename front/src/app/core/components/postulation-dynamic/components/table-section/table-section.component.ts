import {Component, EventEmitter, input, OnInit, Output} from '@angular/core';
import {IRequirement, ISection} from "../../../../models/announcement";
import {FormArray, FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
import {SectionFormComponent} from "../section-form/section-form.component";

@Component({
  selector: 'app-table-section',
  standalone: true,
  imports: [
    SectionFormComponent

  ],
  templateUrl: './table-section.component.html',
  styleUrl: './table-section.component.scss'
})
export class TableSectionComponent implements OnInit {

  section = input.required<ISection>();
  formArray = input.required<FormArray>();
  recordRequirements = input.required<Record<string, string>>()
  @Output() tableFormStateChange = new EventEmitter<{ sectionId: string; activeForm: boolean }>();

  protected formSection: FormGroup;
  protected activeForm: boolean = false;
  protected indexSelect: number = 0;
  protected isEdit: boolean = false;

  constructor(private _fb: FormBuilder,) {
    this.formSection = this._fb.group({})
  }

  public loadFormSection(requirements: IRequirement[]): void {
    this.formSection = this._fb.group({});
    for (const req of requirements) {
      const key = (req.id || 0).toString();
      this.formSection.addControl(
        key,
        new FormControl('', Validators.required)
      );
    }
    this.activeForm = true;
    this.emitState();
  }

  public saveForm() {
    if (!this.formSection.valid) {
      this.formSection.markAllAsTouched()
      return;
    }
    const rowGroup = this._fb.group({...this.formSection.value});
    if (!this.isEdit) {
      this.formArray().push(rowGroup);
    } else {
      this.formArray().setControl(this.indexSelect, rowGroup);
    }
    this.resetForm();
  }

  public resetForm() {
    this.formSection.reset();
    this.activeForm = false;
    this.isEdit = false;
    this.indexSelect = 0;
    this.emitState();
  }

  public deleteRow(index: number): void {
    this.formArray().removeAt(index);
  }

  editRow(item: any, i: number) {
    this.loadFormSection(this.section().requisitos)
    this.indexSelect = i;
    this.formSection.patchValue({...item});
    this.isEdit = true;
    this.activeForm = true;
    this.emitState();
  }

  get rows(): any[] {
    return this.formArray().controls;
  }

  ngOnInit(): void {
    if (this.isFamilyCompositionSection() && this.formArray().length === 0) {
      this.loadFormSection(this.section().requisitos);
    } else {
      this.emitState();
    }
  }

  private emitState(): void {
    this.tableFormStateChange.emit({
      sectionId: (this.section().id || 0).toString(),
      activeForm: this.activeForm,
    });
  }

  private isFamilyCompositionSection(): boolean {
    const sectionText = this.normalizeText(this.section().descripcion || '');
    if (sectionText.includes('composicion familiar')) {
      return true;
    }

    const requirementNames = (this.section().requisitos || [])
      .map((req) => this.normalizeText(req.nombre || ''));

    const expectedFields = ['nombres y apellidos', 'edad', 'estado civil', 'parentesco', 'nivel educativo', 'ocupacion'];
    return expectedFields.every((field) => requirementNames.includes(field));
  }

  private normalizeText(value: string): string {
    return value
      .toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '')
      .replace(/\s+/g, ' ')
      .trim();
  }
}
