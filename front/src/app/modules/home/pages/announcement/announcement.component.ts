import {Component, OnDestroy} from '@angular/core';
import {ModalComponent} from "../../../../core/ui/modal/modal.component";
import {DatePipe, JsonPipe, NgIf} from "@angular/common";
import {CdkAccordionItem} from "@angular/cdk/accordion";
import {IAnnouncement, IRequirement, ISection} from "../../../../core/models/announcement";
import {ManagerService} from "../../../../core/services/manager/manager.service";
import {Subscription} from "rxjs";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {EnvServiceFactory} from "../../../../core/services/env/env.service.provider";
import {FormRequirementComponent} from "./components/form-requirement/form-requirement.component";
import {FormSectionComponent} from "./components/form-section/form-section.component";

@Component({
  selector: 'app-announcement',
  standalone: true,
  imports: [
    ModalComponent,
    NgIf,
    CdkAccordionItem,
    BlockUiComponent,
    DatePipe,
    ReactiveFormsModule,
    ToastComponent,
    FormsModule,
    JsonPipe,
    FormRequirementComponent,
    FormSectionComponent
  ],
  templateUrl: './announcement.component.html',
  styleUrl: './announcement.component.scss',
  providers: [ManagerService, ToastService]
})
export class AnnouncementComponent implements OnDestroy {
  private _subscriptions: Subscription = new Subscription();
  protected typeProcess: string = 'create';
  private reqIndex: number = -1;
  private secIndex: number = -1;
  protected openModal: boolean = false;
  protected openModalSection: boolean = false;
  protected currentSection: number = 0;
  protected view: string = 'list';
  protected announcements: IAnnouncement[] = [];

  public requirementNames: string[] = [];
  protected announcement: IAnnouncement = {
    nombre: '',
    convocatoria_servicio: [],
    fecha_fin: '',
    fecha_inicio: '',
    secciones: [],
    activo: false,
    nota_minima: 0,
    credito_minimo: 0
  };
  protected isLoad: boolean = false;
  protected formAnnouncement: FormGroup;
  protected requirementSelect: IRequirement = {} as IRequirement;
  public sectionSelect: ISection = {} as ISection;

  constructor(
    private _managerService: ManagerService,
    private _toastService: ToastService,
    private _fb: FormBuilder
  ) {
    this._getAnnouncement();

    this.formAnnouncement = this._fb.group({
      name: ['', Validators.required],
      startDate: ['', Validators.required],
      endDate: ['', Validators.required],
      minimum_grade: [0, [Validators.required, Validators.min(0)]],
      minimum_credits: [0, [Validators.required, Validators.min(0)]],
      eat_service: [0, [Validators.required, Validators.min(0)]],
      inter_service: [0, [Validators.required, Validators.min(0)]],
      dentistry: [false, [Validators.required]],
      psychology: [false, [Validators.required]],
      medicine: [false, [Validators.required]],
    });
    this.init();
  }


  ngOnDestroy(): void {
    this._subscriptions.unsubscribe();
    this.init();
  }

  private init(): void {
    this.announcement.nombre = '';
    this.announcement.fecha_fin = '';
    this.announcement.fecha_inicio = '';
    this.announcement.secciones = [
      {
        descripcion: 'Datos Personales',
        type: "form",
        requisitos: EnvServiceFactory().SECTIONS_REQUIREMENTS_ONE
      },
      {
        descripcion: 'Lugar de Nacimiento',
        type: "form",
        requisitos: EnvServiceFactory().SECTIONS_REQUIREMENTS_TWO
      },
      {
        descripcion: 'Lugar de Procedencia',
        type: "form",
        requisitos: EnvServiceFactory().SECTIONS_REQUIREMENTS_THREE
      },
      {
        descripcion: 'Documentos Requeridos',
        type: "form",
        requisitos: EnvServiceFactory().SECTIONS_REQUIREMENTS_FOURTH
      },
      {
        descripcion: 'Datos Socioeconómica',
        type: "form",
        requisitos: EnvServiceFactory().SECTIONS_REQUIREMENTS_FIVE
      },
      {
        descripcion: 'Salud',
        type: "form",
        requisitos: EnvServiceFactory().SECTIONS_REQUIREMENTS_SIX
      },
      {
        descripcion: 'Datos de CERTIJOVEN',
        type: "form",
        requisitos: EnvServiceFactory().SECTIONS_REQUIREMENTS_SEVEN
      },
      {
        descripcion: 'Datos Compocisión familiar',
        type: "table",
        requisitos: EnvServiceFactory().SECTIONS_FAMILY_COMPOSITION
      },
      {
        descripcion: 'Vivienda',
        type: "form",
        requisitos: EnvServiceFactory().SECTIONS_DWELLING
      }
    ];

    this.requirementNames = this.announcement.secciones.flatMap(section =>
      section.requisitos.map(req => req.nombre)
    );

  }

  private _getAnnouncement() {
    this.isLoad = true;
    this._subscriptions.add(
      this._managerService.getAnnouncement().subscribe({
        next: (res: any) => {
          this.isLoad = false;
          if (!res.detalle) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo obtener el listado de convocatorias, intente nuevamente'
            });
            return;
          }

          this.announcements = res.detalle;

          for (const announcement of this.announcements) {
            announcement.activo = announcement.fecha_fin > new Date().toISOString();
          }

          this.announcements.sort((a, b) => {
              return new Date(b.fecha_fin).getTime() - new Date(a.fecha_fin).getTime();
            }
          );
        },
        error: (error: any) => {
          this._toastService.add({type: 'error', message: 'Error al obtener la convocatoria'});
          this.isLoad = false;
        }
      })
    );
  }

  protected createAnnouncement(): void {
    this.view = 'create';
    this.formAnnouncement.enable();
    this.formAnnouncement.reset();
    this.init();
  }

  protected editAnnouncement(item: IAnnouncement): void {
    this.loadAnnouncement(item);
    this.view = 'edit-readonly';
  }

  protected viewAnnouncement(item: IAnnouncement): void {
    this.loadAnnouncement(item);
    this.formAnnouncement.disable();
    this.view = 'view-readonly';
  }

  protected loadAnnouncement(item: IAnnouncement) {
    this.announcement = JSON.parse(JSON.stringify(item));
    const eatServices = this.announcement.convocatoria_servicio.find((item: any) => item.servicio?.id === 1)?.cantidad;
    const interServices = this.announcement.convocatoria_servicio.find((item: any) => item.servicio?.id === 2)?.cantidad;

    this.formAnnouncement.patchValue({
      name: this.announcement.nombre,
      startDate: this.announcement.fecha_inicio,
      endDate: this.announcement.fecha_fin,
      eat_service: eatServices,
      inter_service: interServices,
      minimum_grade: this.announcement.nota_minima,
      minimum_credits: this.announcement.credito_minimo
    });

    const eatService = this.announcement.convocatoria_servicio.find((item: any) => item.id === 2);
    if (eatService) this.formAnnouncement.get('eat_service')?.setValue(eatService.cantidad);

    const interService = this.announcement.convocatoria_servicio.find((item: any) => item.id === 1);
    if (interService) this.formAnnouncement.get('inter_service')?.setValue(interService.cantidad);

  }

  protected createRequirement(index: number): void {
    this.openModal = true;
    this.currentSection = index;
    this.typeProcess = 'create';
  }

  protected editRequirement(req: IRequirement, index: number, section: number): void {
    this.typeProcess = 'edit';
    this.reqIndex = index;
    this.currentSection = section;
    this.requirementSelect = req;
    this.openModal = true;
  }

  protected viewRequirement(req: IRequirement): void {
    this.typeProcess = 'view';
    this.requirementSelect = req;
    this.openModal = true;
  }

  protected submitAnnouncement(): void {
    if (this.formAnnouncement.invalid) {
      this._toastService.add({type: 'error', message: 'Complete los campos correctamente'});
      this.formAnnouncement.markAllAsTouched();
      return;
    }

    this.announcement.nombre = this.formAnnouncement.value.name;
    this.announcement.fecha_inicio = this.formAnnouncement.value.startDate.replace('T', ' ') + ':00';
    this.announcement.fecha_fin = this.formAnnouncement.value.endDate.replace('T', ' ') + ':00';
    this.announcement.activo = true;
    this.announcement.credito_minimo = this.formAnnouncement.value.minimum_credits;
    this.announcement.nota_minima = this.formAnnouncement.value.minimum_grade;

    this.announcement.convocatoria_servicio = [
      {
        cantidad: parseInt(this.formAnnouncement.value.eat_service),
        servicio_id: 2
      },
      {
        cantidad: parseInt(this.formAnnouncement.value.inter_service),
        servicio_id: 1
      }
    ];

    this.isLoad = true;
    this._subscriptions.add(
      this._managerService.createAnnouncement(this.announcement).subscribe({
        next: (res: any) => {
          this.isLoad = false;
          if (!res.detalle) {
            this._toastService.add({type: 'error', message: 'No se pudo crear la convocatoria, intente nuevamente'});
            return;
          }

          this._toastService.add({type: 'success', message: 'Convocatoria creada correctamente'});

          this.announcements.push(this.announcement);
          this.formAnnouncement.reset();
          this.init();
          this._getAnnouncement();
          this.view = 'list';
        },
        error: (error: any) => {
          this._toastService.add({type: 'error', message: 'Error al crear la convocatoria'});
          this.isLoad = false;
        }
      })
    );
  }

  addSection() {
    this.openModalSection = true;
    this.typeProcess = 'create';
  }

  editSection(section: ISection, i: number) {
    this.openModalSection = true;
    this.secIndex = i;
    this.sectionSelect = section;
    this.typeProcess = 'edit';
  }

  onSaveRequirement(requirement: IRequirement) {
    if (this.typeProcess === 'create')
      this.announcement.secciones[this.currentSection].requisitos.push(requirement);

    this.announcement.secciones[this.currentSection].requisitos[this.reqIndex] = requirement;
    return;
  }

  onSaveSection(section: ISection) {
    if (this.typeProcess === 'create')
      this.announcement.secciones.push(section);

    this.announcement.secciones[this.secIndex] = section;
  }

  deleteSection(i: number) {
    this.announcement.secciones.splice(i, 1);
  }

  deleteRequirement(index: number, section: number) {
    this.announcement.secciones[section].requisitos.splice(index, 1);
  }
}
