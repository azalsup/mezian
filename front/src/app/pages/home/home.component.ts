import { Component, inject, computed, signal, HostListener } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';

interface NavCategory {
  slug: string;
  icon: string;
  labelKey: string;
  featured?: boolean;
  subcategories: string[];
}

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
  imports: [CommonModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  readonly lang = inject(LangService);

  readonly cities = ['Casablanca', 'Rabat', 'Marrakech', 'Fès', 'Tanger', 'Agadir', 'Meknès', 'Oujda'];

  openCat = signal<string | null>(null);

  toggleCat(slug: string) {
    this.openCat.set(this.openCat() === slug ? null : slug);
  }

  @HostListener('document:click')
  closeCat() {
    this.openCat.set(null);
  }

  readonly navCategories: NavCategory[] = [
    {
      slug: 'immobilier',
      icon: 'home',
      labelKey: 'catImmobilier',
      featured: true,
      subcategories: ['Appartements', 'Maisons & Villas', 'Terrains', 'Bureaux & Commerces', 'Colocations'],
    },
    {
      slug: 'automobile',
      icon: 'directions_car',
      labelKey: 'catAutomobile',
      featured: true,
      subcategories: ['Voitures', 'Motos & Scooters', 'Camions & Utilitaires', 'Pièces détachées', 'Location'],
    },
    {
      slug: 'emploi',
      icon: 'work',
      labelKey: 'catEmploi',
      featured: true,
      subcategories: ['Offres d\'emploi', 'Demandes d\'emploi', 'Freelance & Mission', 'Stages & Alternance'],
    },
    {
      slug: 'electronique',
      icon: 'smartphone',
      labelKey: 'catElectronique',
      subcategories: ['Téléphones', 'Ordinateurs', 'TV & Audio', 'Jeux vidéo', 'Électroménager'],
    },
    {
      slug: 'maison',
      icon: 'chair',
      labelKey: 'catMaison',
      subcategories: ['Meubles', 'Décoration', 'Jardin', 'Bricolage', 'Électroménager'],
    },
    {
      slug: 'mode',
      icon: 'checkroom',
      labelKey: 'catMode',
      subcategories: ['Vêtements femme', 'Vêtements homme', 'Chaussures', 'Accessoires', 'Montres & Bijoux'],
    },
    {
      slug: 'loisirs',
      icon: 'sports_soccer',
      labelKey: 'catLoisirs',
      subcategories: ['Sport & Loisirs', 'Livres & Musique', 'Enfants & Bébés', 'Collection', 'Animaux'],
    },
    {
      slug: 'services',
      icon: 'build',
      labelKey: 'catServices',
      subcategories: ['Artisans & Travaux', 'Cours & Formations', 'Événements', 'Santé & Beauté', 'Transport'],
    },
  ];

  readonly featuredCategories = this.navCategories.filter(c => c.featured);

  // Placeholder tiles — replaced by real API data later
  readonly adSections: AdSection[] = [
    {
      icon: 'local_fire_department',
      titleKey: 'sectionTrending',
      slug: 'trending',
      ads: [
        { id: 1, title: 'BMW Série 3 2021',      price: '280 000 MAD', location: 'Casablanca', badge: 'TOP',      img: '' },
        { id: 2, title: 'Appartement 3 pièces',   price: '1 500 MAD/m',  location: 'Rabat',       badge: 'NOUVEAU', img: '' },
        { id: 3, title: 'iPhone 15 Pro Max',      price: '12 500 MAD',  location: 'Marrakech',   badge: '',        img: '' },
        { id: 4, title: 'Ingénieur Full-Stack',   price: 'CDI',         location: 'Casablanca',  badge: 'URGENT',  img: '' },
      ],
    },
    {
      icon: 'directions_car',
      titleKey: 'sectionCars',
      slug: 'automobile',
      ads: [
        { id: 5, title: 'Dacia Duster 2022',      price: '195 000 MAD', location: 'Agadir',      badge: '',        img: '' },
        { id: 6, title: 'Renault Clio 5',         price: '135 000 MAD', location: 'Tanger',      badge: 'TOP',     img: '' },
        { id: 7, title: 'Toyota Hilux 4×4',       price: '380 000 MAD', location: 'Fès',         badge: '',        img: '' },
        { id: 8, title: 'Volkswagen Polo 2020',   price: '145 000 MAD', location: 'Rabat',       badge: 'NOUVEAU', img: '' },
      ],
    },
    {
      icon: 'home',
      titleKey: 'sectionRealEstate',
      slug: 'immobilier',
      ads: [
        { id: 9,  title: 'Villa piscine Souissi',   price: '4 500 000 MAD', location: 'Rabat',      badge: 'PREMIUM', img: '' },
        { id: 10, title: 'Studio meublé Gauthier',  price: '6 500 MAD/mois', location: 'Casablanca', badge: '',        img: '' },
        { id: 11, title: 'Terrain 500 m² Bouskoura',price: '1 200 000 MAD', location: 'Casablanca', badge: '',        img: '' },
        { id: 12, title: 'Appartement Guéliz 2 ch', price: '980 000 MAD',   location: 'Marrakech',  badge: 'NOUVEAU', img: '' },
      ],
    },
  ];

  readonly features = computed(() => [
    { icon: 'photo_camera', title: this.lang.t('feat1Title'), desc: this.lang.t('feat1Desc') },
    { icon: 'sms',          title: this.lang.t('feat2Title'), desc: this.lang.t('feat2Desc') },
    { icon: 'storefront',   title: this.lang.t('feat3Title'), desc: this.lang.t('feat3Desc') },
    { icon: 'tune',         title: this.lang.t('feat4Title'), desc: this.lang.t('feat4Desc') },
  ]);
}
