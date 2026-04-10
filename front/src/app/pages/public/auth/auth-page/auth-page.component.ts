import { Component, inject, OnInit, signal, effect } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthFormComponent, AuthScreen } from '../auth-form/auth-form.component';
import { AuthService } from '../../../../core/services/auth.service';

@Component({
  selector: 'app-auth-page',
  imports: [AuthFormComponent],
  templateUrl: './auth-page.component.html',
  styleUrl: './auth-page.component.scss',
})
export class AuthPageComponent implements OnInit {
  private readonly route  = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly auth   = inject(AuthService);

  initialScreen = signal<AuthScreen>('login');

  constructor() {
    // Watch currentUser — redirect as soon as a session is detected,
    // whether it was already in localStorage or just set after login.
    effect(() => {
      if (this.auth.currentUser() !== null) {
        this.router.navigate(['/'], { replaceUrl: true });
      }
    });
  }

  ngOnInit(): void {
    const screen = this.route.snapshot.data['screen'] as AuthScreen | undefined;
    if (screen) this.initialScreen.set(screen);
  }
}
