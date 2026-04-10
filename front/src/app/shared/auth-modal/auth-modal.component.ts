import { Component, inject } from '@angular/core';
import { AuthService } from '../../core/services/auth.service';
import { AuthFormComponent } from '../../pages/public/auth/auth-form/auth-form.component';

@Component({
  selector: 'app-auth-modal',
  imports: [AuthFormComponent],
  templateUrl: './auth-modal.component.html',
  styleUrl: './auth-modal.component.scss',
})
export class AuthModalComponent {
  protected readonly auth = inject(AuthService);

  close(): void {
    this.auth.closeModal();
  }
}
