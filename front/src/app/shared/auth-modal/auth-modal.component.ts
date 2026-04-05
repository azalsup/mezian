import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../core/services/auth.service';
import { LangService } from '../../core/services/lang.service';

type Step = 'phone' | 'otp';

@Component({
  selector: 'app-auth-modal',
  imports: [FormsModule],
  templateUrl: './auth-modal.component.html',
  styleUrl: './auth-modal.component.scss',
})
export class AuthModalComponent {
  protected readonly auth = inject(AuthService);
  protected readonly lang = inject(LangService);

  step    = signal<Step>('phone');
  loading = signal(false);
  error   = signal('');

  phone   = '';
  otp     = '';

  // ── Step 1 : send OTP ────────────────────────────────────────────────────

  submitPhone(): void {
    this.error.set('');
    if (!this.isValidPhone(this.phone)) {
      this.error.set(this.lang.t('errInvalidPhone'));
      return;
    }
    this.loading.set(true);
    this.auth.sendOtp(this.phone).subscribe({
      next: () => {
        this.loading.set(false);
        this.step.set('otp');
      },
      error: () => {
        this.loading.set(false);
        this.error.set(this.lang.t('errNetwork'));
      },
    });
  }

  // ── Step 2 : verify OTP ──────────────────────────────────────────────────

  submitOtp(): void {
    this.error.set('');
    if (this.otp.replace(/\s/g, '').length < 6) {
      this.error.set(this.lang.t('errInvalidOtp'));
      return;
    }
    this.loading.set(true);
    this.auth.verifyOtp(this.phone, this.otp.replace(/\s/g, '')).subscribe({
      next: () => {
        // AuthService.verifyOtp calls setSession → closes modal automatically
        this.loading.set(false);
        this.reset();
      },
      error: () => {
        this.loading.set(false);
        this.error.set(this.lang.t('errInvalidOtp'));
      },
    });
  }

  resend(): void {
    this.otp   = '';
    this.error.set('');
    this.auth.sendOtp(this.phone).subscribe({ error: () => {} });
  }

  goBack(): void {
    this.step.set('phone');
    this.otp = '';
    this.error.set('');
  }

  close(): void {
    this.auth.closeModal();
    this.reset();
  }

  // ── Helpers ──────────────────────────────────────────────────────────────

  private reset(): void {
    this.step.set('phone');
    this.phone = '';
    this.otp   = '';
    this.error.set('');
  }

  private isValidPhone(phone: string): boolean {
    // Accept +212XXXXXXXXX or 0XXXXXXXXX (Morocco)
    return /^(\+212|0)[5-7]\d{8}$/.test(phone.replace(/\s/g, ''));
  }
}
