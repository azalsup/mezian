import { Component, HostListener, inject, ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { LangService } from '../../core/services/lang.service';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss',
})
export class NavbarComponent {
  protected readonly lang   = inject(LangService);
  protected readonly auth   = inject(AuthService);
  private  readonly router  = inject(Router);
  private  readonly el      = inject(ElementRef);

  showProfileMenu = false;

  onSellClick(): void {
    if (!this.auth.isLoggedIn()) {
      this.router.navigate(['/register']);
    }
    // TODO: navigate to ad creation when logged in
  }

  toggleProfileMenu(): void {
    this.showProfileMenu = !this.showProfileMenu;
  }

  openRegister(): void {
    this.showProfileMenu = false;
    this.router.navigate(['/register']);
  }

  openLogin(): void {
    this.showProfileMenu = false;
    this.router.navigate(['/login']);
  }

  logout(): void {
    this.showProfileMenu = false;
    this.auth.logout();
    this.router.navigate(['/']);
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent): void {
    if (!this.el.nativeElement.contains(event.target)) {
      this.showProfileMenu = false;
    }
  }
}
