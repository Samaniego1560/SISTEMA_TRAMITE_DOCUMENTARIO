import { NgIf } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input } from '@angular/core';

@Component({
  selector: 'app-control-error',
  standalone: true,
  imports: [NgIf],
  template: `
    <div *ngIf="textError" class="text-red-400 text-xs font-['Inter']">
      <span>{{ textError }}</span>
    </div>`,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ControlErrorComponent {
  textError = '';

  @Input() set error(value: string) {
    if (value !== this.textError) {
      this.textError = value;
      this.cdr.detectChanges();
    }
  }

  constructor(private cdr: ChangeDetectorRef) {}
}
