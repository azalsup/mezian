import { Component, Output, EventEmitter, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';
import { CategoriesService } from '../../core/services/categories.service';
import type { Category } from '../../sdk';

export interface CategorySelection {
  cat: Category | null;
  sub: Category | null;
}

@Component({
  selector: 'app-category-select',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './category-select.component.html',
})
export class CategorySelectComponent {
  readonly lang       = inject(LangService);
  readonly catService = inject(CategoriesService);

  @Output() readonly selectionChange = new EventEmitter<CategorySelection>();

  readonly selectedCat = signal<Category | null>(null);
  readonly selectedSub = signal<Category | null>(null);

  readonly subcategories = computed<Category[]>(() => {
    const cat = this.selectedCat();
    return cat?.children ?? [];
  });

  setCat(catId: string): void {
    const cat = catId
      ? (this.catService.categories().find(c => c.id === Number(catId)) ?? null)
      : null;
    this.selectedCat.set(cat);
    this.selectedSub.set(null);
    this.selectionChange.emit({ cat, sub: null });
  }

  setSub(subId: string): void {
    const sub = subId
      ? (this.subcategories().find(s => s.id === Number(subId)) ?? null)
      : null;
    this.selectedSub.set(sub);
    this.selectionChange.emit({ cat: this.selectedCat(), sub });
  }

  catLabel(cat: Category): string {
    return this.catService.nameOf(cat);
  }
}
