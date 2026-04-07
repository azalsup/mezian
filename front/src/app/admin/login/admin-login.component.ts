import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-admin-login',
  imports: [FormsModule],
  template: `
    <div class="min-h-screen bg-gray-100 flex items-center justify-center">
      <div class="bg-white rounded-xl shadow-lg w-full max-w-sm p-8">
        <div class="text-center mb-8">
          <h1 class="text-2xl font-bold text-gray-900">Administration</h1>
          <p class="text-sm text-gray-500 mt-1">Daba — accès restreint</p>
        </div>

        @if (error()) {
          <div class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm">
            {{ error() }}
          </div>
        }

        <form (ngSubmit)="submit()">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Email ou téléphone
            </label>
            <input
              type="text"
              [(ngModel)]="identifier"
              name="identifier"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="admin@daba.ma"
              autocomplete="username"
            />
          </div>

          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Mot de passe
            </label>
            <input
              type="password"
              [(ngModel)]="password"
              name="password"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              autocomplete="current-password"
            />
          </div>

          <button
            type="submit"
            [disabled]="loading()"
            class="w-full bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-medium py-2 rounded-lg text-sm transition-colors"
          >
            {{ loading() ? 'Connexion...' : 'Se connecter' }}
          </button>
        </form>
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
        if (res.user.role !== 'admin') {
          this.auth.logout();
          this.error.set('Accès refusé — compte non administrateur.');
          return;
        }
        this.router.navigate(['/admin']);
      },
      error: () => {
        this.loading.set(false);
        this.error.set('Identifiant ou mot de passe incorrect.');
      },
    });
  }
}
