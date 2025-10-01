import {inject, Injectable} from '@angular/core';
import {ActionType, MODULE_PERMISSIONS} from "../../contans/areas/permissions.constant";
import {AuthService} from "../auth/auth.service";

@Injectable({
  providedIn: 'root'
})
export class PermissionsService {

  private authService = inject(AuthService);

  getModuleByName(moduleName: string) {
    const session = this.authService.getSession();
    const modules = MODULE_PERMISSIONS.filter(area => area.user_id === session.id)
    if(modules.length === 0) return undefined;

    return modules.find(module => module.name === moduleName);
  }

  hasPermission(moduleName: string, componentName: string, actionType: ActionType): boolean {
    const module = this.getModuleByName(moduleName)
    if (module === undefined) return false;

    const component = module.components.find(component => component.name === componentName)
    if (component === undefined) return false;

    return component.actions.some(a =>
      a.component_id === component.id &&
      a.type === actionType
    );
  }
}
