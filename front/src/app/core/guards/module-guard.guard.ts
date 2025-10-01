import {CanActivateFn, Router} from "@angular/router";
import {inject} from "@angular/core";
import {AuthService} from "../services/auth/auth.service";
import {MODULE_PERMISSIONS} from "../contans/areas/permissions.constant";

export const moduleGuard: CanActivateFn = (route, state) => {

  const _authService = inject(AuthService);
  const _router = inject(Router);
  const session = _authService.getSession();
  const validateModule: string = route.data['module'];

  const publicRoutes: string[] = [
    '/home/calendar',
    '/home/services',
    '/home/statistics-data',
    '/home/statistics-data-srq',
    '/home/statues-regulations',
    '/home/postulation'
  ];

  if (publicRoutes.includes(state.url)) return true;

  const hasPermission = MODULE_PERMISSIONS.some(module =>  session.id === module.user_id && module.name === validateModule)
  if (_authService.isValidSession() && hasPermission) return true;

  _router.navigateByUrl('/auth/login');
  return false;
};
