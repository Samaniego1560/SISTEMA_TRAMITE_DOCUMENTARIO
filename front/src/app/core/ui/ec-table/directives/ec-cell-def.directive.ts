import {
  Directive,
  Input,
  TemplateRef,
} from "@angular/core";

@Directive({
  selector: '[ecCellDef]',
  host: {},
})
export class EcCellDef {
  @Input() type: string | undefined;

  @Input('ecCellDef') name: any


  constructor(public template: TemplateRef<any>) {}

  getType(): string {
    return this.name!;
  }

}
