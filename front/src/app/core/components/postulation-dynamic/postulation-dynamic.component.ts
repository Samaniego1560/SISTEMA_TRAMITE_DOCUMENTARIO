import {Component, effect, EventEmitter, input, OnInit, Output} from '@angular/core';
import {ModalComponent} from "../../ui/modal/modal.component";
import {CdkAccordionItem} from "@angular/cdk/accordion";
import {NgIf} from "@angular/common";
import {
  PhotoProfileComponent
} from "../../../modules/home/pages/postulation/components/photo-profile/photo-profile.component";
import {FormArray, FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
import {IModalPostulation} from "../../models/modal";
import {IAnnouncement, IRequirement, ISection} from "../../models/announcement";
import {Department, District, Province} from "../../models/ubigeos";
import {IBodyRequest, IDetailRequest, IErrorPostulation, IFileRequest} from "../../models/requests";
import {HttpErrorResponse} from "@angular/common/http";
import {FILE_MAX_SIZE, ValidateFileType} from "../../utils/statics/statics";
import {ToastService} from "../../services/toast/toast.service";
import {Subscription} from "rxjs";
import {ManagerService} from "../../services/manager/manager.service";
import {UbigeoService} from "../../services/ubigeo/ubigeo.service";
import {Student} from "../../models/student";
import {SectionFormComponent} from "./components/section-form/section-form.component";
import {TableSectionComponent} from "./components/table-section/table-section.component";
import {BlockUiComponent} from "../../ui/block-ui/block-ui.component";
import {ToastComponent} from "../../ui/toast/toast.component";
@Component({
  selector: 'app-postulation-dynamic',
  standalone: true,
  imports: [
    CdkAccordionItem,
    PhotoProfileComponent,
    SectionFormComponent,
    TableSectionComponent,
    BlockUiComponent,
    ToastComponent,

  ],
  templateUrl: './postulation-dynamic.component.html',
  styleUrl: './postulation-dynamic.component.scss'
})
export class PostulationDynamicComponent implements OnInit{
  roleId = input.required<number>();
  student = input.required<Student>();
  postulation = input.required<IAnnouncement>();

  @Output('back-event') backEvent: EventEmitter<void> = new EventEmitter<void>();

  public recordRequirements: Record<string, string> = {};

  protected formPostulation: FormGroup;

  public departments: Department[] = [];
  public provinces: Province[] = [];
  public districts: District[] = [];

  protected isLoading: boolean = false;

  protected view: string = 'list';

  private _subscriptions: Subscription = new Subscription();

  protected modal: IModalPostulation = {
    debts: false,
    errors: false,
    manual: false,
    video: false,
    status: false,
    delete: false,
    form: false,
  };

  protected requirementPhotoProfile: IRequirement | null = null;

  constructor(private _fb: FormBuilder,
              private _toastService: ToastService,
              private _managerService: ManagerService,
              ) {
    this.formPostulation = this._fb.group({
    });
  }

  ngOnInit() {
    this.loadForm();
    this.loadDataSource();
  }

  private loadForm(): void {
      this.formPostulation.addControl(
        'dni_student',
        new FormControl(this.student().dni_student))

      this.formPostulation.addControl(
        'type_student',
        new FormControl(this.student().type_student))

      this.formPostulation.addControl(
        'email_student',
        new FormControl(this.student().email_student))

    this.formPostulation.addControl(
      'eat_service',
      new FormControl(this.student().eat_service))

    this.formPostulation.addControl(
      'resident_service',
      new FormControl(this.student().resident_service))

    this.processForm(this.postulation().secciones);
  }

  private loadDataSource(): void {
    const departments = localStorage.getItem('departments');
    const provinces = localStorage.getItem('provinces');
    const districts = localStorage.getItem('districts');

    this.departments = departments ? JSON.parse(departments) : [];
    this.provinces = provinces ? JSON.parse(provinces) : [];
    this.districts = districts ? JSON.parse(districts) : [];
  }

  protected handlePostulation(): void {
    const body: IBodyRequest = {
      convocatoria_id: this.postulation()?.id || 0,
      alumno_id: this.postulation().user_id || 0,
      servicios_solicitados: [],
      detalle_solicitudes: [],
    };

    const keyPhotoProfile = this.recordRequirements['photo-profile'];

    if (keyPhotoProfile) {
      if (this.formPostulation.get(keyPhotoProfile)?.invalid) {
        this._toastService.add({
          type: 'error',
          message: 'Sube una foto de perfil valida!',
        });
        return;
      }
    }

    if (this.formPostulation.invalid) {
      this._toastService.add({
        type: 'error',
        message: 'Complete todos los campos correctamente!',
      });
      this.formPostulation.markAllAsTouched();
      return;
    }

    if (this.formPostulation.value.eat_service) {
      body.servicios_solicitados.push({
        estado: 'pendiente',
        servicio_id: 1,
      });
    }

    if (this.formPostulation.value.resident_service) {
      body.servicios_solicitados.push({
        estado: 'pendiente',
        servicio_id: 2,
      });
    }

    this.isLoading = true;
    const sections = this.postulation().secciones;
    for (const key in this.formPostulation.controls) {
      if (!this.formPostulation.controls.hasOwnProperty(key)) continue;
      const control = this.formPostulation.controls[key];
      if (control instanceof FormArray) {
        const details = this.processFormArray(control, key);
        body.detalle_solicitudes.push(...details);
        continue;
      }

      let req = this.getRequirementByKey(key, sections);
      if (!req) continue;
      const value = control.value.toString();
      body.detalle_solicitudes.push(this.getDetailRequestByRequirementAndControl(key, req, value, 1));
    }

    this._subscriptions.add(
      this._managerService.createRequest(body).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({ type: 'error', message: res.msg });
            return
          }

          this._toastService.add({
            type: 'success',
            message: 'Postulación realizada correctamente',
          });
          this.init();
          this.backEvent.emit()
        },
        error: (err: HttpErrorResponse) => {
          this.isLoading = false;
          this._toastService.add({
            type: 'error',
            message: 'No se pudo realizar la postulación, intente nuevamente',
          });
        },
      })
    );
  }

  protected init(): void {
    this.view = 'list';

    this.modal.form = false;
    this.modal.delete = false;
    this.modal.debts = false;
  }

  private processFormArray(formArray: FormArray, sectionId: string): IDetailRequest[] {
    const details: IDetailRequest[] = [];
    const section = this.postulation().secciones.find((item: any) => item.id?.toString() === sectionId)
    if (!section) return details;
    if (formArray.value.length) {
      let cont = 1;
      for (const item of formArray.value) {
        for (const key in item) {
          if (item.hasOwnProperty(key)) {
            const value = item[key];
            const req = section.requisitos.find(req => req.id?.toString() === key);
            if (!req) continue;
            details.push(this.getDetailRequestByRequirementAndControl(key, req, value, cont));
          }
        }
        cont++;
      }
    }
    return details;
  }

  private getDetailRequestByRequirementAndControl(key: string, req: IRequirement, value: string, order: number): IDetailRequest {
    return {
      respuesta_formulario:
        req.tipo_requisito_id === 3 || req.tipo_requisito_id === 5 || req.tipo_requisito_id === 6 || req.tipo_requisito_id === 7 ? value : null,
      url_documento: [1, 2, 8].includes(req.tipo_requisito_id)
        ? value
        : null,
      opcion_seleccion:
        req.tipo_requisito_id === 4 ? value : null,
      requisito_id: parseInt(key),
      order: order
    }
  }

  private getRequirementByKey(key: string, sections: ISection[]): IRequirement | undefined {
    for (const section of sections) {
      const req = section.requisitos.find((r) => r.id === parseInt(key));
      if (req) return req;
    }
    return undefined;
  }

  protected onProcessFile(e: any): void {
    const {event, key, type} = e;
    const file: File = event.target.files[0];

    if (!ValidateFileType(file, type)) {
      this._toastService.add({
        type: 'error',
        message: 'El archivo debe ser PDF o Imagen',
      });
      return;
    }

    if (file.size > FILE_MAX_SIZE) {
      this._toastService.add({
        type: 'error',
        message: 'El archivo no debe ser mayor a 10MB',
      });
      return;
    }

    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => {
      const body: IFileRequest = {
        id_convocatoria: this.postulation().id || 0,
        dni_alumno: parseInt(this.formPostulation.value.dni_student),
        name_file: file.name,
        file: (reader.result as string).split(',')[1] || '',
      };
      this.handleLoadFile(body, key);
    };
  }

  private handleLoadFile(body: IFileRequest, key: string): void {
    this.isLoading = true;

    this._subscriptions.add(
      this._managerService.uploadRequestFile(body).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({ type: 'error', message: res.msg });
            return;
          }

          this._toastService.add({ type: 'success', message: res.msg });
          this.formPostulation.get(key)?.setValue(res.detalle.url_file);
        },
        error: (err: HttpErrorResponse) => {
          console.error(err);
          this.isLoading = false;
          this._toastService.add({
            type: 'error',
            message: 'No se pudo subir el archivo, intente nuevamente',
          });
        },
      })
    );
  }

  private processForm(sections: ISection[]): void {
    for (const section of sections) {
      const keySection = (section.id || 0).toString();

      if (section.type === 'table') {
        this.handleTableSection(keySection, section);
        continue;
      }

      for (const req of section.requisitos) {
        this.handleRequirement(req);
      }
    }
  }

  private handleTableSection(keySection: string, section: ISection): void {
    this.formPostulation.addControl(keySection, this._fb.array([]));
    this.processRecoverableDetailRequestsTypeTable(keySection, section.requisitos);
  }

  private handleRequirement(req: IRequirement): void {
    const key = (req.id || 0).toString();
    this.recordRequirements[req.nombre] = key;

    if (req.nombre === 'photo-profile') {
      this.requirementPhotoProfile = req;
    }

    const defaultValue = this.getDefaultValue(req);
    const control = new FormControl(defaultValue, this.getValidators(req));

    this.formPostulation.addControl(key, control);

    if (defaultValue) control.disable();
    if (req.is_recoverable) this.addValueRecovery(key);
  }

  private getDefaultValue(req: IRequirement): string {
    switch (req.nombre) {
      case 'Correo institutcional':
        return this.student().email_student;
      case 'Tipo de estudiante':
        return this.student().type_student;
      default:
        return req.default || '';
    }
  }

  private getValidators(req: IRequirement): Validators | null {
    return req.is_dependent && req.show_dependent ? null : Validators.required;
  }

  private addValueRecovery(key: string) {
    const value =  this.getValueRecoveryById(key)

    this.formPostulation.get(key)?.setValue(value);
  }

  private getValueRecoveryById(key: string): string {
    return this.postulation().recoverable_detail_requests?.find(
      (item: any) => item.requisito_id.toString() === key
    )?.respuesta_formulario || '';
  }

  private processRecoverableDetailRequestsTypeTable(sectionId: string, requirements: IRequirement[]): void {
    const requirementsRecoverable = requirements.filter(req => req.is_recoverable);
    const values: { [key: string]: string }[] = [];
    for (const requirement of requirementsRecoverable) {
      const recoverableDetailRequests = this.postulation().recoverable_detail_requests?.filter(
        (item: any) => item.requisito_id === requirement.id
      );
      if (recoverableDetailRequests?.length) {
        for (const recoverableDetailRequest of recoverableDetailRequests) {
          const index = recoverableDetailRequest.order - 1;
          if (!values[index]) {
            values[index] = {};
          }
          const key = (recoverableDetailRequest.requisito_id || 0).toString();
          values[index][key] = recoverableDetailRequest.respuesta_formulario;
        }
      }
    }

    const formArray = this.formPostulation.get(sectionId) as FormArray;
    values.forEach(value => {
      formArray.push(this._fb.group(value));
    });
  }
}
