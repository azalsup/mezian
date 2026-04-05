/**
 * Ressource Boutiques — gestion des boutiques pro.
 * La lecture est publique. Création/modification requièrent un token.
 * Un utilisateur ne peut avoir qu'une seule boutique.
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
   * Retourne une boutique publique avec son profil.
   *
   * @param slug - Slug de la boutique, ex: "auto-occasion-casa"
   */
  get(slug: string): Promise<Shop> {
    return this.http.get(`/shops/${slug}`);
  }

  /**
   * Retourne les annonces actives d'une boutique.
   */
  getAds(slug: string, page = 1, limit = 20): Promise<Paginated<Ad>> {
    return this.http.get(`/shops/${slug}/ads`, { page, limit });
  }

  /**
   * Crée la boutique de l'utilisateur connecté.
   * Chaque utilisateur ne peut créer qu'une seule boutique.
   */
  create(req: CreateShopRequest): Promise<Shop> {
    return this.http.post("/shops", req);
  }

  /**
   * Met à jour la boutique de l'utilisateur connecté.
   */
  update(slug: string, req: UpdateShopRequest): Promise<Shop> {
    return this.http.put(`/shops/${slug}`, req);
  }

  /**
   * Retourne la boutique de l'utilisateur connecté.
   */
  getMyShop(): Promise<Shop> {
    return this.http.get("/users/me/shop");
  }
}
