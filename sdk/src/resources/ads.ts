/**
 * Ads resource — CRUD and search.
 *
 * Public read access: no token required.
 * Create/update: access token required.
 */

import type { HttpClient } from "../utils/http.js";
import type {
  Ad,
  AdFilters,
  CreateAdRequest,
  Paginated,
  UpdateAdRequest,
} from "../types/index.js";

export class AdsResource {
  constructor(private http: HttpClient) {}

  /**
   * Lists ads with filters and pagination.
   * Accessible without login.
   *
   * @example
   * const result = await sdk.ads.list({ category_id: 2, city: "Casablanca", limit: 20 });
   */
  list(filters?: AdFilters): Promise<Paginated<Ad>> {
    return this.http.get("/ads", filters as Record<string, unknown>);
  }

  /**
   * Returns a complete ad (with media and attributes).
   * Increments the view counter.
   *
   * @param slug - The URL-friendly slug of the ad
   */
  get(slug: string): Promise<Ad> {
    return this.http.get(`/ads/${slug}`);
  }

  /**
   * Creates a new ad.
   * Requires an access token.
   *
   * @example
   * const ad = await sdk.ads.create({
   *   category_id: 2,
   *   title: "Renault Clio 2019 – très bon état",
   *   body: "## Description\n\nVoiture en parfait état...",
   *   price: 85000,
   *   city: "Rabat",
   *   attributes: [
   *     { key: "brand", value: "Renault" },
   *     { key: "year", value: "2019" },
   *     { key: "mileage_km", value: "45000" },
   *   ]
   * });
   */
  create(req: CreateAdRequest): Promise<Ad> {
    return this.http.post("/ads", req);
  }

  /**
   * Updates an existing ad (owner only).
   */
  update(slug: string, req: UpdateAdRequest): Promise<Ad> {
    return this.http.put(`/ads/${slug}`, req);
  }

  /**
   * Deletes (archives) an ad.
   */
  delete(slug: string): Promise<{ message: string }> {
    return this.http.delete(`/ads/${slug}`);
  }

  /**
   * Returns the authenticated user's ads (all statuses).
   */
  myAds(page = 1, limit = 20): Promise<Paginated<Ad>> {
    return this.http.get("/users/me/ads", { page, limit });
  }
}
