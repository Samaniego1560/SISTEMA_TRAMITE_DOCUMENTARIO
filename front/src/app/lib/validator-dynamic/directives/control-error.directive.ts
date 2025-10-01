import { ComponentRef, Directive, ElementRef, OnDestroy, OnInit, ViewContainerRef, inject } from '@angular/core';
import { NgControl } from '@angular/forms';
import { EMPTY, Subject, fromEvent, merge, takeUntil } from 'rxjs';
import { FormSubmitDirective } from './form-submit.directive';
import {ControlErrorComponent} from "../components/control-error.component";
import {getFormControlError} from "../utils/form-control-error.util";

@Directive({
  selector: '[formControl], [formControlName]',
  standalone: true
})
export class ControlErrorsDirective implements OnInit, OnDestroy {
  private readonly ngControl = inject(NgControl);
  private readonly form = inject(FormSubmitDirective, { optional: true });
  private readonly destroy$ = new Subject<void>();
  private readonly elementRef: ElementRef<HTMLElement> = inject(ElementRef);
  private readonly vcr = inject(ViewContainerRef);

  private componentRef!: ComponentRef<ControlErrorComponent>;

  private readonly submit$ = this.form?.submit$ ?? EMPTY;
  private readonly blurEvent$ = fromEvent(this.elementRef.nativeElement, 'blur');

  ngOnInit(): void {
    const observables = [this.submit$, this.blurEvent$, this.ngControl.statusChanges!]
    merge(...observables)
      .pipe(takeUntil(this.destroy$))
      .subscribe(this._handleErrorSetText.bind(this));
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private _handleErrorSetText(): void {
    if(this.ngControl.touched) {
      const errorControl = getFormControlError(this.ngControl.control!);
      this._setError(errorControl);
      return;
    }
    this._setError('');
  }

  private _setError(text: string) {
    if (!this.componentRef) {
      this.componentRef = this.vcr.createComponent(ControlErrorComponent);
    }
    this.componentRef.instance.error = text;
  }

  public checkErrors(): void {
    this._handleErrorSetText();
  }
}
