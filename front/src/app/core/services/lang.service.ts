import { Injectable, signal, computed } from '@angular/core';

export type Lang = 'fr' | 'ar';

const TRANSLATIONS: Record<Lang, Record<string, string>> = {
  fr: {
    sell:           'Vendre',
    profile:        'Profil',
    register:       "S'inscrire",
    createProfile:  'Créer mon profil',
    alreadyAccount: "J'ai déjà un compte",
    myProfile:      'Mon profil',
    myShop:         'Ma boutique',
    logout:         'Déconnexion',
    location:       'Casablanca • Particulier',
  },
  ar: {
    sell:           'بيع',
    profile:        'حسابي',
    register:       'تسجيل',
    createProfile:  'إنشاء حسابي',
    alreadyAccount: 'لدي حساب',
    myProfile:      'ملفي الشخصي',
    myShop:         'متجري',
    logout:         'تسجيل الخروج',
    location:       'الدار البيضاء • فرد',
  },
};

@Injectable({ providedIn: 'root' })
export class LangService {
  readonly current = signal<Lang>('fr');
  readonly isRtl   = computed(() => this.current() === 'ar');

  t(key: string): string {
    return TRANSLATIONS[this.current()][key] ?? key;
  }

  toggle(): void {
    const next: Lang = this.current() === 'fr' ? 'ar' : 'fr';
    this.current.set(next);
    document.documentElement.lang = next;
    document.documentElement.dir  = next === 'ar' ? 'rtl' : 'ltr';
  }
}
