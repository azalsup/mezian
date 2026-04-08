import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';
import { SearchBarComponent } from '../../shared/search-bar/search-bar.component';
import { PromoBannerComponent } from '../../shared/promo-banner/promo-banner.component';
import { AdsListComponent } from '../../shared/ads-list/ads-list.component';
import { SiteFooterComponent } from '../../shared/site-footer/site-footer.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    CommonModule,
    SearchBarComponent,
    PromoBannerComponent,
    AdsListComponent,
    SiteFooterComponent,
  ],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  constructor(public lang: LangService) {}
  categories = [
    { icon: 'home',           slug: 'immobilier',   labelFr: 'Immobilier',   labelAr: 'عقارات',      count: '12 400' },
    { icon: 'directions_car', slug: 'automobile',   labelFr: 'Véhicules',    labelAr: 'سيارات',      count: '8 750'  },
    { icon: 'work',           slug: 'emploi',       labelFr: 'Emploi',       labelAr: 'توظيف',       count: '3 200'  },
    { icon: 'smartphone',     slug: 'electronique', labelFr: 'Électronique', labelAr: 'إلكترونيات',  count: '5 600'  },
    { icon: 'chair',          slug: 'maison',       labelFr: 'Maison',       labelAr: 'منزل',        count: '4 100'  },
    { icon: 'checkroom',      slug: 'mode',         labelFr: 'Mode',         labelAr: 'موضة',        count: '2 900'  },
    { icon: 'sports_soccer',  slug: 'loisirs',      labelFr: 'Loisirs',      labelAr: 'ترفيه',       count: '1 800'  },
    { icon: 'build',          slug: 'services',     labelFr: 'Services',      labelAr: 'خدمات',       count: '2 300'  },
  ];

  cities = ['Casablanca', 'Rabat', 'Marrakech', 'Fès', 'Tanger', 'Agadir', 'Meknès', 'Oujda'];

  features = [
    { icon: 'photo_camera', title: 'Photos & Videos',  desc: 'Up to 10 photos and YouTube videos for your real estate listings.' },
    { icon: 'sms',          title: 'Phone Login',       desc: 'No password needed — receive a code by SMS or WhatsApp.' },
    { icon: 'storefront',   title: 'Pro Shops',         desc: 'Create your professional shop and manage all your ads in one place.' },
    { icon: 'tune',         title: 'Advanced Search',   desc: 'Filter by city, price, category, and specific attributes.' },
  ];
}
