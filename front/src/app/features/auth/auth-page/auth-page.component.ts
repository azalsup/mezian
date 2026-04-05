import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AuthFormComponent, AuthScreen } from '../auth-form/auth-form.component';

@Component({
  selector: 'app-auth-page',
  imports: [AuthFormComponent],
  templateUrl: './auth-page.component.html',
  styleUrl: './auth-page.component.scss',
})
export class AuthPageComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  initialScreen = signal<AuthScreen>('login');

  ngOnInit(): void {
    const screen = this.route.snapshot.data['screen'] as AuthScreen | undefined;
    if (screen) this.initialScreen.set(screen);
  }
}
