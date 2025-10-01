import {Component, input} from '@angular/core';
import {ISection} from "../../models/announcement";

@Component({
  selector: 'app-alert-condition',
  standalone: true,
  imports: [],
  templateUrl: './alert-condition.component.html',
  styleUrl: './alert-condition.component.scss'
})
export class AlertConditionComponent {
  title = input.required<string>();
  message = input.required<string>();
}
