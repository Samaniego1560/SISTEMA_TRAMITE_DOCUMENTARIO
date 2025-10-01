import {Component, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {NgForOf} from "@angular/common";
import {
  ApexAxisChartSeries,
  ApexChart, ApexDataLabels, ApexFill, ApexLegend,
  ApexNonAxisChartSeries, ApexPlotOptions,
  ApexResponsive, ApexStroke,
  ApexTitleSubtitle, ApexTooltip, ApexXAxis, ApexYAxis,
  ChartComponent,
  NgApexchartsModule
} from "ng-apexcharts";
import {ManagerService} from "../../../../core/services/manager/manager.service";
import {HttpErrorResponse} from "@angular/common/http";
import {Subscription} from "rxjs";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {IAnnouncement} from "../../../../core/models/announcement";
import {NgIf} from "@angular/common";
import {IStatistics} from "../../../../core/models/statistics";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import { Residence, Submission } from '../../../../core/models/residence';
import { SubmissionService } from '../../../../core/services/residences/submission/submissions.service';
import { FileService } from '../../../../core/services/file/file.service';
import { ResidencesService } from '../../../../core/services/residences/residences/residences.service';

type SchoolCount = {
  [key: string]: number;
};

export type ChartOptions = {
  series: ApexNonAxisChartSeries;
  chart: ApexChart;
  responsive: ApexResponsive[];
  labels: any;
  title: ApexTitleSubtitle;
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
};


@Component({
  selector: 'app-report-residences',
  standalone: true,
  imports: [
    NgApexchartsModule,
    NgIf,
    BlockUiComponent,
    NgForOf
  ],
  providers: [ManagerService, ToastService],
  templateUrl: './report-residences.component.html',
  styleUrl: './report-residences.component.scss'
})
export class ReportResidencesComponent {

  protected isLoading: boolean = false;
  protected submissions: any[] = [];
  protected submissionSelected: any;
  protected students: any[] = [];
  protected totalResidents: number = 0;
  protected totalResidentsAssigned: number = 0;
  protected totalResidentsNotAssigned: number = 0;
  protected facultades: any[] = [];
  zeroWidthSpace: string = '\u200B';
  protected resicencias: Residence[] = [];

  private _subscriptions: Subscription = new Subscription();
  public _idAnnouncementSelect: number = 0;

  public _datasetService: {name: string, value: number} [] = [
    {name: 'Comedor', value: 1},
    {name: 'Residencia', value: 2}
  ];

  public _announcements: {name: string, value: number}[] = [];

  @ViewChild("chart") chart!: ChartComponent;

  public chartBarrasHFacultad: ChartOptionsBarra;
  public chartBarrasHResidencia: ChartOptionsBarra;
  public chartBarrasHGeneroFacultad: ChartOptionsBarra;
  public chartBarrasHDepartamento: ChartOptionsBarra;

  public professionl_school: SchoolCount = {
    "INGENIERIA EN INFORMATICA Y SISTEMAS":0,
    "INGENIERIA AMBIENTAL":0,
    "ECONOMIA":0,
    "INGENIERIA MECANICA ELECTRICA":0,
    "INGENIERIA EN INDUSTRIAS ALIMENTARIAS":0,
    "ZOOTECNIA":0,
    "INGENIERIA FORESTAL":0,
    "INGENIERIA EN CONSERVACION DE SUELOS Y AGUA":0,
    "ADMINISTRACION":0,
    "AGRONOMIA":0,
    "INGENIERIA EN RECURSOS NATURALES RENOVABLES":0,
    "CONTABILIDAD":0,
    "TURISMO Y HOTELERIA":0,
    "INGENIERIA EN CIBERSEGURIDAD":0,
    "INGENIERIA CIVIL":0
  };

  public department: SchoolCount = {
    "AMAZONAS": 0,
    "ANCASH": 0,
    "APURÍMAC": 0,
    "AREQUIPA": 0,
    "AYACUCHO": 0,
    "CAJAMARCA": 0,
    "CALLAO": 0,
    "CUSCO": 0,
    "HUANCAVELICA": 0,
    "HUÁNUCO": 0,
    "ICA": 0,
    "JUNÍN": 0,
    "LA LIBERTAD": 0,
    "LAMBAYEQUE": 0,
    "LIMA": 0,
    "LORETO": 0,
    "MADRE DE DIOS": 0,
    "MOQUEGUA": 0,
    "PASCO": 0,
    "PIURA": 0,
    "PUNO": 0,
    "SAN MARTÍN": 0,
    "TACNA": 0,
    "TUMBES": 0,
    "UCAYALI" : 0
  };
  protected statistics = {
    estados_solicitud: {
      pendiente: 0,
      rechazado: 0,
      aceptado: 0,
      aprobado: 0
    },
    escuelas_profesionales: {
      "ADMINISTRACION": 0,
      "AGRONOMIA": 0,
      "CONTABILIDAD": 0,
      "ECONOMIA": 0,
      "INGENIERIA AMBIENTAL": 0,
      "INGENIERIA EN CONSERVACION DE SUELOS Y AGUA": 0,
      "INGENIERIA EN INDUSTRIAS ALIMENTARIAS": 0,
      "INGENIERIA EN INFORMATICA Y SISTEMAS": 0,
      "INGENIERIA EN RECURSOS NATURALES RENOVABLES": 0,
      "INGENIERIA FORESTAL": 0,
      "INGENIERIA MECANICA ELECTRICA": 0,
      "ZOOTECNIA": 0
    },
    facultades: {
      "FACULTAD DE AGRONOMIA": 0,
      "FACULTAD DE CIENCIAS CONTABLES": 0,
      "FACULTAD DE CIENCIAS ECONOMICAS Y ADMINISTRATIVAS": 0,
      "FACULTAD DE INGENIERIA EN INDUSTRIAS ALIMENTARIAS": 0,
      "FACULTAD DE INGENIERIA EN INFORMATICA Y SISTEMAS": 0,
      "FACULTAD DE INGENIERIA MECANICA ELECTRICA": 0,
      "FACULTAD DE RECURSOS NATURALES RENOVABLES": 0,
      "FACULTAD DE ZOOTECNIA": 0
    },
    residencias: {
      "POMARROSAS":0,
      "TULUMAYOS":0,
      "MARIA ANGOLA":0,
      "MARÍA ANGOLA II":0,
      "SHERATON":0,
      "BAMBU":0,
      "BRITANICO":0,
      "CASTAÑAS":0
    },
    sexo: {
      num_hombres: 0,
      num_mujeres: 0
    }
  }

  constructor(
    private _managerService: ManagerService,
    private _toastService: ToastService,
    private _submissionService: SubmissionService,
    private _fileService: FileService,
    private _residenceService: ResidencesService
  ) {

    this.chartBarrasHFacultad = {
      series: [],
      chart: {
        type: "bar",
        height: 350,
        width: 500
      },
      plotOptions: {
        bar: {
          horizontal: true,
          columnWidth: "55%"
        }
      },
      dataLabels: {
        enabled: false
      },
      stroke: {
        show: true,
        width: 2,
        colors: ["transparent"]
      },
      xaxis: {
        categories: []
      },
      yaxis: {},
      fill: {
        opacity: 1
      },
      tooltip: {},
      legend: {
        position: "bottom",
        horizontalAlign: "left",
        offsetX: 40
      },
      title: {
        align: "center",
        text: "Beneficiarios por Escuela Profesional"
      }
    };


    this.chartBarrasHResidencia = {
      series: [],
      chart: {
        type: "bar",
        height: 350,
        width: 500
      },
      plotOptions: {
        bar: {
          horizontal: true,
          columnWidth: "55%"
        }
      },
      dataLabels: {
        enabled: false
      },
      stroke: {
        show: true,
        width: 2,
        colors: ["transparent"]
      },
      xaxis: {
        categories: [

        ]
      },
      yaxis: {},
      fill: {
        opacity: 1
      },
      tooltip: {},
      legend: {
        position: "bottom",
        horizontalAlign: "left",
        offsetX: 40
      },
      title: {
        align: "center",
        text: "Beneficiarios por Residencia"
      }
    };

    this.chartBarrasHGeneroFacultad = {
      series: [],
      chart: {
        type: "bar",
        height: 350,
        width: 500
      },
      plotOptions: {
        bar: {
          horizontal: true,
          columnWidth: "55%"
        }
      },
      dataLabels: {
        enabled: false
      },
      stroke: {
        show: true,
        width: 2,
        colors: ["transparent"]
      },
      xaxis: {
        categories: [

        ]
      },
      yaxis: {},
      fill: {
        opacity: 1
      },
      tooltip: {},
      legend: {
        position: "bottom",
        horizontalAlign: "center",
        offsetX: 40
      },
      title: {
        align: "center",
        text: "Beneficiarios por facultad y género"
      }
    };

    this.chartBarrasHDepartamento = {
      series: [],
      chart: {
        type: "bar",
        height: 350,
        width: 500
      },
      plotOptions: {
        bar: {
          horizontal: true,
          columnWidth: "55%"
        }
      },
      dataLabels: {
        enabled: false
      },
      stroke: {
        show: true,
        width: 2,
        colors: ["transparent"]
      },
      xaxis: {
        categories: [

        ]
      },
      yaxis: {},
      fill: {
        opacity: 1
      },
      tooltip: {},
      legend: {
        position: "bottom",
        horizontalAlign: "center",
        offsetX: 40
      },
      title: {
        align: "center",
        text: "Beneficiarios por departamento"
      }
    };
    // this._getCurrentAnnouncement();
  }

  public async ngOnInit(): Promise<void> {
    this.getResidences();
  }

  onChangeSubmission(event: any): void {
    debugger
    this.submissionSelected = this.submissions.find((item: any) => item.id == event.target.value);
    this.getStudents();
  }

  protected getAnnoucement(): void {
    try {
      this._managerService.getAnnouncement().subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res) {
            return;
          }
          this.submissions = res.detalle;
          this.submissionSelected = this.submissions.find((item: any) => item.estado === 'Activa');
          if (!this.submissionSelected) {
            this.submissionSelected = this.submissions[this.submissions.length - 1];
          }
          this.submissionSelected = this.submissions[this.submissions.length - 1];
          this.getStudents();
          return;
        },
        error: (err: any) => {
          console.error(err);
          this.isLoading = false;
        }
      });
    } catch (error: any) {
      this.isLoading = false;
      this._toastService.add({ type: 'error', message: error.toString() });
    }
  }

  protected getStudents(): void {
    try {
      this._submissionService.getStudentsByAnnouncement(this.submissionSelected.id, {}).subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res) {
            return;
          }
          this.students = res.data.students;
          this.classifyData();
          const facultades = this.getFacultadesChartData(this.students);
          this.chartBarrasHFacultad.series = [];
          this.chartBarrasHFacultad.series.push(
            {
              name: 'Beneficiario',
              data: Object.values(facultades)
            }
          );
          const residencias = this.getResidenciasChartData(this.students);
          this.chartBarrasHResidencia.series = [];
          this.chartBarrasHResidencia.series.push(
            {
              name: 'Beneficiario',
              data: Object.values(residencias)
            }
          );

          const dataMasculino = this.getGenderSchoolChartDataWithTotals(this.students, 'M');
          const dataFemenino = this.getGenderSchoolChartDataWithTotals(this.students, 'F');
          this.chartBarrasHGeneroFacultad.series = [];
          this.chartBarrasHGeneroFacultad.series.push(
            {
              name: 'Femenino',
              data: Object.values(dataFemenino)
            },
            {
              name: 'Masculino',
              data: Object.values(dataMasculino)
            }
          );

          const dataDepartamento = this.getDepartamentosChartDataWithTotal(this.students);
          this.chartBarrasHDepartamento.series = [];
          this.chartBarrasHDepartamento.series.push(
            {
              name: 'Beneficiario',
              data: Object.values(dataDepartamento)
            }
          );
          return;
        },
        error: (err: any) => {
          console.error(err);
          this.isLoading = false;

        }
      });
    } catch (error: any) {
      this.isLoading = false;
      console.log(error);//
    }
  }

  protected getResidences(): void {
    try {
      this._residenceService.getListResidences().subscribe({
        next: (res: any) => {
          this.isLoading = false;
          if (!res) {
            return;
          }
          this.resicencias = res.data;
          this.getAnnoucement();
          return;
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

  protected classifyData(): void {
    debugger
    this.totalResidents = this.students.length;
    this.totalResidentsAssigned = this.students.filter((item: any) => item.student.room).length;
    this.totalResidentsNotAssigned = this.students.filter((item: any) => !item.student.room).length;
  }

  protected getFacultades(students: any): any {
    return students.reduce((acc:any, obj:any) => { acc[obj.student["professional_school"]] = (acc[obj.student["professional_school"]] || 0) + 1;
      return acc;
    }, {});
  }

  // Transformamos los datos a un formato que ApexCharts acepta
  protected getFacultadesChartData(students: any): { x: string, y: number }[] {
    const facultadesCount = students.reduce((acc: any, obj: any) => {
      const key = obj.student["professional_school"]; // Acceder con comillas
      acc[key] = (acc[key] || 0) + 1;
      return acc;
    }, {});

    // Convertimos el objeto en un array de objetos con { x, y }
    return Object.entries(facultadesCount).map(([key, value]) => ({
      x: key,  // Nombre de la facultad
      y: value as number // Cantidad de estudiantes
    }));
  }

  protected getResidenciasChartData(students: any): { x: string, y: number }[] {
    const residenciesCount = students.reduce((acc: any, obj: any) => {
      const key = obj.student["residence"];
      if (!key) {
        return acc;
      }
      const residenceName = this.resicencias.find((item: any) => item.id == key)?.name;
      if (!residenceName) {
        return acc;
      }
      acc[residenceName] = (acc[residenceName] || 0) + 1;
      return acc;
    }, {});

    // Convertimos el objeto en un array de objetos con { x, y }
    return Object.entries(residenciesCount).map(([key, value]) => ({
      x: key,  // Nombre de la facultad
      y: value as number // Cantidad de estudiantes
    }));
  }

  protected getGenderSchoolChartData(students: any, gender: string): { x: string, y: number }[] {
    const facultadesCount = students.reduce((acc: any, obj: any) => {
      const key = obj.student["professional_school"];

      if (obj.student['sex'] != gender)

        return acc; // Acceder con comillas
      acc[key] = (acc[key] || 0) + 1;
      return acc;
    }, {});

    // Convertimos el objeto en un array de objetos con { x, y }
    return Object.entries(facultadesCount).map(([key, value]) => ({
      x: key,  // Nombre de la facultad
      y: value as number // Cantidad de estudiantes
    }));
  }

  protected getGenderSchoolChartDataWithTotals(students: any, gender: string): { x: string, y: number }[] {
  // Crear una copia del objeto professionl_school para garantizar la completitud
  const facultadesCount = { ...this.professionl_school };

  // Recorrer los estudiantes y actualizar el conteo
  students.forEach((obj: any) => {
    const key = obj.student["professional_school"];

    if (obj.student['sex'] !== gender) return;

    if (key in facultadesCount) {
      facultadesCount[key]++;
    }
  });

  // Convertir el objeto en un array de { x, y }
  return Object.entries(facultadesCount).map(([key, value]) => ({
    x: key,
    y: value
  }));
}

  protected getDepartamentosChartData(students: any): { x: string, y: number }[] {
    const facultadesCount = students.reduce((acc: any, obj: any) => {
      const departamento: string = obj.student["department"];
      const key = departamento.split('/')[0];
      acc[key] = (acc[key] || 0) + 1;
      return acc;
    }, {});

    // Convertimos el objeto en un array de objetos con { x, y }
    return Object.entries(facultadesCount).map(([key, value]) => ({
      x: key,  // Nombre de la facultad
      y: value as number // Cantidad de estudiantes
    }));
  }

  protected getDepartamentosChartDataWithTotal(students: any): { x: string, y: number }[] {

  const department = {...this.department};

  students.reduce((acc: any, obj: any) => {
    const departamento: string = obj.student["department"];
    const key = departamento.split('/')[0].trim().toUpperCase(); // Asegura que coincida con el formato estático
    if (key in department) {
      department[key]++;
    }
    return acc;
  }, {});

  // Devuelve el array completo con los valores existentes o cero
  return Object.entries(department).map(([key, value]) => ({
    x: key,  // Nombre de la facultad
    y: value as number // Cantidad de estudiantes
  }));
}

  protected exportReport(): void {
    this.isLoading = true;
    this._subscriptions.add(
      this._submissionService.exportReportResidence(this.submissionSelected.id).subscribe({
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
                this._fileService.base64ToExcel(res.data, `REPORTE_RESIDENCIAS`);

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

  protected getLastSubmission(): void {

  }
}
