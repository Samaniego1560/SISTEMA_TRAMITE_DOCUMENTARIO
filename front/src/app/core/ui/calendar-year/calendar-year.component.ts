import {Component, forwardRef, OnInit, signal} from '@angular/core';
import {
  AbstractControl,
  ControlValueAccessor,
  FormsModule, NG_VALIDATORS,
  NG_VALUE_ACCESSOR,
  ValidationErrors,
  Validator
} from "@angular/forms";

@Component({
  selector: 'app-calendar-year',
  standalone: true,
  imports: [
    FormsModule
  ],
  templateUrl: './calendar-year.component.html',
  styleUrl: './calendar-year.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => CalendarYearComponent),
      multi: true,
    },
    {
      provide: NG_VALIDATORS,
      useExisting: forwardRef(() => CalendarYearComponent),
      multi: true
    }
  ],
})
export class CalendarYearComponent implements ControlValueAccessor, Validator {
  rangeYear = 10;
  startYearRange: number = 0;
  endYearRange: number = 0;
  years: number[] = [];

  yearValue = signal<string>('');

  private _onChange = (_: any) => {
  };
  private _onTouch = () => {
  };

  onFocusInput(): void {
    const currentYear = new Date().getFullYear();
    this.processRangeYear(currentYear);
  }

  processRangeYear(year: number): void {
    this.years = [];
    const yearsDiffToRangeYear = year % this.rangeYear;
    const maxYear = year + this.rangeYear - yearsDiffToRangeYear;
    for (let i = this.rangeYear; i > 0; i--) {
      this.years.push(maxYear - i);
    }

    this.startYearRange = year - yearsDiffToRangeYear;
    this.endYearRange = maxYear - 1;
  }

  showBackBlock(): void {
    this.processRangeYear(this.startYearRange - this.rangeYear)
  }

  showNextBlock(): void {
    this.processRangeYear(this.startYearRange + this.rangeYear)
  }

  onSelectYear(year: number): void {
    this.yearValue.set(year.toString());
    this._onChange(new Date(year, 0, 1).toString());
    this._onTouch();
  }

  registerOnChange(fn: any): void {
    this._onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this._onTouch = fn;
  }

  writeValue(obj: any): void {
    this.yearValue.set(obj);
  }

  validate(_: AbstractControl): ValidationErrors | null {
    return this.yearValue() ? null : {required: true};
  }
}
