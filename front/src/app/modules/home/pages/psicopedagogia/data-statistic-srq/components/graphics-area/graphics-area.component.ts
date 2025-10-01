import {Component, input} from '@angular/core';
import {NgApexchartsModule} from "ng-apexcharts";
import {ChartOptionsAreaSqr} from "../../models/data-statistic-srq-model";

@Component({
  selector: 'app-graphics-area',
  standalone: true,
  imports: [
    NgApexchartsModule
  ],
  templateUrl: './graphics-area.component.html',
  styleUrl: './graphics-area.component.scss'
})
export class GraphicsAreaComponent {
  chartArea = input.required<ChartOptionsAreaSqr>();
}
