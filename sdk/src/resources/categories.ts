/**
 * Categories resource — read-only, public.
 * Returns the full tree with subcategories and attribute definitions.
 */

import type { HttpClient } from "../utils/http.js";
import type { Category } from "../types/index.js";

export class CategoriesResource {
  constructor(private http: HttpClient) {}

  /**
   * Returns all root categories with their children and attributes.
   *
   * @example
   * const cats = await sdk.categories.list();
   * // cats[0] = { slug: "immobilier", children: [...], ... }
   */
  list(): Promise<Category[]> {
    return this.http.get("/categories");
  }

  /**
   * Returns a specific category with its attributes.
   *
   * @param slug - ex: "voitures", "appartements-vente"
   */
  get(slug: string): Promise<Category> {
    return this.http.get(`/categories/${slug}`);
  }
}
