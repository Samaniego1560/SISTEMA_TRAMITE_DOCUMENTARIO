import {Component, input} from '@angular/core';
import {ChartOptionsPieSqr} from "../../models/data-statistic-srq-model";
import {NgApexchartsModule} from "ng-apexcharts";

@Component({
  selector: 'app-graphics-pie',
  standalone: true,
  imports: [
    NgApexchartsModule
  ],
  templateUrl: './graphics-pie.component.html',
  styleUrl: './graphics-pie.component.scss'
})
export class GraphicsPieComponent {
  chartPie = input.required<ChartOptionsPieSqr>();
}
