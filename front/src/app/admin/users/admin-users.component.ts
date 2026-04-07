import { Component, inject, signal, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AdminApi } from '../../sdk/admin.api';
import { Role, User } from '../../sdk/types';

@Component({
  selector: 'app-admin-users',
  imports: [FormsModule],
  templateUrl: './admin-users.component.html',
})
export class AdminUsersComponent implements OnInit {
  private readonly adminApi = inject(AdminApi);

  users    = signal<User[]>([]);
  total    = signal(0);
  roles    = signal<Role[]>([]);
  loading  = signal(true);
  error    = signal('');
  page     = 1;
  pageSize = 20;

  // Role assignment modal
  showModal    = signal(false);
  editingUser  = signal<User | null>(null);
  selectedRoleIds = signal<Set<number>>(new Set());
  saving       = signal(false);
  modalError   = signal('');

  ngOnInit(): void {
    this.load();
    this.adminApi.listRoles().subscribe({
      next: roles => this.roles.set(roles),
    });
  }

  load(): void {
    this.loading.set(true);
    this.adminApi.listUsers(this.page, this.pageSize).subscribe({
      next: res => {
        this.users.set(res.data);
        this.total.set(res.total);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Erreur lors du chargement des utilisateurs.');
        this.loading.set(false);
      },
    });
  }

  openRoleModal(user: User): void {
    this.editingUser.set(user);
    this.selectedRoleIds.set(new Set(user.roles?.map(r => r.id) ?? []));
    this.modalError.set('');
    this.showModal.set(true);
  }

  closeModal(): void { this.showModal.set(false); }

  toggleRole(id: number): void {
    const s = new Set(this.selectedRoleIds());
    s.has(id) ? s.delete(id) : s.add(id);
    this.selectedRoleIds.set(s);
  }

  isRoleSelected(id: number): boolean {
    return this.selectedRoleIds().has(id);
  }

  saveRoles(): void {
    const user = this.editingUser();
    if (!user) return;
    this.saving.set(true);
    this.adminApi.setUserRoles(user.id, [...this.selectedRoleIds()]).subscribe({
      next: () => { this.saving.set(false); this.showModal.set(false); this.load(); },
      error: () => { this.saving.set(false); this.modalError.set('Erreur lors de la mise à jour.'); },
    });
  }

  get totalPages(): number { return Math.ceil(this.total() / this.pageSize); }

  prevPage(): void { if (this.page > 1) { this.page--; this.load(); } }
  nextPage(): void { if (this.page < this.totalPages) { this.page++; this.load(); } }
}
