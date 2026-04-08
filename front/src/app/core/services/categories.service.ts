import { Injectable, inject, signal, computed, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { CategoriesApi } from '../../sdk';
import type { Category } from '../../sdk';
import { LangService, Lang } from './lang.service';

/** Flat lookup: slug → Category (root or child) */
export type CategoryMap = Map<string, Category>;

@Injectable({ providedIn: 'root' })
export class CategoriesService {
  private readonly api        = inject(CategoriesApi);
  private readonly lang       = inject(LangService);
  private readonly platformId = inject(PLATFORM_ID);

  /** Root categories with children, loaded from API */
  readonly categories = signal<Category[]>([]);
  readonly loading    = signal(true);
  readonly error      = signal<string | null>(null);

  /** Flat map slug → category for O(1) lookup */
  readonly bySlug = computed<CategoryMap>(() => {
    const map = new Map<string, Category>();
    for (const cat of this.categories()) {
      map.set(cat.slug, cat);
      for (const child of cat.children ?? []) {
        map.set(child.slug, child);
      }
    }
    return map;
  });

  /** Localised name for any category based on current language */
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
      this.loadFromApi();
    } else {
      // During SSG prerendering: load from static categories.json in public/
      this.loadFromStaticJson();
    }
  }

  private loadFromApi(): void {
    this.api.getAll().subscribe({
      next: cats => {
        this.categories.set(cats);
        this.loading.set(false);
      },
      error: () => {
        // API unavailable — fall back to static JSON
        this.loadFromStaticJson();
      },
    });
  }

  private loadFromStaticJson(): void {
    if (!isPlatformBrowser(this.platformId)) {
      this.loading.set(false);
      return;
    }
    fetch('/categories.json')
      .then(r => r.json())
      .then((data: { categories: StaticCatEntry[] }) => {
        this.categories.set(data.categories.map(staticToCategory));
        this.loading.set(false);
      })
      .catch(() => {
        this.loading.set(false);
        this.error.set('categories_unavailable');
      });
  }
}

// ── Static JSON → Category conversion ────────────────────────────────────────

interface StaticCatEntry {
  code: string;
  icon?: string;
  sort_order: number;
  name_fr: string;
  name_ar: string;
  name_en: string;
  featured?: boolean;
  subcategories?: StaticSubEntry[];
}

interface StaticSubEntry {
  code: string;
  icon?: string;
  sort_order: number;
  name_fr: string;
  name_ar: string;
  name_en: string;
}

function staticToCategory(e: StaticCatEntry, idx: number): Category {
  return {
    id:         idx + 1,
    slug:       e.code,
    name_fr:    e.name_fr,
    name_ar:    e.name_ar,
    name_en:    e.name_en,
    icon:       e.icon,
    sort_order: e.sort_order,
    is_active:  true,
    children:   (e.subcategories ?? []).map((s, i) => ({
      id:         (idx + 1) * 100 + i,
      slug:       s.code,
      name_fr:    s.name_fr,
      name_ar:    s.name_ar,
      name_en:    s.name_en,
      icon:       s.icon,
      sort_order: s.sort_order,
      is_active:  true,
      parent_id:  idx + 1,
    })),
  };
}
