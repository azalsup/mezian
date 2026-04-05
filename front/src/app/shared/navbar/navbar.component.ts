import { Component, HostListener, inject, ElementRef } from '@angular/core';
import { LangService } from '../../core/services/lang.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss',
})
export class NavbarComponent {
  protected readonly lang = inject(LangService);
  private readonly el    = inject(ElementRef);

  isLoggedIn      = false;
  showProfileMenu = false;
  userPhone       = '+212 6XX XXX XXX';

  onSellClick(): void {
    // TODO: navigate to ad creation
  }

  toggleProfileMenu(): void {
    this.showProfileMenu = !this.showProfileMenu;
  }

  createProfile(): void {
    this.showProfileMenu = false;
    // TODO: navigate to register
  }

  login(): void {
    this.showProfileMenu = false;
    // TODO: navigate to login
  }

  logout(): void {
    this.isLoggedIn      = false;
    this.showProfileMenu = false;
    // TODO: call auth service
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent): void {
    if (!this.el.nativeElement.contains(event.target)) {
      this.showProfileMenu = false;
    }
  }
}
