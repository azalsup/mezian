import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { LangService } from '../../../core/services/lang.service';
import { AdsApi } from '../../../sdk';
import { CategoriesBarComponent } from '../../../shared/categories-bar/categories-bar.component';
import { SiteFooterComponent } from '../../../shared/site-footer/site-footer.component';
import { CategorySelectComponent, type CategorySelection } from '../../../shared/category-select/category-select.component';
import { MOROCCO_CITIES } from '../../../shared/ads-filters/ads-filters.component';
import type { Category } from '../../../sdk';

@Component({
  selector: 'app-post-ad-page',
  standalone: true,
  imports: [CommonModule, CategoriesBarComponent, SiteFooterComponent, CategorySelectComponent],
  templateUrl: './post-ad-page.component.html',
})
export class PostAdPageComponent {
  readonly lang    = inject(LangService);
  private readonly adsApi = inject(AdsApi);
  private readonly router = inject(Router);

  readonly cities = MOROCCO_CITIES;

  // Form state
  readonly selectedCat = signal<Category | null>(null);
  readonly selectedSub = signal<Category | null>(null);
  readonly title       = signal('');
  readonly description = signal('');
  readonly price       = signal('');
  readonly city        = signal('');

  // UI state
  readonly submitting = signal(false);
  readonly error      = signal('');
  readonly success    = signal(false);

  onSelectionChange(sel: CategorySelection): void {
    this.selectedCat.set(sel.cat);
    this.selectedSub.set(sel.sub);
  }

  submit(): void {
    this.error.set('');

    const cat = this.selectedSub() ?? this.selectedCat();
    if (!cat) { this.error.set(this.lang.t('errCatRequired')); return; }
    if (!this.title().trim()) { this.error.set(this.lang.t('errTitleRequired')); return; }
    if (!this.city()) { this.error.set(this.lang.t('errCityRequired')); return; }

    const priceVal = this.price().trim();

    this.submitting.set(true);
    this.adsApi.createAd({
      category_id: cat.id,
      title:       this.title().trim(),
      body:        this.description().trim(),
      price:       priceVal ? Number(priceVal) : undefined,
      currency:    'MAD',
      city:        this.city(),
    }).subscribe({
      next: (ad) => {
        this.submitting.set(false);
        this.success.set(true);
        setTimeout(() => this.router.navigate(['/ads', ad.id]), 1500);
      },
      error: () => {
        this.submitting.set(false);
        this.error.set(this.lang.t('errNetwork'));
      },
    });
  }
}
