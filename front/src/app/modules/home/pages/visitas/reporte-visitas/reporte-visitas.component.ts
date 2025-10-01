import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { ReporteVisitaResidente, LugarProcedencia, VisitasPorDepartamento, ExportDataExcel, ReportePorEscuelaProfesional, ResidentesPorVisitar } from '../../../../../core/models/visita-residente';
import { IAnnouncement } from '../../../../../core/models/announcement';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { Subscription } from 'rxjs';
import { ToastComponent } from '../../../../../core/ui/toast/toast.component';
import { BlockUiComponent } from "../../../../../core/ui/block-ui/block-ui.component";
import { ModalComponent } from '../../../../../core/ui/modal/modal.component';
import { FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import * as XLSX from 'xlsx';
import { saveAs } from 'file-saver';
import {
  ApexAxisChartSeries,
  ApexChart,
  ApexDataLabels,
  ApexFill,
  ApexLegend,
  ApexNonAxisChartSeries,
  ApexPlotOptions,
  ApexResponsive,
  ApexStroke,
  ApexTitleSubtitle,
  ApexTooltip,
  ApexXAxis,
  ApexYAxis,
  ChartComponent,
  NgApexchartsModule
} from "ng-apexcharts";

export type ChartOptions = {
  series: ApexNonAxisChartSeries;
  chart: ApexChart;
  responsive: ApexResponsive[];
  labels: any;
  title: ApexTitleSubtitle;
  colors: string[];
};

export type ChartOptionsBarra = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  dataLabels: ApexDataLabels;
  plotOptions: ApexPlotOptions;
  yaxis: ApexYAxis;
  xaxis: ApexXAxis;
  fill: ApexFill;
  tooltip: ApexTooltip;
  stroke: ApexStroke;
  legend: ApexLegend;
  title: ApexTitleSubtitle;
  colors: string[];
};

@Component({
  selector: 'app-reporte-visitas',
  standalone: true,
  imports: [CommonModule, ToastComponent, BlockUiComponent, NgApexchartsModule, ModalComponent, ReactiveFormsModule],
  templateUrl: './reporte-visitas.component.html',
  styleUrl: './reporte-visitas.component.css',
  providers: [ToastService]
})
export class ReporteVisitasComponent implements OnInit, OnDestroy {
  @ViewChild("chartDonut") chartDonut!: ChartComponent;
  @ViewChild("chartBar") chartBar!: ChartComponent;
  @ViewChild("chartProcedencia") chartProcedencia!: ChartComponent;

  reporteData: ReporteVisitaResidente | null = null;
  ReportePorEscuelaProfesional: ReportePorEscuelaProfesional[] = [];
  reportePorProcedencia: LugarProcedencia[] = [];
  residentesPendientes: ResidentesPorVisitar[] = [];
  isLoading: boolean = false;
  error: string | null = null;
  private _subscriptions: Subscription = new Subscription();
  protected convocatoriaId: string = '';
  protected convocatorias: IAnnouncement[] = [];
  protected isLoadingConvocatorias: boolean = false;

  // Variables para modal de exportación
  protected showExportModal: boolean = false;
  protected departamentoControl: FormControl = new FormControl('', Validators.required);
  protected departamentosDisponibles: string[] = [];
  protected isExporting: boolean = false;

  // Chart configurations
  public chartOptionsDonut: ChartOptions = {
    series: [],
    chart: {
      type: "donut",
      width: 380,
      height: 300
    },
    labels: ["Pendientes", "Verificadas", "Observadas"],
    colors: ["#F59E0B", "#10B981", "#EF4444"],
    title: {
      text: "Estado de Visitas",
      align: "center",
      style: {
        fontSize: "16px",
        fontWeight: "600"
      }
    },
    responsive: [
      {
        breakpoint: 480,
        options: {
          chart: {
            width: 300
          },
          legend: {
            position: "bottom"
          }
        }
      }
    ]
  };

  public chartOptionsBar: ChartOptionsBarra = {
    series: [
      {
        name: "Estudiantes Visitados",
        data: []
      }
    ],
    chart: {
      type: "bar",
      height: 350,
      toolbar: {
        show: true
      }
    },
    colors: ["#3B82F6"],
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: "55%"
      }
    },
    dataLabels: {
      enabled: true
    },
    stroke: {
      show: true,
      width: 2,
      colors: ["transparent"]
    },
    xaxis: {
      categories: []
    },
    yaxis: {
      title: {
        text: "Cantidad de Estudiantes Visitados"
      }
    },
    fill: {
      opacity: 1
    },
    tooltip: {
      y: {
        formatter: function (val) {
          return val + " estudiantes visitados"
        }
      }
    },
    title: {
      text: "Visitas por Escuela Profesional",
      align: "center",
      style: {
        fontSize: "16px",
        fontWeight: "600"
      }
    },
    legend: {
      show: false
    }
  };

  public chartOptionsProcedencia: ChartOptionsBarra = {
    series: [
      {
        name: "Estudiantes Visitados",
        data: []
      }
    ],
    chart: {
      type: "bar",
      height: 350,
      toolbar: {
        show: true
      }
    },
    colors: ["#10B981"],
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: "55%"
      }
    },
    dataLabels: {
      enabled: true
    },
    stroke: {
      show: true,
      width: 2,
      colors: ["transparent"]
    },
    xaxis: {
      categories: []
    },
    yaxis: {
      title: {
        text: "Cantidad de Estudiantes Visitados"
      }
    },
    fill: {
      opacity: 1
    },
    tooltip: {
      y: {
        formatter: function (val) {
          return val + " estudiantes visitados"
        }
      }
    },
    title: {
      text: "Visitas por Lugar de Departamento",
      align: "center",
      style: {
        fontSize: "16px",
        fontWeight: "600"
      }
    },
    legend: {
      show: false
    }
  };

  constructor(
    private _managerService: ManagerService,
    private _toastService: ToastService
  ) {}

  ngOnInit() {
    this.loadConvocatorias();
    this.loadEscuelaProfesionalData();
    this.loadProcedenciaData();
    this.loadReporteData();
  }

  loadEscuelaProfesionalData() {
    this._subscriptions.add(
      this._managerService.getReporteVisitaResidentePorEscuelaProfesional().subscribe({
        next: (response) => {
          if (response.data && Array.isArray(response.data)) {
            this.ReportePorEscuelaProfesional = response.data;
            this.updatePorEscuelaProfesionalChart();
          } else {
            this.ReportePorEscuelaProfesional = [];
          }
        },
        error: (err) => {
          this.ReportePorEscuelaProfesional = [];
        }
      })
    );
  }

  loadProcedenciaData() {
    this._subscriptions.add(
      this._managerService.getLugarProcedencia().subscribe({
        next: (response) => {
          if (response.data && Array.isArray(response.data)) {
            this.reportePorProcedencia = response.data;
            this.updateProcedenciaChart();
          } else {
            this.reportePorProcedencia = [];
          }
        },
        error: (err) => {
          this.reportePorProcedencia = [];
        }
      })
    );
  }

  showToast(type: 'success' | 'error' | 'info' | 'warning', message: string) {
    this._toastService.add({
      key: 'reporte-visitas',
      type: type,
      message: message
    });
  }

  openExportModal(): void {
    if (!this.convocatoriaId) {
      this.showToast('warning', 'Por favor seleccione una convocatoria primero');
      return;
    }
      this.departamentosDisponibles = [
        'AMAZONAS', 'ANCASH', 'APURIMAC', 'AREQUIPA', 'AYACUCHO', 'CAJAMARCA',
        'CALLAO', 'CUSCO', 'HUANCAVELICA', 'HUANUCO', 'ICA', 'JUNIN',
        'LA LIBERTAD', 'LAMBAYEQUE', 'LIMA', 'LORETO', 'MADRE DE DIOS',
        'MOQUEGUA', 'PASCO', 'PIURA', 'PUNO', 'SAN MARTIN', 'TACNA',
        'TUMBES', 'UCAYALI'
      ].sort();
    

    if (this.departamentosDisponibles.length === 0) {
      this.showToast('warning', 'No hay departamentos disponibles para exportar');
      return;
    }
    this.departamentoControl.setValue('');
    this.showExportModal = true;
  }

  closeExportModal(): void {
    this.showExportModal = false;
    this.departamentoControl.setValue('');
  }

  exportToExcel(): void {
    if (this.departamentoControl.invalid) {
      this.showToast('warning', 'Por favor seleccione un departamento');
      return;
    }

    const departamentoSeleccionado = this.departamentoControl.value;
    if (!departamentoSeleccionado || !this.convocatoriaId) {
      this.showToast('error', 'Error: Faltan datos para la exportación');
      return;
    }

    this.isExporting = true;
    
    this._subscriptions.add(
      this._managerService.getAllVisitasPorDepartamento(
        parseInt(this.convocatoriaId), 
        departamentoSeleccionado
      ).subscribe({
        next: (response) => {
          this.isExporting = false;
          
          if (!response.data || !Array.isArray(response.data) || response.data.length === 0) {
            this.showToast('info', `No hay alumnos sin visita para el departamento: ${departamentoSeleccionado}`);
            this.closeExportModal();
            return;
          }

          // Mapear los datos para Excel
          const dataToExport = response.data.map((visita, index) => ({
            'N°': index + 1,
            'Código Estudiante': visita.codigo,
            'Nombre Completo': visita.nombre,
            'Escuela Profesional': visita.escuela_profesional,
            'Departamento': visita.departamento,
            'Provincia': visita.provincia,
            'Distrito': visita.distrito,
            'Dirección': visita.direccion || '',
            'Celular': visita.celular || '',
            'Convocatoria': visita.convocatoria_nombre
          }));

          // Crear workbook y worksheet
          const ws: XLSX.WorkSheet = XLSX.utils.json_to_sheet(dataToExport);
          const wb: XLSX.WorkBook = XLSX.utils.book_new();
          XLSX.utils.book_append_sheet(wb, ws, 'Alumnos Sin Visita');

          // Configurar ancho de columnas
          ws['!cols'] = [
            { wch: 5 },  // N°
            { wch: 15 }, // Código Estudiante
            { wch: 35 }, // Nombre Completo
            { wch: 30 }, // Escuela Profesional
            { wch: 15 }, // Departamento
            { wch: 20 }, // Provincia
            { wch: 20 }, // Distrito
            { wch: 30 }, // Dirección
            { wch: 15 }, // Celular
            { wch: 25 }  // Convocatoria
          ];

          // Generar nombre del archivo
          const convocatoria = this.getSelectedConvocatoria();
          const nombreConvocatoria = convocatoria?.nombre || 'Convocatoria';
          const fechaActual = new Date().toISOString().split('T')[0];
          const nombreArchivo = `alumnos_sin_visita_${departamentoSeleccionado.replace(/\s+/g, '_')}_${nombreConvocatoria.replace(/\s+/g, '_')}_${fechaActual}`;

          // Descargar archivo
          const excelBuffer: any = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });
          this.saveAsExcelFile(excelBuffer, nombreArchivo);
          
          this.showToast('success', `Archivo de alumnos sin visita exportado correctamente para ${departamentoSeleccionado}`);
          this.closeExportModal();
        },
        error: (err) => {
          this.isExporting = false;
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

  exportCompletesDataToExcel(): void {
    if (!this.convocatoriaId) {
      this.showToast('warning', 'Por favor seleccione una convocatoria primero');
      return;
    }

    this.isLoading = true;
    
    this._subscriptions.add(
      this._managerService.getExportDataExcel(parseInt(this.convocatoriaId)).subscribe({
        next: (response) => {
          this.isLoading = false;
          
          if (!response.data || !Array.isArray(response.data) || response.data.length === 0) {
            this.showToast('info', 'No hay datos disponibles para exportar en esta convocatoria');
            return;
          }

          // Mapear los datos para Excel
          const dataToExport = response.data.map((alumno, index) => ({
            'N°': index + 1,
            'Código Estudiante': alumno.codigo,
            'Nombre Completo': alumno.nombre,
            'DNI': alumno.dni,
            'Celular Estudiante': alumno.celular || '',
            'Celular Padre/Madre': alumno.celular_padre || '',
            'Dirección': alumno.direccion || '',
            'Escuela Profesional': alumno.escuela_profesional,
            'Departamento': alumno.departamento,
            'Provincia': alumno.provincia,
            'Distrito': alumno.distrito,
            'ID Solicitud': alumno.solicitud_id,
            'Convocatoria': alumno.convocatoria_nombre,
            'Estado de Visita': alumno.estado_visita
          }));

          // Crear workbook y worksheet
          const ws: XLSX.WorkSheet = XLSX.utils.json_to_sheet(dataToExport);
          const wb: XLSX.WorkBook = XLSX.utils.book_new();
          XLSX.utils.book_append_sheet(wb, ws, 'Datos Completos Alumnos');

          // Configurar ancho de columnas
          ws['!cols'] = [
            { wch: 5 },  // N°
            { wch: 15 }, // Código Estudiante
            { wch: 35 }, // Nombre Completo
            { wch: 12 }, // DNI
            { wch: 15 }, // Celular Estudiante
            { wch: 15 }, // Celular Padre/Madre
            { wch: 30 }, // Dirección
            { wch: 30 }, // Escuela Profesional
            { wch: 15 }, // Departamento
            { wch: 20 }, // Provincia
            { wch: 20 }, // Distrito
            { wch: 12 }, // ID Solicitud
            { wch: 25 }, // Convocatoria
            { wch: 15 }  // Estado de Visita
          ];

          // Generar nombre del archivo
          const convocatoria = this.getSelectedConvocatoria();
          const nombreConvocatoria = convocatoria?.nombre || 'Convocatoria';
          const fechaActual = new Date().toISOString().split('T')[0];
          const nombreArchivo = `datos_completos_alumnos_${nombreConvocatoria.replace(/\s+/g, '_')}_${fechaActual}`;

          // Descargar archivo
          const excelBuffer: any = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });
          this.saveAsExcelFile(excelBuffer, nombreArchivo);
          
          this.showToast('success', `Datos completos de ${response.data.length} alumnos exportados correctamente`);
        },
        error: (err) => {
          this.isLoading = false;
          this.showToast('error', 'Error al exportar los datos completos. Intente nuevamente.');
        }
      })
    );
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

          if (!convocatoriasData || convocatoriasData.length === 0) {
            this.showToast('warning', 'No se encontraron convocatorias disponibles');
            return;
          }

          this.convocatorias = convocatoriasData.filter(conv => conv && conv.id && conv.nombre);

          if (this.convocatorias.length === 0) {
            this.showToast('warning', 'No se encontraron convocatorias válidas');
            return;
          }

          // Seleccion manual
          this.convocatoriaId = '';
        },
        error: err => {
          this.showToast('error', 'Error al cargar las convocatorias');
          this.isLoadingConvocatorias = false;
        }
      })
    );
  }

  onConvocatoriaChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    this.convocatoriaId = target.value;
    
    // Cargar residentes pendientes cuando se selecciona una convocatoria
    if (this.convocatoriaId) {
      this.loadResidentesPendientes();
    } else {
      this.residentesPendientes = [];
    }
  }

  loadReporteData() {
    this.isLoading = true;
    this.error = null;
    this._subscriptions.add(
      this._managerService.getReporteVisitaResidente().subscribe({
        next: (response) => {
          this.isLoading = false;
          if (response.data) {
            this.reporteData = response.data;
            this.updateCharts();
          } else {
            this.reporteData = null;
          }
        },
        error: (err) => {
          this.error = 'Error al cargar el reporte. Por favor, intente nuevamente.';
          this.showToast('error', 'Error al cargar el reporte');
          this.isLoading = false;
          this.reporteData = null;
        }
      })
    );
  }

  updateCharts() {
    if (!this.reporteData) {
      this.chartOptionsDonut.series = [];
      return;
    }

    this.chartOptionsDonut.series = [
      this.reporteData.alumnos_sin_visita,
      this.reporteData.verificadas,
      this.reporteData.observadas
    ];
  }

  updatePorEscuelaProfesionalChart() {
    if (!this.ReportePorEscuelaProfesional || this.ReportePorEscuelaProfesional.length === 0) {
      this.chartOptionsBar.series = [{ name: "Estudiantes Visitados", data: [] }];
      this.chartOptionsBar.xaxis.categories = [];
      return;
    }

    const sortedData = [...this.ReportePorEscuelaProfesional].sort((a, b) => b.total_visitados - a.total_visitados);
    
    this.chartOptionsBar.xaxis.categories = sortedData.map(item => item.escuela_profesional);
    this.chartOptionsBar.series = [{
      name: "Estudiantes Visitados",
      data: sortedData.map(item => item.total_visitados)
    }];
  }

  updateProcedenciaChart() {
    if (!this.reportePorProcedencia || this.reportePorProcedencia.length === 0) {
      this.chartOptionsProcedencia.series = [{ name: "Estudiantes Visitados", data: [] }];
      this.chartOptionsProcedencia.xaxis.categories = [];
      return;
    }

    const sortedData = [...this.reportePorProcedencia].sort((a, b) => b.total_visitados - a.total_visitados);
    
    this.chartOptionsProcedencia.xaxis.categories = sortedData.map(item => item.departamento);
    this.chartOptionsProcedencia.series = [{
      name: "Estudiantes Visitados",
      data: sortedData.map(item => item.total_visitados)
    }];
    
  }

  getSelectedConvocatoria(): IAnnouncement | null {
    return this.convocatorias.find(c => c.id?.toString() === this.convocatoriaId) || null;
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

  ngOnDestroy() {
    this._subscriptions.unsubscribe();
  }
}
