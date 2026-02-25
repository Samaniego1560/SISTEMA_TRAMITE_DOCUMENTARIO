import { Component, Input, input, model, OnChanges, SimpleChanges } from '@angular/core';
import { CdkAccordionItem } from "@angular/cdk/accordion";
import { DatePipe, JsonPipe, NgForOf, NgIf } from "@angular/common";
import { FormControl, ReactiveFormsModule, Validators } from "@angular/forms";
import { ModalComponent } from "../../../../../../../../core/ui/modal/modal.component";
import { IRequest, IUpdateService } from "../../../../../../../../core/models/requests";
import { HttpErrorResponse } from "@angular/common/http";
import { ToastService } from "../../../../../../../../core/services/toast/toast.service";
import { Subscription } from "rxjs";
import { ManagerService } from "../../../../../../../../core/services/manager/manager.service";
import { SafePipePipe } from "../../../../../../../../core/pipes/safe-pipe.pipe";
import { Department, District, Province } from "../../../../../../../../core/models/ubigeos";

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

  @Input({ alias: 'departments', required: true }) departmentsMap: Record<string, string> = {};
  @Input({ alias: 'provinces', required: true }) provincesMap: Record<string, string> = {};
  @Input({ alias: 'districts', required: true }) districtsMap: Record<string, string> = {};

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

  // Propiedades para reemplazo de archivos
  protected showReplaceModal: boolean = false;
  protected selectedFile: File | null = null;
  protected selectedRequirement: any = null;
  protected isUploadingFile: boolean = false;

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
      this._toastService.add({ type: 'warning', message: 'Debe de ingresar el motivo del rechazo!' });
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
            this._toastService.add({ type: 'info', message: res.msg });
            return;
          }
          this._toastService.add({ type: 'info', message: "Estado actualizado correctamente" });

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
        // Manejar tanto arrays como objetos individuales
        const respuesta = Array.isArray(requirement.respuesta)
          ? requirement.respuesta[0]
          : requirement.respuesta;

        if (respuesta && respuesta.url_documento) {
          this.urlProfile = respuesta.url_documento;
        }
      }
    }
  }

  /**
   * Verifica si una sección es de composición familiar
   */
  protected isFamilyCompositionSection(section: any): boolean {
    const normalizedName = this.normalizeText(section?.descripcion || '');
    if (normalizedName.includes('COMPOSICION FAMILIAR')) {
      return true;
    }

    const requirementNames = (section?.requisitos || [])
      .map((req: any) => this.normalizeText(req?.nombre || ''));

    const hasFamilyFields =
      requirementNames.some((name: string) => name.includes('NOMBRES') && name.includes('APELLIDOS')) &&
      requirementNames.some((name: string) => name === 'EDAD') &&
      requirementNames.some((name: string) => name.includes('PARENTESCO'));

    return hasFamilyFields;
  }

  /**
   * Agrupa los requisitos de composición familiar por el campo 'order'
   * Cada grupo representa un integrante familiar con sus 5 campos
   */
  protected groupFamilyMembers(requisitos: any[]): any[] {
    const grouped = new Map<number, any>();

    requisitos.forEach(req => {
      // Manejar tanto arrays de respuestas como objetos individuales
      const respuestas = Array.isArray(req.respuesta) ? req.respuesta : (req.respuesta ? [req.respuesta] : []);
      const key = this.getFamilyFieldKey(req.nombre);
      if (!key) return;

      respuestas.forEach((respuesta: any) => {
        const order = respuesta.order || 0;

        // Crear un nuevo objeto para este order si no existe
        if (!grouped.has(order)) {
          grouped.set(order, {
            order: order,
            family_name_lastname: '-',
            family_age: '-',
            family_marital_status: '-',
            family_educational_level: '-',
            family_occupation: '-'
          });
        }

        // Obtener el objeto del integrante
        const member = grouped.get(order)!;

        // Asignar el valor según el nombre del requisito
        const value = respuesta.respuesta_formulario ||
          respuesta.opcion_seleccion ||
          '-';

        member[key] = value;
      });
    });

    // Convertir el Map a array y filtrar registros vacíos
    return Array.from(grouped.values()).filter(member => {
      // Filtrar si todos los campos están vacíos o son '-'
      return member.family_name_lastname !== '-' ||
        member.family_age !== '-' ||
        member.family_marital_status !== '-' ||
        member.family_educational_level !== '-' ||
        member.family_occupation !== '-';
    });
  }

  private getFamilyFieldKey(requirementName: string): string | null {
    const normalized = this.normalizeText(requirementName);

    if (normalized.includes('NOMBRES') && normalized.includes('APELLIDOS')) {
      return 'family_name_lastname';
    }
    if (normalized === 'EDAD') {
      return 'family_age';
    }
    if (normalized.includes('ESTADO CIVIL')) {
      return 'family_marital_status';
    }
    if (normalized.includes('NIVEL EDUCATIVO')) {
      return 'family_educational_level';
    }
    if (normalized.includes('OCUPACION')) {
      return 'family_occupation';
    }

    return null;
  }

  private normalizeText(value: string): string {
    return (value || '')
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '')
      .replace(/\s+/g, ' ')
      .toUpperCase()
      .trim();
  }

  /**
   * Abre el modal para reemplazar un archivo
   */
  protected openReplaceFileModal(requirement: any): void {
    // Manejar tanto arrays como objetos individuales
    const respuesta = Array.isArray(requirement.respuesta)
      ? requirement.respuesta[0]
      : requirement.respuesta;

    this.selectedRequirement = { ...requirement, respuesta };
    this.showReplaceModal = true;
    this.selectedFile = null;
  }

  /**
   * Maneja la selección de archivo
   */
  protected onFileSelected(event: any): void {
    const file = event.target.files[0];
    if (file) {
      this.selectedFile = file;
    }
  }

  /**
   * Reemplaza el archivo
   */
  protected replaceFile(): void {
    if (!this.selectedFile || !this.selectedRequirement) {
      this._toastService.add({
        type: 'warning',
        message: 'Debe seleccionar un archivo'
      });
      return;
    }

    this.isUploadingFile = true;

    // Convertir archivo a Base64
    const reader = new FileReader();
    reader.onload = () => {
      const base64 = (reader.result as string).split(',')[1];

      const data = {
        detalle_solicitud_id: this.selectedRequirement.respuesta.id,
        id_convocatoria: this.student().convocatoria_id,
        dni_alumno: this.student().alumno.DNI,
        name_file: this.selectedFile!.name,
        file: base64
      };

      this._subscriptions.add(
        this._managerService.replaceRequestFile(data).subscribe({
          next: (res: any) => {
            this.isUploadingFile = false;
            this._toastService.add({
              type: 'success',
              message: 'Archivo reemplazado correctamente'
            });

            // Actualizar la URL en la interfaz
            this.selectedRequirement.respuesta.url_documento = res.detalle.url_file;

            this.showReplaceModal = false;
            this.selectedFile = null;
            this.selectedRequirement = null;
          },
          error: (err: HttpErrorResponse) => {
            this.isUploadingFile = false;
            this._toastService.add({
              type: 'error',
              message: 'Error al reemplazar el archivo'
            });
          }
        })
      );
    };

    reader.readAsDataURL(this.selectedFile);
  }

  /**
   * Cierra el modal de reemplazo
   */
  protected closeReplaceModal(): void {
    this.showReplaceModal = false;
    this.selectedFile = null;
    this.selectedRequirement = null;
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
