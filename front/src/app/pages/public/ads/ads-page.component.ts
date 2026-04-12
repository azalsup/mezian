import { Component, inject, signal, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { LangService } from '../../../core/services/lang.service';
import { AdsApi } from '../../../sdk';
import { SiteFooterComponent } from '../../../shared/site-footer/site-footer.component';
import { AdsSearchBarComponent } from '../../../shared/ads-search-bar/ads-search-bar.component';
import { AdsFiltersComponent } from '../../../shared/ads-filters/ads-filters.component';
import { AdCardComponent } from '../../../shared/ad-card/ad-card.component';
import type { Ad, AdsQuery } from '../../../sdk';

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
export class AdsPageComponent implements OnInit, OnDestroy {
  readonly lang  = inject(LangService);
  private readonly route   = inject(ActivatedRoute);
  private readonly router  = inject(Router);
  private readonly adsApi  = inject(AdsApi);

  readonly currentQuery      = signal<AdsQuery>({});
  readonly mobileFiltersOpen = signal(false);
  readonly sort              = signal<'newest' | 'price_asc' | 'price_desc'>('newest');
  readonly ads               = signal<Ad[]>([]);
  readonly total             = signal(0);
  readonly loading           = signal(false);
  readonly error             = signal(false);

  private sub?: Subscription;

  ngOnInit(): void {
    this.sub = this.route.queryParams.subscribe(p => {
      const q: AdsQuery = {
        q:        p['q']        || undefined,
        cat:      p['cat']      || undefined,
        sub:      p['sub']      || undefined,
        city:     p['city']     || undefined,
        minPrice: p['minPrice'] ? Number(p['minPrice']) : undefined,
        maxPrice: p['maxPrice'] ? Number(p['maxPrice']) : undefined,
      };
      this.currentQuery.set(q);
      this.fetchAds(q);
    });
  }

  ngOnDestroy(): void {
    this.sub?.unsubscribe();
  }

  private fetchAds(q: AdsQuery): void {
    this.loading.set(true);
    this.error.set(false);
    this.adsApi.getAds({ ...q, sort: this.sort() }).subscribe({
      next: res => {
        this.ads.set(res.data);
        this.total.set(res.total);
        this.loading.set(false);
      },
      error: () => {
        this.error.set(true);
        this.loading.set(false);
      },
    });
  }

  onSortChange(sort: 'newest' | 'price_asc' | 'price_desc'): void {
    this.sort.set(sort);
    this.fetchAds({ ...this.currentQuery(), sort });
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
