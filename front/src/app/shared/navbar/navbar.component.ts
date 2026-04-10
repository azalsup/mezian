import { Component, HostListener, inject, ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { LangService, Lang } from '../../core/services/lang.service';
import { AuthService } from '../../core/services/auth.service';
import { ProfileDropdownComponent } from '../profile-dropdown/profile-dropdown.component';

@Component({
  selector: 'app-navbar',
  imports: [CommonModule, RouterModule, ProfileDropdownComponent],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss',
})
export class NavbarComponent {
  protected readonly lang   = inject(LangService);
  protected readonly auth   = inject(AuthService);
  private  readonly router  = inject(Router);
  private  readonly el      = inject(ElementRef);

  showLangMenu    = false;
  showMobileMenu  = false;

  onSellClick(): void {
    if (!this.auth.isLoggedIn()) {
      this.router.navigate(['/register']);
    }
  }

  toggleLangMenu(event: MouseEvent): void {
    event.stopPropagation();
    this.showLangMenu = !this.showLangMenu;
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

  goHome(): void {
    this.router.navigate(['/']);
  }

  @HostListener('document:click')
  onDocumentClick(): void {
    this.showLangMenu    = false;
    this.showMobileMenu  = false;
  }
}
