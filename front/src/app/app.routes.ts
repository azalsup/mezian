import { Routes } from '@angular/router';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from './core/services/auth.service';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/public/home/home.component').then(m => m.HomeComponent),
    title: 'Daba — Petites annonces au Maroc',
  },
  {
    path: 'login',
    loadComponent: () =>
      import('./pages/public/auth/auth-page/auth-page.component').then(m => m.AuthPageComponent),
    title: 'Connexion — Daba',
    data: { screen: 'login' },
  },
  {
    path: 'register',
    loadComponent: () =>
      import('./pages/public/auth/auth-page/auth-page.component').then(m => m.AuthPageComponent),
    title: 'Créer un compte — Daba',
    data: { screen: 'reg-credentials' },
  },
  {
    path: 'ads',
    loadComponent: () =>
      import('./pages/public/ads/ads-page.component').then(m => m.AdsPageComponent),
    title: 'Annonces — Daba',
  },
  {
    path: 'ads/:id',
    loadComponent: () =>
      import('./pages/public/ad-detail/ad-detail-page.component').then(m => m.AdDetailPageComponent),
    title: 'Annonce — Daba',
  },
  {
    path: 'admin',
    loadChildren: () =>
      import('./pages/admin/admin.routes').then(m => m.ADMIN_ROUTES),
    title: 'Administration — Daba',
  },
  {
    path: 'post-ad',
    loadComponent: () =>
      import('./pages/private/post-ad/post-ad-page.component').then(m => m.PostAdPageComponent),
    title: 'Déposer une annonce — Daba',
    canActivate: [() => {
      const auth   = inject(AuthService);
      const router = inject(Router);
      return auth.isLoggedIn() ? true : router.createUrlTree(['/login']);
    }],
  },
  { path: '**', redirectTo: '' },
];
