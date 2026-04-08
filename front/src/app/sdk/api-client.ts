import { Injectable, inject, signal } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environments/environment';
import { ApiLogger } from './api-logger';

/**
 * Core HTTP client for the Daba backend.
 *
 * - Single entry point for every GET / POST (and other verbs) to the API.
 * - Attaches the Bearer token automatically when one is stored.
 * - Logs every request and response at debug level (silent by default).
 *
 * Usage: inject this service in the domain-specific API classes (AuthApi,
 * CategoriesApi…). Never inject it directly in components.
 */
@Injectable({ providedIn: 'root' })
export class ApiClient {
  private readonly http  = inject(HttpClient);
  private readonly log   = inject(ApiLogger);

  private get base(): string {
    if (environment.production) {
      return `https://api.${window.location.hostname}/api/v1`;
    } else {
      return environment.apiBaseUrl;
    }
  }

  /** Current Bearer token — set by AuthService after login / session restore. */
  private readonly _token = signal<string | null>(null);

  setToken(token: string | null): void {
    this._token.set(token);
    this.log.debug('token updated →', token ? '***' : 'null');
  }

  // ── Verbs ─────────────────────────────────────────────────────────────────

  get<T>(path: string): Observable<T> {
    const url = `${this.base}${path}`;
    this.log.debug('→ GET', url);
    return this.http
      .get<T>(url, { headers: this.authHeaders() })
      .pipe(tap(res => this.log.debug('← GET', url, res)));
  }

  post<T>(path: string, body?: unknown): Observable<T> {
    const url = `${this.base}${path}`;
    this.log.debug('→ POST', url, body);
    return this.http
      .post<T>(url, body ?? null, { headers: this.authHeaders() })
      .pipe(tap(res => this.log.debug('← POST', url, res)));
  }

  put<T>(path: string, body?: unknown): Observable<T> {
    const url = `${this.base}${path}`;
    this.log.debug('→ PUT', url, body);
    return this.http
      .put<T>(url, body ?? null, { headers: this.authHeaders() })
      .pipe(tap(res => this.log.debug('← PUT', url, res)));
  }

  patch<T>(path: string, body?: unknown): Observable<T> {
    const url = `${this.base}${path}`;
    this.log.debug('→ PATCH', url, body);
    return this.http
      .patch<T>(url, body ?? null, { headers: this.authHeaders() })
      .pipe(tap(res => this.log.debug('← PATCH', url, res)));
  }

  delete<T>(path: string): Observable<T> {
    const url = `${this.base}${path}`;
    this.log.debug('→ DELETE', url);
    return this.http
      .delete<T>(url, { headers: this.authHeaders() })
      .pipe(tap(res => this.log.debug('← DELETE', url, res)));
  }

  // ── Internal ──────────────────────────────────────────────────────────────

  private authHeaders(): HttpHeaders {
    const token = this._token();
    return token
      ? new HttpHeaders({ Authorization: `Bearer ${token}` })
      : new HttpHeaders();
  }
}
