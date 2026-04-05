import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/home/home.component').then((m) => m.HomeComponent),
    title: 'Mezian — Petites annonces au Maroc',
  },
  // Future routes
  // { path: 'ads', ... }
  // { path: 'ads/:slug', ... }
  // { path: 'auth', ... }
  { path: '**', redirectTo: '' },
];
