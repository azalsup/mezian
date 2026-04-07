import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiClient } from './api-client';
import { AuthResponse, RegisterPayload, User } from './types';

/**
 * Thin wrapper around all /auth/** endpoints.
 * Contains no state — only pure API calls.
 * Business logic (session, signals, modal) lives in AuthService.
 */
@Injectable({ providedIn: 'root' })
export class AuthApi {
  private readonly api = inject(ApiClient);

  sendOtp(phone: string): Observable<{ message: string }> {
    return this.api.post('/auth/send-otp', {
      phone,
      channel: 'sms',
      purpose: 'login',
    });
  }

  verifyOtp(phone: string, code: string): Observable<AuthResponse> {
    return this.api.post('/auth/verify-otp', { phone, code, purpose: 'login' });
  }

  login(identifier: string, password: string): Observable<AuthResponse> {
    return this.api.post('/auth/login', { identifier, password });
  }

  register(payload: RegisterPayload): Observable<AuthResponse> {
    return this.api.post('/auth/register', payload);
  }

  /** Invalidates the refresh token on the server. */
  logout(refreshToken: string): Observable<void> {
    return this.api.post('/auth/logout', { refresh_token: refreshToken });
  }

  /** Returns the authenticated user profile. */
  me(): Observable<User> {
    return this.api.get('/auth/me');
  }
}
