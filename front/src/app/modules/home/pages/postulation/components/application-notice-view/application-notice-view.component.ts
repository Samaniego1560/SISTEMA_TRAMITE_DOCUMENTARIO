import {Component, effect} from '@angular/core';
import {ApplicationNoticeService} from "../../../../../../core/services/application-notice/application-notice.service";
import {ApplicationNoticeModel} from "../../../application-notice/application-notice-models/application-notice-model";
import {take} from "rxjs";
import {IResponse} from "../../../../../../core/models/response";
import {ToastService} from "../../../../../../core/services/toast/toast.service";
import {BlockUiComponent} from "../../../../../../core/ui/block-ui/block-ui.component";
import {ToastComponent} from "../../../../../../core/ui/toast/toast.component";

@Component({
  selector: 'app-application-notice-view',
  standalone: true,
  imports: [
    BlockUiComponent,
    ToastComponent
  ],
  templateUrl: './application-notice-view.component.html',
  styleUrl: './application-notice-view.component.scss',
  providers: [ToastService]
})
export class ApplicationNoticeViewComponent {
  public isLoad: boolean = false;
  public tittle: string = '';
  public subTittle: string = '';
  public applicationNotices: ApplicationNoticeModel[] = [];

  constructor(private _applicationNoticeService: ApplicationNoticeService,
              private _toastService: ToastService,) {
    effect(() => {
      this.loadApplicationNotice().then();
    });
  }

  private async loadApplicationNotice() {
    this.isLoad = true;
    const response = await this.getApplicationNoticePromise()
    this.isLoad = false;
    if (response.error) {
      this.createMessage('error', response.msg)
      return
    }

    if (response.data?.length) {
      this.tittle = response.data.find((item: any) => item.status === 2)?.description || '';
      this.subTittle = response.data.find((item: any) => item.status === 1)?.description || '';
      this.applicationNotices = response.data.filter((item: any) => item.status === 0);
    }
  }

  private getApplicationNoticePromise(): Promise<{ data: ApplicationNoticeModel[] | null, error: boolean, msg: string }> {
    return new Promise((resolve) => {
      this._applicationNoticeService.getApplicationNotice().pipe(take(1)).subscribe({
        next: (res: IResponse<any>) => resolve({ data: res.detalle, error: false, msg: '' }),
        error: () => resolve({ data: null, error: true, msg: 'Error al obtener configuracion de comunicado.' })
      });
    });
  }

  private createMessage(
    type: 'error' | 'success' | 'warning' | 'info',
    message: string
  ): void {
    this._toastService.add({
      type: type,
      message: message
    });
  }
}
