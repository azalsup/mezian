import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../core/services/auth.service';
import { ExternalUsersTabComponent } from './external-users-tab/external-users-tab.component';
import { InternalUsersTabComponent } from './internal-users-tab/internal-users-tab.component';

type Tab = 'external' | 'internal';

@Component({
  selector: 'app-admin-users',
  standalone: true,
  imports: [CommonModule, ExternalUsersTabComponent, InternalUsersTabComponent],
  templateUrl: './admin-users.component.html',
})
export class AdminUsersComponent {
  private readonly auth = inject(AuthService);

  activeTab = signal<Tab>('external');

  /** Role of the currently logged-in admin/moderator */
  get callerRole(): string {
    return this.auth.currentUser()?.role ?? 'moderator';
  }

  setTab(tab: Tab): void { this.activeTab.set(tab); }
}
