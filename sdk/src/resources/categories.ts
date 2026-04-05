/**
 * Ressource Catégories — lecture seule, publique.
 * Retourne l'arbre complet avec sous-catégories et définitions d'attributs.
 */

import type { HttpClient } from "../utils/http.js";
import type { Category } from "../types/index.js";

export class CategoriesResource {
  constructor(private http: HttpClient) {}

  /**
   * Retourne toutes les catégories racines avec leurs enfants et attributs.
   *
   * @example
   * const cats = await sdk.categories.list();
   * // cats[0] = { slug: "immobilier", children: [...], ... }
   */
  list(): Promise<Category[]> {
    return this.http.get("/categories");
  }

  /**
   * Retourne une catégorie spécifique avec ses attributs.
   *
   * @param slug - ex: "voitures", "appartements-vente"
   */
  get(slug: string): Promise<Category> {
    return this.http.get(`/categories/${slug}`);
  }
}
