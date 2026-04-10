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
    role === 'admin' || role === 'moderator';

  if (auth.sessionChecked()) {
    return allow(auth.currentUser()?.role)
      ? true
      : router.createUrlTree(['/admin/login']);
  }

  return toObservable(auth.sessionChecked, { injector }).pipe(
    filter(checked => checked),
    take(1),
    map(() =>
      allow(auth.currentUser()?.role)
        ? true
        : router.createUrlTree(['/admin/login']),
    ),
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
        canActivate: [() => {
          const auth   = inject(AuthService);
          const router = inject(Router);
          const role   = auth.currentUser()?.role;
          return router.createUrlTree([role === 'admin' ? '/admin/roles' : '/admin/users']);
        }],
        loadComponent: () =>
          import('./users/admin-users.component').then(m => m.AdminUsersComponent),
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
