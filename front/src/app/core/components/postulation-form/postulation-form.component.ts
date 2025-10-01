import {Component, EventEmitter, input, model, Output} from '@angular/core';
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {NgForOf, NgIf} from "@angular/common";
import {Subscription} from "rxjs";
import {ToastService} from "../../services/toast/toast.service";
import {ManagerService} from "../../services/manager/manager.service";
import {IDebtsStudent, IValidationUser} from "../../models/user";
import {HttpErrorResponse} from "@angular/common/http";
import {ModalComponent} from "../../ui/modal/modal.component";
import {AlertDebtsComponent} from "./components/alert-debts/alert-debts.component";
import {IAnnouncement} from "../../models/announcement";
import {BlockUiComponent} from "../../ui/block-ui/block-ui.component";
import {Student} from "../../models/student";

@Component({
  selector: 'app-postulation-form',
  standalone: true,
  imports: [
    FormsModule,
    NgForOf,
    NgIf,
    ReactiveFormsModule,
    ModalComponent,
    AlertDebtsComponent,
    BlockUiComponent
  ],
  templateUrl: './postulation-form.component.html',
  styleUrl: './postulation-form.component.scss'
})
export class PostulationFormComponent {
  postulationId = input.required<number>();
  protected formPostulation: FormGroup;
  showModal = model.required<boolean>();
  roleId = input.required<number>();
  public showModalDebt: boolean = false;

  @Output() onPostulation: EventEmitter<{postulation: IAnnouncement, student: Student}> =
    new EventEmitter<{postulation: IAnnouncement, student: Student}>();

  protected isLoading: boolean = false;
  private _subscriptions: Subscription = new Subscription();

  public debts: IDebtsStudent[] = [];

  public typeStudents: { name: string; value: string }[] = [
    { name: 'Estudiante', value: 'Estudiante' },
    { name: 'Tesista', value: 'Tesista' },
    { name: 'Practicante', value: 'Practicante' },
  ];

  constructor(private _fb: FormBuilder,
              private _toastService: ToastService,
              private _managerService: ManagerService,) {
    this.formPostulation = this._fb.group({
      dni_student: ['', Validators.required],
      type_student: ['Estudiante', Validators.required],
      email_student: ['', [Validators.required, Validators.email]],
      eat_service: [false, Validators.required],
      resident_service: [false, Validators.required],
    });
  }

  protected async validateDebts() {
    if (this.formPostulation.invalid) {
      this.createMessage('error', 'Complete todos los campos correctamente!')
      this.formPostulation.markAllAsTouched();
      return;
    }

    if (!this.formPostulation.value.eat_service && !this.formPostulation.value.resident_service) {
      this.createMessage('error','Seleccione al menos un servicio!')
      return;
    }

    const typeStudent = this.formPostulation.value.type_student;
    const dni = this.formPostulation.value.dni_student;
    const email = this.formPostulation.value.email_student

    const responseDebts = await this.validateStudentDebtsPromise(dni);
    if (responseDebts.error) {debugger
      this.createMessage(responseDebts.type, responseDebts.msg)
      return;
    }
debugger
    if (!this.evaluateDebts(responseDebts.data)) return;

        const data: IValidationUser = {
      type_student: typeStudent,
      correo: email,
      DNI: dni,
      announcement_id: this.postulationId()
    };

    if (this.roleId() === 1) {
      const response = await this.validateStudentPermissionPromise(data)
      if (response.error || !response.data) {
        this.createMessage(response.type, response.msg)
        return;
      }
    } else {
      const response = await this.validateStudentThesisOrPracticingPromise(data)
      if (response.error || !response.data) {
        this.createMessage('info', response.msg)
        return;
      }
    }

    /*const responseSign = await this.validateStudentSignAreasPromise(dni);
    debugger
    if (responseSign.error) {
      //this.createMessage(responseStudent.type, responseStudent.msg)
      return;
    }*/

    const responseStudent = await this.getDataAnnouncementAndStudent(this.postulationId(), dni);
    if (responseStudent.error) {
      this.createMessage(responseStudent.type, responseStudent.msg)
      return;
    }

    const eatService = this.formPostulation.value.eat_service
    const residentService = this.formPostulation.value.resident_service

    const student: Student = {
      dni_student: dni,
      type_student: typeStudent,
      email_student: email,
      eat_service: eatService,
      resident_service: residentService
    }
    this.onPostulation.emit({postulation: responseStudent.data, student: student});
  }

  private validateStudentThesisOrPracticingPromise(dataValid: IValidationUser) : Promise<{error: boolean, msg: string, data: any, type: 'error' | 'success' | 'warning' | 'info'}> {
    this.isLoading = true;
    return new Promise<any>(resolve => {
      this._subscriptions.add(
        this._managerService.validateStudentThesisOrPracticingOrStudent(dataValid).subscribe({
          next: async (resp: any) => {
            this.isLoading = false;
            resolve({error: false, msg: resp.msg, data: resp.detalle, type: 'success',});
          },
          error: (err: HttpErrorResponse) => {
            debugger
            this.isLoading = false;
            resolve({ error: true, msg: err.message, data: err.error, type: 'error', code: err.status });
          }
        })
      );
    })
  }

  private validateStudentPermissionPromise(dataValid: IValidationUser) : Promise<{error: boolean, msg: string, data: any, type: 'error' | 'success' | 'warning' | 'info'}> {
    this.isLoading = true;
    return new Promise<any>(resolve => {
      this._subscriptions.add(
        this._managerService.validateStudentPermission(dataValid).subscribe({
          next: async (resp: any) => {
            this.isLoading = false;
            resolve({error: false, msg: resp.msg, data: resp.detalle, type: 'success',});
          },
          error: (err: HttpErrorResponse) => {
            this.isLoading = false;
            resolve({ error: true, msg: 'No se pudo validar el estudiante!', data: null, type: 'error', code: err.status });
          }
        })
      );
    })
  }

  private validateStudentSignAreasPromise(dni: string) : Promise<{error: boolean, msg: string, data: any, type: 'error' | 'success' | 'warning' | 'info'}> {
    this.isLoading = true;
    return new Promise<any>(resolve => {
      this._subscriptions.add(
        this._managerService.validateSignAreas(dni).subscribe({
          next: (resp: any) => {
            this.isLoading = false;
            resolve({error: false, msg: resp.msg, data: resp.detalle, type: 'success',});
          },
          error: (err: HttpErrorResponse) => {
            this.isLoading = false;
            resolve({ error: true, msg: 'No se pudo validar el estudiante!', data: null, type: 'error', code: err.status });
          }
        })
      );
    })
  }

  private getDataAnnouncementAndStudent(postulationId: number, dni: string) : Promise<{error: boolean, msg: string, data: IAnnouncement, type: 'error' | 'success' | 'warning' | 'info'}> {
    this.isLoading = true;
    return new Promise<any>(resolve => {
      this._subscriptions.add(
        this._managerService.getDataAnnouncementAndStudent(postulationId, dni).subscribe({
          next: async (resp: any) => {
            this.isLoading = false;
            if (!resp.detalle)  resolve({error: true, msg: resp.msg, data: resp.detalle, type: 'error',});

            resolve({error: false, msg: resp.msg, data: resp.detalle, type: 'success',});
          },
          error: (err: HttpErrorResponse) => {
            this.isLoading = false;
            resolve({ error: true, msg: 'No se pudo obtener los datos del alumno, intente nuevamente', data: null, type: 'error', code: err.status });
          }
        })
      );
    })
  }

  private evaluateDebts(debts: IDebtsStudent[]): boolean {
    this.debts = [];
    if (debts.length !== 0) {
      const debtsEat = debts.filter(
        (d) => d.concepto_deuda === 'COMEDOR (DEUDA OAEBU)'
        && d.monto_deuda === "0.00"
      );
      if (debtsEat && debtsEat.length) this.debts.push(...debtsEat);

      const debtsInt = debts.filter(
        (d) => d.concepto_deuda === 'INTERNADO UNAS (DEUDA)'
          && d.monto_deuda === "0.00"
      );
      if (debtsInt && debtsInt.length) this.debts.push(...debtsInt);

      if (this.debts.length) {
        this.createMessage('warning', 'El estudiante tiene deudas pendientes, no puede postular!');
        this.resetForm();

        this.showModalDebt = true;
        this.showModal.set(false);
        return true;
      }
    }
    return false;
  }

  private resetForm(): void {
    this.formPostulation.reset({
      type_student: 'student',
      eat_service: false,
      resident_service: false,
    })
  }

  private validateStudentDebtsPromise(dni: string) : Promise<{error: boolean, msg: string, data: IDebtsStudent[], type: 'error' | 'success' | 'warning' | 'info'}> {
    this.isLoading = true;
    return new Promise<any>(resolve => {
      this._subscriptions.add(
        this._managerService.validateStudentDebts(dni).subscribe({
          next: async (resp: any) => {
            this.isLoading = false;
            resolve({error: false, msg: '', data: resp, type: 'success',});
          },
          error: (err: HttpErrorResponse) => {
            this.isLoading = false;
            resolve({ error: true, msg: 'No se pudo validar las deudas del estudiante, intente nuevamente', data: null, type: 'error', code: err.status });
          }
        })
      );
    })
  }

  private createMessage(type: 'error' | 'success' | 'warning' | 'info', message: string): void {
    this._toastService.add({
      type: type,
      message: message
    });
  }

}
