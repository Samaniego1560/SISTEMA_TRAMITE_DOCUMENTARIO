import {Component, input} from '@angular/core';
import {ChartOptionsBarra} from "../../data-statistic.component";
import {NgApexchartsModule} from "ng-apexcharts";

@Component({
  selector: 'app-graphics-statistic',
  standalone: true,
  imports: [
    NgApexchartsModule
  ],
  templateUrl: './graphics-statistic.component.html',
  styleUrl: './graphics-statistic.component.scss'
})
export class GraphicsStatisticComponent {
  chartBars = input.required<ChartOptionsBarra>();
}
