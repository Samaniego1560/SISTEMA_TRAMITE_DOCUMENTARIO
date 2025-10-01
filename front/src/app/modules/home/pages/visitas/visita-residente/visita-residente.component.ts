import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { TableComponent, TableColumn, TableAction, TableFilter } from '../components/table/table.component';
import { VisitaResidente, CreateVisitaResidenteRequest, UpdateVisitaResidenteRequest, ResidentesPorVisitar } from '../../../../../core/models/visita-residente';
import { IAnnouncement } from '../../../../../core/models/announcement';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { Subscription } from 'rxjs';
import { ModalComponent } from "../../../../../core/ui/modal/modal.component";
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from "@angular/forms";
import { BlockUiComponent } from "../../../../../core/ui/block-ui/block-ui.component";

interface VisitaResidenteDisplay extends VisitaResidente {
  fecha_creacion_formatted: string;
  fecha_actualizacion_formatted: string;
  estado_badge: string;
  actions: TableAction[];
}

@Component({
  selector: 'app-visita-residente',
  standalone: true,
  imports: [CommonModule, TableComponent, ModalComponent, ToastComponent, FormsModule, ReactiveFormsModule, BlockUiComponent],
  templateUrl: './visita-residente.component.html',
  styleUrl: './visita-residente.component.css',
  providers: [ToastService]
})
export class VisitaResidenteComponent implements OnInit, OnDestroy {
  visitaData: VisitaResidenteDisplay[] = [];
  residentesPendientes: ResidentesPorVisitar[] = [];
  isLoading: boolean = false;
  error: string | null = null;
  private _subscriptions: Subscription = new Subscription();
  protected deleteModal: boolean = false;
  protected showModal: boolean = false;
  protected modalTitle: string = 'Registrar Visita a Residente';
  protected modalButtonText: string = 'Guardar';
  protected formVisita: FormGroup;
  protected modoFormulario: 'crear' | 'editar' | 'ver' = 'crear';
  protected visitaSeleccionada: VisitaResidente | null = null;
  protected residenteSeleccionado: ResidentesPorVisitar | null = null;
  protected convocatoriaId: string = '';
  protected convocatorias: IAnnouncement[] = [];
  protected isLoadingConvocatorias: boolean = false;
  protected activeTab: 'pendientes' | 'visitados' = 'pendientes';
  
  protected estados = [
    { value: 'verificado', label: 'Verificado', color: 'green'},
    { value: 'observado', label: 'Observado', color: 'red' }
  ];
  
  tableColumns: TableColumn[] = [
    { header: 'Código', field: 'alumno_codigo' },
    { header: 'Nombre', field: 'alumno_nombre' },
    { header: 'Escuela', field: 'escuela_profesional' },
    { header: 'Estado', field: 'estado_badge' },
    { header: 'Fecha Visita', field: 'fecha_creacion_formatted' },
   
  ];
  
  tableActions: TableAction[] = [];
  
  tablePendientesColumns: TableColumn[] = [
    { header: 'Código', field: 'codigo' },
    { header: 'Nombre', field: 'nombre' },
    { header: 'Escuela Profesional', field: 'escuela_profesional' },
    { header: 'Teléfono', field: 'celular' },
  ];
  
  tablePendientesActions: TableAction[] = [
    { label: 'Registrar Visita', icon: 'plus', action: 'register' }
  ];
  
  tablePendientesFilters: TableFilter[] = [
    {
      field: 'escuela_profesional',
      label: 'Escuela Profesional',
      options: []
    }
  ];
  
  tableFilters: TableFilter[] = [
    {
      field: 'estado',
      label: 'Estado',
      options: [
        { value: 'pendiente', label: 'Pendiente' },
        { value: 'verificado', label: 'Verificado' },
        { value: 'observado', label: 'Observado' }
      ]
    },
    {
      field: 'escuela_profesional',
      label: 'Escuela Profesional',
      options: []
    }
  ];

  constructor(
    private _managerService: ManagerService, 
    private _toastService: ToastService,
    private _fb: FormBuilder
  ) {
    this.formVisita = this._fb.group({
      alumno_id: ['', Validators.required],
      estado: ['', Validators.required],
      comentario: [''],
      imagen_url: ['']
    });
  }
 
  ngOnInit() {
    this.loadConvocatorias();
    this.getAllVisitaResidenteData();
  }
  
  showToast(type: 'success' | 'error' | 'info' | 'warning', message: string) {
    this._toastService.add({
      key: 'visita-residente',
      type: type,
      message: message
    });
  }
  
  cambiarConvocatoria(nuevaConvocatoriaId: string) {
    if (!nuevaConvocatoriaId) {
      this.showToast('warning', 'Por favor seleccione una convocatoria válida');
      return;
    }
    
    this.convocatoriaId = nuevaConvocatoriaId;
    this.residentesPendientes = [];
    this.loadResidentesPendientes();
  }
  
  onConvocatoriaChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    this.cambiarConvocatoria(target.value);
  }
  
  loadConvocatorias() {
    this.isLoadingConvocatorias = true;
    this._subscriptions.add(
      this._managerService.getAnnouncement().subscribe({
        next: res => {
          this.isLoadingConvocatorias = false;
          
          let convocatoriasData: IAnnouncement[] = [];
          
          if (res.data && Array.isArray(res.data)) {
            convocatoriasData = res.data;
          } else if (res.detalle && Array.isArray(res.detalle)) {
            convocatoriasData = res.detalle;
          } else if (Array.isArray(res)) {
            convocatoriasData = res as IAnnouncement[];
          }
          
          if(!convocatoriasData || convocatoriasData.length === 0){
            this.showToast('warning', 'No se encontraron convocatorias disponibles');
            return;
          }
          
          this.convocatorias = convocatoriasData.filter(conv => conv && conv.id && conv.nombre);
          
          if (this.convocatorias.length === 0) {
            this.showToast('warning', 'No se encontraron convocatorias válidas');
            return;
          }
          
          this.convocatoriaId = this.convocatorias[this.convocatorias.length - 1].id?.toString() || '';
          
          if (this.convocatoriaId) {
            this.loadResidentesPendientes();
          }
        },
        error: err => {
          this.showToast('error', 'Error al cargar las convocatorias');
          this.isLoadingConvocatorias = false;
        }
      })
    );
  }
  
  loadResidentesPendientes() {
    if (!this.convocatoriaId) {
      return;
    }
    
    this._subscriptions.add(
      this._managerService.getAlumnosPendientesVisitaResidente(parseInt(this.convocatoriaId)).subscribe({
        next: (response) => {
          if (response.data && Array.isArray(response.data)) {
            this.residentesPendientes = response.data;
            this.updatePendientesFilterOptions();
          } else {
            this.residentesPendientes = [];
            this.showToast('info', 'No hay residentes pendientes de visita para esta convocatoria');
          }
        },
        error: err => {
          this.showToast('error', 'Error al cargar residentes pendientes');
          this.residentesPendientes = [];
        }
      })
    );
  }
  
  getAllVisitaResidenteData() {
    this.isLoading = true;
    this.error = null;
    
    this._subscriptions.add(
      this._managerService.getAllVisitaResidente().subscribe({
        next: (response) => {
          if (!response.data || !Array.isArray(response.data)) {
            this.isLoading = false;
            this.visitaData = [];
            this.showToast('info', 'No hay visitas registradas');
            return;
          }
          
          this.visitaData = response.data.map((visita) => {
            const processedVisita: VisitaResidenteDisplay = {
              ...visita,
              fecha_creacion_formatted: visita.fecha_creacion ? 
                new Date(visita.fecha_creacion).toLocaleDateString('es-ES', {
                  year: 'numeric',
                  month: 'short',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit'
                }) : 'Sin fecha',
              fecha_actualizacion_formatted: visita.fecha_actualizacion ? 
                new Date(visita.fecha_actualizacion).toLocaleDateString('es-ES', {
                  year: 'numeric',
                  month: 'short',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit'
                }) : 'Sin fecha',
              estado_badge: this.getEstadoBadgeHtml(visita.estado),
              actions: []
            };
            
            processedVisita.actions = this.getConditionalActions(processedVisita);
            
            return processedVisita;
          });
          
          this.updateEscuelaFilterOptions();
          
          this.isLoading = false;
          
          if (this.visitaData.length === 0) {
            this.showToast('info', 'No hay visitas registradas');
          }
        },
        error: err => {
          this.error = 'Error al cargar los datos. Por favor, intente nuevamente.';
          this.showToast('error', 'Error al obtener la información de las visitas');
          this.isLoading = false;
          this.visitaData = [];
        }
      })
    );
  }
  
  onViewItem(item: VisitaResidente) {
    this.visitaSeleccionada = item;
    this.modoFormulario = 'ver';
    this.modalTitle = 'Ver Detalles de Visita';
    this.modalButtonText = 'Cerrar';
    this.formVisita.patchValue({
      alumno_id: item.alumno_id,
      estado: item.estado,
      comentario: item.comentario,
      imagen_url: item.imagen_url || ''
    });
    this.formVisita.disable();
    this.showModal = true;
  }

  onRegisterItem(residente: ResidentesPorVisitar) {
    this.residenteSeleccionado = residente;
    this.modoFormulario = 'crear';
    this.modalTitle = `Registrar Visita - ${residente.nombre}`;
    this.modalButtonText = 'Registrar Visita';
    
    this.formVisita.patchValue({
      alumno_id: residente.alumno_id,
      estado: 'pendiente',
      comentario: '',
      imagen_url: ''
    });
    this.formVisita.enable();
    this.formVisita.get('alumno_id')?.disable();
    this.showModal = true;
  }

  onEditItem(item: VisitaResidente) {
    this.visitaSeleccionada = item;
    
    this.modoFormulario = 'editar';
    this.modalTitle = 'Editar Visita a Residente';
    this.modalButtonText = 'Guardar Cambios';
    
    this.formVisita.patchValue({
      alumno_id: item.alumno_id,
      estado: item.estado || 'pendiente',
      comentario: item.comentario || '',
      imagen_url: item.imagen_url || ''
    });
    this.formVisita.enable();
    this.formVisita.get('alumno_id')?.disable();
    this.showModal = true;
  }
  
  onDeleteItem(item: VisitaResidente) {
    this.visitaSeleccionada = item;
    this.deleteModal = true;
  }
  
  deleteUser() {
    if (!this.visitaSeleccionada) {
      this.showToast('error', 'No hay visita seleccionada para eliminar');
      return;
    }
    
    this.isLoading = true;
    this.deleteModal = false;
    
    this._subscriptions.add(
      this._managerService.deleteVisitaResidente(this.visitaSeleccionada.id.toString()).subscribe({
        next: () => {
          this.isLoading = false;
          this.showToast('success', 'Visita eliminada correctamente');
          
          this.getAllVisitaResidenteData();
          if (this.convocatoriaId) {
            this.loadResidentesPendientes();
          }
          this.visitaSeleccionada = null;
        },
        error: (err) => {
          this.isLoading = false;
          this.showToast('error', 'Error al eliminar la visita');
          this.visitaSeleccionada = null;
        }
      })
    );
  }
  
  guardarVisita() {
    if (this.formVisita.invalid) {
      this.formVisita.markAllAsTouched();
      this.showToast('warning', 'Por favor complete todos los campos requeridos');
      return;
    }
    
    this.isLoading = true;
    const visitaData = this.formVisita.getRawValue();
    
    if (this.modoFormulario === 'crear') {
      const createData: CreateVisitaResidenteRequest = {
        alumno_id: Number(visitaData.alumno_id),
        estado: visitaData.estado,
        comentario: visitaData.comentario || '',
        imagen_url: visitaData.imagen_url || ''
      };
      
      this._managerService.createVisitaResidente(createData).subscribe({
        next: (response) => {
          this.isLoading = false;
          this.showToast('success', 'Visita registrada correctamente');
          this.showModal = false;
          this.getAllVisitaResidenteData();
          if (this.convocatoriaId) {
            this.loadResidentesPendientes();
          }
        },
        error: (err) => {
          this.isLoading = false;
          this.showToast('error', 'Error al registrar la visita');
        }
      });
    } else if (this.modoFormulario === 'editar' && this.visitaSeleccionada) {
      const updateData: UpdateVisitaResidenteRequest = {
        alumno_id: this.visitaSeleccionada.alumno_id,
        estado: visitaData.estado,
        comentario: visitaData.comentario || '',
        imagen_url: visitaData.imagen_url || ''
      };
      
      const visitaId = this.visitaSeleccionada.id.toString();
      
      this._managerService.updateVisitaResidente(visitaId, updateData).subscribe({
        next: (response) => {
          this.isLoading = false;
          this.showToast('success', 'Visita actualizada correctamente');
          this.showModal = false;
          this.getAllVisitaResidenteData();
          if (this.convocatoriaId) {
            this.loadResidentesPendientes();
          }
        },
        error: (err) => {
          this.isLoading = false;
          this.showToast('error', 'Error al actualizar la visita');
        }
      });
    } else {
      this.isLoading = false;
      this.showToast('error', 'Error: modo de formulario no válido');
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
  
  getEstadoClass(estado: string): string {
    switch (estado) {
      case 'pendiente':
        return 'bg-orange-100 text-orange-800';
      case 'verificado':
        return 'bg-green-100 text-green-800';
      case 'observado':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }
  
  getEstadoText(estado: string): string {
    const estadoObj = this.estados.find(e => e.value === estado);
    return estadoObj ? estadoObj.label : 'Pendiente';
  }
  
  getSelectedConvocatoria(): IAnnouncement | null {
    return this.convocatorias.find(c => c.id?.toString() === this.convocatoriaId) || null;
  }
  
  cambiarTab(tab: 'pendientes' | 'visitados') {
    this.activeTab = tab;
  }
  
  getEstadoSeleccionado() {
    const estadoValue = this.formVisita.get('estado')?.value;
    return this.estados.find(estado => estado.value === estadoValue);
  }
  
  private getEstadoBadgeHtml(estado: string): string {
    switch (estado) {
      case 'verificado':
        return `<div class="inline px-3 py-1 text-sm font-medium select-none rounded-full text-emerald-700 gap-x-2 bg-emerald-100 border border-emerald-200">
          Verificado</div>`;
      case 'observado':
        return `<div class="inline px-3 py-1 text-sm font-medium select-none rounded-full text-red-700 gap-x-2 bg-red-100 border border-red-200">
          Observado</div>`;
      case 'pendiente':
      default:
        return `<div class="inline px-3 py-1 text-sm font-medium select-none rounded-full text-orange-700 gap-x-2 bg-orange-100 border border-orange-200">
          Pendiente</div>`;
    }
  }
  
  private getConditionalActions(visita: VisitaResidente): TableAction[] {
    const actions: TableAction[] = [];
    
    actions.push({ label: 'Ver', icon: 'eye', action: 'view' });
    actions.push({ label: 'Editar', icon: 'pencil', action: 'edit' });
    actions.push({ label: 'Eliminar', icon: 'trash', action: 'delete' });
    
    return actions;
  }
  
  private updateEscuelaFilterOptions() {
    if (!this.visitaData || this.visitaData.length === 0) {
      return;
    }
    
    const escuelas = [...new Set(this.visitaData.map(visita => visita.escuela_profesional).filter(Boolean))];
    const escuelaFilter = this.tableFilters.find(filter => filter.field === 'escuela_profesional');
    if (escuelaFilter) {
      escuelaFilter.options = escuelas.map(escuela => ({
        value: escuela,
        label: escuela
      }));
    }
  }

  private updatePendientesFilterOptions() {
    if (!this.residentesPendientes || this.residentesPendientes.length === 0) {
      return;
    }
    
    // Actualizar opciones del filtro de escuela profesional
    const escuelas = [...new Set(this.residentesPendientes.map(residente => residente.escuela_profesional).filter(Boolean))];
    const escuelaFilter = this.tablePendientesFilters.find(filter => filter.field === 'escuela_profesional');
    if (escuelaFilter) {
      escuelaFilter.options = escuelas.map(escuela => ({
        value: escuela,
        label: escuela
      }));
    }
    
    // Actualizar opciones del filtro de lugar de procedencia
    const lugares = [...new Set(this.residentesPendientes.map(residente => residente.lugar_procedencia).filter(Boolean))];
    const lugarFilter = this.tablePendientesFilters.find(filter => filter.field === 'lugar_procedencia');
    if (lugarFilter) {
      lugarFilter.options = lugares.map(lugar => ({
        value: lugar,
        label: lugar
      }));
    }
  }

  ngOnDestroy() {
    this._subscriptions.unsubscribe();
  }
}
