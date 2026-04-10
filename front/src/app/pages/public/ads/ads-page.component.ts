import { Component, inject, signal, computed, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { LangService } from '../../../core/services/lang.service';
import { SiteFooterComponent } from '../../../shared/site-footer/site-footer.component';
import { AdsSearchBarComponent } from '../../../shared/ads-search-bar/ads-search-bar.component';
import { AdsFiltersComponent } from '../../../shared/ads-filters/ads-filters.component';
import { AdCardComponent } from '../../../shared/ad-card/ad-card.component';
import type { Ad, AdsQuery } from '../../../sdk';

// [TODO] remove mock
const MOCK_ADS: Ad[] = [
  { id: 1,  title: 'BMW Série 3 2020 — 125 000 km',     price: 185000,   city: 'Casablanca', category_slug: 'vehicles',    subcategory_slug: 'cars',       images: [], created_at: '2026-04-08', badge: 'top'     },
  { id: 2,  title: 'Appartement 3 pièces — Maarif',      price: 1450000,  city: 'Casablanca', category_slug: 'real_estate', subcategory_slug: 'sale',       images: [], created_at: '2026-04-07', badge: 'new'     },
  { id: 3,  title: 'iPhone 15 Pro Max 256 Go',            price: 12500,    city: 'Rabat',      category_slug: 'electronics', subcategory_slug: 'phones',     images: [], created_at: '2026-04-09'                   },
  { id: 4,  title: 'Développeur Full-Stack — CDI',        price: null,     city: 'Casablanca', category_slug: 'employment',  subcategory_slug: 'job_offers', images: [], created_at: '2026-04-08', badge: 'urgent'  },
  { id: 5,  title: 'Renault Clio 2019 Diesel',            price: 95000,    city: 'Marrakech',  category_slug: 'vehicles',    subcategory_slug: 'cars',       images: [], created_at: '2026-04-06'                   },
  { id: 6,  title: 'Studio meublé — Hay Riad',            price: 3200,     city: 'Rabat',      category_slug: 'real_estate', subcategory_slug: 'rent',       images: [], created_at: '2026-04-09', badge: 'new'     },
  { id: 7,  title: 'MacBook Pro M3 16" 512 Go',           price: 22000,    city: 'Fès',        category_slug: 'electronics', subcategory_slug: 'computers',  images: [], created_at: '2026-04-05'                   },
  { id: 8,  title: 'Canapé d\'angle en L — Bon état',    price: 3500,     city: 'Agadir',     category_slug: 'home',        subcategory_slug: 'furniture',  images: [], created_at: '2026-04-08'                   },
  { id: 9,  title: 'Dacia Duster 2022 — 45 000 km',      price: 175000,   city: 'Tanger',     category_slug: 'vehicles',    subcategory_slug: 'cars',       images: [], created_at: '2026-04-07', badge: 'top'     },
  { id: 10, title: 'Villa 4 chambres — Route de Fès',     price: 2800000,  city: 'Meknès',     category_slug: 'real_estate', subcategory_slug: 'sale',       images: [], created_at: '2026-04-04', badge: 'premium' },
  { id: 11, title: 'Samsung Galaxy S24 Ultra 512 Go',     price: 9800,     city: 'Oujda',      category_slug: 'electronics', subcategory_slug: 'phones',     images: [], created_at: '2026-04-09'                   },
  { id: 12, title: 'Terrain 500m² — Zone industrielle',   price: 850000,   city: 'Kénitra',    category_slug: 'real_estate', subcategory_slug: 'land',       images: [], created_at: '2026-04-06'                   },
];

@Component({
  selector: 'app-ads-page',
  standalone: true,
  imports: [
    CommonModule,
    SiteFooterComponent,
    AdsSearchBarComponent,
    AdsFiltersComponent,
    AdCardComponent,
  ],
  templateUrl: './ads-page.component.html',
})
export class AdsPageComponent implements OnInit {
  readonly lang  = inject(LangService);
  private readonly route  = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly currentQuery      = signal<AdsQuery>({});
  readonly mobileFiltersOpen = signal(false);
  readonly sort              = signal<'newest' | 'price_asc' | 'price_desc'>('newest');

  readonly ads = computed(() => {
    const q    = this.currentQuery();
    const sort = this.sort();

    let result = MOCK_ADS.filter(ad => {
      if (q.cat      && ad.category_slug    !== q.cat)  return false;
      if (q.sub      && ad.subcategory_slug !== q.sub)  return false;
      if (q.city     && ad.city             !== q.city) return false;
      if (q.minPrice != null && ad.price !== null && ad.price < q.minPrice) return false;
      if (q.maxPrice != null && ad.price !== null && ad.price > q.maxPrice) return false;
      if (q.q) {
        const term = q.q.toLowerCase();
        if (!ad.title.toLowerCase().includes(term)) return false;
      }
      return true;
    });

    if (sort === 'price_asc')  result = [...result].sort((a, b) => (a.price ?? 0) - (b.price ?? 0));
    if (sort === 'price_desc') result = [...result].sort((a, b) => (b.price ?? 0) - (a.price ?? 0));

    return result;
  });

  ngOnInit(): void {
    this.route.queryParams.subscribe(p => {
      this.currentQuery.set({
        q:        p['q']        || undefined,
        cat:      p['cat']      || undefined,
        sub:      p['sub']      || undefined,
        city:     p['city']     || undefined,
        minPrice: p['minPrice'] ? Number(p['minPrice']) : undefined,
        maxPrice: p['maxPrice'] ? Number(p['maxPrice']) : undefined,
      });
    });
  }

  onSearch(q: string): void {
    this.navigate({ ...this.currentQuery(), q: q || undefined });
  }

  onFiltersApply(filters: AdsQuery): void {
    this.navigate({ ...this.currentQuery(), ...filters });
  }

  onFiltersReset(): void {
    this.navigate({ q: this.currentQuery().q });
    this.mobileFiltersOpen.set(false);
  }

  private navigate(query: AdsQuery): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        q:        query.q        || null,
        cat:      query.cat      || null,
        sub:      query.sub      || null,
        city:     query.city     || null,
        minPrice: query.minPrice ?? null,
        maxPrice: query.maxPrice ?? null,
      },
    });
  }
}
