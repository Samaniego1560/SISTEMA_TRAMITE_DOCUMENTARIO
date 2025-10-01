import {Component, OnDestroy, OnInit} from '@angular/core';
import {GraphicsAreaComponent} from "./components/graphics-area/graphics-area.component";
import {GraphicsPieComponent} from "./components/graphics-pie/graphics-pie.component";
import {catchError, of, Subscription} from "rxjs";
import {
  ChartOptionsAreaSqr,
  ChartOptionsPieSqr,
  chartPieSqr,
  chatAreaSqr, cuadros,
  trimestres
} from "./models/data-statistic-srq-model";

import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule} from "@angular/forms";
import {formatDate, NgIf} from "@angular/common";
import { BlockUiComponent } from '../../../../../core/ui/block-ui/block-ui.component';
import { GraphicsStatisticComponent } from '../../data-statistic/components/graphics-statistic/graphics-statistic.component';
import { CalendarYearComponent } from '../../../../../core/ui/calendar-year/calendar-year.component';
import { ControlErrorsDirective } from '../../../../../lib/validator-dynamic/directives/control-error.directive';
import { FormSubmitDirective } from '../../../../../lib/validator-dynamic/directives/form-submit.directive';
import { ChartBarsSchool } from '../../data-statistic/models/data-statistic.contant';
import { ChartOptionsBarra } from '../../data-statistic/data-statistic.component';
import { ManagerService } from '../../../../../core/services/manager/manager.service';
import { ToastService } from '../../../../../core/services/toast/toast.service';
import { Response } from '../../../../../core/models/global.model';
import { NgxValidators } from '../../../../../lib/validator-dynamic/utils/ngx-validators.util';

@Component({
  selector: 'app-data-statistic-srq',
  standalone: true,
  imports: [
    BlockUiComponent,
    GraphicsAreaComponent,
    GraphicsPieComponent,
    GraphicsStatisticComponent,
    CalendarYearComponent,
    ControlErrorsDirective,
    FormSubmitDirective,
    FormsModule,
    ReactiveFormsModule,
    NgIf
  ],
  templateUrl: './data-statistic-srq.component.html',
  styleUrl: './data-statistic-srq.component.scss'
})
export class DataStatisticSrqComponent implements OnInit, OnDestroy {

  private _subscriptions: Subscription = new Subscription();

  reportForm: FormGroup;
  selectedCuadro: string | null = null;
  trues = true;

  protected isLoading: boolean = false;
  public chartPieSqr: ChartOptionsPieSqr = chartPieSqr;
  public chartAreaSqr : ChartOptionsAreaSqr = chatAreaSqr;
  public chartBarsSchool: ChartOptionsBarra = ChartBarsSchool;
  public trimestres = trimestres;
  public reportes = cuadros;

  constructor(
    private _managerService: ManagerService,
    private _fb: FormBuilder,
    private _toastService: ToastService
  ) {
    this.reportForm = this.buildForm();
  }

  ngOnInit() {
    this.getDataForchartArea();
    this.getDataForchartPie();
    this.getDataForchartBarra();
  }
  ngOnDestroy() {
    this._subscriptions.unsubscribe();
  }

  onTrimestreChange(){
    this.getDataForchartPie();
    this.getDataForchartBarra();
    this.getDataForchartArea();
  }

  getDataForchartArea() {
    this.isLoading = true;
    const filterParams = this._getParamsReport();
    this._subscriptions.add(
      this._managerService.getDatachartArea(filterParams).subscribe({
        next: (response: any) => {
          if (response) {

            const dates: string[] = response.detalle.map((item : any) => item.date);
            const countsSqr: number[] = response.detalle.map((item : any) => item.countSrq);
            const counts: number[] = response.detalle.map((item : any) => item.count);
            this.chartAreaSqr.series = [
              {
                name : 'Evaluacion',
                data : countsSqr
              },
              {
                name : 'Consulta psicologica',
                data : counts
              }
            ];
            this.chartAreaSqr.labels = dates;
          }
          this.isLoading = false;
        },
        error: (error: any) => {
          this.isLoading = false;
          console.error('Error al obtener la data para los graficos', error);
        }
      })
    );
  }

  getDataForchartPie() {
    this.isLoading = true;
    const filterParams = this._getParamsReport();
    this._subscriptions.add(
      this._managerService.getDataChartPie(filterParams).subscribe({
        next: (response: any) => {
          if (response) {

            const detalle = Array.isArray(response?.detalle) ? response.detalle : [];

            const labels = detalle.map((item: { estados_evaluacion?: string }) =>
              item.estados_evaluacion ?? ''
            );

            const series = detalle.map((item: { total?: number | string }) =>
              Number(item.total ?? 0)
            );

            this.chartPieSqr.labels = labels;
            this.chartPieSqr.series = series;

          }
          this.isLoading = false;
        },
        error: (error: any) => {
          this.isLoading = false;
          console.error('Error al obtener la data para los graficos Pie', error);
        }
      })
    );
  }

  getDataForchartBarra() {
    this.isLoading = true;
    const filterParams = this._getParamsReport();
    this._subscriptions.add(
      this._managerService.getDataChartBarras(filterParams).subscribe({
        next: (response: any) => {
          if (response) {

            const labels : string[] = response.detalle.map((item : any) => item.escuela);
            const series : number[] = response.detalle.map((item : any) => item.total);

            this.chartBarsSchool.series = [];
            this.chartBarsSchool.series.push(
              {
                name : 'total',
                data : series
              }
            )

            this.chartBarsSchool.xaxis = {
              categories : labels
            }

            this.chartBarsSchool.title = {
              align: "center",
              text: 'Participantes por Facultades'
            }
          }
          this.isLoading = false;
        },
        error: (error: any) => {
          this.isLoading = false;
          console.error('Error al obtener la data para los graficos barras', error);
        }
      })
    );
  }

  onSearchReport(): void {
    const { start_date, end_date, anio, meses  } = this.reportForm.value;
    if (start_date && end_date) {
      this.GenerateReportAttentions();
      return;
    }
    if (anio && meses) {
      this.GenerateReportAttentions();
      return;
    }
    if(this.reportForm.invalid) {
      this.reportForm.markAllAsTouched();
      return
    }

    this.GenerateReportAttentions();
  }

  GenerateReportAttentions() {
    this.isLoading = true;
    const filterParams = this._getParamsReport();
    this._managerService.getReportAttentionsPsicology(filterParams)
      .pipe(catchError(() => (of({error: true} as Response<string>))))
      .subscribe({
        next: (res) => {
          this.isLoading = false;
          if (res.error) {
            this._toastService.add({
              message: 'Error al obtener el reporte',
              type: 'error',
              life: 5000
            });
            return
          }

          console.log(res.data)

          const link = document.createElement('a');
          link.href = `data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,${res.data}`;
          link.download = 'Reporte';
          link.click();
        }
      })
  }

  onCuadroChange() {
    this.selectedCuadro = this.reportForm.get('numero_cuadro')?.value;
    if (this.selectedCuadro === '3') {
      this.reportForm.patchValue({ anio: null, meses: null });
    } else {
      this.reportForm.patchValue({ start_date: null, end_date: null });
    }
  }

  private _getParamsReport() {
    const formValues = this.reportForm.value;
    return {
      ...formValues,
      anio: formValues['anio'] ? formatDate(formValues['anio'], 'yyyy', 'en') : null,
    }
  } 

  buildForm() {
    const currentYear = new Date().getFullYear().toString();
    return this._fb.group({
      anio: [currentYear, [NgxValidators.required()]],
      meses: [null, NgxValidators.required()],
      numero_cuadro: [null, NgxValidators.required()], 
      start_date: [null, NgxValidators.required()], 
      end_date: [null, NgxValidators.required()] 
    });
  }
}
