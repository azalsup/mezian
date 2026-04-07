import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

export type LogLevel = 'debug' | 'info' | 'warn' | 'error' | 'none';

const LEVEL_RANK: Record<LogLevel, number> = {
  debug: 0,
  info:  1,
  warn:  2,
  error: 3,
  none:  99,
};

/**
 * Centralized logger for the Mezian SDK.
 * Log level is controlled by environment.logLevel (default: 'warn').
 * Set to 'debug' to trace every HTTP request/response.
 */
@Injectable({ providedIn: 'root' })
export class ApiLogger {
  private readonly minRank = LEVEL_RANK[environment.logLevel ?? 'warn'];

  debug(...args: unknown[]): void { this.emit('debug', ...args); }
  info (...args: unknown[]): void { this.emit('info',  ...args); }
  warn (...args: unknown[]): void { this.emit('warn',  ...args); }
  error(...args: unknown[]): void { this.emit('error', ...args); }

  private emit(level: LogLevel, ...args: unknown[]): void {
    if (LEVEL_RANK[level] < this.minRank) return;
    const tag = `[mezian-sdk][${level.toUpperCase()}]`;
    switch (level) {
      case 'debug': console.debug(tag, ...args); break;
      case 'info':  console.info (tag, ...args); break;
      case 'warn':  console.warn (tag, ...args); break;
      case 'error': console.error(tag, ...args); break;
    }
  }
}
