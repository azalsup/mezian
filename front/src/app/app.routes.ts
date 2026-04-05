import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/home/home.component').then((m) => m.HomeComponent),
    title: 'Mezian — Petites annonces Maroc',
  },
  // Futures routes
  // { path: 'annonces', ... }
  // { path: 'annonces/:slug', ... }
  // { path: 'auth', ... }
  { path: '**', redirectTo: '' },
];
