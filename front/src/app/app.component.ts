import {Component, OnDestroy} from '@angular/core';
import {CommonModule} from '@angular/common';
import {Router, RouterOutlet} from '@angular/router';
import {HttpClientModule, HttpErrorResponse} from "@angular/common/http";
import {Store} from "@ngrx/store";
import {AppState} from "./core/store/app.reducers";
import {AuthService} from "./core/services/auth/auth.service";
import {controlAuth} from "./core/store/actions/auth.action";
import {interval, Subscription} from "rxjs";
import {UbigeoService} from "./core/services/ubigeo/ubigeo.service";
import {ToastService} from "./core/services/toast/toast.service";
import {ToastComponent} from "./core/ui/toast/toast.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, HttpClientModule, ToastComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  providers: [AuthService, ToastService]
})
export class AppComponent implements OnDestroy {
  title = 'dibuWebApp';
  private _subscriptions: Subscription = new Subscription();

  constructor(
    private _store: Store<AppState>,
    private _authService: AuthService,
    private _router: Router,
    private _ubigeosService: UbigeoService,
    private _toastService: ToastService,
  ) {
    this._store.dispatch(controlAuth({
      auth: {
        isAuth: this._authService.isValidSession(),
        session: this._authService.getSession(),
        token: this._authService.getToken(),
        role: this._authService.getRole()
      }
    }));

    this._store.select('auth').subscribe((auth) => {
      if (auth.isAuth) {
        this._subscriptions = new Subscription();
        this.timerVerifySession();
        return;
      }

      this._subscriptions.unsubscribe();
    });

    if (this._authService.getToken() !== '') this.verifySession();
    this.loadDepartment().then();
    this.loadProvince().then();
    this.loadDistrict().then();
  }

  ngOnDestroy() {
    this._subscriptions.unsubscribe();
  }

  private timerVerifySession(): void {
    this._subscriptions.add(
      interval(60000)
        .subscribe(() => this.verifySession())
    );
  }

  private verifySession(): void {
    if (!this._authService.isValidSession()) {

      this._authService.logout();
      this._router.navigate(['/home/services']);
      return;
    }
  }

  private setLocalStorage(key : string, value: string): void {
    localStorage.setItem(key, value);
  }

  public async loadDepartment() {
    const response = await this.getDepartmentPromise()
    if (response.error) {
      this.createMessage('error', 'Error could not load department');
      return;
    }
    this.setLocalStorage('departments', JSON.stringify(response.data));
  }

  public async loadProvince() {
    const response = await this.getProvincePromise()
    if (response.error) {
      this.createMessage('error', 'Error could not load prvince');
      return;
    }
    this.setLocalStorage('provinces', JSON.stringify(response.data));
  }

  public async loadDistrict() {
    const response = await this.getDistrictPromise()
    if (response.error) {
      this.createMessage('error', 'Error could not load district');
      return;
    }
    this.setLocalStorage('districts', JSON.stringify(response.data));
  }

  private getDepartmentPromise() : Promise<any> {
    return new Promise<any>(resolve => {
      //this._subscriptions.add( //TODO: Corregir el manejo de subscripción
        this._ubigeosService.getUbigeosDepartaments().subscribe({
          next: async (resp: any) => {
            resolve({error: false, msg: '', data: resp, type: 'susses',});
          },
          error: (err: HttpErrorResponse) => {
            debugger
            resolve({ error: true, msg: '', data: null, type: 'error', code: err.status });
          }
        })
      //);
    })
  }

  private getProvincePromise() : Promise<any> {
    return new Promise<any>(resolve => {
      //this._subscriptions.add( //TODO: Corregir el manejo de subscripción
        this._ubigeosService.getUbigeosProvinces().subscribe({
          next: async (resp: any) => {
            resolve({error: false, msg: '', data: resp, type: 'susses',});
          },
          error: (err: HttpErrorResponse) => {
            resolve({ error: true, msg: '', data: null, type: 'error', code: err.status });
          }
        })
      //);
    })
  }

  private getDistrictPromise() : Promise<any> {
    return new Promise<any>(resolve => {
      //this._subscriptions.add( //TODO: Corregir el manejo de subscripción
        this._ubigeosService.getUbigeosDistricts().subscribe({
          next: async (resp: any) => {
            resolve({error: false, msg: '', data: resp, type: 'susses',});
          },
          error: (err: HttpErrorResponse) => {
            resolve({ error: true, msg: '', data: null, type: 'error', code: err.status });
          }
        })
      //);
    })
  }

  private createMessage(type: 'error' | 'success' | 'warning' | 'info', message: string): void {
    this._toastService.add({
      type: type,
      message: message
    });
  }
}
