import { Component, Input, Output, EventEmitter, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { LangService } from '../../../core/services/lang.service';
import { environment } from '../../../../environments/environment';

export type AuthScreen =
  | 'login'
  | 'otp-phone'
  | 'otp-code'
  | 'reg-credentials'
  | 'reg-identity'
  | 'reg-address';

const MOROCCO_CITIES = [
  'Casablanca', 'Rabat', 'Marrakech', 'Fès', 'Tanger', 'Agadir',
  'Meknès', 'Oujda', 'Kénitra', 'Tétouan', 'Safi', 'El Jadida',
  'Béni Mellal', 'Nador', 'Mohammedia', 'Laâyoune',
];

@Component({
  selector: 'app-auth-form',
  imports: [FormsModule],
  templateUrl: './auth-form.component.html',
  styleUrl: './auth-form.component.scss',
})
export class AuthFormComponent implements OnInit {
  /** Initial screen to display */
  @Input() initialScreen: AuthScreen = 'login';
  /** When true, navigates to / after success instead of emitting done */
  @Input() asPage = false;
  /** Emitted when auth succeeds (used by modal to close itself) */
  @Output() done = new EventEmitter<void>();

  protected readonly auth   = inject(AuthService);
  protected readonly lang   = inject(LangService);
  protected readonly router = inject(Router);
  protected readonly cfg    = environment.auth;
  protected readonly cities = MOROCCO_CITIES;

  screen  = signal<AuthScreen>('login');
  loading = signal(false);
  error   = signal('');

  // Login fields
  loginId  = '';
  loginPwd = '';

  // OTP fields
  otpPhone = '';
  otpCode  = '';

  // Register fields
  regPhone      = '';
  regEmail      = '';
  regPassword   = '';
  regConfirm    = '';
  regName       = '';
  regAddress    = '';
  regCity       = '';
  regPostalCode = '';
  regCountry    = environment.auth.defaultCountry;

  ngOnInit(): void {
    this.screen.set(this.initialScreen);
  }

  // ── Login ──────────────────────────────────────────────────────────────────

  submitLogin(): void {
    this.error.set('');
    if (!this.loginId.trim()) { this.error.set(this.lang.t('errIdentifier')); return; }
    if (!this.loginPwd)       { this.error.set(this.lang.t('errPasswordMin')); return; }
    this.loading.set(true);
    this.auth.login(this.loginId.trim(), this.loginPwd).subscribe({
      next:  () => { this.loading.set(false); this.onSuccess(); },
      error: () => { this.loading.set(false); this.error.set(this.lang.t('errCredentials')); },
    });
  }

  // ── OTP login ──────────────────────────────────────────────────────────────

  submitOtpPhone(): void {
    this.error.set('');
    if (!this.isValidPhone(this.otpPhone)) { this.error.set(this.lang.t('errInvalidPhone')); return; }
    this.loading.set(true);
    this.auth.sendOtp(this.otpPhone).subscribe({
      next:  () => { this.loading.set(false); this.screen.set('otp-code'); },
      error: () => { this.loading.set(false); this.error.set(this.lang.t('errNetwork')); },
    });
  }

  submitOtpCode(): void {
    this.error.set('');
    const code = this.otpCode.replace(/\s/g, '');
    if (code.length < 6) { this.error.set(this.lang.t('errInvalidOtp')); return; }
    this.loading.set(true);
    this.auth.verifyOtp(this.otpPhone, code).subscribe({
      next:  () => { this.loading.set(false); this.onSuccess(); },
      error: () => { this.loading.set(false); this.error.set(this.lang.t('errInvalidOtp')); },
    });
  }

  resendOtp(): void {
    this.otpCode = '';
    this.error.set('');
    this.auth.sendOtp(this.otpPhone).subscribe({ error: () => {} });
  }

  // ── Register ───────────────────────────────────────────────────────────────

  submitRegCredentials(): void {
    this.error.set('');
    const phone = this.regPhone.trim();
    const email = this.regEmail.trim();
    if (this.cfg.phoneRequired && !phone)  { this.error.set(this.lang.t('errInvalidPhone')); return; }
    if (this.cfg.emailRequired && !email)  { this.error.set(this.lang.t('errInvalidEmail')); return; }
    if (!phone && !email)                  { this.error.set(this.lang.t('errPhoneOrEmail')); return; }
    if (phone && !this.isValidPhone(phone)){ this.error.set(this.lang.t('errInvalidPhone')); return; }
    if (email && !this.isValidEmail(email)){ this.error.set(this.lang.t('errInvalidEmail')); return; }
    if (this.regPassword.length < 8)       { this.error.set(this.lang.t('errPasswordMin')); return; }
    if (this.regPassword !== this.regConfirm){ this.error.set(this.lang.t('errPasswordMatch')); return; }
    this.screen.set('reg-identity');
  }

  submitRegIdentity(): void {
    this.error.set('');
    if (!this.regName.trim()) { this.error.set(this.lang.t('errDisplayName')); return; }
    this.screen.set('reg-address');
  }

  submitRegAddress(): void {
    this.error.set('');
    this.loading.set(true);
    this.auth.register({
      phone:        this.regPhone.trim()      || undefined,
      email:        this.regEmail.trim()      || undefined,
      password:     this.regPassword,
      display_name: this.regName.trim(),
      address:      this.regAddress.trim()    || undefined,
      city:         this.regCity.trim()       || undefined,
      postal_code:  this.regPostalCode.trim() || undefined,
      country:      this.regCountry           || environment.auth.defaultCountry,
    }).subscribe({
      next: () => { this.loading.set(false); this.onSuccess(); },
      error: (err: { error?: { error?: string } }) => {
        this.loading.set(false);
        const msg = err?.error?.error ?? '';
        if (msg.includes('phone'))      this.error.set(this.lang.t('errInvalidPhone'));
        else if (msg.includes('email')) this.error.set(this.lang.t('errInvalidEmail'));
        else                            this.error.set(this.lang.t('errNetwork'));
      },
    });
  }

  // ── Navigation ─────────────────────────────────────────────────────────────

  goLogin(): void      { this.error.set(''); this.screen.set('login'); }
  goRegister(): void   { this.error.set(''); this.screen.set('reg-credentials'); }
  goOtp(): void        { this.error.set(''); this.screen.set('otp-phone'); }
  goBack(): void {
    const prev: Partial<Record<AuthScreen, AuthScreen>> = {
      'otp-code':      'otp-phone',
      'reg-identity':  'reg-credentials',
      'reg-address':   'reg-identity',
      'otp-phone':     'login',
    };
    const target = prev[this.screen()];
    if (target) { this.error.set(''); this.screen.set(target); }
  }

  get regStep(): number {
    const map: Partial<Record<AuthScreen, number>> = {
      'reg-credentials': 1, 'reg-identity': 2, 'reg-address': 3,
    };
    return map[this.screen()] ?? 0;
  }

  // ── Helpers ────────────────────────────────────────────────────────────────

  private onSuccess(): void {
    if (this.asPage) {
      this.router.navigate(['/']);
    } else {
      this.done.emit();
    }
  }

  private isValidPhone(p: string): boolean {
    return /^(\+212|0)[5-7]\d{8}$/.test(p.replace(/\s/g, ''));
  }

  private isValidEmail(e: string): boolean {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e);
  }
}
