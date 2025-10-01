import {Component, effect} from '@angular/core';
import {DatePipe, KeyValuePipe, NgForOf, NgIf} from "@angular/common";
import {
  AbstractControl,
  FormBuilder,
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators
} from "@angular/forms";
import {ToastService} from "../../../../core/services/toast/toast.service";
import {BlockUiComponent} from "../../../../core/ui/block-ui/block-ui.component";
import {ToastComponent} from "../../../../core/ui/toast/toast.component";
import {ApplicationNoticeModel} from "./application-notice-models/application-notice-model";
import {take} from "rxjs";
import {IResponse} from "../../../../core/models/response";
import {ApplicationNoticeService} from "../../../../core/services/application-notice/application-notice.service";
import {ICalendarSettings} from "../../../../core/models/calendar-setting";

@Component({
  selector: 'app-application-notice',
  standalone: true,
  imports: [
    DatePipe,
    FormsModule,
    NgForOf,
    NgIf,
    ReactiveFormsModule,
    BlockUiComponent,
    ToastComponent,
    KeyValuePipe
  ],
  templateUrl: './application-notice.component.html',
  styleUrl: './application-notice.component.scss',
  providers: [ToastService]
})
export class ApplicationNoticeComponent {
  protected formApplicationNotice: FormGroup;
  public isLoad: boolean = false;
  public applicationNotices: ApplicationNoticeModel[] = [];

  constructor(private _fb: FormBuilder,
              private _toastService: ToastService,
              private _applicationNoticeService: ApplicationNoticeService) {
    this.formApplicationNotice = this._fb.group({
    });

    effect(() => {
      this.loadDefault().then();
    });
  }

  private async loadDefault() {
    const response = await this.getApplicationNoticePromise()
    this.isLoad = false;
    if (response.error) {
      this.createMessage('error', response.msg)
      return
    }

    if (response.data?.length) {
      this.applicationNotices.push(...response.data)
      this.loadApplicationNotice();
    }
  }

  private loadApplicationNotice(): void {
    this.applicationNotices.forEach((item) => {
      this.addFormControl(item.id.toString(), item.description)
    })
  }

  private addFormControl(name: string, value: string): void {
    this.formApplicationNotice.addControl(
      name,
      new FormControl(value, Validators.required)
    );
  }

  get controls() {
    return Object.entries(this.formApplicationNotice.controls);
  }

  public getStatusByKey(key: string) {
    return this.applicationNotices.find((item: any) => item.id === Number(key))?.status;
  }

  public async save() {
    if (this.formApplicationNotice.invalid) {
      this.createMessage('info', 'Complete los campos!')
      return
    }

    this.applicationNotices = this.loadData()

    this.isLoad = true;
    const response = await this.updateApplicationNoticePromise(this.applicationNotices);
    this.isLoad = false;
    if (response.error) {
      this.createMessage('error', 'Error en actualizar noticia: ' + response.msg)
      return ;
    }
    this.createMessage('success', 'Actualizado correctamente!')
  }

  private loadData(): ApplicationNoticeModel[] {
    const applicationNotices: ApplicationNoticeModel[] = [];
    Object.keys(this.formApplicationNotice.controls).forEach((controlName) => {
      const control: AbstractControl = this.formApplicationNotice.get(controlName)!;
      const applicationNotice: ApplicationNoticeModel = {
        id: Number(controlName),
        description: control.value,
        status: this.getStatusByKey(controlName) || 0
      }
      applicationNotices.push(applicationNotice)
    });
    return applicationNotices;
  }

  private updateApplicationNoticePromise(applicationNotices: ApplicationNoticeModel[]): Promise<{ data: ApplicationNoticeModel[] | null, error: boolean, msg: string }> {
    return new Promise((resolve) => {
      this._applicationNoticeService.updateApplicationNotice(applicationNotices).pipe(take(1)).subscribe({
        next: (res: IResponse<any>) => resolve({ data: res.detalle, error: false, msg: '' }),
        error: () => resolve({ data: null, error: true, msg: 'Error al actualizar noticias.' })
      });
    });
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

  public addItem():void {
    const key = this.getMaxKey() + 1;
    this.applicationNotices.push({id: key, status: 0 , description: ''});
    this.addFormControl(key.toString(), '');
  }

  private getMaxKey(): number {
    const keys = Object.keys(this.formApplicationNotice.controls).map(Number);
    return Math.max(...keys);
  }

  public deleteItem(key: string): void {
    this.formApplicationNotice.removeControl(key);
  }
}
