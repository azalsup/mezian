import { Component, HostListener, inject, ElementRef } from '@angular/core';
import { LangService } from '../../core/services/lang.service';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss',
})
export class NavbarComponent {
  protected readonly lang = inject(LangService);
  protected readonly auth = inject(AuthService);
  private  readonly el   = inject(ElementRef);

  showProfileMenu = false;

  onSellClick(): void {
    if (!this.auth.isLoggedIn()) {
      this.auth.openModal('register');
    }
    // TODO: navigate to ad creation
  }

  toggleProfileMenu(): void {
    this.showProfileMenu = !this.showProfileMenu;
  }

  openRegister(): void {
    this.showProfileMenu = false;
    this.auth.openModal('register');
  }

  openLogin(): void {
    this.showProfileMenu = false;
    this.auth.openModal('login');
  }

  logout(): void {
    this.showProfileMenu = false;
    this.auth.logout();
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent): void {
    if (!this.el.nativeElement.contains(event.target)) {
      this.showProfileMenu = false;
    }
  }
}
