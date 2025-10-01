import {Directive, Input, TemplateRef} from "@angular/core";

@Directive({
  selector: '[ecTemplate]',
  standalone: true,
  host: {}
})
export class EcTemplate {
  @Input('ecTemplate') name: string | undefined;

  constructor(public template: TemplateRef<any>) {}

  getType(): string {
    return this.name!;
  }
}
