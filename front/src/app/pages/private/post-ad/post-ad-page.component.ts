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

const MAX_PHOTOS = 5;

@Component({
  selector: 'app-post-ad-page',
  standalone: true,
  imports: [CommonModule, CategoriesBarComponent, SiteFooterComponent, CategorySelectComponent],
  templateUrl: './post-ad-page.component.html',
})
export class PostAdPageComponent {
  readonly lang           = inject(LangService);
  private readonly adsApi = inject(AdsApi);
  private readonly router = inject(Router);

  readonly cities   = MOROCCO_CITIES;
  readonly maxPhotos = MAX_PHOTOS;

  // Step state
  readonly step = signal<1 | 2 | 3>(1);

  // Form state
  readonly selectedCat = signal<Category | null>(null);
  readonly selectedSub = signal<Category | null>(null);
  readonly title       = signal('');
  readonly description = signal('');
  readonly price       = signal('');
  readonly city        = signal('');

  // Photo state
  readonly photos   = signal<File[]>([]);
  readonly previews = signal<string[]>([]);

  // UI state
  readonly submitting = signal(false);
  readonly error      = signal('');
  readonly success    = signal(false);

  onSelectionChange(sel: CategorySelection): void {
    this.selectedCat.set(sel.cat);
    this.selectedSub.set(sel.sub);
  }

  next(): void {
    this.error.set('');
    if (this.step() === 1) {
      if (!(this.selectedSub() ?? this.selectedCat())) {
        this.error.set(this.lang.t('errCatRequired'));
        return;
      }
      this.step.set(2);
    } else if (this.step() === 2) {
      if (!this.title().trim()) {
        this.error.set(this.lang.t('errTitleRequired'));
        return;
      }
      this.step.set(3);
    }
  }

  back(): void {
    this.error.set('');
    if (this.step() === 2) this.step.set(1);
    else if (this.step() === 3) this.step.set(2);
  }

  onFilesSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (!input.files) return;
    const current = this.photos();
    const remaining = MAX_PHOTOS - current.length;
    const newFiles = Array.from(input.files).slice(0, remaining);
    const newPreviews = newFiles.map(f => URL.createObjectURL(f));
    this.photos.set([...current, ...newFiles]);
    this.previews.set([...this.previews(), ...newPreviews]);
    input.value = '';
  }

  removePhoto(index: number): void {
    URL.revokeObjectURL(this.previews()[index]);
    this.photos.update(arr => arr.filter((_, i) => i !== index));
    this.previews.update(arr => arr.filter((_, i) => i !== index));
  }

  submit(): void {
    this.error.set('');
    if (!this.city()) { this.error.set(this.lang.t('errCityRequired')); return; }

    const cat      = this.selectedSub() ?? this.selectedCat()!;
    const priceVal = this.price().trim();

    this.submitting.set(true);
    this.adsApi.createAd({
      category_id: cat.ID,
      title:       this.title().trim(),
      body:        this.description().trim(),
      price:       priceVal ? Number(priceVal) : undefined,
      currency:    'MAD',
      city:        this.city(),
    }, this.photos()).subscribe({
      next: (ad: any) => {
        this.submitting.set(false);
        this.success.set(true);
        setTimeout(() => this.router.navigate(['/ad'], { queryParams: { id: ad.slug } }), 1500);
      },
      error: () => {
        this.submitting.set(false);
        this.error.set(this.lang.t('errNetwork'));
      },
    });
  }
}
