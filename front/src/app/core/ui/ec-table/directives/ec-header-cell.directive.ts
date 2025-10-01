import {Directive, ElementRef, Input} from "@angular/core";

@Directive({
  selector: '[ec-header-cell]',
  host: {}
})
export class EcHeaderCell {
  @Input() type: string | undefined;

  @Input('ec-header-cell') name: string | undefined;

  constructor(public element: ElementRef<any>) {}

  getType(): string {
    return this.name!;
  }
}
