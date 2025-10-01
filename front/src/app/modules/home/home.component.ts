import { Component } from '@angular/core';
import {RouterOutlet} from "@angular/router";
import {HeaderComponent} from "../../core/ui/header/header.component";
import {NgIf} from "@angular/common";
import {UserAgent} from "../../core/utils/functions/userAnget";
import { Router, NavigationEnd } from '@angular/router';
import {MenuComponent} from "../../core/ui/menu/menu.component";

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    RouterOutlet,
    HeaderComponent,
    NgIf,
    MenuComponent
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss'
})
export class HomeComponent {
  showHeaderAndMenu: boolean = true;
  protected openMenu: boolean = !UserAgent.IsMobileDevice();

  constructor(private router: Router) {
    this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        const routesWithoutHeader = ['/home/cuestionario'];
        this.showHeaderAndMenu = !routesWithoutHeader.some(route => event.url.startsWith(route));
      }
    });
  }

  public toggleMenu(): void {
    this.openMenu = !this.openMenu;
  }

}
