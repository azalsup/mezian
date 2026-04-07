import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/home/home.component').then(m => m.HomeComponent),
    title: 'Mezian — Petites annonces au Maroc',
  },
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/auth-page/auth-page.component').then(m => m.AuthPageComponent),
    title: 'Connexion — Mezian',
    data: { screen: 'login' },
  },
  {
    path: 'register',
    loadComponent: () =>
      import('./features/auth/auth-page/auth-page.component').then(m => m.AuthPageComponent),
    title: 'Créer un compte — Mezian',
    data: { screen: 'reg-credentials' },
  },
  {
    path: 'admin',
    loadChildren: () =>
      import('./admin/admin.routes').then(m => m.ADMIN_ROUTES),
    title: 'Administration — Mezian',
  },
  { path: '**', redirectTo: '' },
];
