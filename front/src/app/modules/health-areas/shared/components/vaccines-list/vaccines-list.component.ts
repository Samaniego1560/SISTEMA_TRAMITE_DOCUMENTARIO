import {Component, input, output} from '@angular/core';
import {VaccineListCard, Vacuna} from "../../../../../core/models/areas/nursing/nursing.model";
import {DatePipe} from "@angular/common";

@Component({
  selector: 'app-vaccines-list',
  standalone: true,
  imports: [
    DatePipe
  ],
  templateUrl: './vaccines-list.component.html',
  styleUrl: './vaccines-list.component.scss'
})
export class VaccinesListComponent {
  isModeCreation = input.required<boolean>()
  vaccines = input.required<VaccineListCard[]>()

  editRecord = output<Vacuna>()
  deleteRecord = output<string>()
}
