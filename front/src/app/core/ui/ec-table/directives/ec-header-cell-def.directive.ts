import {
  Directive,
  Input,
  TemplateRef,
} from "@angular/core";


@Directive({
  selector: '[ecHeaderCellDef]',
  host: {},
})
export class EcHeaderCellDef{
  @Input('ecHeader') name: string | undefined;


  constructor(public template: TemplateRef<any>) {}

  getType(): string {
    return this.name!;
  }

}
