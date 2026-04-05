import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../core/services/auth.service';
import { LangService } from '../../core/services/lang.service';
import { environment } from '../../../environments/environment';

/** Possible screens inside the modal */
type Screen =
  | 'login'          // identifier + password
  | 'otp-phone'      // phone entry before OTP send
  | 'otp-code'       // OTP code entry
  | 'reg-credentials'// step 1: phone / email / password
  | 'reg-identity'   // step 2: display name
  | 'reg-address';   // step 3: address (optional)

const MOROCCO_CITIES = [
  'Casablanca','Rabat','Marrakech','Fès','Tanger','Agadir',
  'Meknès','Oujda','Kénitra','Tétouan','Safi','El Jadida',
  'Béni Mellal','Nador','Mohammedia','Laâyoune',
];

@Component({
  selector: 'app-auth-modal',
  imports: [FormsModule],
  templateUrl: './auth-modal.component.html',
  styleUrl: './auth-modal.component.scss',
})
export class AuthModalComponent {
  protected readonly auth    = inject(AuthService);
  protected readonly lang    = inject(LangService);
  protected readonly cfg     = environment.auth;
  protected readonly cities  = MOROCCO_CITIES;

  screen  = signal<Screen>('login');
  loading = signal(false);
  error   = signal('');

  // ── Login fields ─────────────────────────────────────────────────────────
  loginId   = '';
  loginPwd  = '';

  // ── OTP fields ───────────────────────────────────────────────────────────
  otpPhone  = '';
  otpCode   = '';

  // ── Register fields ──────────────────────────────────────────────────────
  regPhone       = '';
  regEmail       = '';
  regPassword    = '';
  regConfirm     = '';
  regName        = '';
  regAddress     = '';
  regCity        = '';
  regPostalCode  = '';
  regCountry     = environment.auth.defaultCountry;

  // ── Login ────────────────────────────────────────────────────────────────

  submitLogin(): void {
    this.error.set('');
    if (!this.loginId.trim()) { this.error.set(this.lang.t('errIdentifier')); return; }
    if (!this.loginPwd)       { this.error.set(this.lang.t('errPasswordMin')); return; }
    this.loading.set(true);
    this.auth.login(this.loginId.trim(), this.loginPwd).subscribe({
      next:  () => { this.loading.set(false); this.reset(); },
      error: () => { this.loading.set(false); this.error.set(this.lang.t('errCredentials')); },
    });
  }

  // ── OTP login ────────────────────────────────────────────────────────────

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
      next:  () => { this.loading.set(false); this.reset(); },
      error: () => { this.loading.set(false); this.error.set(this.lang.t('errInvalidOtp')); },
    });
  }

  resendOtp(): void {
    this.otpCode = '';
    this.error.set('');
    this.auth.sendOtp(this.otpPhone).subscribe({ error: () => {} });
  }

  // ── Register ─────────────────────────────────────────────────────────────

  submitRegCredentials(): void {
    this.error.set('');
    const phone = this.regPhone.trim();
    const email = this.regEmail.trim();

    if (this.cfg.phoneRequired && !phone) { this.error.set(this.lang.t('errInvalidPhone')); return; }
    if (this.cfg.emailRequired && !email) { this.error.set(this.lang.t('errInvalidEmail')); return; }
    if (!phone && !email)                 { this.error.set(this.lang.t('errPhoneOrEmail')); return; }
    if (phone && !this.isValidPhone(phone)){ this.error.set(this.lang.t('errInvalidPhone')); return; }
    if (email && !this.isValidEmail(email)){ this.error.set(this.lang.t('errInvalidEmail')); return; }
    if (this.regPassword.length < 8)      { this.error.set(this.lang.t('errPasswordMin')); return; }
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
      phone:        this.regPhone.trim() || undefined,
      email:        this.regEmail.trim() || undefined,
      password:     this.regPassword,
      display_name: this.regName.trim(),
      address:      this.regAddress.trim() || undefined,
      city:         this.regCity.trim() || undefined,
      postal_code:  this.regPostalCode.trim() || undefined,
      country:      this.regCountry || environment.auth.defaultCountry,
    }).subscribe({
      next:  () => { this.loading.set(false); this.reset(); },
      error: (err: { error?: { error?: string } }) => {
        this.loading.set(false);
        const msg: string = err?.error?.error ?? '';
        if (msg.includes('phone')) this.error.set(this.lang.t('errInvalidPhone'));
        else if (msg.includes('email')) this.error.set(this.lang.t('errInvalidEmail'));
        else this.error.set(this.lang.t('errNetwork'));
      },
    });
  }

  // ── Navigation helpers ───────────────────────────────────────────────────

  switchToLogin(): void {
    this.error.set('');
    this.screen.set('login');
    this.auth.modalMode.set('login');
  }

  switchToRegister(): void {
    this.error.set('');
    this.screen.set('reg-credentials');
    this.auth.modalMode.set('register');
  }

  switchToOtp(): void {
    this.error.set('');
    this.screen.set('otp-phone');
  }

  backFromOtpCode(): void {
    this.otpCode = '';
    this.error.set('');
    this.screen.set('otp-phone');
  }

  backFromIdentity(): void {
    this.error.set('');
    this.screen.set('reg-credentials');
  }

  backFromAddress(): void {
    this.error.set('');
    this.screen.set('reg-identity');
  }

  close(): void {
    this.auth.closeModal();
    this.reset();
  }

  // ── Validation helpers ───────────────────────────────────────────────────

  private isValidPhone(phone: string): boolean {
    return /^(\+212|0)[5-7]\d{8}$/.test(phone.replace(/\s/g, ''));
  }

  private isValidEmail(email: string): boolean {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  }

  private reset(): void {
    this.screen.set('login');
    this.loginId = ''; this.loginPwd = '';
    this.otpPhone = ''; this.otpCode = '';
    this.regPhone = ''; this.regEmail = ''; this.regPassword = ''; this.regConfirm = '';
    this.regName = ''; this.regAddress = ''; this.regCity = '';
    this.regPostalCode = ''; this.regCountry = environment.auth.defaultCountry;
    this.error.set('');
  }

  // ── Step indicator for register ──────────────────────────────────────────
  get regStep(): number {
    const map: Partial<Record<Screen, number>> = {
      'reg-credentials': 1,
      'reg-identity':    2,
      'reg-address':     3,
    };
    return map[this.screen()] ?? 0;
  }
}
