import { Component } from '@angular/core';
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {ExportComponent} from "../../../../core/utils/export/export.component";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {ModalComponent} from "../../../../core/ui/modal/modal.component";
import {NgForOf, NgIf} from "@angular/common";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {HouseholdTargeting} from "../../../../core/models/household-targeting";
import * as XLSX from "xlsx";

@Component({
  selector: 'app-household-targeting',
  standalone: true,
  imports: [
    BlockUiComponent,
    ExportComponent,
    FormsModule,
    ModalComponent,
    NgForOf,
    NgIf,
    ReactiveFormsModule,
    ToastComponent
  ],
  templateUrl: './household-targeting.component.html',
  styleUrl: './household-targeting.component.scss'
})
export class HouseholdTargetingComponent {

  public householdsTargeting: HouseholdTargeting[] = []

  uploadData(event: any) {
    const file = event.target.files[0];
    const reader = new FileReader();

    reader.onload = (e: any) => {
      const binaryStr = e.target.result;
      const workbook = XLSX.read(binaryStr, { type: 'binary' });

      // first field
      const worksheet = workbook.Sheets[workbook.SheetNames[0]];
      const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });

      const headers = jsonData[0];
      const data = jsonData.slice(1); // Los datos son el resto de las filas

      // Opcional: Convertir cada fila en un objeto usando los encabezados
      const dataFormat = data.map((row: any) => {
        // @ts-ignore
        return headers.reduce((acc: any, header: string, index: number) => {
          acc[header] = row[index]; // Asigna cada valor a su respectivo encabezado
          return acc;
        }, {});
      });
      this.householdsTargeting = dataFormat.map((row: any) => ({
        dni: row.dni,
        status: row.estado,
        current_date: row["fecha vigente"],
      }));

    };

    reader.readAsBinaryString(file); // Lee el archivo como cadena binaria
  }

  findRequest($event: Event) {

  }
}
