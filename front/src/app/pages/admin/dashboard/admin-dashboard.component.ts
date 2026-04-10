import { Component, inject, computed } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';

interface QuickLink {
  label:       string;
  description: string;
  icon:        string;
  path:        string;
  adminOnly:   boolean;
}

@Component({
  selector: 'app-admin-dashboard',
  standalone: true,
  imports: [RouterLink],
  template: `
    <div class="max-w-3xl">

      <!-- Welcome banner -->
      <div class="bg-gradient-to-r from-indigo-600 to-indigo-500 rounded-2xl px-8 py-7 text-white mb-8 flex items-center gap-6">
        <div class="w-14 h-14 rounded-full bg-white/20 flex items-center justify-center shrink-0">
          <i class="fa-solid fa-user-shield text-2xl"></i>
        </div>
        <div>
          <p class="text-indigo-200 text-sm font-medium mb-0.5">Bienvenue,</p>
          <h1 class="text-2xl font-bold">{{ user()?.display_name || user()?.phone }}</h1>
          <span class="inline-flex items-center gap-1.5 mt-2 text-xs font-semibold uppercase tracking-widest bg-white/20 px-3 py-1 rounded-full">
            <i class="fa-solid {{ isAdmin() ? 'fa-shield-halved' : 'fa-user-pen' }} text-[11px]"></i>
            {{ isAdmin() ? 'Administrateur' : 'Modérateur' }}
          </span>
        </div>
      </div>

      <!-- Quick links -->
      <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-widest mb-3">Accès rapide</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        @for (link of visibleLinks; track link.path) {
          <a [routerLink]="link.path"
            class="flex items-start gap-4 bg-white border border-gray-200 hover:border-indigo-300 hover:shadow-sm rounded-xl px-5 py-4 transition-all group">
            <div class="w-10 h-10 rounded-lg bg-indigo-50 group-hover:bg-indigo-100 flex items-center justify-center shrink-0 transition-colors">
              <i class="fa-solid {{ link.icon }} text-indigo-500 text-[16px]"></i>
            </div>
            <div>
              <p class="font-semibold text-gray-900 text-sm">{{ link.label }}</p>
              <p class="text-xs text-gray-500 mt-0.5">{{ link.description }}</p>
            </div>
          </a>
        }
      </div>

      <!-- Account info -->
      <h2 class="text-sm font-semibold text-gray-500 uppercase tracking-widest mt-8 mb-3">Mon compte</h2>
      <div class="bg-white border border-gray-200 rounded-xl divide-y divide-gray-100 text-sm">
        @for (row of accountRows; track row.label) {
          <div class="flex items-center gap-3 px-5 py-3">
            <span class="text-gray-400 w-28 shrink-0">{{ row.label }}</span>
            <span class="text-gray-800 font-medium">{{ row.value || '—' }}</span>
          </div>
        }
      </div>

    </div>
  `,
})
export class AdminDashboardComponent {
  private readonly auth = inject(AuthService);

  protected readonly user    = this.auth.currentUser;
  protected readonly isAdmin = computed(() => this.auth.currentUser()?.role === 'admin');

  readonly links: QuickLink[] = [
    { label: 'Utilisateurs',  description: 'Gérer les comptes, bannissements',     icon: 'fa-users',           path: '/admin/users',      adminOnly: false },
    { label: 'Catégories',    description: 'Arborescence, vedettes, activation',    icon: 'fa-tags',            path: '/admin/categories', adminOnly: false },
    { label: 'Rôles',         description: 'Permissions et rôles personnalisés',    icon: 'fa-shield-halved',   path: '/admin/roles',      adminOnly: true  },
  ];

  get visibleLinks(): QuickLink[] {
    return this.links.filter(l => !l.adminOnly || this.isAdmin());
  }

  get accountRows(): { label: string; value: string }[] {
    const u = this.user();
    return [
      { label: 'Nom',       value: u?.display_name ?? '' },
      { label: 'Email',     value: u?.email ?? '' },
      { label: 'Téléphone', value: u?.phone ?? '' },
      { label: 'Rôle',      value: u?.role ?? '' },
    ];
  }
}
