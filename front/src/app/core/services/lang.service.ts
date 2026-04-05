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

    // Auth modal — shared
    back:           'Retour',
    next:           'Suivant',
    validate:       'Valider',
    loading:        'Chargement…',
    orSeparator:    'ou',

    // Auth modal — login
    loginTitle:     'Connexion',
    loginSubtitle:  'Connectez-vous avec votre numéro ou e-mail',
    identifierLabel: 'Téléphone ou e-mail',
    identifierPlaceholder: '+212 6XX XXX XXX ou email@exemple.ma',
    passwordLabel:  'Mot de passe',
    passwordPlaceholder: '••••••••',
    loginBtn:       'Se connecter',
    noAccount:      'Pas encore de compte ?',
    createAccount:  'Créer un compte',
    loginWithOtp:   'Se connecter par SMS',

    // Auth modal — register step 1 (credentials)
    registerTitle:     'Créer un compte',
    registerSubtitle:  'Renseignez vos informations de connexion',
    phoneLabel:        'Numéro de téléphone',
    phonePlaceholder:  '+212 6XX XXX XXX',
    emailLabel:        'Adresse e-mail',
    emailPlaceholder:  'email@exemple.ma',
    optional:          '(optionnel)',
    confirmPasswordLabel: 'Confirmer le mot de passe',
    confirmPasswordPlaceholder: '••••••••',
    hasAccount:        'Déjà un compte ?',
    login:             'Se connecter',

    // Auth modal — register step 2 (identity)
    identityTitle:    'Votre identité',
    identitySubtitle: 'Comment souhaitez-vous être affiché ?',
    displayNameLabel: 'Nom affiché',
    displayNamePlaceholder: 'Ex : Ahmed Benali',

    // Auth modal — register step 3 (address)
    addressTitle:      'Votre adresse',
    addressSubtitle:   'Optionnel — utilisée pour localiser vos annonces',
    addressLabel:      'Adresse',
    addressPlaceholder: '12 rue Mohammed V',
    cityLabel:         'Ville',
    cityPlaceholder:   'Casablanca',
    postalCodeLabel:   'Code postal',
    postalCodePlaceholder: '20000',
    countryLabel:      'Pays',
    finishBtn:         "Créer mon compte",

    // OTP step
    sendCode:         'Recevoir le code',
    otpTitle:         'Vérification',
    otpSentTo:        'Code envoyé au',
    otpLabel:         'Code à 6 chiffres',
    otpPlaceholder:   '• • • • • •',
    resendCode:       'Renvoyer le code',

    // Errors
    errInvalidPhone:    'Numéro de téléphone invalide',
    errInvalidEmail:    'Adresse e-mail invalide',
    errPasswordMin:     'Le mot de passe doit contenir au moins 8 caractères',
    errPasswordMatch:   'Les mots de passe ne correspondent pas',
    errDisplayName:     'Veuillez saisir un nom',
    errIdentifier:      'Téléphone ou e-mail requis',
    errInvalidOtp:      'Code incorrect ou expiré',
    errNetwork:         'Erreur réseau, réessayez',
    errPhoneOrEmail:    'Téléphone ou e-mail requis',
    errCredentials:     'Identifiants incorrects',
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

    // Auth modal — shared
    back:           'رجوع',
    next:           'التالي',
    validate:       'تأكيد',
    loading:        'جاري التحميل…',
    orSeparator:    'أو',

    // Auth modal — login
    loginTitle:     'تسجيل الدخول',
    loginSubtitle:  'سجّل الدخول برقم هاتفك أو بريدك',
    identifierLabel: 'الهاتف أو البريد الإلكتروني',
    identifierPlaceholder: '+212 6XX XXX XXX أو email@exemple.ma',
    passwordLabel:  'كلمة المرور',
    passwordPlaceholder: '••••••••',
    loginBtn:       'تسجيل الدخول',
    noAccount:      'ليس لديك حساب؟',
    createAccount:  'إنشاء حساب',
    loginWithOtp:   'الدخول عبر SMS',

    // Auth modal — register step 1 (credentials)
    registerTitle:    'إنشاء حساب',
    registerSubtitle: 'أدخل معلومات تسجيل الدخول',
    phoneLabel:       'رقم الهاتف',
    phonePlaceholder: '+212 6XX XXX XXX',
    emailLabel:       'البريد الإلكتروني',
    emailPlaceholder: 'email@exemple.ma',
    optional:         '(اختياري)',
    confirmPasswordLabel: 'تأكيد كلمة المرور',
    confirmPasswordPlaceholder: '••••••••',
    hasAccount:       'لديك حساب بالفعل؟',
    login:            'تسجيل الدخول',

    // Auth modal — register step 2 (identity)
    identityTitle:    'هويتك',
    identitySubtitle: 'كيف تريد أن تظهر للآخرين؟',
    displayNameLabel: 'الاسم المعروض',
    displayNamePlaceholder: 'مثال: أحمد بنعلي',

    // Auth modal — register step 3 (address)
    addressTitle:      'عنوانك',
    addressSubtitle:   'اختياري — يُستخدم لتحديد موقع إعلاناتك',
    addressLabel:      'العنوان',
    addressPlaceholder: '12 شارع محمد الخامس',
    cityLabel:         'المدينة',
    cityPlaceholder:   'الدار البيضاء',
    postalCodeLabel:   'الرمز البريدي',
    postalCodePlaceholder: '20000',
    countryLabel:      'البلد',
    finishBtn:         'إنشاء حسابي',

    // OTP step
    sendCode:         'إرسال الرمز',
    otpTitle:         'التحقق',
    otpSentTo:        'تم إرسال الرمز إلى',
    otpLabel:         'رمز مكون من 6 أرقام',
    otpPlaceholder:   '• • • • • •',
    resendCode:       'إعادة إرسال الرمز',

    // Errors
    errInvalidPhone:    'رقم الهاتف غير صالح',
    errInvalidEmail:    'البريد الإلكتروني غير صالح',
    errPasswordMin:     'يجب أن تحتوي كلمة المرور على 8 أحرف على الأقل',
    errPasswordMatch:   'كلمتا المرور غير متطابقتين',
    errDisplayName:     'الرجاء إدخال اسم',
    errIdentifier:      'الهاتف أو البريد الإلكتروني مطلوب',
    errInvalidOtp:      'الرمز غير صحيح أو منتهي الصلاحية',
    errNetwork:         'خطأ في الشبكة، حاول مرة أخرى',
    errPhoneOrEmail:    'الهاتف أو البريد الإلكتروني مطلوب',
    errCredentials:     'بيانات الاعتماد غير صحيحة',
  },
};

@Injectable({ providedIn: 'root' })
export class LangService {
  readonly current = signal<Lang>('fr');
  readonly isRtl   = computed(() => this.current() === 'ar');

  t(key: string): string {
    const lang: Lang = this.current();
    return TRANSLATIONS[lang][key] ?? key;
  }

  toggle(): void {
    const next: Lang = this.current() === 'fr' ? 'ar' : 'fr';
    this.current.set(next);
    document.documentElement.lang = next;
    document.documentElement.dir  = next === 'ar' ? 'rtl' : 'ltr';
  }
}
