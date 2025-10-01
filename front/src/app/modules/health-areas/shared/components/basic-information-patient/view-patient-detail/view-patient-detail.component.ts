import {Component, input, output} from '@angular/core';
import {Patient} from "../../../../../../core/models/areas/patient.model";
import {EcDialogComponent} from "../../../../../../core/ui/ec-dialog/ec-dialog.component";
import {EcTemplate} from "../../../../../../core/directives/ec-template.directive";
import {ToggleSwitchComponent} from "../../../../../../core/ui/toggle/toggle-switch.component";
import {EcTableModule} from "../../../../../../core/ui/ec-table/ec-table.module";

@Component({
  selector: 'app-view-patient-detail',
  standalone: true,
  imports: [
    EcDialogComponent,
    EcTemplate,
    ToggleSwitchComponent,
    EcTableModule
  ],
  templateUrl: './view-patient-detail.component.html',
  styleUrl: './view-patient-detail.component.scss'
})
export class ViewPatientDetailComponent {
  patient = input.required<Patient>()

  close = output<void>()
}
