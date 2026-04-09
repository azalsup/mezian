import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/home/home.component').then(m => m.HomeComponent),
    title: 'Daba — Petites annonces au Maroc',
  },
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/auth-page/auth-page.component').then(m => m.AuthPageComponent),
    title: 'Connexion — Daba',
    data: { screen: 'login' },
  },
  {
    path: 'register',
    loadComponent: () =>
      import('./features/auth/auth-page/auth-page.component').then(m => m.AuthPageComponent),
    title: 'Créer un compte — Daba',
    data: { screen: 'reg-credentials' },
  },
  {
    path: 'ads',
    loadComponent: () =>
      import('./pages/ads/ads-page.component').then(m => m.AdsPageComponent),
    title: 'Annonces — Daba',
  },
  {
    path: 'ads/:id',
    loadComponent: () =>
      import('./pages/ad-detail/ad-detail-page.component').then(m => m.AdDetailPageComponent),
    title: 'Annonce — Daba',
  },
  {
    path: 'admin',
    loadChildren: () =>
      import('./admin/admin.routes').then(m => m.ADMIN_ROUTES),
    title: 'Administration — Daba',
  },
  { path: '**', redirectTo: '' },
];
