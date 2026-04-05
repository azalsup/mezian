/**
 * @mezian/sdk — SDK TypeScript officiel pour l'API Mezian
 *
 * Usage basique :
 * ```ts
 * import { MezianClient } from "@mezian/sdk";
 *
 * const sdk = new MezianClient({ baseUrl: "https://api.mezian.ma/api/v1" });
 *
 * // Browse ads without login
 * const result = await sdk.ads.list({ city: "Casablanca", limit: 20 });
 *
 * // Log in with WhatsApp OTP
 * await sdk.auth.sendOtp({ phone: "+212600000000", channel: "whatsapp" });
 * const { tokens, user } = await sdk.auth.verifyOtp({ phone: "+212600000000", code: "123456" });
 * sdk.setTokens(tokens);
 *
 * // Publish an ad
 * const ad = await sdk.ads.create({ category_id: 1, title: "Appartement...", ... });
 * ```
 */

import { HttpClient, type HttpClientOptions } from "./utils/http.js";
import { AuthResource } from "./resources/auth.js";
import { AdsResource } from "./resources/ads.js";
import { CategoriesResource } from "./resources/categories.js";
import { MediaResource } from "./resources/media.js";
import { ShopsResource } from "./resources/shops.js";
import type { AuthTokens } from "./types/index.js";

export * from "./types/index.js";
export { MezianApiError } from "./utils/http.js";

export interface MezianClientOptions {
  /** URL de base de l'API, ex: "http://localhost:8080/api/v1" */
  baseUrl: string;
  /**
   * Initial tokens (e.g. restored from localStorage).
   * Can be updated later with setTokens().
   */
  initialTokens?: AuthTokens;
  /**
   * Appelé quand les tokens sont refreshés automatiquement.
   * Utilisez-le pour persister les nouveaux tokens.
   */
  onTokensRefreshed?: (tokens: AuthTokens) => void;
  /**
   * Appelé quand le refresh échoue (session expirée).
   * Use it to redirect to the login page.
   */
  onAuthFailure?: () => void;
}

/**
 * Main entry point for the Mezian SDK.
 * Instanciez une seule fois dans votre application.
 */
export class MezianClient {
  /** Access authentication operations */
  readonly auth: AuthResource;
  /** Access ads operations */
  readonly ads: AdsResource;
  /** Access categories operations */
  readonly categories: CategoriesResource;
  /** Access media operations (upload images, YouTube) */
  readonly media: MediaResource;
  /** Access professional shops operations */
  readonly shops: ShopsResource;

  private tokens: AuthTokens | null;
  private readonly http: HttpClient;

  constructor(opts: MezianClientOptions) {
    this.tokens = opts.initialTokens ?? null;

    const httpOpts: HttpClientOptions = {
      baseUrl: opts.baseUrl,
      getAccessToken: () => this.tokens?.access_token ?? null,
      onTokensRefreshed: (t) => {
        this.tokens = t;
        opts.onTokensRefreshed?.(t);
      },
      onAuthFailure: () => {
        this.tokens = null;
        opts.onAuthFailure?.();
      },
    };

    this.http = new HttpClient(httpOpts);

    this.auth = new AuthResource(this.http);
    this.ads = new AdsResource(this.http);
    this.categories = new CategoriesResource(this.http);
    this.media = new MediaResource(this.http);
    this.shops = new ShopsResource(this.http);
  }

  /**
   * Updates tokens (after login or refresh).
   * Store tokens in localStorage/AsyncStorage and restore them
   * via initialTokens à la prochaine initialisation.
   */
  setTokens(tokens: AuthTokens): void {
    this.tokens = tokens;
  }

  /** Clears tokens (client-side logout). */
  clearTokens(): void {
    this.tokens = null;
  }

  /** Returns true if the user has an access token. */
  get isAuthenticated(): boolean {
    return this.tokens !== null;
  }
}
