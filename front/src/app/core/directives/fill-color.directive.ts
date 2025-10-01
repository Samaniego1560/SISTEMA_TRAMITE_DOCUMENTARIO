import {Directive, ElementRef, HostListener, Input} from '@angular/core';

@Directive({
  selector: '[appFillColor]',
  standalone: true
})
export class FillColorDirective {
  @Input() changeColor: string = 'blue'; // Color por defecto

  constructor(private el: ElementRef) {}

  @HostListener('click') onClick() {
    this.applyChangeColor(this.changeColor);
  }

  private applyChangeColor(color: string) {
    this.el.nativeElement.style.fill = color;
  }
}
