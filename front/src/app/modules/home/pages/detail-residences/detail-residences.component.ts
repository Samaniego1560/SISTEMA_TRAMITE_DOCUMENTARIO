import { Component, EventEmitter, Input, Output, Pipe, PipeTransform } from '@angular/core';
import { Subscription } from 'rxjs';
import { ManagerService } from '../../../../core/services/manager/manager.service';
import { ToastService } from '../../../../core/services/toast/toast.service';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators, FormsModule} from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { ToastComponent } from '../../../../core/ui/toast/toast.component';
import { BlockUiComponent } from '../../../../core/ui/block-ui/block-ui.component';
import { ModalComponent } from '../../../../core/ui/modal/modal.component';
import { NgFor, NgIf, CommonModule } from '@angular/common';
import { Residence, StudentInfo } from '../../../../core/models/residence';
import { DatasetSexo } from '../../../../core/contans/postulation';
import { Submission } from '../../../../core/models/residence';
import {IAnnouncement} from '../../../../core/models/announcement';
import { ResidencesService } from '../../../../core/services/residences/residences/residences.service';
import { RoomsService } from '../../../../core/services/residences/rooms/rooms.service';
import { SubmissionService } from '../../../../core/services/residences/submission/submissions.service';
import { FilterStudentsPipe } from '../../../../core/ui/pipes/filter-students.pipe';
import { StudentDebtInfo } from '../../../../core/models/residence';
import {Store} from "@ngrx/store";
import { AppState } from '../../../../core/store/app.reducers';

@Component({
  selector: 'app-residents',
  standalone: true,
  imports: [ToastComponent, BlockUiComponent, ModalComponent, ReactiveFormsModule, NgIf, NgFor,FormsModule, CommonModule],
  templateUrl: './detail-residences.component.html',
  styleUrl: './detail-residences.component.scss'
})

export class DetailResidencesComponent {

  @Input() view: string = 'list';//residents
  @Input() residence: Residence = {} as Residence;
  @Input() residenceView: string = '';
  @Output() cerrarHijo: EventEmitter<void> = new EventEmitter<void>();

  Math = Math;

  protected showModal: boolean = false;
  protected showRoomModal: boolean = false;
  protected showUpdateRoomAssignationModal: boolean = false;
  protected showRoomAssignation: boolean = false;
  protected modalTitle: string = 'Crear residencia';
  protected deleteModal: boolean = false;
  protected balanceModal: boolean = false;
  protected isLoading: boolean = false;
  protected isUpdateRoom: boolean = false;
  protected isUpdateStudent: boolean = false;
  protected formResident: FormGroup;
  protected formFilter: FormGroup;
  protected formRoom: FormGroup;
  protected formStudent: FormGroup;
  protected formStudentAssignation: FormGroup;
  protected formRoomAssignation: FormGroup;
  protected typeListSelected: string = 'cuartos'; //cuartos
  protected annoucements: any;
  protected announcementSelected = {} as any;
  protected roomsResidence: any = [];
  protected studentsResidence: any = [];
  protected roomResidenceSelected: any = {};
  protected studentResidenceSelected: any;
  protected studentsRoomSelected: any = [];
  protected students: any = [];
  protected studentEnabledForResidence: StudentInfo[] = [];
  protected reasonUnassing: string = '';
  protected statusUnassing: string = '';
  protected infoDeleteAssignation: {student_id: number, room_id: string} = {student_id: 0, room_id: ''};
  //Pagination
  protected rooms: any = [];
  protected page: number = 1;
  protected totalItems: number = 0;
  protected pages = 0;
  protected itemsPerPage: number = 10;
  protected currentPage: number = 1;
  protected leftLimit: number = 0;
  protected rightLimit: number = 10;
  protected filterStudent: string = '';
  protected filter = '';

  //Filter students
  dropdownOpen: boolean = false;
  filteredStudents: StudentInfo[] = [];
  protected selectedStudentForm: string = '';
  protected studentdDebts: StudentDebtInfo = {} as StudentDebtInfo;
  public filterStudentResidence:any[] = [];

  //GUARD
  protected isAuth: boolean = false;
  protected role: number = 0;

  public DatasetSexo: {name: string, value: string}[] = DatasetSexo;

  constructor(
      private _fb: FormBuilder,
      private _toastService: ToastService,
      private _managerService: ManagerService,
      private _residenceService: ResidencesService,
      private _roomService: RoomsService,
      private _submissionService: SubmissionService,
      private _store: Store<AppState>
    ) {
    this.formResident = this._fb.group({
      name: ['', Validators.required],
      student_code: ['', Validators.required],
      type_residence: ['', Validators.required],
      residence: ['', Validators.required],
      status: ['', Validators.required]
    });
    this.formFilter = this._fb.group({
      filter: ['', Validators.required],
      type_residence: ['', Validators.required],
      residence: ['', Validators.required]
    });
    this.formRoom = this._fb.group({
      capacity: ['', Validators.required],
      status: ['', Validators.required]
    });
    this.formStudent = this._fb.group({
      student: ['', Validators.required]
    });
    this.formRoomAssignation = this._fb.group({
      student: ['', Validators.required]
    });
    this.formStudentAssignation = this._fb.group({
      student: ['', Validators.required]
    });
    this._store.select('auth').subscribe((auth) => {
      this.isAuth = auth.isAuth;
      this.role = auth.role;
    });
  }

  public async ngOnInit(): Promise<void> {
    this.getAnnoucement();
  }

  protected onBackResidences(): void {
    this.view = 'list';

    this.cerrarHijo.emit();
  }

  protected findResidents(): void {
    if (this.studentsResidence.length > 0) {
      this.filterStudentResidence = this.studentsResidence.filter((item: any) => item.student.full_name.toLowerCase().includes(this.filter.toLocaleLowerCase()));
    }
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
          this.announcementSelected = this.annoucements.find((item: any) => item.estado === 'Activa');
          if (!this.announcementSelected) {
            this.announcementSelected = this.annoucements[this.annoucements.length - 1];
          }
          this.getRoomsResidence(this.residence!.id);
          this.getStudents(this.announcementSelected.id);
        },
        error: (err: any) => {
          console.error(err);
          this.isLoading = false;
        }
      });
    } catch (error: any) {
      this.isLoading = false;
      this._toastService.add({type: 'error', message: error.message});
    }
  }

  protected getResidence(idResidence: string): void {
    this.isLoading = true;
    this._residenceService.getResidence(idResidence).subscribe({
      next: (res: any) => {
        this.isLoading = false;
        if (res && res.error) {
          this._toastService.add({type: 'error', message: 'No se pudo obtener los cuartos, intente nuevamente'});
          return;
        }

      },
      error: (err: any) => {
        console.error(err);
        this._toastService.add({type: 'error', message: 'No se pudo obtener los cuartos, intente nuevamente'});
        this.isLoading = false;
      }
    });
  }

  protected getRoomsResidence(idResidence: string): void {
    this.isLoading = true;
    this._residenceService.getRoomsResidence(idResidence, {page: this.currentPage, limit: this.itemsPerPage, submission_id: this.announcementSelected.id, gender: this.residence.gender}).subscribe({
      next: (res: any) => {
        this.isLoading = false;
        if (res && res.error) {
          this._toastService.add({type: 'error', message: 'No se pudo obtener los cuartos, intente nuevamente'});
          return;
        }
        this.rooms = res.data;
        this.totalItems = this.rooms.length;
        this.pages = this.Math.ceil(this.totalItems / this.itemsPerPage);
        this.getEnabledRooms();
        this.getStudentsEnabled(this.announcementSelected.id);
      },
      error: (err: any) => {
        console.error(err);
        this._toastService.add({type: 'error', message: 'No se pudo obtener los cuartos, intente nuevamente'});
        this.isLoading = false;
      }
    });
  }

  protected getEnabledRooms(): any {
    this.roomResidenceSelected = this.roomsResidence.filter((item: any) => item.status === 'habilitado');
  }

  protected getStudents(annoucementId: number): void {
    try {
      const data = {page: this.currentPage, limit: this.itemsPerPage, gender: this.residence.gender, filter: this.filter};
      this._submissionService.getStudentsByAnnouncement(annoucementId, data).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res) {
            return;
          }
          this.students = res.data.students;
          this.filterStudentResidence = res.data.students;
          this.studentEnabledForResidence = this.students.filter((item: any) => !item.student.room);
        },
        error: (err: any) => {
          console.error(err);
          this.isLoading = false;
        }
      });
    } catch (error: any) {
      this.isLoading = false;
      this._toastService.add({type: 'error', message: error.message});
    }
  }

  protected getStudentsEnabled(annoucementId: number): void {
    this.isLoading = true;
    try {
      const data = {gender: this.residence.gender};
      this._submissionService.getStudentsByAnnouncement(annoucementId, data).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res) {
            this.isLoading = false;
            return;
          }
          this.students = res.data.students;
          this.filterStudentResidence = res.data.students;
          this.studentEnabledForResidence = this.students.filter((item: any) => !item.student.room);
          this.isLoading = false;
        },
        error: (err: any) => {
          console.error(err);
          this.isLoading = false;
        }
      });
    } catch (error: any) {
      this.isLoading = false;
      this._toastService.add({type: 'error', message: error.message});
    }
  }

  protected onChangeAnnouncement(event: any): void {
    debugger
    this.announcementSelected = this.annoucements.find((item: Submission) => item.id === Number(event.target.value));
    if (this.typeListSelected === 'alumnos') {
      this.getStudents(this.announcementSelected.id);
    } else {
      this.getRoomsResidence(this.residence!.id);
    }
  }

  protected onEditResident(): void {
    this.modalTitle = 'Editar asignación de residencia';
    this.showModal = true;
  }

  protected onDeleteResident(): void {
    this.deleteModal = true;
  }

  protected onBalanceModal(): void {
    this.balanceModal = true;
  }

  protected sortResidenceAssignation(): void {
    this.isLoading = true;
    try {
      this._residenceService.sortResidenceAssignation(this.residence?.id).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (res && res.error) {
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }
          this._toastService.add({type: 'success', message: 'Balance realizado correctamente'});
          this.balanceModal = false;
        },
        error: (err: any) => {
          this._toastService.add({type: 'error', message: 'No se pudo realizar el balance, intente nuevamente'});
          this.isLoading = false;
        }
      });
    } catch (error: any) {
      this.isLoading = false;
      this._toastService.add({type: 'error', message: 'No se pudo realizar el balance, intente nuevamente'});
    }
  }

  protected assignRoomStudent( idRoom: string,idStudent: string): void {
    this.isLoading = true;
    try {
      this._roomService.assignRoomStudent({submission_id: this.announcementSelected.id,student_id: Number(idStudent),room_id:idRoom}).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (res && res.error) {
            this._toastService.add({type: 'error', message: 'No se pudo realizar la asignación, intente nuevamente'});
            return;
          }
          this._toastService.add({type: 'success', message: 'Asignación realizada correctamente'});
          this.getRoomsResidence(this.residence!.id);
          this.getStudents(this.announcementSelected?.id ? this.announcementSelected?.id :0);
          if(this.studentEnabledForResidence.find((item: any) => item.student.id == idStudent)!.student) {
            this.studentsRoomSelected.push(this.studentEnabledForResidence.find((item: any) => item.student.id == idStudent)!.student);
          }
          this.showRoomAssignation = false;
          this.selectedStudentForm = '';
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

  protected onChangeListType(event: any): void{
    this.typeListSelected = event.target.value;
    if (this.typeListSelected === 'cuartos') {
      this.getRoomsResidence(this.residence!.id);
    } else {
      this.getStudentsByResidence();
    }
  }

  protected getStudentsByResidence(): void {
    this.isLoading = true;
    try {
      console.log(this.residence);
      console.log(this.announcementSelected);
      let filters = {page: this.currentPage, limit: this.itemsPerPage, submission_id: this.announcementSelected.id};
      this._residenceService.getStudentsResidence(this.residence.id, filters).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (res && res.error) {
            this._toastService.add({type: 'error', message: 'No se pudo obtener la información, intente nuevamente'});
            return;
          }
          this.studentsResidence = res.data.students;
          this.filterStudentResidence = res.data.students;
          this.isLoading = false;
        },
        error: (err: any) => {
          console.error(err);
          this._toastService.add({type: 'error', message: 'No se pudo obtener la información, intente nuevamente'});
          this.isLoading = false;
        }
      });
    } catch (error: any) {
      this.isLoading = false;
      this._toastService.add({type: 'error', message: 'No se pudo obtener la información, intente nuevamente'});
    }
  }

  protected onDeleteAssignation(studentInfo: any): void {
    this.isLoading = true;
    try {
      let data = {
          student_id: studentInfo.student_id,
          submission_id: this.announcementSelected.id,
          room_id:studentInfo.room_id,
          observation: this.reasonUnassing,
          status: this.statusUnassing
      };
      this._roomService.deleteStudentAssignation(data).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (res && res.error) {
            this._toastService.add({type: 'error', message: 'No se pudo eliminar la asignación, intente nuevamente'});
            return;
          }
          this._toastService.add({type: 'success', message: 'Asignación eliminada correctamente'});
          this.getRoomsResidence(this.residence!.id);
          this.getStudents(this.announcementSelected?.id ? this.announcementSelected?.id :0);
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

  protected getRoomIdByNumber(room_number: string): number {
    return this.rooms.find((item: any) => item.number == room_number).id;
  }

  protected onShowStudentAssignation(student: any): void {
    this.showModal = !this.showModal;
    this.studentResidenceSelected = student;
    console.log("Estudiante seleccionado",this.studentResidenceSelected);
    this.isUpdateStudent = false;
    this.getStudentDebts(this.studentResidenceSelected.student.number_identification);
  }

  protected onShowRoomAssignation(room: any): void {
    this.showRoomModal = !this.showRoomModal;

    this.roomResidenceSelected = {...room};

    this.studentsRoomSelected = room.students ? room.students : [];
    this.formRoom.patchValue({capacity: room.capacity, status: room.status});
    this.isUpdateRoom = false;
  }

  protected onUpdateRoomAssignation(room: any): void {
    this.showUpdateRoomAssignationModal = !this.showUpdateRoomAssignationModal;

    this.roomResidenceSelected = {...room};

    this.studentsRoomSelected = room.students ? [...room.students] : [];
    this.isUpdateRoom = false;
  }

  protected onUpdateRoom(room: any): void {
    this.showRoomModal = !this.showRoomModal;
    this.roomResidenceSelected = {...room};
    this.isUpdateRoom = true;
    this.formRoom.patchValue({capacity: room.capacity, status: room.status});
  }

  protected updateRoom(): void {
    if (this.formRoom.invalid) {
      this._toastService.add({type: 'warning', message: 'Por favor, complete los campos requeridos'});
      return;
    }
    this.isLoading = true;
    let params = {capacity: this.formRoom.value.capacity, status: this.formRoom.value.status, id: this.roomResidenceSelected.id};
    try {
      this._roomService.updateRoom(params).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (res && res.error) {
            this._toastService.add({type: 'error', message: 'No se pudo actualizar el cuarto, intente nuevamente'});
            return;
          }
          this.formRoom.reset();
          this.showRoomModal = false;
          this._toastService.add({type: 'success', message: 'Cuarto actualizado correctamente'});
          this.getRoomsResidence(this.residence!.id);
        },
        error: (err: any) => {
          console.error(err);
          this._toastService.add({type: 'error', message: 'No se pudo actualizar el cuarto, intente nuevamente'});
          this.isLoading = false;
        }
      });
    } catch (error: any) {
      this.isLoading = false;
      this._toastService.add({type: 'error', message: 'No se pudo actualizar el cuarto, intente nuevamente'});
    }
  }

  protected onDeleteAssignationStudent(student: any): void {
    try {
      this.studentResidenceSelected = {...student};
      this.isUpdateStudent = true;
      this.onDeleteAssignation(student);
      this.getStudentsByResidence();
    } catch (error: any) {
      this._toastService.add({type: "error", message: error.message});
    }
  }

  protected closeModal(): void {
    this.showModal = false;
  }
  protected closeRoomModal(): void {
    this.showRoomModal = false;
    this.formRoom.reset();
  }

  protected closeUpdateRoomAssignatioModal(): void {
    this.showUpdateRoomAssignationModal = false;
  }

  protected getCurrentAnnoucement(): void {
     return this.annoucements.find((item: Submission) => item.state === true);
  }

  protected changePage(newPage: number) {
    this.currentPage = newPage;
  }

  protected paginate() : void{
    debugger
    if (this.typeListSelected === 'alumnos') {
      this.getStudents(this.announcementSelected?.id ? this.announcementSelected?.id :0);
      this.totalItems = this.students.length;
    } else if (this.typeListSelected === 'cuartos') {
      this.getRoomsResidence(this.residence!.id);
      this.totalItems = this.rooms.length;
    }
    //this.totalItems = this.students.length;
    //return this.students.slice(this.leftLimit, this.rightLimit);
  }

  protected nextPage(): void {
    this.currentPage++;
    this.leftLimit = this.currentPage * this.itemsPerPage;
    this.rightLimit = this.leftLimit + this.itemsPerPage;
    this.paginate();
  }

  protected previousPage(): void {
    this.currentPage--;
    this.leftLimit = this.leftLimit - this.itemsPerPage;
    this.rightLimit = this.leftLimit + this.itemsPerPage;
    this.paginate();
  }

  get paginatedRooms() {
    this.totalItems = this.rooms.length;
    this.rooms = [...this.roomsResidence];

    const startIndex = (this.currentPage - 1) * this.itemsPerPage;
    const endIndex = startIndex + this.itemsPerPage;
    return this.rooms.slice(startIndex, endIndex);
  }

  getRowIndex(index: number): number {
    return index + 1 + (this.currentPage * this.itemsPerPage) - this.itemsPerPage;
  }

  showRoomAssignationCreator(): void {
    this.showRoomAssignation = true;
    this.formRoomAssignation.reset();
    //Falta filtrar por sexo
    this.studentEnabledForResidence = this.students.filter((item: any) => !item.student.room);
    this.filteredStudents = [...this.studentEnabledForResidence];
  }

  protected cancelRoomAssignationCreator(): void {
    this.showRoomAssignation = false;
    this.formRoomAssignation.reset();
    //this.studentsRoomSelected = [];
    this.selectedStudentForm = '';
  }

  protected addStudentRoom(): void {
    if (this.formRoomAssignation.invalid) {
      this._toastService.add({type: 'warning', message: 'Por favor, complete los campos requeridos'});
      return;
    }
    try {
      this.showRoomAssignation = false;
      let studentId = this.formRoomAssignation.value.student;
      let studentInfo = this.students.find((item: any) => item.student.id == studentId);
      this.assignRoomStudent(this.roomResidenceSelected.id,studentId);
      //this.studentsRoomSelected.push(studentInfo.student);
      this.formRoomAssignation.reset();
    } catch (error: any) {

    }
  }

  protected removeStudentRoom(): void {
    if (this.statusUnassing === '' || this.reasonUnassing === '') {
      this._toastService.add({type: 'warning', message: 'Por favor, ingrese el motivo para eliminar la asignación del cuarto'});
      return;
    }
    const _student = this.studentsRoomSelected.filter((item: any) => item.id !== this.infoDeleteAssignation.student_id);
    let index = this.studentsRoomSelected.indexOf(_student);
    debugger
    try {
      //if(this.checkItemInArray(this.roomResidenceSelected.students,_student)) {
        this.onDeleteAssignation(this.infoDeleteAssignation);
      //}
      this.studentsRoomSelected = this.studentsRoomSelected.splice(index, 1);
      this.showUpdateRoomAssignationModal = false;
      this.deleteModal = false;
      this.reasonUnassing = '';
    } catch (error: any) {
      this._toastService.add({type: 'error', message: error.toString()});
    }
  }

  protected checkItemInArray( arr: any[], obj: any): boolean {
    return arr.includes(obj);
  }

  filterStudents(event: any) {
  const _filter = event.target.value?.toLowerCase();
  if (!_filter || _filter === '') {
    this.filteredStudents = [...this.studentEnabledForResidence];
    return;
  }

  this.filteredStudents = this.studentEnabledForResidence.filter(
    student => {
      return student.student.code.toLowerCase().includes(_filter) ||
             student.student.full_name.toLowerCase().includes(_filter);
    }
  );
}


  studentsEnabled(){
    return  this.filteredStudents;
  }

  protected onSetStudent(student_id: any): void {
    let _student = this.students.find((item: any) => item.student.id == student_id);
    this.formRoomAssignation.controls['student'].setValue(_student.student.id);
    this.selectedStudentForm = _student.student.full_name;
    this.dropdownOpen = false;
  }

  protected toggleDropdown() {
    this.dropdownOpen = !this.dropdownOpen;
  }

  protected getStudentDebts(identification: string) {
    this.isLoading = true;
    this._residenceService.validateStudentDebts(identification).subscribe({
      next: (res: any) => {
        if (!res) {
          this.isLoading = false;
          this._toastService.add({type: 'error', message: 'No se pudo obtener las deudas, intente nuevamente'});
          return;
        }
        this.studentdDebts = res;
        this.studentdDebts.historial_pagos = res.historial_pagos.filter((item: any) => item.monto_restante_deuda !== '0.00');
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

  public onDeleteAssignationRoom(student_id: number, room_id: string): void {
    this.infoDeleteAssignation = {student_id, room_id};
    debugger
    this.deleteModal = true;
  }

  public onCancelAssignationRoom(): void {
    this.deleteModal = false;
  }
}
