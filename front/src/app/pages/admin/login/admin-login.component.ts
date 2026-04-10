import { Component, inject, signal, effect } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-admin-login',
  standalone: true,
  imports: [FormsModule],
  template: `
    <div class="min-h-screen bg-[#f3f4f6] flex items-center justify-center p-8">
      <div class="w-[min(460px,100%)] bg-white rounded-3xl shadow-[0_8px_40px_rgba(0,0,0,.10)] overflow-hidden">

        <!-- Header -->
        <div class="flex items-center gap-2 px-8 pt-6 pb-2">
          <span class="text-3xl">🇲🇦</span>
          <span class="font-pacifico text-3xl text-[#006233]">daba</span>
        </div>
        <div class="px-8 pb-2 pt-1">
          <span class="inline-flex items-center gap-1.5 text-xs font-semibold uppercase tracking-widest text-gray-400">
            <span class="material-icons-round text-[14px]">admin_panel_settings</span>
            Administration
          </span>
        </div>

        <!-- Form -->
        <div class="px-8 pb-8 pt-3">
          <h2 class="text-[1.35rem] font-bold text-[#1a1a2e]">Connexion</h2>
          <p class="text-gray-500 text-[.9rem] mt-1 mb-4">Accès réservé aux administrateurs et modérateurs.</p>

          <label class="block text-sm font-medium text-gray-700 mb-1.5">Email ou téléphone</label>
          <input
            type="text"
            [(ngModel)]="identifier"
            name="identifier"
            placeholder="admin@daba.ma"
            autocomplete="username"
            (keydown.enter)="submit()"
            class="w-full h-12 px-4 border-[1.5px] border-gray-300 rounded-xl text-base outline-none focus:border-[#006233] transition-colors font-[inherit]"
          />

          <label class="block text-sm font-medium text-gray-700 mb-1.5 mt-4">Mot de passe</label>
          <input
            type="password"
            [(ngModel)]="password"
            name="password"
            autocomplete="current-password"
            (keydown.enter)="submit()"
            class="w-full h-12 px-4 border-[1.5px] border-gray-300 rounded-xl text-base outline-none focus:border-[#006233] transition-colors font-[inherit]"
          />

          @if (error()) {
            <div class="mt-3 flex items-start gap-2 p-3 bg-red-50 border border-red-200 rounded-xl">
              <span class="material-icons-round text-red-500 text-[18px] mt-px shrink-0">error_outline</span>
              <p class="text-[.85rem] text-red-600 leading-snug">{{ error() }}</p>
            </div>
          }

          <button
            (click)="submit()"
            [disabled]="loading()"
            class="w-full h-12 bg-[#006233] hover:bg-[#005229] disabled:opacity-60 disabled:cursor-not-allowed text-white text-base font-semibold rounded-full flex items-center justify-center gap-2 border-none cursor-pointer transition-colors mt-6"
          >
            @if (loading()) {
              <span class="material-icons-round animate-spin">sync</span>
            } @else {
              Se connecter
            }
          </button>
        </div>

      </div>
    </div>
  `,
})
export class AdminLoginComponent {
  private readonly auth   = inject(AuthService);
  private readonly router = inject(Router);

  identifier = '';
  password   = '';
  loading    = signal(false);
  error      = signal('');

  constructor() {
    effect(() => {
      if (!this.auth.sessionChecked()) return;
      const user = this.auth.currentUser();
      if (!user) return;
      const role = user.role;
      const dest = (role === 'admin' || role === 'moderator')
        ? (role === 'admin' ? '/admin/roles' : '/admin/users')
        : '/';
      this.router.navigate([dest], { replaceUrl: true });
    });
  }

  submit(): void {
    this.error.set('');
    if (!this.identifier.trim() || !this.password) {
      this.error.set('Identifiant et mot de passe requis.');
      return;
    }
    this.loading.set(true);
    this.auth.login(this.identifier.trim(), this.password).subscribe({
      next: res => {
        this.loading.set(false);
        if (res.user.role !== 'admin' && res.user.role !== 'moderator') {
          this.auth.logout();
          this.error.set('Accès refusé — ce compte ne dispose pas des droits d\'administration.');
          return;
        }
        this.router.navigate([res.user.role === 'admin' ? '/admin/roles' : '/admin/users']);
      },
      error: () => {
        this.loading.set(false);
        this.error.set('Identifiant ou mot de passe incorrect.');
      },
    });
  }
}
