import { Component } from '@angular/core';
import {NgIf} from "@angular/common";
import {HouseholdTargeting} from "../../../../core/models/household-targeting";
import * as XLSX from "xlsx";
import {ModalComponent} from "../../../../core/ui/modal/modal.component";
import {
  ScholarshipStudentsFormComponent
} from "./components/scholarship-students-form/scholarship-students-form.component";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ManagerService} from "../../../../core/services/manager/manager.service";
import {ScholarshipStudentsModel} from "./models/scholarship-students.model";
import {
  ScholarshipStudentsTableComponent
} from "./components/scholarship-students-table/scholarship-students-table.component";
import {IRequest} from "../../../../core/models/requests";

@Component({
  selector: 'app-scholarship-students',
  standalone: true,
  imports: [
    NgIf,
    ModalComponent,
    ScholarshipStudentsFormComponent,
    BlockUiComponent,
    ToastComponent,
    ScholarshipStudentsTableComponent
  ],
  templateUrl: './scholarship-students.component.html',
  styleUrl: './scholarship-students.component.scss',
  providers: [ToastService],
})
export class ScholarshipStudentsComponent {
  public formScholarshipStudents: FormGroup;
  public scholarshipStudents: ScholarshipStudentsModel[] = []
  public modalTitle: string = '';
  public showModal: boolean = false;
  public isLoading: boolean = false;


  constructor(private _fb: FormBuilder,
              private _toastService: ToastService,
              ) {
    this.formScholarshipStudents = this._fb.group({
      dni: ['', Validators.required],
      names: ['', Validators.required],
      career: ['', Validators.required],
      program: ['', Validators.required],
      year_announcement: ['', Validators.required],
      status: ['', Validators.required],
      condition: ['', Validators.required],
    })
  }

  uploadData(event: any) {
    const file = event.target.files[0];
    const reader = new FileReader();

    reader.onload = (e: any) => {
      const binaryStr = e.target.result;
      const workbook = XLSX.read(binaryStr, { type: 'binary' });
      const worksheet = workbook.Sheets[workbook.SheetNames[0]];
      const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });

      const headers = jsonData[0];
      const data = jsonData.slice(1);

      const dataFormat = data.map((row: any) => {
        // @ts-ignore
        return headers.reduce((acc: any, header: string, index: number) => {
          acc[header] = row[index];
          return acc;
        }, {});
      });
      this.scholarshipStudents = dataFormat.map((row: any) => ({
        dni: row["DNI"],
        names: row['NOMBRES COMPLETOS'],
        program: row['PROGRAMA'],
        year_announcement: row['CONVOC.'],
        career: row['CARRERA'],
        status: row['ESTADO'],
        condition: row["CONDICION"],
      }));

    };
    reader.readAsBinaryString(file);
  }

  public onSave(): void {
    if (this.formScholarshipStudents.invalid) {
      this.formScholarshipStudents.markAllAsTouched();
      this.eventMessage('info', 'Complete el formulario!');
    }

  }

  public addScholarshipStudent() {
    this.modalTitle = 'Registrar un becado';
    this.showModal = true;
  }

  public eventMessage(type: "error" | "success" | "warning" | "info", message: string): void {
    this._toastService.add({
      type: type,
      message: message,
    });
  }
}
