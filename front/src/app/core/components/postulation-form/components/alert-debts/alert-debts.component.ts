import {Component, input, model} from '@angular/core';
import {ModalComponent} from "../../../../ui/modal/modal.component";
import {FormGroup} from "@angular/forms";
import {IDebtsStudent} from "../../../../models/user";
import {DatePipe} from "@angular/common";

@Component({
  selector: 'app-alert-debts',
  standalone: true,
  imports: [
    ModalComponent,
    DatePipe
  ],
  templateUrl: './alert-debts.component.html',
  styleUrl: './alert-debts.component.scss'
})
export class AlertDebtsComponent {
  showModal = model<boolean>(false);
  debts = input.required<IDebtsStudent[]>();
}
