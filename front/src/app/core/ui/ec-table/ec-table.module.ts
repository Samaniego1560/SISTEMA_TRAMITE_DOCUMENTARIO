import {NgModule} from "@angular/core";
import {EcTableComponent} from "./ec-table.component";
import {EcCellDef, EcColumnDef, EcHeaderCell, EcHeaderCellDef} from "./directives";
import {NgForOf, NgTemplateOutlet} from "@angular/common";

@NgModule({
  declarations: [
    EcTableComponent,
    EcHeaderCell,
    EcHeaderCellDef,
    EcCellDef,
    EcColumnDef
  ],
  imports: [
    NgForOf,
    NgTemplateOutlet
  ],
  exports: [
    EcTableComponent,
    EcHeaderCell,
    EcHeaderCellDef,
    EcCellDef,
    EcColumnDef
  ]
})
export class EcTableModule {}
