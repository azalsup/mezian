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
    { icon: '🏠', slug: 'immobilier',  labelFr: 'Immobilier',  labelAr: 'عقارات',       count: '12 400' },
    { icon: '🚗', slug: 'automobile',  labelFr: 'Automobile',   labelAr: 'سيارات',       count: '8 750'  },
    { icon: '💼', slug: 'emploi',      labelFr: 'Emploi',       labelAr: 'توظيف',        count: '3 200'  },
    { icon: '📱', slug: 'electronique',labelFr: 'Électronique', labelAr: 'إلكترونيات',   count: '5 600'  },
    { icon: '🛋️', slug: 'maison',      labelFr: 'Maison',       labelAr: 'منزل',         count: '4 100'  },
    { icon: '👗', slug: 'mode',        labelFr: 'Mode',         labelAr: 'موضة',         count: '2 900'  },
    { icon: '⚽', slug: 'loisirs',     labelFr: 'Loisirs',      labelAr: 'ترفيه',        count: '1 800'  },
    { icon: '🔧', slug: 'services',    labelFr: 'Services',     labelAr: 'خدمات',        count: '2 300'  },
  ];

  cities = ['Casablanca', 'Rabat', 'Marrakech', 'Fès', 'Tanger', 'Agadir', 'Meknès', 'Oujda'];

  features = [
    { icon: '📸', title: 'Photos & vidéos',   desc: 'Jusqu\'à 10 photos et des vidéos YouTube pour vos annonces immobilier.' },
    { icon: '📱', title: 'Connexion par téléphone', desc: 'Pas besoin de mot de passe — recevez un code par SMS ou WhatsApp.' },
    { icon: '🏪', title: 'Boutiques pro',     desc: 'Créez votre boutique professionnelle et gérez vos annonces en un seul endroit.' },
    { icon: '🔍', title: 'Recherche avancée', desc: 'Filtrez par ville, prix, catégorie et attributs spécifiques.' },
  ];
}
