import {Component, computed, OnInit} from '@angular/core';
import {ThesisPracticing} from '../../../../core/models/ThesisPracticing';
import {FormBuilder, FormGroup, ReactiveFormsModule, Validators} from '@angular/forms';
import {ToastComponent} from '../../../../core/ui/toast/toast.component';
import {BlockUiComponent} from '../../../../core/ui/block-ui/block-ui.component';
import {ModalComponent} from '../../../../core/ui/modal/modal.component';
import {ToastService} from '../../../../core/services/toast/toast.service';
import {NgFor, NgIf} from '@angular/common';
import {Subscription} from 'rxjs';
import {ManagerService} from '../../../../core/services/manager/manager.service';
import {HttpErrorResponse} from '@angular/common/http';
import {
  DatasetEscuelaProfesional,
  DatasetFaculad,
  DatasetModalidadIngeso,
  DatasetSexo
} from '../../../../core/contans/postulation';
import * as XLSX from "xlsx";
import {ExportComponent} from "../../../../core/utils/export/export.component";

@Component({
  selector: 'app-thesis-practicing',
  standalone: true,
  imports: [ToastComponent, BlockUiComponent, ModalComponent, ReactiveFormsModule, NgIf, NgFor, ExportComponent],
  templateUrl: './thesis-practicing.component.html',
  styleUrl: './thesis-practicing.component.scss',
  providers: [ToastService]
})
export class ThesisPracticingComponent implements OnInit {
  private _subscriptions: Subscription = new Subscription();
  private typeModal: string = '';
  protected modalTitle: string = 'Crear Tesistas o practicantes';
  protected thesisPracticings: ThesisPracticing[] = [];
  protected formThesisPracticing: FormGroup;
  protected isLoading: boolean = false;
  protected thesisPracticingSelected!: ThesisPracticing;
  public _announcements: { name: string, value: number }[] = [];
  public typeStudents: { name: string, value: string }[] = [
    {name: 'Estudiante', value: 'Estudiante'},
    {name: 'Tesista', value: 'Tesista'},
    {name: 'Practicante', value: 'Practicante'},
  ];
  public DatasetSexo: { name: string, value: string }[] = DatasetSexo;
  public DatasetFaculad: { name: string, value: string }[] = DatasetFaculad;
  public DatasetEscuelaProfesional: { name: string, value: string }[] = DatasetEscuelaProfesional;
  public DatasetModalidadIngeso: { name: string, value: string }[] = DatasetModalidadIngeso;
  public showModal: boolean = false;
  public showDeleteModal: boolean = false;
  public showExportModal: boolean = false;
  protected orderStatus: Record<string, boolean> = {}

  constructor(private _fb: FormBuilder, private _toastService: ToastService, private _managerService: ManagerService,) {
    this.formThesisPracticing = this._fb.group({
      type_student: ['', Validators.required],
      cod_student: ['', Validators.required],
      name: ['', Validators.required],
      appaterno: ['', Validators.required],
      apmaterno: ['', Validators.required],
      sexo: ['', Validators.required],
      age: ['', Validators.required],
      facultad: ['', Validators.required],
      escuela_profesional: ['', Validators.required],
      mod_ingreso: ['', Validators.required],
      dni: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      convocatoria_id: ['', Validators.required],
      status_id: ['']
    });
  }

  ngOnInit() {
    this._getAnnouncementAll();
  }

  private _getAnnouncementAll() {
    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.getAnnouncement().subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this.isLoading = false;
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }

          this._announcements = res.detalle.map((item: any) => ({
            name: item.nombre,
            value: item.id !== undefined ? item.id : 0
          }));
          // this._announcement = res.detalle;
          this.getThesisPracticing();
        },
        error: (err: HttpErrorResponse) => {
          this.isLoading = false;
          this._toastService.add({type: 'error', message: 'No se pudo obtener la convocatoria'});
          console.error(err)
        }
      })
    );
  }

  protected onCreateThesisPracticing(): void {
    this.modalTitle = 'Crear Tesista o practicante'
    this.formThesisPracticing.reset();
    this.formThesisPracticing.enable();
    this.typeModal = 'create';
    this.showModal = true;
  }

  protected onUpdateThesisPracticing(usr: ThesisPracticing): void {
    this.thesisPracticingSelected = usr;
    this.modalTitle = 'Actualizar Usuario'
    this.formThesisPracticing.patchValue({
      type_student: usr.type_student,
      cod_student: usr.cod_student,
      apmaterno: usr.apmaterno,
      appaterno: usr.appaterno,
      age: usr.age,
      sexo: usr.sexo,
      facultad: usr.facultad,
      mod_ingreso: usr.mod_ingreso,
      escuela_profesional: usr.escuela_profesional,
      name: usr.name,
      dni: usr.dni,
      email: usr.email,
      convocatoria_id: usr.convocatoria_id,
      status_id: usr.status_id,
    })
    this.formThesisPracticing.enable();
    this.typeModal = 'update';
    this.showModal = true;
  }

  protected saveThesisPracticing(): void {
    if (this.formThesisPracticing.invalid) {
      this._toastService.add({type: 'error', message: 'Complete todos los campos correctamente.'});
      this.formThesisPracticing.markAllAsTouched();
      return;
    }

    if (this.typeModal === 'create') {
      this.createThesisPracticing();
      return;
    }

    this.updateThesisPracticing();
  }

  protected createThesisPracticing(): void {
    const data: ThesisPracticing = {
      cod_student: this.formThesisPracticing.value.cod_student,
      apmaterno: this.formThesisPracticing.value.apmaterno,
      appaterno: this.formThesisPracticing.value.appaterno,
      age: this.formThesisPracticing.value.age,
      sexo: this.formThesisPracticing.value.sexo,
      facultad: this.formThesisPracticing.value.facultad,
      mod_ingreso: this.formThesisPracticing.value.mod_ingreso,
      escuela_profesional: this.formThesisPracticing.value.escuela_profesional,
      type_student: this.formThesisPracticing.value.type_student,
      name: this.formThesisPracticing.value.name,
      dni: this.formThesisPracticing.value.dni,
      email: this.formThesisPracticing.value.email,
      convocatoria_id: this.formThesisPracticing.value.convocatoria_id,
      status_id: 1,
    };

    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.createThesisOrPracticing(data).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }

          this.getThesisPracticing();
          this._toastService.add({type: 'success', message: 'Tesista o practicante agregado correctamente.'});
          this.typeModal = '';
          this.showModal = false;

        },
        error: (err: any) => {
          this.isLoading = false;
          this._toastService.add({type: 'error', message: 'No se pudo crear el usuario, intente nuevamente'});
        }
      })
    );
  }

  protected updateThesisPracticing(): void {
    const data: ThesisPracticing = {
      cod_student: this.formThesisPracticing.value.cod_student,
      apmaterno: this.formThesisPracticing.value.apmaterno,
      appaterno: this.formThesisPracticing.value.appaterno,
      age: this.formThesisPracticing.value.age,
      sexo: this.formThesisPracticing.value.sexo,
      facultad: this.formThesisPracticing.value.facultad,
      mod_ingreso: this.formThesisPracticing.value.mod_ingreso,
      escuela_profesional: this.formThesisPracticing.value.escuela_profesional,
      type_student: this.formThesisPracticing.value.type_student,
      name: this.formThesisPracticing.value.name,
      dni: this.formThesisPracticing.value.dni,
      email: this.formThesisPracticing.value.email,
      convocatoria_id: this.formThesisPracticing.value.convocatoria_id,
      status_id: this.formThesisPracticing.value.status_id,
    };

    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.updateThesisOrPracticing(data, this.thesisPracticingSelected.id!).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }

          this.getThesisPracticing();
          this._toastService.add({
            type: 'success',
            message: 'Practicante o tesista e actualizo correctamente!'
          });
          this.showModal = false;
          this.typeModal = '';

        },
        error: (err: any) => {
          this.isLoading = false;
          console.error(err);
          this._toastService.add({
            type: 'error',
            message: 'No se pudo actualizar el tesistas o practicantes, intente nuevamente'
          });
        }
      })
    );
  }

  private getThesisPracticing(): void {
    this.isLoading = true;

    this._subscriptions.add(
      this._managerService.getThesisPracticing().subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo obtner los tesistas o practicantes, intente nuevamente'
            });
            return;
          }

          this.thesisPracticings = res.detalle;
        },
        error: (err: any) => {
          console.error(err);
          this._toastService.add({
            type: 'error',
            message: 'No se pudo obtener los tesistas o practicantes, intente nuevamente'
          });
          this.isLoading = false;
        }
      })
    );
  }

  protected deleteThesisPracticing(): void {

    this.isLoading = false;
    this._subscriptions.add(
      this._managerService.deleteThesisOrPracticing(this.thesisPracticingSelected.id!).subscribe({
        next: (res: any) => {
          this.isLoading = false;

          if (!res.detalle) {
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }

          this.getThesisPracticing();
          this.typeModal = 'delete';
        },
        error: (err: any) => {
          console.error(err);
          this._toastService.add({type: 'error', message: 'No se pudo borrar al usuario, intente nuevamente.'});
          this.isLoading = false;
        }
      })
    );

  }


  public getNameTypeStudentByValue(value: string) {
    return this.typeStudents.find((item) => item.value === value)?.name;
  }

  public getNameConvocatoriaById(id: number) {
    return this._announcements.find((item) => item.value == id)?.name;
  }

  exportToExcel() {
    this.showExportModal = true;
  }


  findRequest($event: Event) {

  }

  orderByField(field: keyof ThesisPracticing) {
    if (this.orderStatus[field] === undefined) {
      this.orderStatus[field] = true;
    }

    this.thesisPracticings = this.thesisPracticings.sort((a, b) => {
      const fieldA = a[field] ? a[field].toString().toLowerCase() : '';
      const fieldB = b[field] ? b[field].toString().toLowerCase() : '';

      let comparison = 0;
      if (fieldA < fieldB) {
        comparison = -1;
      } else if (fieldA > fieldB) {
        comparison = 1;
      }
      return this.orderStatus[field] ? comparison : -comparison;
    });
    this.orderStatus[field] = !this.orderStatus[field];
  }
}

