import { Component, inject, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { LangService } from '../../../core/services/lang.service';
import { CategoriesBarComponent } from '../../../shared/categories-bar/categories-bar.component';
import { SiteFooterComponent } from '../../../shared/site-footer/site-footer.component';

interface AdTile {
  id: number;
  title: string;
  price: string;
  location: string;
  badge?: string;
  img: string;
}

interface AdSection {
  icon: string;
  titleKey: string;
  slug: string;
  ads: AdTile[];
}

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, CategoriesBarComponent, SiteFooterComponent],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  readonly lang   = inject(LangService);
  private readonly router = inject(Router);

  goToAds(query: string): void {
    this.router.navigate(['/ads'], query.trim() ? { queryParams: { q: query.trim() } } : {});
  }

  // ── Why Daba features ──────────────────────────────────────────────────────
  readonly features = computed(() => [
    { icon: 'photo_camera', title: this.lang.t('feat1Title'), desc: this.lang.t('feat1Desc') },
    { icon: 'sms',          title: this.lang.t('feat2Title'), desc: this.lang.t('feat2Desc') },
    { icon: 'storefront',   title: this.lang.t('feat3Title'), desc: this.lang.t('feat3Desc') },
    { icon: 'tune',         title: this.lang.t('feat4Title'), desc: this.lang.t('feat4Desc') },
  ]);

  // ── [TODO] Ad tile sections (placeholder — replaced by API data later) ───────────
  readonly adSections: AdSection[] = [
    {
      icon: 'local_fire_department',
      titleKey: 'sectionTrending',
      slug: 'trending',
      ads: [
        { id: 1, title: 'BMW Série 3 2021',       price: '280 000 MAD',    location: 'Casablanca', badge: 'TOP',     img: '' },
        { id: 2, title: 'Appartement 3 pièces',    price: '1 500 MAD/m²',  location: 'Rabat',      badge: 'NOUVEAU', img: '' },
        { id: 3, title: 'iPhone 15 Pro Max',       price: '12 500 MAD',    location: 'Marrakech',  badge: '',        img: '' },
        { id: 4, title: 'Ingénieur Full-Stack CDI',price: 'CDI',           location: 'Casablanca', badge: 'URGENT',  img: '' },
      ],
    },
    {
      icon: 'directions_car',
      titleKey: 'sectionCars',
      slug: 'automobiles',
      ads: [
        { id: 5, title: 'Dacia Duster 2022',     price: '195 000 MAD', location: 'Agadir',  badge: '',        img: '' },
        { id: 6, title: 'Renault Clio 5',        price: '135 000 MAD', location: 'Tanger',  badge: 'TOP',     img: '' },
        { id: 7, title: 'Toyota Hilux 4×4',      price: '380 000 MAD', location: 'Fès',     badge: '',        img: '' },
        { id: 8, title: 'Volkswagen Polo 2020',  price: '145 000 MAD', location: 'Rabat',   badge: 'NOUVEAU', img: '' },
      ],
    },
    {
      icon: 'home',
      titleKey: 'sectionRealEstate',
      slug: 'immobilier',
      ads: [
        { id: 9,  title: 'Villa piscine Souissi',    price: '4 500 000 MAD',  location: 'Rabat',      badge: 'PREMIUM', img: '' },
        { id: 10, title: 'Studio meublé Gauthier',   price: '6 500 MAD/mois', location: 'Casablanca', badge: '',        img: '' },
        { id: 11, title: 'Terrain 500 m² Bouskoura', price: '1 200 000 MAD',  location: 'Casablanca', badge: '',        img: '' },
        { id: 12, title: 'Appart Guéliz 2 chambres', price: '980 000 MAD',    location: 'Marrakech',  badge: 'NOUVEAU', img: '' },
      ],
    },
  ];
}
