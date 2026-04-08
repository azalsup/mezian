import { Injectable, inject, signal, computed, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { CategoriesApi } from '../../sdk';
import type { Category } from '../../sdk';
import { LangService, Lang } from './lang.service';
import staticData from '../../../../public/categories.json';

/** Flat lookup: slug → Category (root or child) */
export type CategoryMap = Map<string, Category>;

@Injectable({ providedIn: 'root' })
export class CategoriesService {
  private readonly api        = inject(CategoriesApi);
  private readonly lang       = inject(LangService);
  private readonly platformId = inject(PLATFORM_ID);

  /** Root categories — pre-populated from categories.json, updated from API when available */
  readonly categories = signal<Category[]>(staticData.categories.map(staticToCategory));
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
    // Defaults already in signal — attempt API refresh only in browser
    if (isPlatformBrowser(this.platformId)) {
      this.loadFromApi();
    } else {
      // SSG: defaults are sufficient, mark loading done
      this.loading.set(false);
    }
  }

  private loadFromApi(): void {
    this.api.getAll().subscribe({
      next: cats => {
        this.categories.set(cats);
        this.loading.set(false);
      },
      error: () => {
        // API unavailable — categories.json defaults remain
        this.loading.set(false);
      },
    });
  }
}

// ── Static JSON → Category conversion ────────────────────────────────────────

type StaticEntry = (typeof staticData.categories)[number];
type StaticSub   = NonNullable<StaticEntry['subcategories']>[number];

function staticToCategory(e: StaticEntry, idx: number): Category {
  return {
    id:         idx + 1,
    slug:       e.code,
    name_fr:    e.name_fr,
    name_ar:    e.name_ar,
    name_en:    e.name_en,
    icon:       e.icon,
    sort_order: e.sort_order,
    is_active:  true,
    children:   (e.subcategories ?? []).map((s: StaticSub, i: number) => ({
      id:         (idx + 1) * 100 + i + 1,
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
