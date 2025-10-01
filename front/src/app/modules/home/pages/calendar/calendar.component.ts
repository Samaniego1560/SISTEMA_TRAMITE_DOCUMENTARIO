import {Component, OnInit} from '@angular/core';
import {NgIf} from "@angular/common";
import {ICalendarSettings} from "../../../../core/models/calendar-setting";
import {take} from "rxjs";
import {IResponse} from "../../../../core/models/response";
import {CalendarSettingService} from "../../../../core/services/calendar-setting/calendar-setting.service";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ToastService} from "../../../../core/services/toast/toast.service";

@Component({
  selector: 'app-calendar',
  standalone: true,
  imports: [
    NgIf,
    BlockUiComponent,
    ToastComponent
  ],
  templateUrl: './calendar.component.html',
  styleUrl: './calendar.component.scss',
  providers: [ToastService]
})
export class CalendarComponent implements OnInit {

  public calendars: ICalendarSettings[] = [];
  public isLoading: boolean = false;

  constructor(private _calendarSettingService: CalendarSettingService,
              private _toastService: ToastService,) {
  }

  ngOnInit() {
    this.getCalendarSettings().then()
  }

  private async getCalendarSettings() {
    this.isLoading = true;
    const response = await this.getCalendarSettingPromise();
    this.isLoading = false;
    if (response.error) {
      this.eventMessage('error', response.msg)
      return
    }

    if (response.data?.length) {
      this.calendars = response.data;
    }
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

  public monthByDateString(dateString: string) {
    const date = new Date(dateString + 'T00:00:00');
    return date.toLocaleString('default', { month: 'short' }).toUpperCase();
  }

  public dayByDateString(dateString: string) {
    const date = new Date(dateString + 'T00:00:00');
    return date.toLocaleString('default', { day: '2-digit' });
  }

  public yearByDateString(dateString: string) {
    const date = new Date(dateString + 'T00:00:00');
    return date.getFullYear();
  }

}
