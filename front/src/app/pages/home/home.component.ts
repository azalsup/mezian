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
    { icon: '🏠', slug: 'immobilier',  labelFr: 'Real estate',  labelAr: 'عقارات',       count: '12 400' },
    { icon: '🚗', slug: 'automobile',  labelFr: 'Automotive',   labelAr: 'سيارات',       count: '8 750'  },
    { icon: '💼', slug: 'emploi',      labelFr: 'Jobs',       labelAr: 'توظيف',        count: '3 200'  },
    { icon: '📱', slug: 'electronique',labelFr: 'Electronics', labelAr: 'إلكترونيات',   count: '5 600'  },
    { icon: '🛋️', slug: 'maison',      labelFr: 'Home',       labelAr: 'منزل',         count: '4 100'  },
    { icon: '👗', slug: 'mode',        labelFr: 'Fashion',         labelAr: 'موضة',         count: '2 900'  },
    { icon: '⚽', slug: 'loisirs',     labelFr: 'Leisure',      labelAr: 'ترفيه',        count: '1 800'  },
    { icon: '🔧', slug: 'services',    labelFr: 'Services',     labelAr: 'خدمات',        count: '2 300'  },
  ];

  cities = ['Casablanca', 'Rabat', 'Marrakech', 'Fès', 'Tanger', 'Agadir', 'Meknès', 'Oujda'];

  features = [
    { icon: '📸', title: 'Photos & videos',   desc: 'Jusqu\'à 10 photos et des vidéos YouTube pour vos annonces immobilier.' },
    { icon: '📱', title: 'Phone login', desc: 'No password needed — receive a code by SMS or WhatsApp.' },
    { icon: '🏪', title: 'Pro shops',     desc: 'Create your professional shop and manage your ads in one place.' },
    { icon: '🔍', title: 'Advanced search', desc: 'Filter by city, price, category, and specific attributes.' },
  ];
}
