import {
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output,
  forwardRef,
  input,
  OnChanges,
  SimpleChanges
} from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

@Component({
  selector: 'app-toggle-switch',
  templateUrl: './toggle-switch.component.html',
  styleUrls: ['./toggle-switch.component.scss'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => ToggleSwitchComponent),
      multi: true,
    },
  ],
  standalone: true
})
export class ToggleSwitchComponent implements OnChanges, OnInit, ControlValueAccessor {
  trueValue = input<any>(true);
  falseValue = input<any>(false)
  public currentValue: any;
  @Input() label: string = '';
  @Input() id: string = 'id';
  @Input() value: any = '';
  @Input() set disabled(isDisabled: boolean) {
    this.isDisabled = isDisabled;
  }

  @Output() change: EventEmitter<any> = new EventEmitter<any>();

  public isDisabled: boolean = false;

  constructor() {}

  ngOnChanges(changes: SimpleChanges) {
    if (changes['value'] && changes['value'].currentValue !== undefined) {
      this.writeValue(changes['value'].currentValue);
    }
  }

  ngOnInit(): void {
    this.writeValue(this.falseValue());
  }

  writeValue(obj: any): void {
    this.currentValue = obj;
  }

  registerOnChange(fn: any): void {
    this._onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this._onTouch = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this.isDisabled = isDisabled;
  }

  private _onChange = (_: any) => {};
  private _onTouch = () => {};

  public handleChange(event: Event): void {
    event.stopPropagation();
    this._updateModel()
    this._onChange(this.currentValue);
    this._onTouch();
    this.change.emit(this.currentValue);
  }

  private _updateModel(){
    this.currentValue = this.currentValue === this.trueValue() ? this.falseValue() : this.trueValue()
  }

  public get checked(): boolean {
    return this.currentValue === this.trueValue();
  }
}
