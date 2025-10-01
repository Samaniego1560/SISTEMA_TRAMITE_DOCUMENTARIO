import {Component, effect, Input, OnDestroy, OnInit} from '@angular/core';
import {DatePipe, NgClass, NgFor, NgIf, NgStyle} from '@angular/common';
import { ModalComponent } from '../../../../core/ui/modal/modal.component';
import { CdkAccordionItem } from '@angular/cdk/accordion';
import {
  IAnnouncement,
} from '../../../../core/models/announcement';
import { ManagerService } from '../../../../core/services/manager/manager.service';
import { ToastService } from '../../../../core/services/toast/toast.service';
import { Subscription } from 'rxjs';
import { BlockUiComponent } from '../../../../core/ui/block-ui/block-ui.component';
import { HttpErrorResponse } from '@angular/common/http';
import {
  FormControl,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ToastComponent } from '../../../../core/ui/toast/toast.component';
import {
  Estudiante,
  IStatusRequest,
} from '../../../../core/models/requests';
import { IModalPostulation } from '../../../../core/models/modal';
import { FormsModule } from '@angular/forms';
import {SisfohComponent} from "./components/sisfoh/sisfoh.component";
import {SisComponent} from "./components/sis/sis.component";
import {AlertConditionComponent} from "../../../../core/components/alert-condition/alert-condition.component";
import {PhotoProfileComponent} from "./components/photo-profile/photo-profile.component";
import {PostulationFormComponent} from "../../../../core/components/postulation-form/postulation-form.component";
import {
  PostulationDynamicComponent
} from "../../../../core/components/postulation-dynamic/postulation-dynamic.component";
import {Student} from "../../../../core/models/student";
import {ApplicationNoticeViewComponent} from "./components/application-notice-view/application-notice-view.component";
@Component({
  selector: 'app-postulation',
  standalone: true,
  imports: [
    NgIf,
    NgFor,
    ModalComponent,
    CdkAccordionItem,
    BlockUiComponent,
    ReactiveFormsModule,
    ToastComponent,
    DatePipe,
    FormsModule,
    NgClass,
    SisfohComponent,
    SisComponent,
    NgStyle,
    AlertConditionComponent,
    PhotoProfileComponent,
    PostulationFormComponent,
    PostulationDynamicComponent,
    ApplicationNoticeViewComponent,
  ],
  templateUrl: './postulation.component.html',
  styleUrl: './postulation.component.scss',
  providers: [ManagerService, ToastService],
})
export class PostulationComponent implements OnInit, OnDestroy {
  private _subscriptions: Subscription = new Subscription();
  protected view: string = 'list';
  protected announcement: IAnnouncement = {
    nombre: '',
    convocatoria_servicio: [],
    fecha_fin: '',
    fecha_inicio: '',
    secciones: [],
    activo: false,
    credito_minimo: 0,
    nota_minima: 0
  };
  protected postulation: IAnnouncement = {
    nombre: '',
    convocatoria_servicio: [],
    fecha_fin: '',
    fecha_inicio: '',
    secciones: [],
    activo: false,
    nota_minima: 0,
    credito_minimo: 0
  };
  protected student: Student = {} as Student;

  protected isLoading: boolean = false;
  protected modal: IModalPostulation = {
    debts: false,
    errors: false,
    manual: false,
    video: false,
    status: false,
    delete: false,
    form: false,
  };


  protected dniForm: FormControl;
  protected statusRequest: IStatusRequest[] = [];
  protected title: string = 'LISTA DE REQUISITOS PARA LA CONVOCATORIA';


  constructor(
    private _managerService: ManagerService,
    private _toastService: ToastService,
  ) {
    this.dniForm = new FormControl<string>('', Validators.required);
  }

  ngOnInit() {
      this.getCurrentAnnouncement();
  }

  ngOnDestroy() {
    this._subscriptions.unsubscribe();
    this.resetValidateStatus();
    //this.formPostulation.reset();
    this.modal = {
      debts: false,
      errors: false,
      manual: false,
      video: false,
      status: false,
      delete: false,
      form: false,
    };
  }

  protected init(): void {
    this.view = 'list';
    this.postulation = {
      nombre: '',
      convocatoria_servicio: [],
      fecha_fin: '',
      fecha_inicio: '',
      secciones: [],
      activo: false,
      credito_minimo: 0,
      nota_minima: 0
    };

    this.modal.form = false;
    this.modal.delete = false;
    this.modal.debts = false;
  }

  private getCurrentAnnouncement(): void {
    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.getCurrentAnnouncement().subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this.title = 'Convocatoria no disponible';
            this._toastService.add({ type: 'error', message: res.msg });
            return;
          }
          this.announcement = res.detalle;
        },
        error: (err: HttpErrorResponse) => {
          this.isLoading = false;
          this._toastService.add({
            type: 'error',
            message: 'No se pudo obtener la convocatoria actual',
          });
          console.error(err);
        },
      })
    );
  }

  protected validateStatusRequest(): void {
    if (this.dniForm.invalid) {
      this._toastService.add({
        type: 'warning',
        message: 'Debe de ingresar un DNI válido!',
      });
      this.dniForm.markAllAsTouched();
      return;
    }

    this.isLoading = true;

    this._subscriptions.add(
      this._managerService.getRequestStatus(this.dniForm.value).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({ type: 'info', message: res.msg });
            return;
          }

          console.log(res.detalle);
          const studet: Estudiante = res.detalle;
          if (studet.convocatorias_participadas.length > 0) {
            const ultimoObjeto = studet.convocatorias_participadas[studet.convocatorias_participadas.length - 1];
            this.statusRequest = ultimoObjeto.solicitudes;
        }

        },
        error: (err: HttpErrorResponse) => {
          console.error(err);
          this.isLoading = false;
          this._toastService.add({
            type: 'error',
            message:
              'No se pudo obtener la solicitud del estudiante, intente de nuevo',
          });
        },
      })
    );
  }

  protected resetValidateStatus(): void {
    this.modal.status = false;
    this.statusRequest = [];
    this.dniForm.reset();
  }

  onPostulation(event: { postulation: IAnnouncement, student: Student }) {
    this.postulation = event.postulation;
    this.student = event.student;
    this.modal.form = false;
    this.view = 'form';
  }
}
