import { Component, inject, signal, computed, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { LangService } from '../../../core/services/lang.service';
import { CategoriesService } from '../../../core/services/categories.service';
import { AdsApi } from '../../../sdk';
import { CategoriesBarComponent } from '../../../shared/categories-bar/categories-bar.component';
import { SiteFooterComponent } from '../../../shared/site-footer/site-footer.component';
import type { Ad, Category } from '../../../sdk';

@Component({
  selector: 'app-ad-detail-page',
  standalone: true,
  imports: [CommonModule, RouterLink, CategoriesBarComponent, SiteFooterComponent],
  templateUrl: './ad-detail-page.component.html',
})
export class AdDetailPageComponent implements OnInit {
  readonly lang       = inject(LangService);
  readonly catService = inject(CategoriesService);
  private readonly route  = inject(ActivatedRoute);
  private readonly adsApi = inject(AdsApi);

  readonly ad          = signal<Ad | null>(null);
  readonly loading     = signal(true);
  readonly error       = signal(false);
  readonly activeImage = signal(0);

  // Resolve category from preloaded relation or from categories service
  readonly category = computed<Category | null>(() => {
    const ad = this.ad();
    if (!ad) return null;
    if (ad.category) return ad.category;
    const slug = ad.category_slug;
    return slug ? (this.catService.categories().find(c => c.slug === slug) ?? null) : null;
  });

  readonly subcategory = computed<Category | null>(() => {
    const ad  = this.ad();
    const cat = this.category();
    if (!ad || !cat) return null;
    // subcategory is the category itself if it has a parent_id
    if (cat.parent_id) return cat;
    const subSlug = ad.subcategory_slug;
    return subSlug ? (cat.children?.find(c => c.slug === subSlug) ?? null) : null;
  });

  // Unified image list: prefer media URLs, fall back to legacy images array
  readonly images = computed<string[]>(() => {
    const ad = this.ad();
    if (!ad) return [];
    if (ad.media?.length) return ad.media.map(m => m.url);
    return ad.images ?? [];
  });

  readonly body = computed(() => this.ad()?.body ?? this.ad()?.description ?? '');

  readonly sellerName = computed(() =>
    this.ad()?.user?.display_name ?? this.ad()?.seller_name ?? null
  );

  readonly viewCount = computed(() =>
    this.ad()?.view_count ?? this.ad()?.views ?? null
  );

  readonly categorySlug = computed(() =>
    this.ad()?.category?.slug ?? this.ad()?.category_slug ?? null
  );

  ngOnInit(): void {
    const slug = this.route.snapshot.queryParamMap.get('id') ?? '';
    this.adsApi.getBySlug(slug).subscribe({
      next:  ad => { this.ad.set(ad); this.loading.set(false); },
      error: ()  => { this.error.set(true); this.loading.set(false); },
    });
  }

  catLabel(cat: Category): string {
    return this.catService.nameOf(cat);
  }

  formatDate(iso: string): string {
    return new Date(iso).toLocaleDateString(
      this.lang.current() === 'ar' ? 'ar-MA' : 'fr-MA',
      { day: 'numeric', month: 'long', year: 'numeric' },
    );
  }

  badgeCls(badge: Ad['badge']): Record<string, boolean> {
    return {
      'bg-[#006233] text-white': badge === 'top',
      'bg-blue-600 text-white':  badge === 'new',
      'bg-red-600 text-white':   badge === 'urgent',
      'bg-amber-600 text-white': badge === 'premium',
    };
  }

  badgeLabel(badge: Ad['badge']): string {
    const map: Record<string, string> = { top: 'TOP', new: 'NOUVEAU', urgent: 'URGENT', premium: 'PREMIUM' };
    return badge ? (map[badge] ?? badge.toUpperCase()) : '';
  }
}
