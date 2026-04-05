/**
 * @mezian/sdk — SDK TypeScript officiel pour l'API Mezian
 *
 * Usage basique :
 * ```ts
 * import { MezianClient } from "@mezian/sdk";
 *
 * const sdk = new MezianClient({ baseUrl: "https://api.mezian.ma/api/v1" });
 *
 * // Parcourir les annonces sans connexion
 * const result = await sdk.ads.list({ city: "Casablanca", limit: 20 });
 *
 * // Se connecter par OTP WhatsApp
 * await sdk.auth.sendOtp({ phone: "+212600000000", channel: "whatsapp" });
 * const { tokens, user } = await sdk.auth.verifyOtp({ phone: "+212600000000", code: "123456" });
 * sdk.setTokens(tokens);
 *
 * // Publier une annonce
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
   * Tokens initiaux (ex: restaurés depuis localStorage).
   * Peuvent être mis à jour ultérieurement avec setTokens().
   */
  initialTokens?: AuthTokens;
  /**
   * Appelé quand les tokens sont refreshés automatiquement.
   * Utilisez-le pour persister les nouveaux tokens.
   */
  onTokensRefreshed?: (tokens: AuthTokens) => void;
  /**
   * Appelé quand le refresh échoue (session expirée).
   * Utilisez-le pour rediriger vers la page de connexion.
   */
  onAuthFailure?: () => void;
}

/**
 * Point d'entrée principal du SDK Mezian.
 * Instanciez une seule fois dans votre application.
 */
export class MezianClient {
  /** Accès aux opérations d'authentification */
  readonly auth: AuthResource;
  /** Accès aux annonces */
  readonly ads: AdsResource;
  /** Accès aux catégories */
  readonly categories: CategoriesResource;
  /** Accès aux médias (upload images, YouTube) */
  readonly media: MediaResource;
  /** Accès aux boutiques pro */
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
   * Met à jour les tokens (après login ou refresh).
   * Stockez les tokens dans localStorage/AsyncStorage et restaurez-les
   * via initialTokens à la prochaine initialisation.
   */
  setTokens(tokens: AuthTokens): void {
    this.tokens = tokens;
  }

  /** Supprime les tokens (déconnexion côté client). */
  clearTokens(): void {
    this.tokens = null;
  }

  /** Retourne true si l'utilisateur a un token d'accès. */
  get isAuthenticated(): boolean {
    return this.tokens !== null;
  }
}
