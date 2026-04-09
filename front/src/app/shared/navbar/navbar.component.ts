import { Component, HostListener, inject, ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { LangService, Lang } from '../../core/services/lang.service';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-navbar',
  imports: [CommonModule, RouterModule],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss',
})
export class NavbarComponent {
  protected readonly lang   = inject(LangService);
  protected readonly auth   = inject(AuthService);
  private  readonly router  = inject(Router);
  private  readonly el      = inject(ElementRef);

  showProfileMenu = false;
  showLangMenu    = false;
  showMobileMenu  = false;

  onSellClick(): void {
    if (!this.auth.isLoggedIn()) {
      this.router.navigate(['/register']);
    }
  }

  toggleProfileMenu(event: MouseEvent): void {
    event.stopPropagation();
    this.showProfileMenu = !this.showProfileMenu;
    this.showLangMenu = false;
  }

  toggleLangMenu(event: MouseEvent): void {
    event.stopPropagation();
    this.showLangMenu = !this.showLangMenu;
    this.showProfileMenu = false;
  }

  toggleMobileMenu(event: MouseEvent): void {
    event.stopPropagation();
    this.showMobileMenu = !this.showMobileMenu;
  }

  selectLang(code: Lang, event: MouseEvent): void {
    event.stopPropagation();
    this.lang.setLang(code);
    this.showLangMenu = false;
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

  goHome(): void {
    this.router.navigate(['/']);
  }

  @HostListener('document:click')
  onDocumentClick(): void {
    this.showProfileMenu = false;
    this.showLangMenu    = false;
    this.showMobileMenu  = false;
  }
}
