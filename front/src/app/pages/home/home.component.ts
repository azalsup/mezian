import { Component, inject, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  readonly lang = inject(LangService);

  readonly categories = computed(() => [
    { icon: 'home',           slug: 'immobilier',   label: this.lang.t('catImmobilier'),   count: '12 400' },
    { icon: 'directions_car', slug: 'automobile',   label: this.lang.t('catAutomobile'),   count: '8 750'  },
    { icon: 'work',           slug: 'emploi',       label: this.lang.t('catEmploi'),       count: '3 200'  },
    { icon: 'smartphone',     slug: 'electronique', label: this.lang.t('catElectronique'), count: '5 600'  },
    { icon: 'chair',          slug: 'maison',       label: this.lang.t('catMaison'),       count: '4 100'  },
    { icon: 'checkroom',      slug: 'mode',         label: this.lang.t('catMode'),         count: '2 900'  },
    { icon: 'sports_soccer',  slug: 'loisirs',      label: this.lang.t('catLoisirs'),      count: '1 800'  },
    { icon: 'build',          slug: 'services',     label: this.lang.t('catServices'),     count: '2 300'  },
  ]);

  readonly features = computed(() => [
    { icon: 'photo_camera', title: this.lang.t('feat1Title'), desc: this.lang.t('feat1Desc') },
    { icon: 'sms',          title: this.lang.t('feat2Title'), desc: this.lang.t('feat2Desc') },
    { icon: 'storefront',   title: this.lang.t('feat3Title'), desc: this.lang.t('feat3Desc') },
    { icon: 'tune',         title: this.lang.t('feat4Title'), desc: this.lang.t('feat4Desc') },
  ]);

  readonly cities = ['Casablanca', 'Rabat', 'Marrakech', 'Fès', 'Tanger', 'Agadir', 'Meknès', 'Oujda'];
}
