import { Component, OnInit, OnDestroy, HostListener } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { TableComponent, TableColumn, TableAction } from '../components/table/table.component';
import { IVisitaGeneral, TipoUsuario, Genero, Departamento, Provincia, Distrito } from '../../../../../core/models/visita-general';
import { ESCUELAS_PROFESIONALES, EscuelaProfesional } from '../models/escuela-profesional';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { Subscription } from 'rxjs';
import { ModalComponent } from "../../../../../core/ui/modal/modal.component";
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators, AbstractControl } from "@angular/forms";
import { NgClass, NgIf } from "@angular/common";
import { BlockUiComponent } from "../../../../../core/ui/block-ui/block-ui.component";
import * as XLSX from 'xlsx';
import { saveAs } from 'file-saver';

@Component({
  selector: 'app-visita-general',
  standalone: true,
  imports: [CommonModule, TableComponent, ModalComponent, ToastComponent, FormsModule, ReactiveFormsModule, NgIf, BlockUiComponent],
  templateUrl: './visita-general.component.html',
  styleUrl: './visita-general.component.css',
  providers: [ToastService]
})
export class VisitaGeneralComponent implements OnInit, OnDestroy { 

  visitaData: IVisitaGeneral[] = [];
  isLoading: boolean = false;
  error: string | null = null;
  private _subscriptions: Subscription = new Subscription();
  protected deleteModal: boolean = false;
  protected showModal: boolean = false;
  protected modalTitle: string = 'Registrar Visita';
  protected formVisita: FormGroup;
  protected modoFormulario: 'crear' | 'editar' | 'ver' = 'crear';
  protected visitaSeleccionada: IVisitaGeneral | null = null;
  
  protected tiposUsuario = [
    { value: TipoUsuario.Alumno, label: 'Alumno' },
    { value: TipoUsuario.Docente, label: 'Docente' },
    { value: TipoUsuario.Administrativo, label: 'Administrativo' }
  ];
  
  protected generos = [
    { value: Genero.F, label: 'F' },
    { value: Genero.M, label: 'M' }
  ];
  
  protected escuelasProfesionales: EscuelaProfesional[] = ESCUELAS_PROFESIONALES;
  
  protected departamentos: Departamento[] = [];
  protected provincias: Provincia[] = [];
  protected distritos: Distrito[] = [];
  
  protected isExporting: boolean = false;
  protected dropdownOpen: boolean = false;
  private isLoadingData: boolean = false; // Bandera para controlar la carga de datos
  
  protected motivosAtencion = [
    { value: 'social', label: 'Social' },
    { value: 'academico', label: 'Académico' },
    { value: 'salud', label: 'Salud' },
    { value: 'economico', label: 'Económico' }
  ];
  
  tableColumns: TableColumn[] = [
    { header: 'Nombre Completo', field: 'nombre_completo' },
    { header: 'Tipo Usuario', field: 'tipo_usuario', filterable: true },
    { header: 'Género', field: 'genero' },
    { header: 'Edad', field: 'edad' },
    { header: 'Motivo Atención', field: 'motivo_atencion', filterable: true },
    { header: 'Lugar Atención', field: 'lugar_atencion' },
  ];
  
  tableActions: TableAction[] = [
    { label: 'Ver', icon: 'eye', action: 'view' },
    { label: 'Editar', icon: 'pencil', action: 'edit' },
    { label: 'Eliminar', icon: 'trash', action: 'delete' }
  ];
  
  tableFilters: {
    field: string;
    label: string;
    options: { value: string; label: string; }[];
    customFilter?: (data: any[], filterValue: string) => any[];
    dependsOn?: string;
    childFilters?: { [key: string]: { value: string; label: string; }[] };
  }[] = [
    {
      field: 'tipo_usuario',
      label: 'Tipo de Usuario',
      options: [
        { value: TipoUsuario.Alumno, label: 'Alumno' },
        { value: TipoUsuario.Docente, label: 'Docente' },
        { value: TipoUsuario.Administrativo, label: 'Administrativo' }
      ]
    },
    {
      field: 'motivo_atencion',
      label: 'Motivo de Atención',
      options: [
        { value: 'social', label: 'Social' },
        { value: 'academico', label: 'Académico' },
        { value: 'salud', label: 'Salud' },
        { value: 'economico', label: 'Económico' }
      ],
      customFilter: (data: IVisitaGeneral[], filterValue: string) => this.filterMotivoAtencion(data, filterValue)
    }
  ];

  constructor(
    private _managerService: ManagerService, 
    private _toastService: ToastService,
    private _fb: FormBuilder
  ) {
    this.formVisita = this._fb.group({
      tipo_usuario: ['', Validators.required],
      codigo_estudiante: [''],
      dni: [''],
      nombre_completo: ['', Validators.required],
      genero: ['', Validators.required],
      edad: ['', [Validators.min(1), Validators.max(120)]],
      escuela: [''],
      area: [''],
      departamento: [''],
      provincia: [''],
      distrito: [''],
      motivo_atencion: [[], (control: AbstractControl) => {
        const value = control.value;
        return value && value.length > 0 ? null : { required: true };
      }],
      descripcion_motivo: ['', Validators.required],
      url_imagen: ['', [this.validarUrlOpcional]],
      lugar_atencion: ['', Validators.required]
    });
    
    this.formVisita.get('tipo_usuario')?.valueChanges.subscribe(tipo => {
      if (tipo === TipoUsuario.Alumno) {
        this.formVisita.get('codigo_estudiante')?.setValidators([Validators.required]);
        this.formVisita.get('dni')?.clearValidators();
        this.formVisita.get('escuela')?.setValidators([Validators.required]);
      } else {
        this.formVisita.get('codigo_estudiante')?.clearValidators();
        this.formVisita.get('dni')?.setValidators([Validators.required]);
        this.formVisita.get('escuela')?.clearValidators();
      }
      
      this.formVisita.get('codigo_estudiante')?.updateValueAndValidity();
      this.formVisita.get('dni')?.updateValueAndValidity();
      this.formVisita.get('escuela')?.updateValueAndValidity();
    });

    this.formVisita.get('departamento')?.valueChanges.subscribe(departamentoNombre => {
      // Solo actuar si no estamos cargando datos existentes
      if (!this.isLoadingData && departamentoNombre) {
        const departamento = this.departamentos.find(d => d.name === departamentoNombre);
        if (departamento) {
          this.cargarProvincias(departamento.id);
        }
        this.formVisita.get('provincia')?.setValue('');
        this.formVisita.get('distrito')?.setValue('');
        this.distritos = [];
      }
    });

    this.formVisita.get('provincia')?.valueChanges.subscribe(provinciaNombre => {
      // Solo actuar si no estamos cargando datos existentes
      if (!this.isLoadingData && provinciaNombre) {
        const provincia = this.provincias.find(p => p.name === provinciaNombre);
        if (provincia) {
          this.cargarDistritos(provincia.id);
        }
        this.formVisita.get('distrito')?.setValue('');
      }
    });
  }
 
  ngOnInit() {
    this.GetAllVisitaGeneralData();
    this.cargarDepartamentos();
  }
  
  ngOnDestroy() {
    this._subscriptions.unsubscribe();
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: Event) {
    const target = event.target as HTMLElement;
    const dropdownContainer = target.closest('[data-dropdown-container]');
    if (!dropdownContainer) {
      this.dropdownOpen = false;
    }
  }
  
  showToast(type: 'success' | 'error' | 'info' | 'warning', message: string) {
    this._toastService.add({
      key: 'visita-general',
      type: type,
      message: message
    });
  }
  
  abrirModalCrear() {
    this.modoFormulario = 'crear';
    this.modalTitle = 'Registrar Visita';
    this.isLoadingData = false; // Resetear bandera
    this.formVisita.reset();
    this.formVisita.get('motivo_atencion')?.setValue([]);
    this.formVisita.enable();
    this.provincias = [];
    this.distritos = [];
    this.dropdownOpen = false;
    this.showModal = true;
  }
  
  GetAllVisitaGeneralData() {
    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.getAllVisitaGeneral().subscribe({
        next: res => {
          this.isLoading = false;
          if(!res.data){
            this.showToast('error', 'No se pudo obtener la información de las visitas generales');
            return;
          }
          this.visitaData = res.data;
        },
        error: err => {
          this.showToast('error', 'Error al obtener la información de las visitas generales');
          this.isLoading = false;
        }
      })
    );
  }
  
  
  onViewItem(item: IVisitaGeneral) {
    this._managerService.getVisitaGeneralById(item.id).subscribe({
      next: (response) => {
        if (response.data) {
          this.visitaSeleccionada = response.data;
          this.modoFormulario = 'ver';
          this.modalTitle = 'Ver Detalles de Visita';
          
          // Activar bandera para prevenir interferencia de valueChanges
          this.isLoadingData = true;
          
          // Cargar datos geográficos de forma secuencial
          this.cargarDatosGeograficos(this.visitaSeleccionada).then(() => {
            const datosProcessados = this.procesarDatosVisita(this.visitaSeleccionada!);
            this.formVisita.patchValue(datosProcessados);
            this.formVisita.disable();
            
            // Desactivar bandera después de cargar los datos
            this.isLoadingData = false;
            this.showModal = true;
          }).catch(() => {
            this.isLoadingData = false;
            this.showModal = true;
          });
        }
      },
      error: (err) => {
        this.isLoadingData = false;
        this.showToast('error', 'Error al obtener los detalles del registro');
      }
    });
  }
  
  onEditItem(item: IVisitaGeneral) {
    this._managerService.getVisitaGeneralById(item.id).subscribe({
      next: (response) => {
        if (response.data) {
          this.visitaSeleccionada = response.data;
          this.modoFormulario = 'editar';
          this.modalTitle = 'Editar Visita';
          
          // Activar bandera para prevenir interferencia de valueChanges
          this.isLoadingData = true;
          
          // Cargar datos geográficos de forma secuencial
          this.cargarDatosGeograficos(this.visitaSeleccionada).then(() => {
            const datosProcessados = this.procesarDatosVisita(this.visitaSeleccionada!);
            this.formVisita.patchValue(datosProcessados);
            this.formVisita.enable();
            
            // Desactivar bandera después de cargar los datos
            this.isLoadingData = false;
            this.showModal = true;
          }).catch(() => {
            this.isLoadingData = false;
            this.showModal = true;
          });
        }
      },
      error: (err) => {
        this.isLoadingData = false;
        this.showToast('error', 'Error al obtener los detalles del registro');
      }
    });
  }
  
  private selectedItem: IVisitaGeneral | null = null;
  
  onDeleteItem(item: IVisitaGeneral) {
    this.selectedItem = item;
    this.deleteModal = true;
  }

  deleteUser() {
    if (this.selectedItem) {
      this._managerService.deleteVisitaGeneral(this.selectedItem.id).subscribe({
        next: (response) => {
          this.showToast('success', 'Registro eliminado correctamente');
          this.GetAllVisitaGeneralData(); 
          this.deleteModal = false; 
          this.selectedItem = null; 
        },
        error: (err) => {
          this.showToast('error', 'Error al eliminar el registro');
        }
      });
    }
  }
  
  guardarVisita() {
    if (this.formVisita.invalid) {
      this.formVisita.markAllAsTouched();
      this.showToast('warning', 'Por favor complete todos los campos requeridos');
      return;
    }
    
    const formValue = this.formVisita.value;
    
    const visitaData = {
      ...formValue,
      motivo_atencion: Array.isArray(formValue.motivo_atencion) 
        ? formValue.motivo_atencion.join(',') 
        : formValue.motivo_atencion
    };
    
    if (this.modoFormulario === 'crear') {
      this._managerService.createVisitaGeneral(visitaData).subscribe({
        next: (response) => {
          this.showToast('success', 'Visita registrada correctamente');
          this.showModal = false;
          this.GetAllVisitaGeneralData();
        },
        error: (err) => {
          this.showToast('error', 'Error al registrar la visita');
        }
      });
    } else if (this.modoFormulario === 'editar' && this.visitaSeleccionada) {
      this._managerService.updateVisitaGeneral(visitaData, this.visitaSeleccionada.id).subscribe({
        next: (response) => {
          this.showToast('success', 'Visita actualizada correctamente');
          this.showModal = false;
          this.GetAllVisitaGeneralData();
        },
        error: (err) => {
          this.showToast('error', 'Error al actualizar la visita');
        }
      });
    }
  }

  esCampoInvalido(nombreCampo: string): boolean {
    const control = this.formVisita.get(nombreCampo);
    return control !== null && control.invalid && (control.dirty || control.touched);
  }
  tieneCampoError(nombreCampo: string, tipoError: string): boolean {
    const control = this.formVisita.get(nombreCampo);
    return control !== null && control.hasError(tipoError) && (control.dirty || control.touched);
  }

  cargarDepartamentos() {
    this._subscriptions.add(
      this._managerService.geDepartamentos().subscribe({
        next: (response) => {
          if (response.data) {
            this.departamentos = response.data;
          }
        },
        error: (err) => {
          this.showToast('error', 'Error al cargar los departamentos');
        }
      })
    );
  }

  // Método para cargar datos geográficos de forma secuencial
  async cargarDatosGeograficos(visitaData: IVisitaGeneral): Promise<void> {
    try {
      // Limpiar arrays
      this.provincias = [];
      this.distritos = [];
      
      if (visitaData.departamento) {
        const departamento = this.departamentos.find(d => d.name === visitaData.departamento);
        if (departamento) {
          // Cargar provincias y esperar
          await this.cargarProvinciasAsync(departamento.id);
          
          if (visitaData.provincia) {
            const provincia = this.provincias.find(p => p.name === visitaData.provincia);
            if (provincia) {
              // Cargar distritos y esperar
              await this.cargarDistritosAsync(provincia.id);
            }
          }
        }
      }
    } catch (error) {
      console.error('Error al cargar datos geográficos:', error);
    }
  }

  cargarProvincias(departamentoId: string) {
    this._subscriptions.add(
      this._managerService.getProvinciasByDepartamento(departamentoId).subscribe({
        next: (response) => {
          if (response.data) {
            this.provincias = response.data;
          }
        },
        error: (err) => {
          this.showToast('error', 'Error al cargar las provincias');
        }
      })
    );
  }

  cargarDistritos(provinciaId: string) {
    this._subscriptions.add(
      this._managerService.getDistritosByProvincia(provinciaId).subscribe({
        next: (response) => {
          if (response.data) {
            this.distritos = response.data;
          }
        },
        error: (err) => {
          this.showToast('error', 'Error al cargar los distritos');
        }
      })
    );
  }

  // Métodos asíncronos para cargar datos geográficos
  private cargarProvinciasAsync(departamentoId: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this._subscriptions.add(
        this._managerService.getProvinciasByDepartamento(departamentoId).subscribe({
          next: (response) => {
            if (response.data) {
              this.provincias = response.data;
              resolve();
            } else {
              reject('No se pudieron cargar las provincias');
            }
          },
          error: (err) => {
            this.showToast('error', 'Error al cargar las provincias');
            reject(err);
          }
        })
      );
    });
  }

  private cargarDistritosAsync(provinciaId: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this._subscriptions.add(
        this._managerService.getDistritosByProvincia(provinciaId).subscribe({
          next: (response) => {
            if (response.data) {
              this.distritos = response.data;
              resolve();
            } else {
              reject('No se pudieron cargar los distritos');
            }
          },
          error: (err) => {
            this.showToast('error', 'Error al cargar los distritos');
            reject(err);
          }
        })
      );
    });
  }

  // Métodos para el dropdown multi-select
  toggleDropdown() {
    this.dropdownOpen = !this.dropdownOpen;
  }

  getSelectedMotivos(): string[] {
    const motivosActuales = this.formVisita.get('motivo_atencion')?.value || [];
    return motivosActuales.map((motivo: string) => {
      const motivoObj = this.motivosAtencion.find(m => m.value === motivo);
      return motivoObj ? motivoObj.label : motivo;
    });
  }

  removeMotivo(motivoLabel: string, event: Event) {
    event.stopPropagation();
    const motivoObj = this.motivosAtencion.find(m => m.label === motivoLabel);
    if (motivoObj) {
      const motivosActuales = this.formVisita.get('motivo_atencion')?.value || [];
      const index = motivosActuales.indexOf(motivoObj.value);
      if (index > -1) {
        motivosActuales.splice(index, 1);
        this.formVisita.get('motivo_atencion')?.setValue(motivosActuales);
      }
    }
  }

  onMotivoChange(motivo: string, event: any) {
    const motivosActuales = this.formVisita.get('motivo_atencion')?.value || [];
    
    if (event.target.checked) {
      if (!motivosActuales.includes(motivo)) {
        motivosActuales.push(motivo);
      }
    } else {
      const index = motivosActuales.indexOf(motivo);
      if (index > -1) {
        motivosActuales.splice(index, 1);
      }
    }
    
    this.formVisita.get('motivo_atencion')?.setValue(motivosActuales);
  }

  isMotivoSelected(motivo: string): boolean {
    const motivosActuales = this.formVisita.get('motivo_atencion')?.value || [];
    return motivosActuales.includes(motivo);
  }

  getMotivoAtencionLabel(): string {
    const motivosActuales = this.formVisita.get('motivo_atencion')?.value || [];
    if (motivosActuales.length === 0) return 'Seleccione motivos de atención';
    
    const labels = motivosActuales.map((motivo: string) => {
      const motivoObj = this.motivosAtencion.find(m => m.value === motivo);
      return motivoObj ? motivoObj.label : motivo;
    });
    
    return labels.join(', ');
  }

  private procesarDatosVisita(visitaData: IVisitaGeneral) {
    let motivosArray: string[] = [];
    if (typeof visitaData.motivo_atencion === 'string') {
      motivosArray = visitaData.motivo_atencion.split(',').map(m => m.trim()).filter(m => m.length > 0);
    } else if (Array.isArray(visitaData.motivo_atencion)) {
      motivosArray = visitaData.motivo_atencion;
    }

    return {
      ...visitaData,
      motivo_atencion: motivosArray
    };
  }

  filterMotivoAtencion(data: IVisitaGeneral[], filterValue: string): IVisitaGeneral[] {
    if (!filterValue) {
      return data;
    }
    
    return data.filter(item => {
      if (!item.motivo_atencion) {
        return false;
      }
      
      let motivosItem: string[] = [];
      if (typeof item.motivo_atencion === 'string') {
        motivosItem = item.motivo_atencion.split(',').map(m => m.trim());
      } else if (Array.isArray(item.motivo_atencion)) {
        motivosItem = item.motivo_atencion;
      }
      return motivosItem.includes(filterValue);
    });
  }

  exportToExcel(): void {
    this.isExporting = true;
    
    this._subscriptions.add(
      this._managerService.getAllVisitaGeneral().subscribe({
        next: (response) => {
          this.isExporting = false;
          
          if (!response.data || !Array.isArray(response.data) || response.data.length === 0) {
            this.showToast('info', 'No hay datos disponibles para exportar');
            return;
          }

          const dataToExport = response.data.map((visita, index) => ({
            'N°': index + 1,
            'Tipo Usuario': this.formatTipoUsuario(visita.tipo_usuario),
            'Código Estudiante': visita.codigo_estudiante || '',
            'DNI': visita.dni || '',
            'Nombre Completo': visita.nombre_completo,
            'Género': visita.genero,
            'Edad': visita.edad || '',
            'Escuela/Área': visita.escuela || visita.area || '',
            'Motivo Atención': this.formatMotivoAtencion(visita.motivo_atencion),
            'Descripción Motivo': visita.descripcion_motivo,
            'Departamento': visita.departamento || '',
            'Provincia': visita.provincia || '',
            'Distrito': visita.distrito || '',
            'Lugar Atención': visita.lugar_atencion,
            'Fecha Creación': visita.created_at ? new Date(visita.created_at).toLocaleDateString('es-PE') : '',
            'Última Actualización': visita.updated_at ? new Date(visita.updated_at).toLocaleDateString('es-PE') : ''
          }));

          const ws: XLSX.WorkSheet = XLSX.utils.json_to_sheet(dataToExport);
          const wb: XLSX.WorkBook = XLSX.utils.book_new();
          XLSX.utils.book_append_sheet(wb, ws, 'Visitas Generales');

          ws['!cols'] = [
            { wch: 5 },
            { wch: 15 },
            { wch: 15 },
            { wch: 12 },
            { wch: 35 },
            { wch: 8 },
            { wch: 8 },
            { wch: 20 },
            { wch: 20 },
            { wch: 40 },
            { wch: 15 },
            { wch: 15 },
            { wch: 15 },
            { wch: 30 },
            { wch: 12 },
            { wch: 12 }
          ];

          const fechaActual = new Date().toISOString().split('T')[0];
          const nombreArchivo = `visitas_generales_${fechaActual}`;

          const excelBuffer: any = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });
          this.saveAsExcelFile(excelBuffer, nombreArchivo);
          
          this.showToast('success', `Datos de ${response.data.length} visitas exportados correctamente`);
        },
        error: (err) => {
          this.isExporting = false;
          this.showToast('error', 'Error al exportar los datos. Intente nuevamente.');
        }
      })
    );
  }

  private formatTipoUsuario(tipo: string): string {
    const tipoObj = this.tiposUsuario.find(t => t.value === tipo);
    return tipoObj ? tipoObj.label : tipo;
  }

  private formatMotivoAtencion(motivo: string | string[]): string {
    let motivosArray: string[] = [];
    
    if (typeof motivo === 'string') {
      motivosArray = motivo.split(',').map(m => m.trim()).filter(m => m.length > 0);
    } else if (Array.isArray(motivo)) {
      motivosArray = motivo;
    }

    return motivosArray.map(m => {
      const motivoObj = this.motivosAtencion.find(ma => ma.value === m);
      return motivoObj ? motivoObj.label : m;
    }).join(', ');
  }

  private saveAsExcelFile(buffer: any, fileName: string): void {
    const EXCEL_TYPE = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8';
    const EXCEL_EXTENSION = '.xlsx';

    const data: Blob = new Blob([buffer], { type: EXCEL_TYPE });
    saveAs(data, fileName + EXCEL_EXTENSION);
  }

  // Validador personalizado para URL opcional
  validarUrlOpcional(control: AbstractControl): {[key: string]: any} | null {
    if (!control.value || control.value.trim() === '') {
      return null; // Es opcional, permite valores vacíos
    }
    
    const urlPattern = /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/;
    const isValid = urlPattern.test(control.value);
    
    return isValid ? null : { 'invalidUrl': { value: control.value } };
  }

  // Método para copiar URL al portapapeles
  copiarUrl() {
    const url = this.formVisita.get('url_imagen')?.value;
    if (url) {
      navigator.clipboard.writeText(url).then(() => {
        this.showToast('success', 'URL copiada al portapapeles');
      }).catch(() => {
        this.showToast('error', 'Error al copiar la URL');
      });
    }
  }
}
