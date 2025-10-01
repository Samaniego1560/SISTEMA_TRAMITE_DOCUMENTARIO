import {Directive, ElementRef, HostListener, Input, OnChanges, SimpleChanges} from '@angular/core';

@Directive({
  selector: '[appFillColorDefault]',
  standalone: true
})
export class FillColorDefaultDirective implements OnChanges {
  @Input({alias: 'appFillColorDefault'}) colors: Record<string, string> = {}

  constructor(private el: ElementRef) {}

  ngOnChanges(changes: SimpleChanges) {
    if (changes['colors'] ?.currentValue) {
      const hasID = this.colors.hasOwnProperty(this.el.nativeElement.id)
      if(!hasID) return

      this.applyChangeColor(this.colors[this.el.nativeElement.id])
    }
  }

  private applyChangeColor(color: string) {
    this.el.nativeElement.style.fill = color;
  }
}
