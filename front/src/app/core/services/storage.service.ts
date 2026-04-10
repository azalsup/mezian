import { Injectable, inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';

export type StorageType = 'localStorage' | 'sessionStorage' | 'cookie';

@Injectable({ providedIn: 'root' })
export class StorageService {
  private readonly platformId = inject(PLATFORM_ID);

  setItem(key: string, value: string, type: StorageType = 'localStorage'): void {
    if (!isPlatformBrowser(this.platformId)) return;

    switch (type) {
      case 'localStorage':
        localStorage.setItem(key, value);
        break;
      case 'sessionStorage':
        sessionStorage.setItem(key, value);
        break;
      case 'cookie':
        this.setCookie(key, value);
        break;
    }
  }

  getItem(key: string, type: StorageType = 'localStorage'): string | null {
    if (!isPlatformBrowser(this.platformId)) return null;

    switch (type) {
      case 'localStorage':
        return localStorage.getItem(key);
      case 'sessionStorage':
        return sessionStorage.getItem(key);
      case 'cookie':
        return this.getCookie(key);
      default:
        return null;
    }
  }

  removeItem(key: string, type: StorageType = 'localStorage'): void {
    if (!isPlatformBrowser(this.platformId)) return;

    switch (type) {
      case 'localStorage':
        localStorage.removeItem(key);
        break;
      case 'sessionStorage':
        sessionStorage.removeItem(key);
        break;
      case 'cookie':
        this.removeCookie(key);
        break;
    }
  }

  private setCookie(key: string, value: string, days: number = 30): void {
    const expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
    document.cookie = `${key}=${value};expires=${expires.toUTCString()};path=/;SameSite=Strict`;
  }

  private getCookie(key: string): string | null {
    const name = key + '=';
    const decodedCookie = decodeURIComponent(document.cookie);
    const ca = decodedCookie.split(';');
    for (let c of ca) {
      c = c.trim();
      if (c.indexOf(name) === 0) {
        return c.substring(name.length);
      }
    }
    return null;
  }

  private removeCookie(key: string): void {
    document.cookie = `${key}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;SameSite=Strict`;
  }
}