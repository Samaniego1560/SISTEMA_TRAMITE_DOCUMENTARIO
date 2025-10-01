import {Component, effect, input} from '@angular/core';
import {ModalComponent} from "../../../../../../core/ui/modal/modal.component";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {GobService} from "../../../../../../core/services/gob/gob.service";
import {Subscription} from "rxjs";
import {HttpErrorResponse} from "@angular/common/http";
import {NgIf} from "@angular/common";
import {DomSanitizer, SafeResourceUrl} from '@angular/platform-browser';

@Component({
  selector: 'app-sisfoh',
  standalone: true,
  imports: [
    ModalComponent,
    ReactiveFormsModule,
    NgIf,
    FormsModule
  ],
  templateUrl: './sisfoh.component.html',
  styleUrl: './sisfoh.component.scss'
})
export class SisfohComponent {
  isShow: boolean = true;
  dni = input.required<string>();

  formSisfoh: FormGroup;

  private _subscription = new Subscription();
  captchaUrl: string = '';
  valueCaptcha: string = '';
  currentIndex: number = 1;
  constructor(private _fb: FormBuilder,
              private _gobService: GobService) {
    this.formSisfoh = this._fb.group({
      dni_student: ['', Validators.required],
      email_student: ['', [Validators.required, Validators.email]],
      eat_service: [false, Validators.required],
      resident_service: [false, Validators.required],
    });

    effect(() => {
      this.getCaptcha().then();
    });
  }

  async getCaptcha() {
    const response = await this.getCaptchaValuePromise();
    if (response.error) {
      console.log(response.error);
    }
    this.captchaUrl = response.data;
  }

  private getCaptchaValuePromise() : Promise<any> {
    return new Promise<any>(resolve => {
      this._subscription.add(
        this._gobService.getCaptchaValue().subscribe({
          next: async (resp: any) => {
            resolve(
              {
                error: false,
                msg: '',
                data: resp.data,
                type: 'susses',
              });
          },
          error: (err: HttpErrorResponse) => {
            resolve({
              error: true,
              msg: '',
              data: null,
              type: 'error',
              code: err.status
            });
            console.error(err);
          }
        })
      );
    })
  }

  nextProcess() {
    if (this.currentIndex === 1) {
      this.validateDni().then();
    }
    this.currentIndex++;
  }

  backProcess() {
    this.currentIndex--;
  }

  private async validateDni() {

    this._gobService.validaDNI('76797001', this.valueCaptcha).subscribe(
      response => {
        console.log('Respuesta del servidor:', response);
      },
      error => {
        console.error('Error en la petición:', error);
      }
    );
  }

}
