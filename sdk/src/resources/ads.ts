/**
 * Ressource Annonces — CRUD et recherche.
 *
 * Lecture publique : aucun token requis.
 * Création / modification : access token requis.
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
   * Liste les annonces avec filtres et pagination.
   * Accessible sans connexion.
   *
   * @example
   * const result = await sdk.ads.list({ category_id: 2, city: "Casablanca", limit: 20 });
   */
  list(filters?: AdFilters): Promise<Paginated<Ad>> {
    return this.http.get("/ads", filters as Record<string, unknown>);
  }

  /**
   * Retourne une annonce complète (avec médias et attributs).
   * Incrémente le compteur de vues.
   *
   * @param slug - Le slug URL-friendly de l'annonce
   */
  get(slug: string): Promise<Ad> {
    return this.http.get(`/ads/${slug}`);
  }

  /**
   * Crée une nouvelle annonce.
   * Requiert un access token.
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
   * Met à jour une annonce existante (propriétaire uniquement).
   */
  update(slug: string, req: UpdateAdRequest): Promise<Ad> {
    return this.http.put(`/ads/${slug}`, req);
  }

  /**
   * Supprime (archive) une annonce.
   */
  delete(slug: string): Promise<{ message: string }> {
    return this.http.delete(`/ads/${slug}`);
  }

  /**
   * Retourne les annonces de l'utilisateur connecté (tous statuts).
   */
  myAds(page = 1, limit = 20): Promise<Paginated<Ad>> {
    return this.http.get("/users/me/ads", { page, limit });
  }
}
