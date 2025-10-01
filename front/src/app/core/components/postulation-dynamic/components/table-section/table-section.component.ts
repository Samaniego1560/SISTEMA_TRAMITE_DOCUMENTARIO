import {Component, input} from '@angular/core';
import {IRequirement, ISection} from "../../../../models/announcement";
import {FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
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
export class TableSectionComponent {

  section = input.required<ISection>();
  arrayValue = input.required<any[]>();
  recordRequirements = input.required<Record<string, string>>()

  protected formSection: FormGroup;
  protected activeForm: boolean = false;
  protected indexSelect: number = 0;
  protected isEdit: boolean = false;

  constructor(private _fb: FormBuilder,) {
    this.formSection = this._fb.group({})
  }

  public loadFormSection(requirements: IRequirement[]): void {
    for (const req of requirements) {
      const key = (req.id || 0).toString();
      this.formSection.addControl(
        key,
        new FormControl('', Validators.required)
      );
    }
    this.activeForm = true;
  }

  public saveForm() {
    if (!this.formSection.valid) {
      this.formSection.markAllAsTouched()
      return;
    }
    if (!this.isEdit) {
      this.arrayValue().push(this.formSection.value);
    } else {
      this.arrayValue()[this.indexSelect] = this.formSection.value;
    }
    this.resetForm();
  }

  public resetForm() {
    this.formSection.reset();
    this.activeForm = false;
  }

  public deleteRow(index: number): void {
    this.arrayValue().splice(index, 1);
  }

  editRow(item: any, i: number) {
    this.loadFormSection(this.section().requisitos)
    this.indexSelect = i;
    this.formSection.patchValue({...item});
    this.isEdit = true;
    this.activeForm = true;
  }
}
