import { Component, OnInit } from '@angular/core';
import { Residence, Dataset, DatasetResidence, Floors, Room, ConfigurationRules} from '../../../../core/models/residence';
import { Subscription } from 'rxjs';
import { ManagerService } from '../../../../core/services/manager/manager.service';
import { ToastService } from '../../../../core/services/toast/toast.service';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { ToastComponent } from '../../../../core/ui/toast/toast.component';
import { BlockUiComponent } from '../../../../core/ui/block-ui/block-ui.component';
import { ModalComponent } from '../../../../core/ui/modal/modal.component';
import { CommonModule, NgFor, NgIf } from '@angular/common';
import { DetailResidencesComponent } from '../detail-residences/detail-residences.component';
import { DatasetSexo } from '../../../../core/contans/postulation';
import { v4 as uuidv4 } from 'uuid';
import {ResidencesService} from '../../../../core/services/residences/residences/residences.service';
import { SubmissionService } from '../../../../core/services/residences/submission/submissions.service';
import {Store} from "@ngrx/store";
import { AppState } from '../../../../core/store/app.reducers';

@Component({
  selector: 'app-residences',
  standalone: true,
  imports: [ToastComponent, BlockUiComponent, ModalComponent, ReactiveFormsModule, NgIf, NgFor,DetailResidencesComponent,CommonModule, FormsModule],
  templateUrl: './residences.component.html',
  styleUrl: './residences.component.scss',
  providers: [ToastService]
})
export class ResidencesComponent implements OnInit{



  protected typeModal: string = 'create';
  private typeView: string[] = ['residence', 'residents'];
  private _subscriptions: Subscription = new Subscription();
  public deleteModal: boolean = false;

  protected modalTitle: string = 'Crear residencia';
  protected thesisPracticingSelected!: 1;
  protected residences: Residence[] = [];
  protected residencesFiltered: Residence[] = [];
  protected residenceSelected!: Residence;
  protected isLoading: boolean = false;
  protected formResidence: FormGroup;
  protected formFloor: FormGroup;
  protected formRulesAssgination: FormGroup;
  protected showModal: boolean = false;
  protected showRulesModal: boolean = false;
  protected view: string = 'list';//residents
  protected showFloorCreator: boolean = false;
  protected mostrarHijo: boolean = true;
  protected floors: Floors[] = [];
  protected annoucements: any[] = [];
  protected annoucementSeleted: any;
  protected filterGender: string = 'todos';
  protected filterName: string = '';
  //GUARD
  protected isAuth: boolean = false;
  protected role: number = 0;

  public DatasetSexo: {name: string, value: string}[] = DatasetSexo;

  public _announcements: any[] = [];

  constructor(
      private _fb: FormBuilder, private _toastService: ToastService ,private _residenceService: ResidencesService,
      private _managerService: ManagerService, private _submissionService: SubmissionService, private _store: Store<AppState>) {
    this.formResidence = this._fb.group({
      name: ['', Validators.required],
      gender: ['', Validators.required],
      status: [''],
      description : ['', Validators.required],
      address: ['', Validators.required]
    });
    this.formFloor = this._fb.group({
      floor: ['', Validators.required],
      q_rooms: ['', Validators.required],
      space_rooms: ['', Validators.required]
    });
    this.formRulesAssgination = this._fb.group({
      percentage_fcea: ['', Validators.required],
      percentage_engineering: ['', Validators.required],
      minimum_grade_fcea: ['', Validators.required],
      minimum_grade_engineering: ['', Validators.required],
      is_newbie: ['', Validators.required],
    });
    this._store.select('auth').subscribe((auth) => {
      this.isAuth = auth.isAuth;
      this.role = auth.role;
    });
  }

  ngOnInit() {
    this.getAnnoucement();
    this.getListResidences();
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
          this.annoucementSeleted = this.annoucements.find((item: any) => item.estado === 'Activa');
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

  protected onCreateResidence(): void {
    this.typeModal = 'create';
    this.modalTitle = 'Crear Residencia'
      this.formResidence.reset();
      this.formResidence.enable();
      this.showModal = true;
  }

  protected onUpdateResidence(residences: Residence): void {
    this.residenceSelected = residences;
    this.typeModal = 'update';
    this.showModal = true;
    this.formResidence.patchValue({
      name: this.residenceSelected.name,
      gender: this.residenceSelected.gender,
      status: this.residenceSelected.status,
      description : this.residenceSelected.description,
      address: this.residenceSelected.address,
    });
    this.floors = this.residenceSelected.floors;
  }

  protected onListResidents(resident: any): void {
    this.view = 'residents';
    this.residenceSelected = resident;
  }

  public toggleFloorCreator(): void {
     this.showFloorCreator = true;
  }

  protected cancelFloorCreator(): void {
    this.showFloorCreator = false;
  }

  protected addFloor(): void {
    if (this.formFloor.invalid) {
      this.formFloor.markAllAsTouched();
      this._toastService.add({type: 'error', message: 'Por favor complete los campos de piso requeridos'});
      return;
    }
    const floor = this.formFloor.value;
    this.floors.push({
      floor: floor.floor,
      q_rooms: floor.q_rooms,
      space_rooms: floor.space_rooms,
      rooms: []
    });
    this.formFloor.reset();
    this.showFloorCreator = false;
  }

  protected hiddenViewResidents(): void {
    this.mostrarHijo = false;
    this.view = 'list';
  }

  protected closeModal(): void {
    this.showModal = !this.showModal;
    this.residenceSelected = {} as Residence;
    this.floors = [];
  }

  protected closeRulesModal(): void {
    this.showRulesModal = !this.showRulesModal;
  }

  protected onDeleteFloor(floor: Floors): void {
    this.floors = this.floors.filter((item: any) => item.floor !== floor.floor);
  }

  protected getListResidences(): void {
    this.isLoading = true;

      this._subscriptions.add(
        this._residenceService.getListResidences().subscribe({
          next: (res: any) => {
            this.isLoading = false;
            if (res && res.error) {
              this._toastService.add({type: 'error', message: 'No se pudo obtener las residencias, intente nuevamente'});
              return;
            }
            this.residences = res.data;
            this.residencesFiltered = [...this.residences];
          },
          error: (err: any) => {
            console.error(err);
            this._toastService.add({type: 'error', message: 'No se pudo obtener las residencias, intente nuevamente'});
            this.isLoading = false;
          }
        })
      );
  }

  public getNameConvocatoriaById(id: number) {
    return this._announcements.find((item) => item.value == id)?.name;
  }

  protected saveResidence(): void {
    if (this.formResidence.invalid) {
      this.formResidence.markAllAsTouched();
      this._toastService.add({type: 'error', message: 'Por favor complete los campos requeridos'});
      return;
    }

    this.isLoading = true;
    let residence = this.formResidence.value;
    residence.floors = this.floors;
    if (this.typeModal === 'create') {
      if (this.floors.length === 0) {
        this._toastService.add({type: 'error', message: 'Por favor agregue al menos un piso'});
        this.isLoading = false;
        return;
      }
      residence.id = uuidv4();
      residence.status = 'habilitado';
      residence.floors = this.onBuildRooms(residence.floors);
      this._subscriptions.add(
        this._residenceService.createResidence(residence).subscribe({
          next: (res: any) => {
            this.isLoading = false;
            if (res && res.error) {
              this._toastService.add({type: 'error', message: res.msg});
              return;
            }
            this.getListResidences();
            this._toastService.add({type: 'success', message: 'Residencia creada correctamente'});
            this.showModal = false;
          },
          error: (err: any) => {
            console.error(err);
            this._toastService.add({type: 'error', message: 'No se pudo crear la residencia, intente nuevamente'});
            this.isLoading = false;
          }
        })
      );
    } else {
      this._subscriptions.add(
        //residence.id = this.residenceSelected.id;
        this._residenceService.updateResidence(residence, this.residenceSelected.id).subscribe({
          next: (res: any) => {
            this.isLoading = false;
            if (res && res.error) {
              this._toastService.add({type: 'error', message: res.msg});
              return;
            }
            this.getListResidences();
            this.showModal = false;
          },
          error: (err: any) => {
            console.error(err);
            this._toastService.add({type: 'error', message: 'No se pudo actualizar la residencia, intente nuevamente'});
            this.isLoading = false;
          }
        })
      );
    }
  }

  public onBuildRooms(floors: Floors[]): Floors[] {
    for (let i = 0; i < floors.length; i++) {
      let rooms = [];
      for (let j = 0; j < floors[i].q_rooms!; j++) {
        rooms.push({
          id: uuidv4(),
          number: j + 1,
          capacity: floors[i]!.space_rooms,
          status: 'habilitado'
        });
      }
      floors[i].rooms = rooms;
      delete floors[i].q_rooms;
      delete floors[i].space_rooms;
    }
    return floors;
  }

  protected saveRulesModal(): void {
      if(this.formRulesAssgination.invalid) {
        this.formRulesAssgination.markAllAsTouched();
        this._toastService.add({type: 'error', message: 'Por favor complete los campos requeridos'});
        return;
      }
      this.isLoading = true;
      let rules = this.formRulesAssgination.value;
      rules.id = this.residenceSelected.id;
      rules.is_newbie = rules.is_newbie === 'true'? true : false;
      if (this.residenceSelected.configuration === rules) {
        this._toastService.add({type: 'info', message: 'No se realizaron cambios'});
        this.isLoading = false;
        this.showRulesModal = false;
        return;
      }
      this._subscriptions.add(
        this._residenceService.updateRulesResidence(this.residenceSelected.id,rules).subscribe({
          next: (res: any) => {
            this.isLoading = false;
            if (res && res.error) {
              this._toastService.add({type: 'error', message: res.msg});
              return;
            }
            this.getListResidences();
            this.showRulesModal = false;
          },
          error: (err: any) => {
            console.error(err);
            this._toastService.add({type: 'error', message: 'No se pudo actualizar las reglas de asignación, intente nuevamente'});
            this.isLoading = false;
          }
        })
      );
  }

  protected onUpdateRules(residence: any): void {
    this.residenceSelected = residence;
    this.showRulesModal = true;
    this.formRulesAssgination.patchValue({
      percentage_fcea: residence.configuration.percentage_fcea,
      percentage_engineering: residence.configuration.percentage_engineering,
      minimum_grade_fcea: residence.configuration.minimum_grade_fcea,
      minimum_grade_engineering: residence.configuration.minimum_grade_engineering,
      is_newbie: residence.configuration.is_newbie
    });
  }

  protected toogleDeleteModal(residence: Residence): void {
    this.residenceSelected = residence;
    this.deleteModal = true;
  }

  protected cancelDeleteResidence(): void {
    this.deleteModal = false;
    this.residenceSelected = {} as Residence;
  }

  protected deleteResidence(): void {
    this.isLoading = false;
    this._subscriptions.add(
      this._residenceService.deleteResidence(this.residenceSelected.id).subscribe({
        next: (res: any) => {
          this.isLoading = false;

          if (res.error) {
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }

          this.getListResidences();
          this.deleteModal = false;
        },
        error: (err: any) => {
          console.error(err);
          this._toastService.add({type: 'error', message: 'No se pudo borrar la residencia, intente nuevamente.'});
          this.isLoading = false;
        }
      })
    );

  }

  protected exportReport(): void {
    this.isLoading = true;
    this._subscriptions.add(
      this._submissionService.exportReportResidence(this.annoucementSeleted.id).subscribe({
        next: (res: any) => {
            this.isLoading = false;
            if (!res || res.error) {
                this._toastService.add({ type: 'error', message: res?.msg || 'Error en la exportación' });
                return;
            }

            if (!res.data) {
                this._toastService.add({ type: 'error', message: 'El archivo no es válido' });
                return;
            }

            try {
                console.log("Base64 recibido:", res.data.slice(0, 50) + "..."); // Verificar los primeros caracteres

                // Convertir Base64 a Blob
                const byteCharacters = atob(res.data);
                const byteArray = new Uint8Array(byteCharacters.length);

                for (let i = 0; i < byteCharacters.length; i++) {
                    byteArray[i] = byteCharacters.charCodeAt(i);
                }

                const blob = new Blob([byteArray], { type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" });

                console.log("Tamaño del archivo generado:", blob.size); // Depuración

                if (blob.size === 0) {
                    throw new Error("El archivo generado está vacío o corrupto.");
                }

                // Crear un enlace para descargar el archivo
                const link = document.createElement("a");
                const url = URL.createObjectURL(blob);

                link.href = url;
                link.download = `Reporte_Residencias_${new Date().toISOString().slice(0, 10)}.xlsx`;

                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);

                URL.revokeObjectURL(url); // Liberar memoria

                this._toastService.add({ type: 'success', message: 'Reporte exportado correctamente' });

            } catch (error: any) {
                console.error("Error al convertir el archivo:", error);
                this._toastService.add({ type: 'error', message: 'No se pudo procesar el archivo' });
            }
        },
        error: (err: any) => {
            console.error(err);
            this._toastService.add({ type: 'error', message: 'No se pudo exportar el reporte, intente nuevamente' });
            this.isLoading = false;
        }
      })
    );
  }

  protected onFilterGender(): void {
    if (this.filterName !== '') {
      this.residencesFiltered = this.residences.filter((item: any) => (item.gender === this.filterGender || this.filterGender === 'todos') && item.name.toLowerCase().includes(this.filterName));
      return;
    }
    this.residencesFiltered = this.residences.filter((item: any) => item.gender === this.filterGender || this.filterGender === 'todos');
  }

  protected onFilterName(): void {
    debugger
    this.residencesFiltered = this.residences.filter((item: any) => item.name.toLowerCase().includes(this.filterName.toLocaleLowerCase()) && (item.gender === this.filterGender || this.filterGender === 'todos'));
  }
}
