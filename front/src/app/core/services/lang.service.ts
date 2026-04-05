import { Injectable, signal, computed } from '@angular/core';

export type Lang = 'fr' | 'ar';

const TRANSLATIONS: Record<Lang, Record<string, string>> = {
  fr: {
    // Navbar
    sell:           'Vendre',
    profile:        'Profil',
    register:       "S'inscrire",
    createProfile:  'Créer mon profil',
    alreadyAccount: "J'ai déjà un compte",
    myProfile:      'Mon profil',
    myShop:         'Ma boutique',
    logout:         'Déconnexion',
    location:       'Casablanca • Particulier',

    // Auth modal — phone step
    loginTitle:       'Connexion',
    registerTitle:    'Créer un compte',
    phoneLabel:       'Numéro de téléphone',
    phonePlaceholder: '+212 6XX XXX XXX',
    sendCode:         'Recevoir le code',
    channelSms:       'SMS',
    channelWhatsapp:  'WhatsApp',

    // Auth modal — OTP step
    otpTitle:         'Vérification',
    otpSentTo:        'Code envoyé au',
    otpLabel:         'Code à 6 chiffres',
    otpPlaceholder:   '• • • • • •',
    validate:         'Valider',
    resendCode:       'Renvoyer le code',
    back:             'Retour',

    // Errors
    errInvalidPhone:  'Numéro de téléphone invalide',
    errInvalidOtp:    'Code incorrect ou expiré',
    errNetwork:       'Erreur réseau, réessayez',
  },
  ar: {
    // Navbar
    sell:           'بيع',
    profile:        'حسابي',
    register:       'تسجيل',
    createProfile:  'إنشاء حسابي',
    alreadyAccount: 'لدي حساب',
    myProfile:      'ملفي الشخصي',
    myShop:         'متجري',
    logout:         'تسجيل الخروج',
    location:       'الدار البيضاء • فرد',

    // Auth modal — phone step
    loginTitle:       'تسجيل الدخول',
    registerTitle:    'إنشاء حساب',
    phoneLabel:       'رقم الهاتف',
    phonePlaceholder: '+212 6XX XXX XXX',
    sendCode:         'إرسال الرمز',
    channelSms:       'SMS',
    channelWhatsapp:  'واتساب',

    // Auth modal — OTP step
    otpTitle:         'التحقق',
    otpSentTo:        'تم إرسال الرمز إلى',
    otpLabel:         'رمز مكون من 6 أرقام',
    otpPlaceholder:   '• • • • • •',
    validate:         'تأكيد',
    resendCode:       'إعادة إرسال الرمز',
    back:             'رجوع',

    // Errors
    errInvalidPhone:  'رقم الهاتف غير صالح',
    errInvalidOtp:    'الرمز غير صحيح أو منتهي الصلاحية',
    errNetwork:       'خطأ في الشبكة، حاول مرة أخرى',
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
