import {Component, input} from '@angular/core';
import {FormBuilder, FormGroup, ReactiveFormsModule} from "@angular/forms";
import {NgForOf, NgIf} from "@angular/common";
import {DatasetEscuelaProfesional} from "../../../../../../core/contans/postulation";
import {DatasetCondition, DatasetStatus} from "../../models/scholarship-students.contant";

@Component({
  selector: 'app-scholarship-students-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    NgIf,
    NgForOf
  ],
  templateUrl: './scholarship-students-form.component.html',
  styleUrl: './scholarship-students-form.component.scss'
})
export class ScholarshipStudentsFormComponent {
  formScholarshipStudents = input.required<FormGroup>();
  protected readonly DatasetProfessionalSchool = DatasetEscuelaProfesional;
  protected readonly DatasetStatus = DatasetStatus;
  protected readonly DatasetCondition = DatasetCondition;
}
