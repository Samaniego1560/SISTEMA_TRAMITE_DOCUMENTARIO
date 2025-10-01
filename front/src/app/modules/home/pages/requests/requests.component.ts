import {Component, OnDestroy, OnInit} from '@angular/core';
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {ModalComponent} from "../../../../core/ui/modal/modal.component";
import {FormBuilder, FormControl, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {ManagerService} from "../../../../core/services/manager/manager.service";
import {Subscription} from "rxjs";
import {DatePipe, KeyValuePipe, NgFor, NgIf, NgClass} from "@angular/common";
import {CdkAccordionItem} from "@angular/cdk/accordion";
import {SafePipePipe} from "../../../../core/pipes/safe-pipe.pipe";
import {AuthService} from "../../../../core/services/auth/auth.service";
import {FilterService} from "../../../../core/services/filter/filter.service";
import {IAnnouncement} from "../../../../core/models/announcement";
import { PostulationComponent } from '../postulation/postulation.component';
import {SearchStudentComponent} from "./components/search-student/search-student.component";
import {PostulationFormComponent} from "../../../../core/components/postulation-form/postulation-form.component";
import {RequestsReceivedComponent} from "./components/requests-received/requests-received.component";
import {Department, District, Province} from "../../../../core/models/ubigeos";
import {
  RequestEvaluationComponent
} from "./components/requests-received/components/request-evaluation/request-evaluation.component";

@Component({
  selector: 'app-requests',
  standalone: true,
  imports: [
    BlockUiComponent,
    ReactiveFormsModule,
    ToastComponent,
    DatePipe,
    SearchStudentComponent,
    RequestsReceivedComponent,
  ],
  templateUrl: './requests.component.html',
  styleUrl: './requests.component.scss',
  providers: [ToastService, ManagerService, AuthService, FilterService]
})
export class RequestsComponent implements OnInit, OnDestroy {

  private _subscriptions: Subscription = new Subscription();


  protected isLoading: boolean = false;


  protected view: string = 'list';

  protected showModal: boolean = false;
  protected role: number = 0;
  protected messageService: FormControl;
  protected announcements: IAnnouncement[] = [];
  protected announcementSelect: IAnnouncement = {} as IAnnouncement;


  protected formPostulation: FormGroup;
  protected dniForm: FormControl;
  protected postulation: IAnnouncement = {} as IAnnouncement;

  public departmentsMap: Record<string, string> = {};
  public provincesMap: Record<string, string> = {};
  public districtsMap: Record<string, string> = {};

  constructor(
    private _toastService: ToastService,
    private _managerService: ManagerService,
    private _authService: AuthService,
    private _fb: FormBuilder,
  ) {
    this._getAnnouncement();
    this.role = this._authService.getRole();
    this.messageService = new FormControl<string>('', Validators.required);

    this.formPostulation = this._fb.group({
      type_student: ['', Validators.required],
      dni_student: ['', Validators.required],
      email_student: ['', [Validators.required, Validators.email]],
      eat_service: [false, Validators.required],
      resident_service: [false, Validators.required],
    });
    this.dniForm = new FormControl<string>('', Validators.required);
  }

  ngOnInit() {
    this.loadDataSource()
  }

  ngOnDestroy() {
    this._subscriptions.unsubscribe();
  }

  private loadDataSource() {
    const storageData = {
      departments: localStorage.getItem('departments'),
      provinces: localStorage.getItem('provinces'),
      districts: localStorage.getItem('districts'),
    };

    const departments: Department[] = JSON.parse(storageData.departments || '[]');
    const provinces: Province[] = JSON.parse(storageData.provinces || '[]');
    const districts: District[] = JSON.parse(storageData.districts || '[]');

    departments.forEach(d => this.departmentsMap[d.id] = d.name);
    provinces.forEach(p => this.provincesMap[p.id] = p.name);
    districts.forEach(d => this.districtsMap[d.id] = d.name);
  }

  private _getAnnouncement() {
    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.getAnnouncement().subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res.detalle) {
            this._toastService.add({
              type: 'error',
              message: 'No se pudo obtener el listado de convocatorias, intente nuevamente'
            });
            return;
          }

          this.announcements = res.detalle;

          this.announcements.sort((a, b) => {
              return new Date(b.fecha_fin).getTime() - new Date(a.fecha_fin).getTime();
            }
          );
        },
        error: (error: any) => {
          this._toastService.add({type: 'error', message: 'Error al obtener las convocatorias'});
          this.isLoading = false;
        }
      })
    );
  }

  protected selectAnnouncement(item: IAnnouncement): void {
    this.announcementSelect = item;
    this.view = 'request';
  }

}
