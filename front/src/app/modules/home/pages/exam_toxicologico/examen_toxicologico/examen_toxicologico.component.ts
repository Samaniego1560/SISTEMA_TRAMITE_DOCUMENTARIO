import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Subscription } from 'rxjs';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { ModalComponent } from '../../../../../core/ui/modal/modal.component';
import { BlockUiComponent } from '../../../../../core/ui/block-ui/block-ui.component';
import { ExamenToxicologico, Estado, CreateExamenToxicologicoRequest, UpdateExamenToxicologicoRequest } from '../../../../../core/models/exam_toxicologico';
import { IAnnouncement } from '../../../../../core/models/announcement';
import * as XLSX from 'xlsx';
import { saveAs } from 'file-saver';
import { TableComponent } from '../../visitas/components/table/table.component';

interface TableColumn {
  header: string;
  field: string;
}

interface TableAction {
  label: string;
  icon: string;
  action: string;
}

interface ExamenToxicologicoDisplay extends ExamenToxicologico {
  nombre_completo: string;
  fecha_examen_formatted: string;
  estado_badge: string;
  actions: TableAction[];
}

@Component({
  selector: 'app-examen-toxicologico',
  standalone: true,
  imports: [CommonModule,TableComponent, ModalComponent, ToastComponent, FormsModule, ReactiveFormsModule, BlockUiComponent],
  templateUrl: './examen_toxicologico.component.html',
  styleUrl: './examen_toxicologico.component.css',
  providers: [ToastService]
})
export class ExamenToxicologicoComponent implements OnInit, OnDestroy { 
  examenData: ExamenToxicologicoDisplay[] = [];
  isLoading: boolean = false;
  error: string | null = null;
  private _subscriptions: Subscription = new Subscription();
  
  protected estadisticas = {
    total: 0,
    registrados: 0,
    sinRegistrar: 0,
    porcentajeCompletado: 0
  };
  protected deleteModal: boolean = false;
  protected showModal: boolean = false;
  protected modalTitle: string = 'Registrar Examen Toxicológico';
  protected modalButtonText: string = 'Guardar';
  protected formExamen: FormGroup;
  protected modoFormulario: 'crear' | 'editar' | 'ver' = 'crear';
  protected examenSeleccionado: ExamenToxicologico | null = null;
  protected examenAEliminar: ExamenToxicologico | null = null;
  protected convocatoriaId: string = '';
  protected convocatorias: IAnnouncement[] = [];
  protected isLoadingConvocatorias: boolean = false;
  
  cambiarConvocatoria(nuevaConvocatoriaId: string) {
    if (!nuevaConvocatoriaId) {
      this.showToast('warning', 'Por favor seleccione una convocatoria válida');
      this.examenData = []; 
      this.calcularEstadisticas(); 
      return;
    }
    this.convocatoriaId = nuevaConvocatoriaId;
    this.formExamen.patchValue({
      convocatoria_id: this.convocatoriaId
    });
    this.getAllExamenToxicologicoData();
  }
  
  onConvocatoriaChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    this.cambiarConvocatoria(target.value);
  }
  
  protected estados = [
    { value: Estado.Observado, label: 'Observado' },
    { value: Estado.Verificado, label: 'Verificado' }
  ];
  
  tableColumns: TableColumn[] = [
    { header: 'Código Estudiante', field: 'codigo_estudiante' },
    { header: 'Nombre Completo', field: 'nombre_completo' },
    { header: 'Escuela Profesional', field: 'escuela_profesional' },
    { header: 'Estado', field: 'estado_badge' },
    { header: 'Fecha Examen', field: 'fecha_examen_formatted' }
  ];
  
  tableActions: TableAction[] = [];
  
  tableFilters: {
    field: string;
    label: string;
    options: { value: string; label: string; }[];
  }[] = [
    {
      field: 'estado',
      label: 'Estado',
      options: [
        { value: Estado.Pendiente, label: 'Pendiente' },
        { value: Estado.Observado, label: 'Observado' },
        { value: Estado.Verificado, label: 'Verificado' }
      ]
    },
    {
      field: 'escuela_profesional',
      label: 'Escuela Profesional',
      options: []
    }
  ];

  private mapEstadoFromApi(apiEstado: string): Estado {
    switch (apiEstado?.toLowerCase()) {
      case 'observado':
        return Estado.Observado;
      case 'verificado':
        return Estado.Verificado;
      case 'sin registro':
      case 'pendiente':
      default:
        return Estado.Pendiente;
    }
  }

  constructor(
    private _managerService: ManagerService, 
    private _toastService: ToastService,
    private _fb: FormBuilder
  ) {
    this.formExamen = this._fb.group({
      alumno_id: ['', Validators.required],
      convocatoria_id: ['', Validators.required],
      estado: ['', Validators.required],
      comentario: ['']
    });
  }
 
  ngOnInit() {
    this.loadConvocatorias();
  }
  
  showToast(type: 'success' | 'error' | 'info' | 'warning', message: string) {
    this._toastService.add({
      key: 'examen-toxicologico',
      type: type,
      message: message
    });
  }

  exportToExcel(): void {
    if (!this.convocatoriaId) {
      this.showToast('warning', 'No hay convocatoria seleccionada para exportar');
      return;
    }

    this.isLoading = true;
    
    this._subscriptions.add(
      this._managerService.getExamenToxicologicoByConvocatoria(this.convocatoriaId).subscribe({
        next: (res) => {
          this.isLoading = false;
          
          if (!res.data || !Array.isArray(res.data) || res.data.length === 0) {
            this.showToast('info', 'No hay datos para exportar');
            return;
          }

          const dataToExport = res.data.map((examen, index) => ({
            'N°': index + 1,
            'Código Estudiante': examen.codigo_estudiante,
            'Nombres': examen.nombres,
            'Apellido Paterno': examen.apellido_paterno,
            'Apellido Materno': examen.apellido_materno,
            'Nombre Completo': `${examen.nombres} ${examen.apellido_paterno} ${examen.apellido_materno}`.trim(),
            'Escuela Profesional': examen.escuela_profesional   ,
            'Estado': examen.estado,
            'Comentario': examen.comentario || '',
            'Fecha Examen': examen.fecha_examen ? new Date(examen.fecha_examen).toLocaleDateString('es-ES', {
              year: 'numeric',
              month: '2-digit',
              day: '2-digit',
              hour: '2-digit',
              minute: '2-digit'
            }) : 'Sin fecha',
            'Usuario': examen.usuario_nombre || ''
          }));
          const ws: XLSX.WorkSheet = XLSX.utils.json_to_sheet(dataToExport);
          const wb: XLSX.WorkBook = XLSX.utils.book_new();
          XLSX.utils.book_append_sheet(wb, ws, 'Exámenes Toxicológicos');

          // Configurar ancho de columnas
          ws['!cols'] = [
            { wch: 5 },  // N°
            { wch: 15 }, // Código Estudiante
            { wch: 20 }, // Nombres
            { wch: 20 }, // Apellido Paterno
            { wch: 20 }, // Apellido Materno
            { wch: 35 }, // Nombre Completo
            { wch: 40 }, // Escuela Profesional
            { wch: 12 }, // Estado
            { wch: 30 }, // Comentario
            { wch: 20 }, // Fecha Examen
            { wch: 15 }  // Usuario
          ];

          // Generar nombre del archivo
          const convocatoria = this.getSelectedConvocatoria();
          const nombreConvocatoria = convocatoria?.nombre || 'Convocatoria';
          const fechaActual = new Date().toISOString().split('T')[0];
          const nombreArchivo = `examenes_toxicologicos_${nombreConvocatoria.replace(/\s+/g, '_')}_${fechaActual}`;

          // Descargar archivo
          const excelBuffer: any = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });
          this.saveAsExcelFile(excelBuffer, nombreArchivo);
          
          this.showToast('success', 'Archivo exportado correctamente');
        },
        error: (err) => {
          this.isLoading = false;
          this.showToast('error', 'Error al exportar los datos. Intente nuevamente.');
        }
      })
    );
  }

  private saveAsExcelFile(buffer: any, fileName: string): void {
    const EXCEL_TYPE = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8';
    const EXCEL_EXTENSION = '.xlsx';

    const data: Blob = new Blob([buffer], { type: EXCEL_TYPE });
    saveAs(data, fileName + EXCEL_EXTENSION);
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
          
          this.convocatoriaId = '';
        },
        error: err => {
          this.showToast('error', 'Error al cargar las convocatorias');
          this.isLoadingConvocatorias = false;
        }
      })
    );
  }
  
  getAllExamenToxicologicoData() {
    if (!this.convocatoriaId) {
      this.showToast('warning', 'No hay convocatoria seleccionada');
      return;
    }
    
    this.isLoading = true;
    this.error = null;
    
    this._subscriptions.add(
      this._managerService.getExamenToxicologicoByConvocatoria(this.convocatoriaId).subscribe({
        next: (displayDataRes) => {
          if (!displayDataRes.data || !Array.isArray(displayDataRes.data) || displayDataRes.data.length === 0) {
            this.isLoading = false;
            this.examenData = [];
            this.calcularEstadisticas();
            this.showToast('info', 'No hay estudiantes registrados para esta convocatoria');
            return;
          }
          
          this._managerService.getAllExamenToxicologico().subscribe({
            next: (allExamsRes) => {
              const examIdMap = new Map<number, any>();
              if (allExamsRes.data && Array.isArray(allExamsRes.data)) {
                allExamsRes.data.forEach(exam => {
                  if (exam.alumno_id && exam.convocatoria_id && 
                      exam.convocatoria_id.toString() === this.convocatoriaId) {
                    examIdMap.set(exam.alumno_id, exam);
                  }
                });
              }
              
              this.examenData = (displayDataRes.data || []).map((displayExam, index) => {
                const examWithId = examIdMap.get(displayExam.alumno_id);
                
                const processedExam: ExamenToxicologicoDisplay = {
                  ...displayExam,
                  id: examWithId?.id || undefined,
                  fecha_examen: displayExam.fecha_examen ? new Date(displayExam.fecha_examen) : null,
                  nombre_completo: `${displayExam.nombres} ${displayExam.apellido_paterno} ${displayExam.apellido_materno}`.trim(),
                  fecha_examen_formatted: displayExam.fecha_examen ? 
                    new Date(displayExam.fecha_examen).toLocaleDateString('es-ES', {
                      year: 'numeric',
                      month: 'short',
                      day: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit'
                    }) : 'Sin fecha',
                  estado: this.mapEstadoFromApi(displayExam.estado as string),
                  estado_badge: this.getEstadoBadgeHtml(this.mapEstadoFromApi(displayExam.estado as string)),
                  convocatoria_id: parseInt(this.convocatoriaId),
                  actions: []
                };
                
                processedExam.actions = this.getConditionalActions(processedExam);
                
                return processedExam;
              });
              
              this.updateEscuelaProfesionalFilterOptions();
              this.calcularEstadisticas();
              
              this.isLoading = false;
            
            },
            error: (err) => {
              this.isLoading = false;
              this.examenData = (displayDataRes.data || []).map(displayExam => {
                const processedExam: ExamenToxicologicoDisplay = {
                  ...displayExam,
                  id: undefined,
                  nombre_completo: `${displayExam.nombres} ${displayExam.apellido_paterno} ${displayExam.apellido_materno}`.trim(),
                  fecha_examen: displayExam.fecha_examen ? new Date(displayExam.fecha_examen) : null,
                  fecha_examen_formatted: displayExam.fecha_examen ? 
                    new Date(displayExam.fecha_examen).toLocaleDateString('es-ES', {
                      year: 'numeric',
                      month: 'short', 
                      day: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit'
                    }) : 'Sin fecha',
                  estado: this.mapEstadoFromApi(displayExam.estado as string),
                  estado_badge: this.getEstadoBadgeHtml(this.mapEstadoFromApi(displayExam.estado as string)),
                  convocatoria_id: parseInt(this.convocatoriaId),
                  actions: []
                };
                
                processedExam.actions = this.getConditionalActions(processedExam);
                return processedExam;
              });
              
              this.updateEscuelaProfesionalFilterOptions();
              this.calcularEstadisticas();
            }
          });
        },
        error: err => {
          this.error = 'Error al cargar los datos. Por favor, intente nuevamente.';
          this.showToast('error', 'Error al obtener la información de los exámenes toxicológicos');
          this.isLoading = false;
          this.examenData = [];
          this.calcularEstadisticas();
        }
      })
    );
  }
  
  onViewItem(item: ExamenToxicologico) {
    this.examenSeleccionado = item;
    this.modoFormulario = 'ver';
    this.modalTitle = 'Ver Detalles de Examen Toxicológico';
    this.modalButtonText = 'Cerrar';
    this.formExamen.patchValue({
      alumno_id: item.alumno_id,
      convocatoria_id: this.convocatoriaId,
      estado: item.estado,
      comentario: item.comentario
    });
    this.formExamen.disable();
    this.showModal = true;
  }

  onRegisterItem(item: ExamenToxicologico) {
    this.examenSeleccionado = item;
    this.modoFormulario = 'crear';
    this.modalTitle = 'Registrar Examen Toxicológico';
    this.modalButtonText = 'Registrar';
    
    this.formExamen.patchValue({
      alumno_id: item.alumno_id,
      convocatoria_id: this.convocatoriaId,
      estado: Estado.Pendiente,
      comentario: ''
    });
    this.formExamen.enable();
    this.formExamen.get('alumno_id')?.disable();
    this.formExamen.get('convocatoria_id')?.disable();
    this.showModal = true;
  }

  onEditItem(item: ExamenToxicologico) {
    this.examenSeleccionado = item;
    
    this.modoFormulario = 'editar';
    this.modalTitle = 'Editar Examen Toxicológico';
    this.modalButtonText = 'Guardar Cambios';
    
    this.formExamen.patchValue({
      alumno_id: item.alumno_id,
      convocatoria_id: this.convocatoriaId,
      estado: item.estado || Estado.Pendiente,
      comentario: item.comentario || ''
    });
    this.formExamen.enable();
    this.formExamen.get('alumno_id')?.disable();
    this.formExamen.get('convocatoria_id')?.disable();
    this.showModal = true;
  }
  
  onDeleteItem(item: ExamenToxicologico) {
    this.examenAEliminar = item;
    this.deleteModal = true;
  }

  confirmarEliminacion() {
    if (!this.examenAEliminar || !this.examenAEliminar.id) {
      this.showToast('error', 'Error: No se puede eliminar un examen sin ID válido');
      return;
    }

    this.isLoading = true;
    const examId = this.examenAEliminar.id.toString();

    this._subscriptions.add(
      this._managerService.deleteExamenToxicologico(examId).subscribe({
        next: (response) => {
          this.isLoading = false;
          this.deleteModal = false;
          this.examenAEliminar = null;
          this.showToast('success', 'Examen toxicológico eliminado correctamente');
          this.getAllExamenToxicologicoData();
        },
        error: (err) => {
          this.isLoading = false;
          this.showToast('error', 'Error al eliminar el examen toxicológico');
        }
      })
    );
  }

  cancelarEliminacion() {
    this.deleteModal = false;
    this.examenAEliminar = null;
  }
  
  guardarExamen() {
    const estado = this.formExamen.get('estado')?.value;
    const comentario = this.formExamen.get('comentario')?.value?.trim();
    
    if ((estado === Estado.Observado || estado === Estado.Verificado) && !comentario) {
      this.showToast('warning', 'El comentario es requerido para estados Observado y Verificado');
      return;
    }
    
    if (this.formExamen.invalid) {
      this.formExamen.markAllAsTouched();
      this.showToast('warning', 'Por favor complete todos los campos requeridos');
      return;
    }
    
    this.isLoading = true;
    const examenData = this.formExamen.getRawValue();
    
    if (this.modoFormulario === 'crear') {
      const createData: CreateExamenToxicologicoRequest = {
        alumno_id: Number(examenData.alumno_id),
        convocatoria_id: Number(examenData.convocatoria_id),
        estado: examenData.estado,
        comentario: examenData.comentario || ''
      };
      
      this._managerService.createOrUpdateExamenToxicologico(createData).subscribe({
        next: (response) => {
          this.isLoading = false;
          this.showToast('success', 'Examen toxicológico verificado correctamente');
          this.showModal = false;
          this.getAllExamenToxicologicoData();
        },
        error: (err) => {
          this.isLoading = false;
          this.showToast('error', 'Error al verificar el examen toxicológico');
        }
      });
    } else if (this.modoFormulario === 'editar' && this.examenSeleccionado) {
      if (!this.examenSeleccionado.id) {
        this.showToast('error', 'Error: No se puede editar un examen sin ID válido');
        this.isLoading = false;
        return;
      }
      
      const updateData: UpdateExamenToxicologicoRequest = {
        estado: examenData.estado,
        comentario: examenData.comentario || ''
      };
      
      const examId = this.examenSeleccionado.id.toString();
      
      this._managerService.updateExamenToxicologico(examId, updateData).subscribe({
        next: (response) => {
          this.isLoading = false;
          this.showToast('success', 'Examen toxicológico actualizado correctamente');
          this.showModal = false;
          this.getAllExamenToxicologicoData();
        },
        error: (err) => {
          this.isLoading = false;
          this.showToast('error', 'Error al actualizar el examen toxicológico');
        }
      });
    } else {
      this.isLoading = false;
      this.showToast('error', 'Error: modo de formulario no válido');
    }
  }
  
  esCampoInvalido(nombreCampo: string): boolean {
    const control = this.formExamen.get(nombreCampo);
    return control !== null && control.invalid && (control.dirty || control.touched);
  }
  
  tieneCampoError(nombreCampo: string, tipoError: string): boolean {
    const control = this.formExamen.get(nombreCampo);
    return control !== null && control.hasError(tipoError) && (control.dirty || control.touched);
  }
  
  getEstadoClass(estado: Estado): string {
    switch (estado) {
      case Estado.Pendiente:
        return 'bg-yellow-100 text-yellow-800';
      case Estado.Observado:
        return 'bg-red-100 text-red-800';
      case Estado.Verificado:
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }
  
  getEstadoText(estado: Estado): string {
    const estadoObj = this.estados.find(e => e.value === estado);
    return estadoObj ? estadoObj.label : 'Pendiente';
  }
  
  getSelectedConvocatoria(): IAnnouncement | null {
    return this.convocatorias.find(c => c.id?.toString() === this.convocatoriaId) || null;
  }
  
  private getEstadoBadgeHtml(estado: Estado): string {
    switch (estado) {
      case Estado.Verificado:
        return `<div class="inline px-3 py-1 text-sm font-normal select-none rounded-full text-emerald-500 gap-x-2 bg-emerald-100/60">Verificado</div>`;
      case Estado.Observado:
        return `<div class="inline px-3 py-1 text-sm font-normal select-none rounded-full text-red-500 gap-x-2 bg-red-100/60">Observado</div>`;
      case Estado.Pendiente:
      default:
        return `<div class="inline px-3 py-1 text-sm font-normal select-none rounded-full text-orange-500 gap-x-2 bg-orange-100/60">Pendiente</div>`;
    }
  }
  
  private getConditionalActions(examen: ExamenToxicologico): TableAction[] {
    const actions: TableAction[] = [];
    
    actions.push({ label: 'Ver', icon: 'eye', action: 'view' });
    
    if (examen.id) {
      actions.push({ label: 'Editar', icon: 'pencil', action: 'edit' });
      actions.push({ label: 'Eliminar', icon: 'trash', action: 'delete' });
    } else {
      actions.push({ label: 'Registrar', icon: 'plus', action: 'register' });
    }
    
    return actions;
  }

  private updateEscuelaProfesionalFilterOptions() {
    if (!this.examenData || this.examenData.length === 0) {
      return;
    }

    const escuelas = [...new Set(this.examenData.map(exam => exam.escuela_profesional).filter(Boolean))];
    const escuelaFilter = this.tableFilters.find(filter => filter.field === 'escuela_profesional');
    if (escuelaFilter) {
      escuelaFilter.options = escuelas.map(escuela => ({
        value: escuela,
        label: escuela
      }));
    }
  }

  private calcularEstadisticas() {
    if (!this.examenData || this.examenData.length === 0) {
      this.estadisticas = {
        total: 0,
        registrados: 0,
        sinRegistrar: 0,
        porcentajeCompletado: 0
      };
      return;
    }

    const total = this.examenData.length;
    const registrados = this.examenData.filter(exam => 
      exam.estado === Estado.Verificado || exam.estado === Estado.Observado
    ).length;
    const sinRegistrar = this.examenData.filter(exam => 
      exam.estado === Estado.Pendiente
    ).length;
    const porcentajeCompletado = total > 0 ? Math.round((registrados / total) * 100) : 0;

    this.estadisticas = {
      total,
      registrados,
      sinRegistrar,
      porcentajeCompletado
    };
  }

  ngOnDestroy() {
    this._subscriptions.unsubscribe();
  }
}
