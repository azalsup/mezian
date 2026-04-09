import {
  Component, Input, Output, EventEmitter,
  inject, signal, HostListener, OnChanges
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AdminApi } from '../../../sdk/admin.api';
import type { Role, User } from '../../../sdk/types';

export type UserAction = 'updated' | 'banned' | 'unbanned' | 'deleted' | 'rolesChanged';

@Component({
  selector: 'app-user-tile',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './user-tile.component.html',
  host: { style: 'display: block; position: relative;' },
})
export class UserTileComponent implements OnChanges {
  private readonly adminApi = inject(AdminApi);

  @Input({ required: true }) user!: User;
  /** 'admin' → all actions; 'moderator' → ban only */
  @Input() callerRole: string = 'admin';
  /** Pass available roles for RBAC assignment (optional, fetched if not provided) */
  @Input() allRoles: Role[] = [];

  @Output() readonly action = new EventEmitter<UserAction>();

  // ── Dropdown ────────────────────────────────────────────────────────────────
  menuOpen = signal(false);

  toggleMenu(e: MouseEvent): void { e.stopPropagation(); this.menuOpen.set(!this.menuOpen()); }

  @HostListener('document:click')
  closeMenu(): void { this.menuOpen.set(false); }

  // ── Edit user modal ─────────────────────────────────────────────────────────
  showEdit    = signal(false);
  editForm    = signal({ display_name: '', phone: '', email: '', address: '',
                         city: '', postal_code: '', country: '', is_verified: false, role: 'user' });
  saving      = signal(false);
  editError   = signal('');

  openEdit(): void {
    const u = this.user;
    this.editForm.set({
      display_name: u.display_name,
      phone:        u.phone,
      email:        u.email        ?? '',
      address:      u.address      ?? '',
      city:         u.city         ?? '',
      postal_code:  u.postal_code  ?? '',
      country:      u.country      ?? '',
      is_verified:  u.is_verified,
      role:         u.role,
    });
    this.editError.set('');
    this.showEdit.set(true);
  }

  patchEdit(field: string, value: unknown): void {
    this.editForm.set({ ...this.editForm(), [field]: value });
  }

  saveEdit(): void {
    const f = this.editForm();
    if (!f.display_name.trim()) { this.editError.set('Le nom est requis.'); return; }
    if (!f.phone.trim())        { this.editError.set('Le téléphone est requis.'); return; }
    this.saving.set(true);
    this.adminApi.updateUser(this.user.id, {
      display_name: f.display_name.trim(),
      phone:        f.phone.trim(),
      email:        f.email        || null,
      address:      f.address      || null,
      city:         f.city         || null,
      postal_code:  f.postal_code  || null,
      country:      f.country      || null,
      is_verified:  f.is_verified,
      role:         f.role,
    }).subscribe({
      next: () => { this.saving.set(false); this.showEdit.set(false); this.action.emit('updated'); },
      error: () => { this.saving.set(false); this.editError.set('Erreur lors de la sauvegarde.'); },
    });
  }

  // ── Confirm (ban / unban / delete) ──────────────────────────────────────────
  showConfirm   = signal(false);
  confirmAct    = signal<'ban' | 'unban' | 'delete' | null>(null);
  confirming    = signal(false);
  confirmError  = signal('');

  openConfirm(act: 'ban' | 'unban' | 'delete'): void {
    this.confirmAct.set(act); this.confirmError.set(''); this.showConfirm.set(true);
  }

  doConfirm(): void {
    const act = this.confirmAct();
    if (!act) return;
    this.confirming.set(true);
    const obs = act === 'ban'   ? this.adminApi.banUser(this.user.id)
              : act === 'unban' ? this.adminApi.unbanUser(this.user.id)
              :                   this.adminApi.deleteUser(this.user.id);
    obs.subscribe({
      next: () => {
        this.confirming.set(false); this.showConfirm.set(false);
        this.action.emit(act === 'ban' ? 'banned' : act === 'unban' ? 'unbanned' : 'deleted');
      },
      error: () => { this.confirming.set(false); this.confirmError.set('Erreur lors de l\'opération.'); },
    });
  }

  // ── Reset password modal ────────────────────────────────────────────────────
  showReset   = signal(false);
  newPassword = signal('');
  resetting   = signal(false);
  resetError  = signal('');

  openReset(): void { this.newPassword.set(''); this.resetError.set(''); this.showReset.set(true); }

  doReset(): void {
    if (this.newPassword().length < 6) { this.resetError.set('Minimum 6 caractères.'); return; }
    this.resetting.set(true);
    this.adminApi.resetUserPassword(this.user.id, this.newPassword()).subscribe({
      next: () => { this.resetting.set(false); this.showReset.set(false); },
      error: () => { this.resetting.set(false); this.resetError.set('Erreur lors de la réinitialisation.'); },
    });
  }

  // ── RBAC roles modal ────────────────────────────────────────────────────────
  showRoles       = signal(false);
  selectedRoleIds = signal<Set<number>>(new Set());
  savingRoles     = signal(false);
  rolesError      = signal('');
  loadedRoles     = signal<Role[]>([]);

  openRoles(): void {
    this.selectedRoleIds.set(new Set(this.user.roles?.map(r => r.id) ?? []));
    this.rolesError.set('');
    if (this.allRoles.length > 0) {
      this.loadedRoles.set(this.allRoles);
    } else {
      this.adminApi.listRoles().subscribe({ next: r => this.loadedRoles.set(r) });
    }
    this.showRoles.set(true);
  }

  toggleRole(id: number): void {
    const s = new Set(this.selectedRoleIds());
    s.has(id) ? s.delete(id) : s.add(id);
    this.selectedRoleIds.set(s);
  }

  saveRoles(): void {
    this.savingRoles.set(true);
    this.adminApi.setUserRoles(this.user.id, [...this.selectedRoleIds()]).subscribe({
      next: () => { this.savingRoles.set(false); this.showRoles.set(false); this.action.emit('rolesChanged'); },
      error: () => { this.savingRoles.set(false); this.rolesError.set('Erreur lors de la mise à jour.'); },
    });
  }

  // ── Computed helpers ────────────────────────────────────────────────────────
  get canBan():   boolean { return this.callerRole === 'admin' || this.callerRole === 'moderator'; }
  get canEdit():  boolean { return this.callerRole === 'admin'; }
  get canDelete():boolean { return this.callerRole === 'admin'; }
  get canReset(): boolean { return this.callerRole === 'admin'; }
  get canRoles(): boolean { return this.callerRole === 'admin'; }

  get initials(): string {
    return this.user.display_name.split(' ').map(w => w[0]).join('').slice(0, 2).toUpperCase();
  }

  ngOnChanges(): void {
    if (this.allRoles.length && this.showRoles()) this.loadedRoles.set(this.allRoles);
  }
}
