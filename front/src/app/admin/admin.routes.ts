import { Routes } from '@angular/router';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../core/services/auth.service';

/** Guard: user must be authenticated AND have role=admin or role=moderator. */
const adminGuard = () => {
  const auth   = inject(AuthService);
  const router = inject(Router);
  const user   = auth.currentUser();
  if (user?.role === 'admin' || user?.role === 'moderator') return true;
  return router.createUrlTree(['/admin/login']);
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
    ],
  },
];
