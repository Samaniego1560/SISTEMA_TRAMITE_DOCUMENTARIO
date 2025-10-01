import {Component, EventEmitter, Input, input, OnChanges, Output, SimpleChanges} from '@angular/core';
import {ModalComponent} from "../../ui/modal/modal.component";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {KeyValuePipe, NgForOf, NgIf} from "@angular/common";
import * as XLSX from "xlsx";
import {saveAs} from "file-saver";

@Component({
  selector: 'app-export',
  standalone: true,
  imports: [
    ModalComponent,
    FormsModule,
    NgIf,
    NgForOf,
    KeyValuePipe,
    ReactiveFormsModule
  ],
  templateUrl: './export.component.html',
  styleUrl: './export.component.scss'
})
export class ExportComponent implements OnChanges{
  @Input() public show: boolean = false;
  @Output() public showChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  titleModal = input.required<string>();
  @Input() public dataExport: any = [false];
  protected headers: Record<string, boolean> = {}
  public allSelected: boolean = false;

  constructor() {
  }

  ngOnChanges(changes: SimpleChanges) {
    if (this.dataExport.length) {
      const headers = Object.keys(this.dataExport[0]);
      headers.forEach(header => {
        this.headers[header] = true;
      })
    }
  }

  public exportData(): void {
    const objectsExport = this.getFilteredData();

    if (objectsExport && objectsExport.length > 0) {
      const originalHeaders = Object.keys(objectsExport[0]);
      const headers = originalHeaders.map(header =>
        header
          .split('_')
          .map(word => word.charAt(0).toUpperCase() + word.slice(1))
          .join(' ')
      );
      const worksheetData = [
        headers,
        ...objectsExport.map((item: any) =>
          originalHeaders.map(originalHeader => item[originalHeader])
        )
      ];

      const worksheet: XLSX.WorkSheet = XLSX.utils.aoa_to_sheet(worksheetData);
      const workbook: XLSX.WorkBook = { Sheets: { 'Datos': worksheet }, SheetNames: ['Datos'] };

      const excelBuffer: any = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });

      this.saveAsExcelFile(excelBuffer, 'Practicantes y tesistas');
    }
  }

  saveAsExcelFile(buffer: any, fileName: string): void {
    const EXCEL_TYPE = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8';
    const EXCEL_EXTENSION = '.xlsx';

    const data: Blob = new Blob([buffer], { type: EXCEL_TYPE });
    saveAs(data, fileName + '_export_' + new Date().getTime() + EXCEL_EXTENSION);
  }


  getFilteredData(): any[] {
    return this.dataExport.map((item: any) => {
      const filteredItem: any = {};

      for (const header in this.headers) {
        if (this.headers[header] && item.hasOwnProperty(header)) {
          filteredItem[header] = item[header];
        }
      }

      return filteredItem;
    }).filter((item: any) => Object.keys(item).length > 0);
  }

  toggleSelectAll() {
    Object.keys(this.headers).forEach(key => {
      this.headers[key] = this.allSelected;
    });
  }

  changeValue($event: any, key: string) {
    this.headers[key] = $event;
  }
}
