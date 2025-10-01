import { Component } from '@angular/core';
import {ModalComponent} from "../../../../../../core/ui/modal/modal.component";
import {DomSanitizer, SafeResourceUrl} from "@angular/platform-browser";

@Component({
  selector: 'app-sis',
  standalone: true,
    imports: [
        ModalComponent
    ],
  templateUrl: './sis.component.html',
  styleUrl: './sis.component.scss'
})
export class SisComponent {
  isShow: boolean = true;

  safeUrl: SafeResourceUrl;

  constructor(private sanitizer: DomSanitizer) {
    // Asegúrate de sanitizar la URL
    this.safeUrl = this.sanitizer.bypassSecurityTrustResourceUrl('http://app3.sis.gob.pe/SisConsultaEnLinea/Consulta/frmResultadoEnLinea.aspx');
  }
}
