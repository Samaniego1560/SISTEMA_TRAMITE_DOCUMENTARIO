import { Component } from '@angular/core';
import { PostulationDynamicComponent } from 'ruta/correcta/postulation-dynamic.component';

export const REQUEST_STATUS = {
  PENDING: 'PENDIENTE',
  APPROVED: 'APROBADO',
  ACCEPTED: 'ACEPTADO',
  REJECTED: 'RECHAZADO',
  RETIRED: 'RETIRADO'
} as const;

export const REQUEST_STATUS_ICON = {
  PENDIENTE: 'hourglass_empty',
  APROBADO: 'check_circle',
  ACEPTADO: 'done_all',
  RECHAZADO: 'cancel',
  RETIRADO: 'cancel'
} as const;

// Si necesitas la lógica del switch, ponla en una función:
export function getStatusLogic(status: string) {
  switch (status) {
    case 'PENDIENTE':
    case 'APROBADO':
    case 'ACEPTADO':
    case 'RECHAZADO':
    case 'RETIRADO':
      // lógica aquí
      break;
    default:
      // manejar estado desconocido
  }
}

@Component({
  selector: 'app-postulation-dynamic',
  standalone: true, // <-- esto debe estar presente
  // ...
})
export class PostulationDynamicComponent { ... }import { Component } from '@angular/core';

@Component({
  selector: 'app-postulation-dynamic',
  standalone: true,
  // imports: [CommonModule, ...] // agregar imports necesarios
  templateUrl: './postulation-dynamic.component.html',
  styleUrls: ['./postulation-dynamic.component.scss']
})
export class PostulationDynamicComponent {
  // ...existing code...
}