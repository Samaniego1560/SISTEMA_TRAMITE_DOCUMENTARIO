import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { ManagerService } from '../../../../core/services/manager/manager.service';
import { ToastService } from '../../../../core/services/toast/toast.service';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators,  } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { ToastComponent } from '../../../../core/ui/toast/toast.component';
import { BlockUiComponent } from '../../../../core/ui/block-ui/block-ui.component';
import { ModalComponent } from '../../../../core/ui/modal/modal.component';
import { NgFor, NgIf, CommonModule } from '@angular/common';
import { Residence, StudentInfo, StudentDebtInfo} from '../../../../core/models/residence';
import { DatasetSexo } from '../../../../core/contans/postulation';
import { Submission } from '../../../../core/models/residence';
import {IAnnouncement} from '../../../../core/models/announcement';
import { RoomsService } from '../../../../core/services/residences/rooms/rooms.service';
import { ResidencesService } from '../../../../core/services/residences/residences/residences.service';
import { SubmissionService } from '../../../../core/services/residences/submission/submissions.service';
import { IDebtsStudent } from '../../../../core/models/user';
import {Store} from "@ngrx/store";
import { AppState } from '../../../../core/store/app.reducers';

@Component({
  selector: 'app-students-residences',
  standalone: true,
  imports: [ToastComponent, BlockUiComponent, ModalComponent, ReactiveFormsModule, NgIf, NgFor,CommonModule, FormsModule],
  templateUrl: './students-residences.component.html',
  styleUrl: './students-residences.component.scss',
  providers: [ManagerService, ToastService]
})
export class StudentsResidencesComponent {
  @Input() view: string = 'list';//residents
    @Input() residenceView: string = '';
    @Output() cerrarHijo: EventEmitter<void> = new EventEmitter<void>();

    Math = Math;
    private _subscriptions: Subscription = new Subscription();
    protected showModal: boolean = false;
    protected showRoomModal: boolean = false;
    protected showRoomAssignation: boolean = false;
    protected modalTitle: string = 'Crear residencia';
    protected deleteModal: boolean = false;
    protected balanceModal: boolean = false;
    protected isLoading: boolean = false;
    protected isUpdateRoom: boolean = false;
    protected isUpdateStudent: boolean = false;
    protected formFilter: FormGroup;
    protected formStudentAssignation: FormGroup;
    protected typeListSelected: string = 'cuartos'; //cuartos
    protected annoucements: any;
    protected annoucementSeleted = {} as any;
    protected roomsResidence: any = [];
    protected studentsResidence: any = [];
    protected roomResidenceSelected: any = {};
    protected studentSelected: StudentInfo;
    protected studentsRoomSelected: any = [];
    protected students: StudentInfo [] = [];
    protected studentsFiltered: StudentInfo [] = [];
    protected residences: Residence|any = [];
    protected cuartos: any = [];
    protected residence: any = {};
    //Pagination
    protected rooms: any = [];
    protected page: number = 1;
    protected totalItems: number = 0;
    protected pages = 0;
    protected itemsPerPage: number = 10;
    protected currentPage: number = 1;
    protected leftLimit: number = 0;
    protected rightLimit: number = 10;
    protected studentdDebts: StudentDebtInfo = {} as StudentDebtInfo;
    protected reasonUnassing: string = '';
    protected statusUnassing: string = '';
    protected infoDeleteAssignation: {student_id: number, room_id: string} = {student_id: 0, room_id: ''};

    public DatasetSexo: {name: string, value: string}[] = DatasetSexo;

    //GUARD
    protected isAuth: boolean = false;
    protected role: number = 0;

    constructor(
      private _fb: FormBuilder,
      private _toastService: ToastService,
      private _managerService: ManagerService,
      private router: Router,
      private _roomService: RoomsService,
      private _residencesService: ResidencesService,
      private _submissionService: SubmissionService,
      private _store: Store<AppState>
    ) {

      this.formFilter = this._fb.group({
        filter: ['', Validators.required]
      });
      this.formStudentAssignation = this._fb.group({
        student: ['', Validators.required],
        residence: ['', Validators.required],
        room: ['', Validators.required],

      });

      this._store.select('auth').subscribe((auth) => {
        this.isAuth = auth.isAuth;
        this.role = auth.role;
      });

      this.studentSelected = {
        student: {
          full_name: '',
          code: '',
          professional_school: '',
          residence: '',
          number_identification: '',
          sex: '',
          room: {
            number: 0,
            id: ''
          },
          admission_date: '',
          id: 0
        },
        assigned_goods: [],
        sanctions: [],
        room_mates: []
      };
    }

    public async ngOnInit(): Promise<void> {
      this.getAnnoucement();
      this.getListResidences();
    }

    protected onBackResidences(): void {
      this.router.navigate(['home/residences']);
    }

    protected onShowStudentAssignation(student: any): void {
      //this.studentdDebts = null;
      this.showModal = !this.showModal;
      this.studentSelected = student;
      this.isUpdateStudent = false;
      this.getStudentDebts(student.student.number_identification);
    }

    protected getStudentDebts(identification: string) {
      this.isLoading = true;
      this._residencesService.validateStudentDebts(identification).subscribe({
        next: (res: any) => {
          if (!res) {
            this.isLoading = false;
            this._toastService.add({type: 'error', message: 'No se pudo obtener las deudas, intente nuevamente'});
            return;
          }
          this.studentdDebts = res;
          this.studentdDebts.historial_pagos = res.historial_pagos.filter((item: any) =>
            item.monto_restante_deuda !== '0.00' &&
            (
                item.desc_deuda.toLowerCase().includes('comedor') ||
                item.desc_deuda.toLowerCase().includes('internado')
            )
        );
          console.log(this.studentdDebts);

          this.isLoading = false;
        },
        error: (err: any) => {
          console.error(err);
          this._toastService.add({type: 'error', message: 'No se pudo obtener las deudas, intente nuevamente'});
          this.isLoading = false;
        }
      });
    }

    protected findResidents(): void {
      const _studentsFiltered = this.students.filter((item: any) => {
        return (item.student.full_name.toLowerCase().includes(this.formFilter.value.filter.toLowerCase()) ||
        item.student.code.toLowerCase().includes(this.formFilter.value.filter.toLowerCase())) || this.formFilter.value.filter.toLowerCase() === ''
      });
      this.totalItems = _studentsFiltered.length;
      this.studentsFiltered = _studentsFiltered.slice(this.leftLimit, this.rightLimit);
      this.pages = this.Math.ceil(this.totalItems / this.itemsPerPage);
      //this.getStudents(this.annoucementSeleted?.id ? this.annoucementSeleted?.id :0);
    }

    protected getListResidences(): void {
      this.isLoading = true;
          this._residencesService.getListResidences().subscribe({
            next: (res: any) => {
              this.isLoading = false;
              if (res && res.error) {
                this._toastService.add({type: 'error', message: 'No se pudo obtener las residencias, intente nuevamente'});
                return;
              }
              this.residences = res.data;
            },
            error: (err: any) => {
              console.error(err);
              this._toastService.add({type: 'error', message: 'No se pudo obtener las residencias, intente nuevamente'});
              this.isLoading = false;
            }
          })
    }

    protected getAnnoucement(): void {
      try {
        this._managerService.getAnnouncement().subscribe({
          next: (res: any) => {
            this.isLoading = false;
            if (!res) {
              return;
            }
            this.annoucements = res.detalle;
            this.balanceModal = false;
            this.annoucementSeleted = this.annoucements.find((item: any) => item.estado === 'Activa');
            if (!this.annoucementSeleted) {
              this.annoucementSeleted = this.annoucements[this.annoucements.length - 1];
            }
            this.getStudents(this.annoucementSeleted?.id ? this.annoucementSeleted?.id :0);
          },
          error: (err: any) => {
            console.error(err);
            this.isLoading = false;
          }
        });
      } catch (error: any) {
        this.isLoading = false;
        console.log(error);
      }
    }

    protected getStudents(annoucement_id: number): void {
      this.isLoading = true;
      let filters = {filter: this.formFilter.value.filter};
      try {
          this._submissionService.getStudentsByAnnouncement(annoucement_id, filters).subscribe({
            next: (res: any) => {
              if (!res) {
                this.isLoading = false;
                return;
              }
              if (res && res.error) {
                return;
              }
              this.students = res.data.students;
              this.totalItems = this.students.length;
              this.pages = this.Math.ceil(this.totalItems / this.itemsPerPage);
              this.studentsFiltered = this.students.slice(this.leftLimit, this.rightLimit);
              this.isLoading = false;
            },
            error: (err: any) => {
              console.error(err);
              this.isLoading = false;
            }
          });
      } catch (error: any) {
        this.isLoading = false;
      }
    }

    protected getRoomsResidence(idResidence: string): void {
      this.isLoading = true;
          this._residencesService
          .getRoomsResidence(idResidence, {page: 1, limit: 100, submission_id: this.annoucementSeleted.id}).subscribe({
            next: (res: any) => {
              this.isLoading = false;
              if (res && res.error) {
                this._toastService.add({type: 'error', message: 'No se pudo obtener los cuartos, intente nuevamente'});
                return;
              }
              this.roomsResidence = res.data;
              this.getEnabledRooms();
            },
            error: (err: any) => {
              console.error(err);
              this._toastService.add({type: 'error', message: 'No se pudo obtener los cuartos, intente nuevamente'});
              this.isLoading = false;
            }
          })
    }

    protected onToogleResidence(event: any): void {
      this.residence = this.residences.find((item: Residence) => item.id === event.target.value);
      this.getRoomsResidence(this.residence.id);
    }

    protected getEnabledRooms(): any {
      this.cuartos = this.roomsResidence.filter((item: any) => item.status === 'habilitado');
    }

    protected updateAssignation(): void {
      if(this.formStudentAssignation.invalid) {
        this._toastService.add({type: 'error', message: 'Complete los campos requeridos'});
        return;
      }
      let studentAssignation = this.formStudentAssignation.value;
      this.assignRoomStudent(studentAssignation.room, studentAssignation.student);
      this.showRoomAssignation = false;
    }

    protected onChangeAnnouncement(event: any): void {
      this.annoucementSeleted = this.annoucements.find((item: Submission) => item.id == event.target.value);
      this.getStudents(this.annoucementSeleted?.id ? this.annoucementSeleted?.id :0);
    }


    protected assignRoomStudent( idRoom: string,idStudent: string): void {
      this.isLoading = true;
      let params = {room_id:idRoom, submission_id: this.annoucementSeleted.id, student_id: idStudent};
      try {
        this._roomService.assignRoomStudent(params).subscribe({
          next: (res: any) => {
            this.isLoading = false;
            if (res && res.error) {
              this._toastService.add({type: 'error', message: 'No se pudo realizar la asignación, intente nuevamente'});
              return;
            }
            this._toastService.add({type: 'success', message: 'Asignación realizada correctamente'});
            this.showModal = false;
            this.getStudents(this.annoucementSeleted?.id ? this.annoucementSeleted?.id :0);
          },
          error: (err: any) => {
            console.error(err);
            this._toastService.add({type: 'error', message: 'No se pudo realizar la asignación, intente nuevamente'});
            this.isLoading = false;
          }
        });
      } catch (error: any) {
        this.isLoading = false;
        this._toastService.add({type: 'error', message: 'No se pudo realizar la asignación, intente nuevamente'});
      }
    }

    protected getResidenceName(id: string): string {
      let residence = this.residences.find((item: Residence) => item.id === id);
      return residence ? residence.name : '';
    }

    getEnabledResidences(gender:string): any[] {
      console.log('Genero filtrado: ',gender);
      console.log('Residencias: ',this.residences);
      return this.residences.filter((item: Residence) => (gender === 'F' && item.gender === 'femenino') || (gender === 'M' && item.gender === 'masculino'));
      //return this.residences;
    }

    protected assignStudentRoom(student: any): void {
      this.studentSelected = student;
      this.showRoomAssignation = !this.showRoomAssignation;
      this.formStudentAssignation.patchValue({student:this.studentSelected.student.id});
      if (this.studentSelected.student.residence && this.studentSelected.student.residence !== '')
          this.getRoomsResidence(this.studentSelected.student.residence);
    }

    protected closeModalAssignation(): void {
      this.showRoomAssignation = false;
      this.formStudentAssignation.reset();
    }

    protected closeModal(): void {
      this.showModal = false;
    }

    protected getCurrentAnnoucement(): void {
       return this.annoucements.find((item: Submission) => item.state === true);
    }

    protected paginate() : void{
      //this.getStudents(this.annoucementSeleted?.id ? this.annoucementSeleted?.id :0);
      this.studentsFiltered = this.students.slice(this.leftLimit, this.rightLimit);
      this.totalItems = this.students.length;
      //return this.students.slice(this.leftLimit, this.rightLimit);
    }

    protected nextPage(): void {
      this.leftLimit = this.currentPage * this.itemsPerPage;
      this.rightLimit = this.leftLimit + this.itemsPerPage;
      this.currentPage++;
      this.paginate();
    }

    protected previousPage(): void {
      this.currentPage--;
      this.leftLimit = this.leftLimit - this.itemsPerPage;
      this.rightLimit = this.leftLimit + this.itemsPerPage;
      this.paginate();
    }

    getRowIndex(index: number): number {
      return index + 1 + (this.currentPage * this.itemsPerPage) - this.itemsPerPage;
    }

    public onDeleteAssignationRoom(student_id: number, room_id: string): void {
      this.infoDeleteAssignation = {student_id, room_id};
      debugger
      this.deleteModal = true;
    }

    public onCancelAssignationRoom(): void {
      this.deleteModal = false;
      this.reasonUnassing = '';
      this.statusUnassing = '';
    }

    protected deleteAssignation(): void {
      this.isLoading = true;
      if (this.reasonUnassing === '' || this.statusUnassing === '') {
        this._toastService.add({type: 'error', message: 'Complete la información requerida'});
        this.isLoading = false;
        return;
      }
      let data = {
          student_id:this.infoDeleteAssignation.student_id,
          submission_id:this.annoucementSeleted.id,
          room_id: this.infoDeleteAssignation.room_id,
          observation: this.reasonUnassing,
          status: this.statusUnassing
        };
      try {
        this._roomService.deleteStudentAssignation(data).subscribe({
          next: (res: any) => {
            this.isLoading = false;
            if (res && res.error) {
              this._toastService.add({type: 'error', message: 'No se pudo eliminar la asignación, intente nuevamente'});
              return;
            }
            this._toastService.add({type: 'success', message: 'Asignación eliminada correctamente'});
            this.getStudents(this.annoucementSeleted?.id ? this.annoucementSeleted?.id :0);
            this.reasonUnassing = '';
            this.statusUnassing = '';
            this.deleteModal = false;
          },
          error: (err: any) => {
            console.error(err);
            this._toastService.add({type: 'error', message: 'No se pudo eliminar la asignación, intente nuevamente'});
            this.isLoading = false;
          }
        });
      } catch (error: any) {
        this.isLoading = false;
        this._toastService.add({type: 'error', message: 'No se pudo eliminar la asignación, intente nuevamente'});
      }
    }
}
