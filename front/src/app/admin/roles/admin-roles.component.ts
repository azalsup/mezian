import { Component, inject, signal, computed, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { KeyValuePipe } from '@angular/common';
import { AdminApi, RolePayload } from '../../sdk/admin.api';
import { Permission, Role } from '../../sdk/types';

@Component({
  selector: 'app-admin-roles',
  imports: [FormsModule, KeyValuePipe],
  templateUrl: './admin-roles.component.html',
})
export class AdminRolesComponent implements OnInit {
  private readonly adminApi = inject(AdminApi);

  roles       = signal<Role[]>([]);
  permissions = signal<Permission[]>([]);
  loading     = signal(true);
  error       = signal('');

  // ── Form state ─────────────────────────────────────────────────────────────
  showForm    = signal(false);
  saving      = signal(false);
  formError   = signal('');
  editingId   = signal<number | null>(null);

  formName        = '';
  formSlug        = '';
  formDescription = '';
  formPermIds     = signal<Set<number>>(new Set());

  // Groups permissions by their group field
  permGroups = computed(() => {
    const map = new Map<string, Permission[]>();
    for (const p of this.permissions()) {
      const list = map.get(p.group) ?? [];
      list.push(p);
      map.set(p.group, list);
    }
    return map;
  });

  ngOnInit(): void {
    this.load();
  }

  private load(): void {
    this.loading.set(true);
    this.adminApi.listRoles().subscribe({
      next: roles => {
        this.roles.set(roles);
        this.adminApi.listPermissions().subscribe({
          next: perms => { this.permissions.set(perms); this.loading.set(false); },
          error: () => { this.error.set('Erreur lors du chargement des permissions.'); this.loading.set(false); },
        });
      },
      error: () => { this.error.set('Erreur lors du chargement des rôles.'); this.loading.set(false); },
    });
  }

  // ── Form ───────────────────────────────────────────────────────────────────

  openCreate(): void {
    this.editingId.set(null);
    this.formName = '';
    this.formSlug = '';
    this.formDescription = '';
    this.formPermIds.set(new Set());
    this.formError.set('');
    this.showForm.set(true);
  }

  openEdit(role: Role): void {
    this.editingId.set(role.id);
    this.formName        = role.name;
    this.formSlug        = role.slug;
    this.formDescription = role.description;
    this.formPermIds.set(new Set(role.permissions?.map(p => p.id) ?? []));
    this.formError.set('');
    this.showForm.set(true);
  }

  closeForm(): void { this.showForm.set(false); }

  togglePerm(id: number): void {
    const s = new Set(this.formPermIds());
    s.has(id) ? s.delete(id) : s.add(id);
    this.formPermIds.set(s);
  }

  isPermSelected(id: number): boolean {
    return this.formPermIds().has(id);
  }

  slugify(): void {
    this.formSlug = this.formName
      .toLowerCase()
      .normalize('NFD').replace(/[\u0300-\u036f]/g, '')
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '');
  }

  submit(): void {
    this.formError.set('');
    if (!this.formName.trim() || !this.formSlug.trim()) {
      this.formError.set('Nom et slug sont requis.');
      return;
    }
    const payload: RolePayload = {
      name:           this.formName.trim(),
      slug:           this.formSlug.trim(),
      description:    this.formDescription.trim(),
      permission_ids: [...this.formPermIds()],
    };
    this.saving.set(true);
    const id = this.editingId();
    const req = id
      ? this.adminApi.updateRole(id, payload)
      : this.adminApi.createRole(payload);

    req.subscribe({
      next: () => { this.saving.set(false); this.showForm.set(false); this.load(); },
      error: () => { this.saving.set(false); this.formError.set('Erreur lors de la sauvegarde.'); },
    });
  }

  deleteRole(role: Role): void {
    if (!confirm(`Supprimer le rôle « ${role.name} » ?`)) return;
    this.adminApi.deleteRole(role.id).subscribe({
      next: () => this.load(),
      error: () => this.error.set('Impossible de supprimer ce rôle.'),
    });
  }
}
