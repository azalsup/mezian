import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiClient } from './api-client';
import { Category } from './types';

/**
 * Thin wrapper around all /categories/** endpoints.
 */
@Injectable({ providedIn: 'root' })
export class CategoriesApi {
  private readonly api = inject(ApiClient);

  /** Returns all root categories with their children and attribute definitions. */
  getAll(): Observable<Category[]> {
    return this.api.get('/categories');
  }

  /** Returns a single category (with children and attributes) by slug. */
  getBySlug(slug: string): Observable<Category> {
    return this.api.get(`/categories/${slug}`);
  }
}
