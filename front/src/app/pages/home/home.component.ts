import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent {
  categories = [
    { icon: 'home',           slug: 'immobilier',   labelFr: 'Real Estate',   labelAr: 'عقارات',      count: '12 400' },
    { icon: 'directions_car', slug: 'automobile',   labelFr: 'Automotive',    labelAr: 'سيارات',      count: '8 750'  },
    { icon: 'work',           slug: 'emploi',       labelFr: 'Jobs',          labelAr: 'توظيف',       count: '3 200'  },
    { icon: 'smartphone',     slug: 'electronique', labelFr: 'Electronics',   labelAr: 'إلكترونيات',  count: '5 600'  },
    { icon: 'chair',          slug: 'maison',       labelFr: 'Home & Living', labelAr: 'منزل',        count: '4 100'  },
    { icon: 'checkroom',      slug: 'mode',         labelFr: 'Fashion',       labelAr: 'موضة',        count: '2 900'  },
    { icon: 'sports_soccer',  slug: 'loisirs',      labelFr: 'Leisure',       labelAr: 'ترفيه',       count: '1 800'  },
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
