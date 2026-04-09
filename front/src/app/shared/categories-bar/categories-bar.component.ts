import { Component, inject, computed, signal, HostListener, ViewChild, ElementRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { LangService } from '../../core/services/lang.service';
import { CategoriesService } from '../../core/services/categories.service';
import type { Category } from '../../sdk';

@Component({
  selector: 'app-categories-bar',
  imports: [CommonModule, RouterModule],
  templateUrl: './categories-bar.component.html',
  styleUrl: './categories-bar.component.scss',
})
export class CategoriesBarComponent {
  readonly lang       = inject(LangService);
  readonly catService = inject(CategoriesService);

  @ViewChild('barEl') barEl!: ElementRef<HTMLDivElement>;

  openCat       = signal<string | null>(null);
  expandedInAll = signal<string | null>(null);
  dropdownLeft  = signal(24);

  readonly categories = this.catService.categories;

  /** First 6 categories shown in the nav bar */
  readonly navCats = computed(() => this.categories().slice(0, 6));

  /** Data for the open per-category dropdown (null when 'all' or closed) */
  readonly openCatData = computed(() => {
    const slug = this.openCat();
    if (!slug || slug === 'all') return null;
    return this.categories().find(c => c.slug === slug) ?? null;
  });

  label(cat: Category): string {
    return this.catService.nameOf(cat);
  }

  toggleCat(slug: string, event: MouseEvent): void {
    event.stopPropagation();
    const next = this.openCat() === slug ? null : slug;
    this.openCat.set(next);
    if (next !== 'all') this.expandedInAll.set(null);
    if (next) {
      const wrapper = (event.currentTarget as HTMLElement).parentElement!;
      const barRect = this.barEl.nativeElement.getBoundingClientRect();
      this.dropdownLeft.set(wrapper.getBoundingClientRect().left - barRect.left);
    }
  }

  toggleAllSub(slug: string, event: MouseEvent): void {
    event.stopPropagation();
    this.expandedInAll.set(this.expandedInAll() === slug ? null : slug);
  }

  @HostListener('document:click')
  closeDropdowns(): void {
    this.openCat.set(null);
    this.expandedInAll.set(null);
  }

  catBtnCls(cat: Category): Record<string, boolean> {
    const isOpen     = this.openCat() === cat.slug;
    const isFeatured = cat.sort_order <= 3;
    return {
      'text-white/75':   !isFeatured && !isOpen,
      'text-[#6ee7a8]':  isFeatured && !isOpen,
      'font-semibold':   isFeatured,
      'text-white':      isOpen,
      'bg-white/[.06]':  isOpen,
    };
  }

  /** Responsive visibility class for each nav category by index:
   *  0-1 → always visible, 2 → lg+, 3-5 → xl+ */
  navCatCls(i: number): Record<string, boolean> {
    return {
      'hidden':   i >= 2,
      'lg:block': i === 2,
      'xl:block': i >= 3,
    };
  }
}
