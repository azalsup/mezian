import { Component, Input, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { LangService } from '../../core/services/lang.service';
import type { Ad } from '../../sdk';

@Component({
  selector: 'app-ad-card',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './ad-card.component.html',
})
export class AdCardComponent {
  readonly lang = inject(LangService);

  @Input({ required: true }) ad!: Ad;

  coverUrl(): string | null {
    return this.ad.media?.[0]?.url ?? this.ad.images?.[0] ?? null;
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

  formatDate(iso: string): string {
    return new Date(iso).toLocaleDateString(this.lang.current() === 'ar' ? 'ar-MA' : 'fr-MA', {
      day: '2-digit', month: 'short',
    });
  }
}
