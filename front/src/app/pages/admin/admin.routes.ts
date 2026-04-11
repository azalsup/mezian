import { Routes } from '@angular/router';
import { inject, Injector } from '@angular/core';
import { Router } from '@angular/router';
import { filter, take, map } from 'rxjs';
import { toObservable } from '@angular/core/rxjs-interop';
import { AuthService } from '../../core/services/auth.service';

/**
 * Protects admin routes.
 * - If session is already resolved → evaluate synchronously.
 * - If session is still loading  → wait for sessionChecked, then evaluate.
 */
const adminGuard = () => {
  const auth     = inject(AuthService);
  const router   = inject(Router);
  const injector = inject(Injector);

  const allow = (role: string | undefined) =>
    role?.toLowerCase() !== 'user';

  const redirect = (user: any) => {
    if (!user) return router.createUrlTree(['/login']);
    if (!allow(user.role)) return router.createUrlTree(['/']);
    return true;
  };

  if (auth.sessionChecked()) {
    return redirect(auth.currentUser());
  }

  return toObservable(auth.sessionChecked, { injector }).pipe(
    filter(checked => checked),
    take(1),
    map(() => redirect(auth.currentUser())),
  );
};

export const ADMIN_ROUTES: Routes = [
  {
    path: 'login',
    loadComponent: () =>
      import('./login/admin-login.component').then(m => m.AdminLoginComponent),
  },
  {
    path: '',
    canActivate: [adminGuard],
    loadComponent: () =>
      import('./shell/admin-shell.component').then(m => m.AdminShellComponent),
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./admin-main.component').then(m => m.AdminMainComponent),
      },
      {
        path: 'roles',
        loadComponent: () =>
          import('./roles/admin-roles.component').then(m => m.AdminRolesComponent),
      },
      {
        path: 'users',
        loadComponent: () =>
          import('./users/admin-users.component').then(m => m.AdminUsersComponent),
      },
      {
        path: 'categories',
        loadComponent: () =>
          import('./categories/admin-categories.component').then(m => m.AdminCategoriesComponent),
      },
    ],
  },
];
