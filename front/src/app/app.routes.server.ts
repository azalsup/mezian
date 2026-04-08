import { RenderMode, ServerRoute } from '@angular/ssr';

export const serverRoutes: ServerRoute[] = [
  // Static pages — prerendered to HTML at build time
  { path: '',         renderMode: RenderMode.Prerender },
  { path: 'login',   renderMode: RenderMode.Prerender },
  { path: 'register', renderMode: RenderMode.Prerender },

  // Admin panel — auth-gated, client-side only
  { path: 'admin/login', renderMode: RenderMode.Prerender },
  { path: 'admin/roles', renderMode: RenderMode.Prerender },
  { path: 'admin/users', renderMode: RenderMode.Prerender },
  { path: 'admin',       renderMode: RenderMode.Prerender },

  // Catch-all — client-side
  { path: '**',          renderMode: RenderMode.Prerender },
];
