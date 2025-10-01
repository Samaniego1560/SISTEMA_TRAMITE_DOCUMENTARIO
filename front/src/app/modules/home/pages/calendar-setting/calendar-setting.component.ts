import {Component, effect} from '@angular/core';
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {ExportComponent} from "../../../../core/utils/export/export.component";
import {FormArray, FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {ModalComponent} from "../../../../core/ui/modal/modal.component";
import {DatePipe, NgForOf, NgIf} from "@angular/common";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {CalendarSettingService} from "../../../../core/services/calendar-setting/calendar-setting.service";
import {take} from "rxjs";
import {IResponse} from "../../../../core/models/response";
import {ICalendarSettings} from "../../../../core/models/calendar-setting";

@Component({
  selector: 'app-calendar-setting',
  standalone: true,
  imports: [
    BlockUiComponent,
    ExportComponent,
    FormsModule,
    ModalComponent,
    NgForOf,
    NgIf,
    ReactiveFormsModule,
    ToastComponent,
    DatePipe
  ],
  templateUrl: './calendar-setting.component.html',
  styleUrl: './calendar-setting.component.scss',
  providers: [ToastService]
})
export class CalendarSettingComponent {
  protected formCalendar: FormGroup;
  public isLoading: boolean = false;

  constructor(private _fb: FormBuilder,
             private _toastService: ToastService,
              private _calendarSettingService: CalendarSettingService,) {
    this.formCalendar = this._fb.group({
      items: this._fb.array([]),
    });


    effect(() => {
      this.createDefault().then()
    });
  }

  private async createDefault() {
    this.isLoading = true;
    const response = await this.getCalendarSettingPromise();
    this.isLoading = false;
    if (response.error) {
      this.eventMessage('error', response.msg)
      return
    }

    if (response.data?.length) {
      response.data.forEach(value => {
        const formGroup = this.createFormItem();
        formGroup.patchValue(value);
        this.items.push(formGroup);
      });
    }

  }

  public addItem() {
    const items = this.formCalendar.get('items') as FormArray;
    items.push(this.createFormItem());
  }


  get items() : FormArray {
    return this.formCalendar.get('items') as FormArray;
  }

  public createFormItem(): FormGroup {
    return this._fb.group({
      start_date: ['', [Validators.required]],
      end_date: [null, []],
      title: ['', [Validators.required]],
      description: ['', [Validators.required]],
    });
  }

  public deleteItem(index: number) {
    const items = this.formCalendar.get('items') as FormArray;
    if (items.length > index && index >= 0) {
      items.removeAt(index);
    }
  }

  async saveForm() {
    if (this.formCalendar.invalid) {
      this.eventMessage('info', 'Complete los campos correctamente!.')
      return
    }

    const items = this.formCalendar.get('items') as FormArray;
    const data: ICalendarSettings[] = items.value.map((item: any) => ({
      start_date: item.start_date,
      end_date: item.end_date === '' ? null : item.end_date,
      title: item.title,
      description: item.description
    }));

    this.isLoading = true;
    const response = await this.updateCalendarSettingPromise(data);
    this.isLoading = false;
    if (response.error) {
      this.eventMessage('error', response.msg)
      return
    }

    this.eventMessage('success', 'Actualizado correctamente!')
  }

  private updateCalendarSettingPromise(calendarSetting: ICalendarSettings[]): Promise<{ data: ICalendarSettings[] | null, error: boolean, msg: string }> {
    return new Promise((resolve) => {
      this._calendarSettingService.updateCalendarSetting(calendarSetting).pipe(take(1)).subscribe({
        next: (res: IResponse<any>) => resolve({ data: res.detalle, error: false, msg: '' }),
        error: () => resolve({ data: null, error: true, msg: 'Error al actualizar calendario.' })
      });
    });
  }

  private getCalendarSettingPromise(): Promise<{ data: ICalendarSettings[] | null, error: boolean, msg: string }> {
    return new Promise((resolve) => {
      this._calendarSettingService.getCalendarSetting().pipe(take(1)).subscribe({
        next: (res: IResponse<any>) => resolve({ data: res.detalle, error: false, msg: '' }),
        error: () => resolve({ data: null, error: true, msg: 'Error al obtener configuracion de calendario.' })
      });
    });
  }

  public eventMessage(type: "error" | "success" | "warning" | "info", message: string): void {
    this._toastService.add({
      type: type,
      message: message,
    });
  }
}
