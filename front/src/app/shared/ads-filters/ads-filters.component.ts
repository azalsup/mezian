import {
  Component, Input, Output, EventEmitter,
  OnInit, OnChanges, SimpleChanges,
  inject, signal, computed,
} from '@angular/core';
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
export class AdsFiltersComponent implements OnInit, OnChanges {
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

  readonly categories    = this.catService.categories;
  readonly subcategories = computed(() =>
    this.categories().find(c => c.slug === this.cat())?.children ?? []
  );

  readonly activeCount = computed(() =>
    [this.cat(), this.sub(), this.city(), this.minPrice(), this.maxPrice()]
      .filter(Boolean).length
  );

  ngOnInit(): void {
    this.syncFromInput(this.initial);
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['initial'] && !changes['initial'].firstChange) {
      this.syncFromInput(changes['initial'].currentValue as AdsQuery);
    }
  }

  private syncFromInput(q: AdsQuery): void {
    this.cat.set(q.cat      ?? '');
    this.sub.set(q.sub      ?? '');
    this.city.set(q.city    ?? '');
    this.minPrice.set(q.minPrice != null ? String(q.minPrice) : '');
    this.maxPrice.set(q.maxPrice != null ? String(q.maxPrice) : '');
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
