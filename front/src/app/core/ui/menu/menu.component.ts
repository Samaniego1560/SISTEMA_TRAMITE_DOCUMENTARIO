import {Component} from '@angular/core';
import {RouterLink, RouterLinkActive} from "@angular/router";
import {AuthService} from "../../services/auth/auth.service";
import {NgClass, NgFor, NgIf} from "@angular/common";
import {HttpClientModule} from "@angular/common/http";
import {Store} from "@ngrx/store";
import {AppState} from "../../store/app.reducers";
import { EnvServiceFactory } from '../../services/env/env.service.provider';
import { menuRoutes } from './menuRoutes';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Component({
  selector: 'app-menu',
  standalone: true,
  imports: [
    RouterLink,
    RouterLinkActive,
    NgIf,
    NgFor,
    HttpClientModule,
    NgClass
  ],
  templateUrl: './menu.component.html',
  styleUrl: './menu.component.scss',
  providers: [AuthService]
})
export class MenuComponent {
  protected isAuth: boolean = false;
  protected role: number = 0;
  menu: any[] = [];
  rolePsicopedagogia = EnvServiceFactory().ROLE_PSICOPEDAGOGIA;
  showVisitOptions: boolean = false;


  constructor(private _store: Store<AppState>,  private sanitizer: DomSanitizer) {
    this._store.select('auth').subscribe((auth) => {
      this.isAuth = auth.isAuth;
      this.role = auth.role;
    });
    if(this.role === this.rolePsicopedagogia || this.role === 1){
      this.menu = menuRoutes.filter(item => item.roles.includes(this.role));
    }
  }

  showSubOptions = false;
  showSubMenuResidence = false;
  showPsicopedagogiaOptions = false;
  showSubMenuNursing = false;
  showSubMenuMedicine = false;
  showSubMenuDentistry = false;

  getSafeHtml(html: string): SafeHtml {
    return this.sanitizer.bypassSecurityTrustHtml(html);
  }

  onToogleResidence() {
    this.showSubMenuResidence = !this.showSubMenuResidence;
  }

  toggleSubOptions() {
    this.showSubOptions = !this.showSubOptions;

  }

  toggleVisitOptions(): void {
    this.showVisitOptions = !this.showVisitOptions;
  }

  togglePsicopedagogiaOptions(): void {
    this.showPsicopedagogiaOptions = !this.showPsicopedagogiaOptions;
  }
}
