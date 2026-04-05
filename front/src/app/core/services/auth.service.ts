import { Injectable, inject, signal, computed } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  token_type: 'Bearer';
  expires_in: number;
}

export interface User {
  id: number;
  phone: string;
  email?: string;
  is_verified: boolean;
  display_name: string;
  avatar_url?: string;
  city?: string;
  role: 'user' | 'admin';
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  tokens: AuthTokens;
  user: User;
}

const STORAGE_KEY = 'mezian_tokens';
const API = environment.apiBaseUrl;

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);

  readonly currentUser = signal<User | null>(null);
  readonly isLoggedIn  = computed(() => this.currentUser() !== null);

  // Modal visibility + mode
  readonly modalOpen = signal(false);
  readonly modalMode = signal<'login' | 'register'>('login');

  private accessToken: string | null  = null;
  private refreshToken: string | null = null;

  constructor() {
    this.restoreSession();
  }

  // ── Modal ────────────────────────────────────────────────────────────────

  openModal(mode: 'login' | 'register'): void {
    this.modalMode.set(mode);
    this.modalOpen.set(true);
  }

  closeModal(): void {
    this.modalOpen.set(false);
  }

  // ── Auth flows ───────────────────────────────────────────────────────────

  sendOtp(phone: string): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${API}/auth/send-otp`, {
      phone,
      channel: 'sms',
    });
  }

  verifyOtp(phone: string, code: string): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${API}/auth/verify-otp`, { phone, code }).pipe(
      tap(res => this.setSession(res)),
    );
  }

  logout(): void {
    if (this.refreshToken && this.accessToken) {
      this.http
        .post(
          `${API}/auth/logout`,
          { refresh_token: this.refreshToken },
          { headers: new HttpHeaders({ Authorization: `Bearer ${this.accessToken}` }) },
        )
        .subscribe({ error: () => {} });
    }
    this.clearSession();
  }

  getAuthHeader(): HttpHeaders {
    return new HttpHeaders({ Authorization: `Bearer ${this.accessToken ?? ''}` });
  }

  // ── Session helpers ──────────────────────────────────────────────────────

  private setSession(res: AuthResponse): void {
    this.accessToken  = res.tokens.access_token;
    this.refreshToken = res.tokens.refresh_token;
    localStorage.setItem(STORAGE_KEY, JSON.stringify(res.tokens));
    this.currentUser.set(res.user);
    this.closeModal();
  }

  private clearSession(): void {
    this.accessToken  = null;
    this.refreshToken = null;
    localStorage.removeItem(STORAGE_KEY);
    this.currentUser.set(null);
  }

  private restoreSession(): void {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return;

    try {
      const tokens: AuthTokens = JSON.parse(raw);
      this.accessToken  = tokens.access_token;
      this.refreshToken = tokens.refresh_token;
      this.loadMe();
    } catch {
      localStorage.removeItem(STORAGE_KEY);
    }
  }

  private loadMe(): void {
    this.http
      .get<User>(`${API}/auth/me`, { headers: this.getAuthHeader() })
      .subscribe({
        next:  user => this.currentUser.set(user),
        error: ()   => this.clearSession(),
      });
  }
}
