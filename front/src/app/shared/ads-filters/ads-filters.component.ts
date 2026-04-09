import { Component, Input, Output, EventEmitter, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';
import { CategoriesService } from '../../core/services/categories.service';
import type { AdsQuery, Category } from '../../sdk';

export const MOROCCO_CITIES = [
  'Casablanca', 'Rabat', 'Marrakech', 'Fès', 'Tanger', 'Agadir',
  'Meknès', 'Oujda', 'Kénitra', 'Tétouan', 'Safi', 'El Jadida',
  'Béni Mellal', 'Nador', 'Mohammedia', 'Laâyoune',
];

@Component({
  selector: 'app-ads-filters',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './ads-filters.component.html',
})
export class AdsFiltersComponent implements OnInit {
  readonly lang       = inject(LangService);
  readonly catService = inject(CategoriesService);

  @Input() initial: AdsQuery = {};
  @Output() readonly apply = new EventEmitter<AdsQuery>();
  @Output() readonly reset = new EventEmitter<void>();

  readonly CITIES = MOROCCO_CITIES;

  cat      = signal('');
  sub      = signal('');
  city     = signal('');
  minPrice = signal('');
  maxPrice = signal('');

  readonly categories   = this.catService.categories;
  readonly subcategories = computed(() =>
    this.categories().find(c => c.slug === this.cat())?.children ?? []
  );

  ngOnInit(): void {
    this.cat.set(this.initial.cat      ?? '');
    this.sub.set(this.initial.sub      ?? '');
    this.city.set(this.initial.city    ?? '');
    this.minPrice.set(this.initial.minPrice != null ? String(this.initial.minPrice) : '');
    this.maxPrice.set(this.initial.maxPrice != null ? String(this.initial.maxPrice) : '');
  }

  setCat(value: string): void {
    this.cat.set(value);
    this.sub.set('');
  }

  doApply(): void {
    this.apply.emit({
      cat:      this.cat()      || undefined,
      sub:      this.sub()      || undefined,
      city:     this.city()     || undefined,
      minPrice: this.minPrice() ? Number(this.minPrice()) : undefined,
      maxPrice: this.maxPrice() ? Number(this.maxPrice()) : undefined,
    });
  }

  doReset(): void {
    this.cat.set('');
    this.sub.set('');
    this.city.set('');
    this.minPrice.set('');
    this.maxPrice.set('');
    this.reset.emit();
  }

  label(cat: Category): string {
    return this.catService.nameOf(cat);
  }
}
