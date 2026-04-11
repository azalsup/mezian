import { Injectable, inject, signal, computed, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { AuthApi, ApiClient } from '../../sdk';
import { StorageService } from './storage.service';
import type { AuthTokens, AuthResponse, RegisterPayload, User } from '../../sdk';

// Re-export types so existing consumers keep their import path working.
export type { AuthTokens, AuthResponse, RegisterPayload, User };

/** Exposes auth config from environment so templates can read it. */
export const authConfig = environment.auth;

const STORAGE_KEY = 'auth_tokens';

/**
 * AuthService — session management, UI signals, modal state.
 *
 * Does NOT make HTTP calls directly. All backend communication is
 * delegated to AuthApi (which uses ApiClient).
 */
@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly authApi    = inject(AuthApi);
  private readonly apiClient  = inject(ApiClient);
  private readonly storage    = inject(StorageService);
  private readonly platformId = inject(PLATFORM_ID);

  readonly currentUser    = signal<User | null>(null);
  readonly isLoggedIn     = computed(() => this.currentUser() !== null);
  readonly isStaff        = computed(() => {
    const role = this.currentUser()?.role?.toLowerCase();
    return !!(role && role !== 'user');
  });
  readonly isAdmin        = computed(() => this.currentUser()?.role?.toLowerCase() === 'administrator' || this.currentUser()?.role?.toLowerCase() === 'admin');
  readonly sessionChecked = signal(false);

  /** Controls the auth modal */
  readonly modalOpen = signal(false);
  readonly modalMode = signal<'login' | 'register'>('login');

  private refreshToken: string | null = null;

  constructor() {
    if (isPlatformBrowser(this.platformId)) {
      this.restoreSession();
    } else {
      // SSR: no localStorage, no session to restore — mark as checked immediately.
      this.sessionChecked.set(true);
    }
  }

  // ── Modal ────────────────────────────────────────────────────────────────

  openModal(mode: 'login' | 'register'): void {
    this.modalMode.set(mode);
    this.modalOpen.set(true);
  }

  closeModal(): void {
    this.modalOpen.set(false);
  }

  // ── OTP flow ─────────────────────────────────────────────────────────────

  sendOtp(phone: string): Observable<{ message: string }> {
    return this.authApi.sendOtp(phone);
  }

  verifyOtp(phone: string, code: string): Observable<AuthResponse> {
    return this.authApi.verifyOtp(phone, code).pipe(tap(r => this.setSession(r)));
  }

  // ── Password flow ────────────────────────────────────────────────────────

  login(identifier: string, password: string): Observable<AuthResponse> {
    return this.authApi.login(identifier, password).pipe(tap(r => this.setSession(r)));
  }

  register(payload: RegisterPayload): Observable<AuthResponse> {
    return this.authApi.register(payload).pipe(tap(r => this.setSession(r)));
  }

  // ── Logout ───────────────────────────────────────────────────────────────

  logout(): void {
    if (this.refreshToken) {
      this.authApi.logout(this.refreshToken).subscribe({ error: () => {} });
    }
    this.clearSession();
  }

  // ── Session helpers ──────────────────────────────────────────────────────

  private setSession(res: AuthResponse): void {
    this.apiClient.setToken(res.tokens.access_token);
    this.refreshToken = res.tokens.refresh_token;
    this.storage.setItem(STORAGE_KEY, JSON.stringify(res.tokens));
    this.currentUser.set(res.user);
  }

  private clearSession(): void {
    this.apiClient.setToken(null);
    this.refreshToken = null;
    this.storage.removeItem(STORAGE_KEY);
    this.currentUser.set(null);
  }

  private restoreSession(): void {
    const raw = this.storage.getItem(STORAGE_KEY);
    if (!raw) {
      this.sessionChecked.set(true);
      return;
    }
    try {
      const tokens: AuthTokens = JSON.parse(raw);
      this.apiClient.setToken(tokens.access_token);
      this.refreshToken = tokens.refresh_token;
      this.authApi.me().subscribe({
        next:  user => { this.currentUser.set(user); this.sessionChecked.set(true); },
        error: ()   => { this.clearSession(); this.sessionChecked.set(true); },
      });
    } catch {
      this.storage.removeItem(STORAGE_KEY);
      this.sessionChecked.set(true);
    }
  }
}
