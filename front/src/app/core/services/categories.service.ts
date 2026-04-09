import { Injectable, inject, signal, computed, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { CategoriesApi } from '../../sdk';
import type { Category } from '../../sdk';
import { LangService, Lang } from './lang.service';

/** Shape of /categories.json served from the public/ folder */
interface StaticSubEntry {
  code: string; icon: string; sort_order: number;
  name_fr: string; name_ar: string; name_en: string;
}
interface StaticCatEntry extends StaticSubEntry {
  featured?: boolean;
  subcategories?: StaticSubEntry[];
}

@Injectable({ providedIn: 'root' })
export class CategoriesService {
  private readonly api        = inject(CategoriesApi);
  private readonly http       = inject(HttpClient);
  private readonly lang       = inject(LangService);
  private readonly platformId = inject(PLATFORM_ID);

  /** Root categories — populated immediately from bundled JSON, then replaced by API data */
  readonly categories = signal<Category[]>([]);
  readonly loading    = signal(true);
  readonly error      = signal<string | null>(null);

  /** Flat slug → category map for O(1) lookups */
  readonly bySlug = computed<Map<string, Category>>(() => {
    const map = new Map<string, Category>();
    for (const cat of this.categories()) {
      map.set(cat.slug, cat);
      for (const child of cat.children ?? []) {
        map.set(child.slug, child);
      }
    }
    return map;
  });

  nameOf(cat: Category): string {
    return this.nameForLang(cat, this.lang.current());
  }

  nameForLang(cat: Category, lang: Lang): string {
    if (lang === 'ar') return cat.name_ar || cat.name_fr;
    if (lang === 'en') return cat.name_en || cat.name_fr;
    return cat.name_fr;
  }

  constructor() {
    if (isPlatformBrowser(this.platformId)) {
      // Load bundled JSON first for instant display, then override with live API data
      this.loadStaticFallback();
      this.loadFromApi();
    } else {
      this.loading.set(false);
    }
  }

  // ── API load ───────────────────────────────────────────────────────────────

  private loadFromApi(): void {
    this.api.getAll().subscribe({
      next: (res) => {
        // Backend wraps response in { data: [...] } via respondOK()
        const raw: unknown = res;
        let list: Category[];
        if (Array.isArray(raw)) {
          list = raw;
        } else if (raw && typeof raw === 'object' && 'data' in raw) {
          list = ((raw as { data: unknown }).data as Category[]) ?? [];
        } else {
          list = [];
        }
        // Only replace static data when the API actually returns categories
        if (list.length > 0) {
          this.categories.set(list.map(c => this.normalizeApiCategory(c)));
        }
        this.loading.set(false);
      },
      error: () => {
        // Backend unavailable — keep the static fallback already displayed
        this.loading.set(false);
      },
    });
  }

  /**
   * Normalize a raw category from the backend.
   * gorm.Model serializes its primary key as "ID" (uppercase), so we remap it to "id".
   */
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private normalizeApiCategory(raw: any): Category {
    return {
      ...raw,
      id: raw['id'] ?? raw['ID'] ?? 0,
      children: (raw['children'] ?? []).map((c: unknown) => this.normalizeApiCategory(c)),
    };
  }

  // ── Static fallback ────────────────────────────────────────────────────────

  /** Load the bundled /categories.json (from public/) for instant display */
  private loadStaticFallback(): void {
    this.http.get<{ categories: StaticCatEntry[] }>('/categories.json').subscribe({
      next: ({ categories }) => {
        // Only apply if the API hasn't already populated the signal
        if (this.categories().length === 0) {
          this.categories.set(this.transformStatic(categories));
        }
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      },
    });
  }

  /** Convert the bundled JSON shape (code/subcategories) to the Category interface */
  private transformStatic(entries: StaticCatEntry[]): Category[] {
    return entries.map((e, i) => ({
      id:         -(i + 1),          // negative so they never collide with DB ids
      slug:       e.code,
      name_fr:    e.name_fr,
      name_ar:    e.name_ar,
      name_en:    e.name_en,
      icon:       e.icon ?? '',
      sort_order: e.sort_order,
      is_active:  true,
      featured:   e.featured ?? false,
      children: (e.subcategories ?? []).map((s, j) => ({
        id:         -(i * 100 + j + 1),
        slug:       s.code,
        name_fr:    s.name_fr,
        name_ar:    s.name_ar,
        name_en:    s.name_en,
        icon:       s.icon ?? '',
        sort_order: s.sort_order,
        is_active:  true,
        featured:   false,
        children:   [],
      })),
    }));
  }
}
