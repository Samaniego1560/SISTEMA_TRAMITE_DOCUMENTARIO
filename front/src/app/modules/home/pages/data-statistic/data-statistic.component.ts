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
import {GraphicsStatisticComponent} from "./components/graphics-statistic/graphics-statistic.component";
import {ChartBarsFaculty, ChartBarsSchool, ChartBarsSchoolGender} from "./models/data-statistic.contant";

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
  selector: 'app-data-statistic',
  standalone: true,
  imports: [
    NgApexchartsModule,
    NgIf,
    BlockUiComponent,
    NgForOf,
    GraphicsStatisticComponent
  ],
  templateUrl: './data-statistic.component.html',
  styleUrl: './data-statistic.component.scss',
  providers: [ManagerService, ToastService]
})
export class DataStatisticComponent implements OnInit, OnDestroy {

  private _subscriptions: Subscription = new Subscription();
  public _idAnnouncementSelect: number = 0;
  private _announcement: IAnnouncement = {
    nombre: '',
    convocatoria_servicio: [],
    fecha_fin: '',
    fecha_inicio: '',
    secciones: [],
    activo: false,
    credito_minimo: 0,
    nota_minima: 0
  };

  public _datasetService: {name: string, value: number} [] = [
    {name: 'Comedor', value: 1},
    {name: 'Residencia', value: 2}
  ];

  public _announcements: {name: string, value: number}[] = [];

  @ViewChild("chart") chart!: ChartComponent;

  public chartBarsFaculty: ChartOptionsBarra = ChartBarsFaculty;
  public chartBarsSchool: ChartOptionsBarra = ChartBarsSchool;
  public chartBarsSchoolGender: ChartOptionsBarra = ChartBarsSchoolGender;

  protected isLoading: boolean = false;

  protected statistics: IStatistics = {
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
    sexo: {
      num_hombres: 0,
      num_mujeres: 0
    },
    escuela_profesionales_genero: {
      Femenino: {
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
      Masculino: {
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
      }
    }
  }

  constructor(
    private _managerService: ManagerService,
    private _toastService: ToastService
  ) {
    // this._getCurrentAnnouncement();
  }
  ngOnInit() {
    this._getAnnouncementAll();
  }

  onChangeOfSubjects(event: Event): void {
    const selectedValue = (event.target as HTMLSelectElement).value;
    this._idAnnouncementSelect = Number(selectedValue);
    this._getStatistic(Number(selectedValue));
  }

  onChangeOfService(event: Event): void {
    const selectedService = (event.target as HTMLSelectElement).value;
    if (this._idAnnouncementSelect === 0) {
      return;
    }
    this._getStatisticAndService(this._idAnnouncementSelect, Number(selectedService));
  }

  ngOnDestroy() {
    this._subscriptions.unsubscribe();
  }

  private _getAnnouncementAll() {
    this.isLoading = true;
    this._subscriptions.add(
      this._managerService.getAnnouncement().subscribe({
        next: (res: any) => {
          this.isLoading = false;
         if (!res.detalle) {
            this.isLoading = false;
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }

          this._announcements = res.detalle.map((item: any) => ({
            name: item.nombre,
            value: item.id !== undefined ? item.id : 0
          }));
          // this._announcement = res.detalle;

          // this._getStatistic();
        },
        error: (err: HttpErrorResponse) => {
          this.isLoading = false;
          this._toastService.add({type: 'error', message: 'No se pudo obtener la convocatoria actual'});
          console.error(err)
        }
      })
    );
  }

  private _getCurrentAnnouncement(): void {
    this.isLoading = true;

    this._subscriptions.add(
      this._managerService.getCurrentAnnouncement().subscribe({
        next: (res: any) => {
         if (!res.detalle) {
            this.isLoading = false;
            this._toastService.add({type: 'error', message: res.msg});
            return;
          }

          this._announcement = res.detalle;

         // this._getStatistic();
        },
        error: (err: HttpErrorResponse) => {
          this.isLoading = false;
          this._toastService.add({type: 'error', message: 'No se pudo obtener la convocatoria actual'});
          console.error(err)
        }
      })
    );
  }

  private _getStatistic(code: number): void {
    this.isLoading = true;
    this._managerService.getStatistics(code).subscribe({
      next: (res: any) => {
        this.isLoading = false;
        if (!res.detalle) {
          this._toastService.add({type: 'error', message: res.msg});
          return;
        }

        this.statistics = res.detalle;

        this.chartBarsFaculty.series = [];
        this.chartBarsFaculty.series.push(
          {
            name: 'Postulantes',
            data: Object.values(this.statistics.facultades)
          }
        );

        this.chartBarsSchool.series = [];
        this.chartBarsSchool.series.push(
          {
            name: 'Postulantes',
            data: Object.values(this.statistics.escuelas_profesionales)
          }
        );

        this.chartBarsSchoolGender.series = [];
        this.chartBarsSchoolGender.series.push(
          {
            name: 'Femenino',
            data: Object.values(this.statistics.escuela_profesionales_genero.Femenino)
          },
          {
            name: 'Masculino',
            data: Object.values(this.statistics.escuela_profesionales_genero.Masculino)
          }
        );

      },
      error: (err: HttpErrorResponse) => {
        this._toastService.add({type: 'error', message: 'No se pudo obtener las estadísticas'});
        console.log(err);
        this.isLoading = false;
      }
    });
  }

  private _getStatisticAndService(code: number, typeService: number): void {
    this.isLoading = true;
    this._managerService.getStatisticsAndService(code, typeService).subscribe({
      next: (res: any) => {
        this.isLoading = false;
        if (!res.detalle) {
          this._toastService.add({type: 'error', message: res.msg});
          return;
        }

        this.statistics = res.detalle;

        this.chartBarsFaculty.series = [];
        this.chartBarsFaculty.series.push(
          {
            name: 'Postulantes',
            data: Object.values(this.statistics.facultades)
          }
        );

        this.chartBarsSchool.series = [];
        this.chartBarsSchool.series.push(
          {
            name: 'Postulantes',
            data: Object.values(this.statistics.escuelas_profesionales)
          }
        );

        console.log(this.chartBarsSchool)
      },
      error: (err: HttpErrorResponse) => {
        this._toastService.add({type: 'error', message: 'No se pudo obtener las estadísticas'});
        console.log(err);
        this.isLoading = false;
      }
    });
  }
}
