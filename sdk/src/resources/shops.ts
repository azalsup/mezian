/**
 * Ressource Shops — gestion des boutiques pro.
 * Read access is public. Create/update require a token.
 * Each user can only have one shop.
 */

import type { HttpClient } from "../utils/http.js";
import type {
  Ad,
  CreateShopRequest,
  Paginated,
  Shop,
  UpdateShopRequest,
} from "../types/index.js";

export class ShopsResource {
  constructor(private http: HttpClient) {}

  /**
   * Returns a public shop profile.
   *
   * @param slug - Shop slug, e.g. "auto-occasion-casa"
   */
  get(slug: string): Promise<Shop> {
    return this.http.get(`/shops/${slug}`);
  }

  /**
   * Returns a shop's active ads.
   */
  getAds(slug: string, page = 1, limit = 20): Promise<Paginated<Ad>> {
    return this.http.get(`/shops/${slug}/ads`, { page, limit });
  }

  /**
   * Creates the authenticated user's shop.
   * Each user can create only one shop.
   */
  create(req: CreateShopRequest): Promise<Shop> {
    return this.http.post("/shops", req);
  }

  /**
   * Updates the authenticated user's shop.
   */
  update(slug: string, req: UpdateShopRequest): Promise<Shop> {
    return this.http.put(`/shops/${slug}`, req);
  }

  /**
   * Returns the authenticated user's shop.
   */
  getMyShop(): Promise<Shop> {
    return this.http.get("/users/me/shop");
  }
}
