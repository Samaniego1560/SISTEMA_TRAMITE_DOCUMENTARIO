import {Component, Input, input, model, OnChanges, SimpleChanges} from '@angular/core';
import {CdkAccordionItem} from "@angular/cdk/accordion";
import {DatePipe, JsonPipe, NgForOf, NgIf} from "@angular/common";
import {FormControl, ReactiveFormsModule, Validators} from "@angular/forms";
import {ModalComponent} from "../../../../../../../../core/ui/modal/modal.component";
import {IRequest, IUpdateService} from "../../../../../../../../core/models/requests";
import {HttpErrorResponse} from "@angular/common/http";
import {ToastService} from "../../../../../../../../core/services/toast/toast.service";
import {Subscription} from "rxjs";
import {ManagerService} from "../../../../../../../../core/services/manager/manager.service";
import {SafePipePipe} from "../../../../../../../../core/pipes/safe-pipe.pipe";
import {Department, District, Province} from "../../../../../../../../core/models/ubigeos";

@Component({
  selector: 'app-request-evaluation',
  standalone: true,
  imports: [
    CdkAccordionItem,
    DatePipe,
    NgIf,
    ModalComponent,
    ReactiveFormsModule,
    SafePipePipe,
  ],
  templateUrl: './request-evaluation.component.html',
  styleUrl: './request-evaluation.component.scss'
})
export class RequestEvaluationComponent implements OnChanges {
  student = input.required<IRequest>();
  view = model.required<string>();

  @Input({alias: 'departments', required: true}) departmentsMap: Record<string, string> = {};
  @Input({alias: 'provinces', required: true}) provincesMap: Record<string, string> = {};
  @Input({alias: 'districts', required: true}) districtsMap: Record<string, string> = {};

  public urlProfile: string = '';
  protected messageService: FormControl;

  protected action: string = '';

  private serviceID: number = 0;

  protected alertModal: boolean = false;

  protected role: number = 0;

  protected fileUrl: string = '';

  protected showModal: boolean = false;
  protected showModalService: boolean = false;

  protected isLoading: boolean = false;

  private _subscriptions: Subscription = new Subscription();

  constructor(private _toastService: ToastService,
              private _managerService: ManagerService,
  ) {
    this.messageService = new FormControl<string>('', Validators.required);
  }

  ngOnChanges(changes: SimpleChanges) {
    this.searchPhotoProfile()
  }

  protected updateStatusService(service: number, status: string): void {
    this.action = status;
    this.serviceID = service;
    this.alertModal = true;

    if (status !== 'rechazado') {
      this.messageService.clearValidators();
      this.messageService.updateValueAndValidity();
      this.messageService.setValue('-');
    } else {
      this.messageService.setValidators(Validators.required);
      this.messageService.updateValueAndValidity();
    }
  }

  protected backToList(): void {
    this.view.set('list')
  }

  protected viewFileAnnexe(url: string): void {
    this.fileUrl = url;
    this.showModal = true;
  }

  protected handleUpdateStatus(): void {

    if (this.messageService.invalid) {
      this._toastService.add({type: 'warning', message: 'Debe de ingresar el motivo del rechazo!'});
      this.messageService.markAllAsTouched();
      return;
    }

    const data: IUpdateService = {
      solicitud_id: this.student().id,
      servicios: [
        {
          servicio_id: this.serviceID,
          estado: this.action,
          detalle_rechazo: this.messageService.value.trim()
        }
      ]
    }

    this.isLoading = true;

    this._subscriptions.add(
      this._managerService.updateStatusService(data).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({type: 'info', message: res.msg});
            return;
          }
          this._toastService.add({type: 'info', message: "Estado actualizado correctamente"});

          const index = this.student().servicios_solicitados.findIndex(service => service.id === this.serviceID);
          if (index !== -1) this.student().servicios_solicitados[index].estado = this.action;

          this.alertModal = false;
        },
        error: (err: HttpErrorResponse) => {
          this.isLoading = false;
          this._toastService.add({
            type: 'error',
            message: 'No se pudo actualizar el estado del servicio, intente nuevamente!'
          });
        }
      })
    );
  }

  private searchPhotoProfile() {
    const application = this.student().detalle_solicitudes.find(detail =>
      detail.requisitos.some(req => req.nombre === 'photo-profile')
    );

    if (application) {
      const requirement = application.requisitos.find(req => req.nombre === 'photo-profile');
      if (requirement && requirement.respuesta) {
        this.urlProfile = requirement.respuesta.url_documento;
      }
    }
  }

  saveForm() {
    this.student().servicios_solicitados.push({
      estado: 'Pendiente',
      id: 1,
      servicio_id: 2,
      solicitud_id: 2,
      servicio: {
        id: 1,
        nombre: 'Internado',
        descripcion: 'Internado'
      }
    })
    this.showModalService = false;
  }
}
