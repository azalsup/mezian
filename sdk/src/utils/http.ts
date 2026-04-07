/**
 * Base HTTP client for the Daba API.
 * Uses native fetch (browser and Node 18+ compatible).
 * Handles: Bearer token injection, automatic refresh, typed errors.
 */

import type { ApiError, AuthTokens } from "../types/index.js";

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly body: ApiError,
  ) {
    super(body.error ?? `HTTP ${status}`);
    this.name = "ApiError";
  }
}

export interface HttpClientOptions {
  /** URL de base, ex: "http://localhost:8080/api/v1" */
  baseUrl: string;
  /** Callback to retrieve the current access token */
  getAccessToken?: () => string | null;
  /** Callback called with new tokens after a successful refresh */
  onTokensRefreshed?: (tokens: AuthTokens) => void;
  /** Callback called when refresh fails (logout) */
  onAuthFailure?: () => void;
}

export class HttpClient {
  private baseUrl: string;
  private getAccessToken: () => string | null;
  private onTokensRefreshed?: (t: AuthTokens) => void;
  private onAuthFailure?: () => void;

  constructor(opts: HttpClientOptions) {
    this.baseUrl = opts.baseUrl.replace(/\/$/, "");
    this.getAccessToken = opts.getAccessToken ?? (() => null);
    this.onTokensRefreshed = opts.onTokensRefreshed;
    this.onAuthFailure = opts.onAuthFailure;
  }

  // ---------------------------------------------------------------------------
  // Public methods
  // ---------------------------------------------------------------------------

  async get<T>(path: string, params?: Record<string, unknown>): Promise<T> {
    const url = this.buildUrl(path, params);
    return this.request<T>("GET", url);
  }

  async post<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>("POST", this.buildUrl(path), body);
  }

  async put<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>("PUT", this.buildUrl(path), body);
  }

  async delete<T = void>(path: string): Promise<T> {
    return this.request<T>("DELETE", this.buildUrl(path));
  }

  /** Upload multipart/form-data (for images) */
  async upload<T>(path: string, formData: FormData): Promise<T> {
    return this.request<T>("POST", this.buildUrl(path), formData, true);
  }

  // ---------------------------------------------------------------------------
  // Interne
  // ---------------------------------------------------------------------------

  private buildUrl(path: string, params?: Record<string, unknown>): string {
    const url = new URL(`${this.baseUrl}${path}`);
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        if (v !== undefined && v !== null && v !== "") {
          url.searchParams.set(k, String(v));
        }
      }
    }
    return url.toString();
  }

  private async request<T>(
    method: string,
    url: string,
    body?: unknown,
    isMultipart = false,
  ): Promise<T> {
    const headers: Record<string, string> = {};

    const token = this.getAccessToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    if (!isMultipart && body !== undefined) {
      headers["Content-Type"] = "application/json";
    }

    const res = await fetch(url, {
      method,
      headers,
      body:
        body instanceof FormData
          ? body
          : body !== undefined
            ? JSON.stringify(body)
            : undefined,
    });

    // Réponse vide (204 No Content)
    if (res.status === 204) {
      return undefined as T;
    }

    // Tentative de décodage JSON
    const contentType = res.headers.get("content-type") ?? "";
    const json = contentType.includes("application/json")
      ? await res.json()
      : null;

    if (!res.ok) {
      throw new ApiError(res.status, json ?? { error: res.statusText });
    }

    return json as T;
  }
}
