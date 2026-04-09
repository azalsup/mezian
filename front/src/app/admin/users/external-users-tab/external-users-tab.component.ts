import { Component, Input, inject, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AdminApi } from '../../../sdk/admin.api';
import { UserTileComponent } from '../user-tile/user-tile.component';
import type { User } from '../../../sdk/types';

@Component({
  selector: 'app-external-users-tab',
  standalone: true,
  imports: [CommonModule, UserTileComponent],
  templateUrl: './external-users-tab.component.html',
})
export class ExternalUsersTabComponent implements OnInit {
  private readonly adminApi = inject(AdminApi);

  @Input() callerRole = 'admin';

  users    = signal<User[]>([]);
  total    = signal(0);
  loading  = signal(true);
  error    = signal('');
  page     = 1;
  pageSize = 24;

  ngOnInit(): void { this.load(); }

  load(): void {
    this.loading.set(true);
    this.adminApi.listUsers(this.page, this.pageSize, 'external').subscribe({
      next: res => { this.users.set(res.data); this.total.set(res.total); this.loading.set(false); },
      error: () => { this.error.set('Erreur lors du chargement.'); this.loading.set(false); },
    });
  }

  onAction(): void { this.load(); }

  get totalPages(): number { return Math.ceil(this.total() / this.pageSize); }
  prevPage(): void { if (this.page > 1) { this.page--; this.load(); } }
  nextPage(): void { if (this.page < this.totalPages) { this.page++; this.load(); } }
}
