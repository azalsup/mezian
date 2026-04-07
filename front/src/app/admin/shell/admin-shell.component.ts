import { Component, inject } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { Router } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

interface NavItem {
  label: string;
  icon:  string;
  path:  string;
}

@Component({
  selector: 'app-admin-shell',
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  template: `
    <div class="flex h-screen bg-gray-50 overflow-hidden">

      <!-- ── Sidebar ──────────────────────────────────────────────────────── -->
      <aside class="w-60 flex-shrink-0 bg-gray-900 text-gray-100 flex flex-col">

        <!-- Logo -->
        <div class="px-6 py-5 border-b border-gray-700">
          <span class="text-lg font-bold tracking-wide text-white">Mezian</span>
          <span class="ml-2 text-xs bg-indigo-600 text-white px-2 py-0.5 rounded-full">Admin</span>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
          @for (item of navItems; track item.path) {
            <a
              [routerLink]="item.path"
              routerLinkActive="bg-gray-700 text-white"
              [routerLinkActiveOptions]="{ exact: false }"
              class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-gray-300 hover:bg-gray-700 hover:text-white transition-colors"
            >
              <i class="fa-solid {{ item.icon }} w-4 text-center"></i>
              {{ item.label }}
            </a>
          }
        </nav>

        <!-- User + logout -->
        <div class="px-4 py-4 border-t border-gray-700">
          <p class="text-xs text-gray-400 truncate mb-3">
            {{ auth.currentUser()?.display_name }}
          </p>
          <button
            (click)="logout()"
            class="flex items-center gap-2 text-xs text-gray-400 hover:text-white transition-colors"
          >
            <i class="fa-solid fa-arrow-right-from-bracket"></i>
            Déconnexion
          </button>
        </div>
      </aside>

      <!-- ── Main content ──────────────────────────────────────────────────── -->
      <div class="flex-1 flex flex-col overflow-hidden">

        <!-- Top bar -->
        <header class="bg-white border-b border-gray-200 px-8 py-4 flex-shrink-0">
          <h1 class="text-xl font-semibold text-gray-800">Panneau d'administration</h1>
        </header>

        <!-- Page content -->
        <main class="flex-1 overflow-y-auto p-8">
          <router-outlet />
        </main>
      </div>
    </div>
  `,
})
export class AdminShellComponent {
  protected readonly auth   = inject(AuthService);
  private   readonly router = inject(Router);

  readonly navItems: NavItem[] = [
    { label: 'Rôles',         icon: 'fa-shield-halved', path: '/admin/roles' },
    { label: 'Utilisateurs',  icon: 'fa-users',         path: '/admin/users' },
  ];

  logout(): void {
    this.auth.logout();
    this.router.navigate(['/admin/login']);
  }
}
