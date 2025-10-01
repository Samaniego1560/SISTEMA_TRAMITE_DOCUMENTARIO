import {Component, input, model} from '@angular/core';
import {IErrorPostulation} from "../../../../models/requests";
import {ModalComponent} from "../../../../ui/modal/modal.component";

@Component({
  selector: 'app-postulation-error',
  standalone: true,
  imports: [
    ModalComponent
  ],
  templateUrl: './postulation-error.component.html',
  styleUrl: './postulation-error.component.scss'
})
export class PostulationErrorComponent {
  showModal = model<boolean>(false);
  errorsPostulation = input.required<IErrorPostulation[]>();
}
