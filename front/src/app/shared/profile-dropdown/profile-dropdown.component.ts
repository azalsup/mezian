import { Component, inject, Input, HostListener, computed } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-profile-dropdown',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './profile-dropdown.component.html',
  styleUrls: ['./profile-dropdown.component.scss'],
})
export class ProfileDropdownComponent {
  protected readonly lang = inject(LangService);
  protected readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  @Input() isMobile = false;

  readonly isStaff = computed(() => this.auth.isStaff());

  showProfileMenu = false;

  toggleProfileMenu(event: MouseEvent): void {
    event.stopPropagation();
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

  close(): void {
    this.showProfileMenu = false;
  }

  logout(): void {
    this.showProfileMenu = false;
    this.auth.logout();
    this.router.navigate(['/']);
  }

  @HostListener('document:click')
  onDocumentClick(): void {
    this.showProfileMenu = false;
  }
}