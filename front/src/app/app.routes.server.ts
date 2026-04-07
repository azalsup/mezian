import { RenderMode, ServerRoute } from '@angular/ssr';

export const serverRoutes: ServerRoute[] = [
  // Static pages — prerendered to HTML at build time
  { path: '',         renderMode: RenderMode.Prerender },
  { path: 'login',   renderMode: RenderMode.Prerender },
  { path: 'register', renderMode: RenderMode.Prerender },

  // Admin panel — auth-gated, client-side only
  { path: 'admin/login', renderMode: RenderMode.Client },
  { path: 'admin/roles', renderMode: RenderMode.Client },
  { path: 'admin/users', renderMode: RenderMode.Client },
  { path: 'admin',       renderMode: RenderMode.Client },

  // Catch-all — client-side
  { path: '**',          renderMode: RenderMode.Client },
];
