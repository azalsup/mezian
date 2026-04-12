import { Component, inject, signal, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AdminApi, CategoryCreatePayload, CategoryPayload } from '../../../sdk/admin.api';
import { Category } from '../../../sdk/types';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-admin-categories',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './admin-categories.component.html',
})
export class AdminCategoriesComponent implements OnInit {
  private readonly adminApi = inject(AdminApi);
  protected readonly auth   = inject(AuthService);

  categories = signal<Category[]>([]);
  loading    = signal(true);
  error      = signal('');

  expanded   = signal<Set<number>>(new Set());

  // ── Edit modal ─────────────────────────────────────────────────────────────
  showEdit   = signal(false);
  saving     = signal(false);
  editError  = signal('');
  editingCat = signal<Category | null>(null);

  fNameFr    = '';
  fNameAr    = '';
  fNameEn    = '';
  fIcon      = '';
  fSortOrder = 0;
  fFeatured  = false;
  fIsActive  = true;

  // ── Create modal ───────────────────────────────────────────────────────────
  showCreate    = signal(false);
  createError   = signal('');
  creating      = signal(false);
  createParent  = signal<Category | null>(null); // null = root category

  cSlug      = '';
  cNameFr    = '';
  cNameAr    = '';
  cNameEn    = '';
  cIcon      = '';
  cSortOrder = 0;
  cFeatured  = false;
  cIsActive  = true;

  get isAdmin(): boolean { return this.auth.currentUser()?.role === 'admin'; }

  ngOnInit(): void { this.load(); }

  private load(): void {
    this.loading.set(true);
    this.adminApi.listAdminCategories().subscribe({
      next:  cats => { this.categories.set(cats); this.loading.set(false); },
      error: ()   => { this.error.set('Erreur lors du chargement.'); this.loading.set(false); },
    });
  }

  // ── Expand/collapse subcategories ──────────────────────────────────────────

  toggleExpand(id: number): void {
    const s = new Set(this.expanded());
    s.has(id) ? s.delete(id) : s.add(id);
    this.expanded.set(s);
  }

  isExpanded(id: number): boolean { return this.expanded().has(id); }

  // ── Featured quick-toggle ──────────────────────────────────────────────────

  toggleFeatured(cat: Category): void {
    const next = !cat.featured;
    this.adminApi.updateCategory(cat.ID, { featured: next }).subscribe({
      next: () => {
        this.categories.update(list =>
          list.map(c => c.ID === cat.ID ? { ...c, featured: next } : c)
        );
      },
      error: () => this.error.set('Impossible de modifier le statut featured.'),
    });
  }

  // ── Active quick-toggle ───────────────────────────────────────────────────

  toggleActive(cat: Category, parentId?: number): void {
    const next = !cat.is_active;
    this.adminApi.updateCategory(cat.ID, { is_active: next }).subscribe({
      next: () => {
        this.categories.update(list => list.map(parent => {
          if (!parentId && parent.ID === cat.ID) return { ...parent, is_active: next };
          if (parentId && parent.ID === parentId) {
            return {
              ...parent,
              children: parent.children?.map(c =>
                c.ID === cat.ID ? { ...c, is_active: next } : c
              ),
            };
          }
          return parent;
        }));
      },
      error: () => this.error.set('Impossible de modifier le statut.'),
    });
  }

  // ── Edit ───────────────────────────────────────────────────────────────────

  openEdit(cat: Category): void {
    this.editingCat.set(cat);
    this.fNameFr    = cat.name_fr;
    this.fNameAr    = cat.name_ar;
    this.fNameEn    = cat.name_en;
    this.fIcon      = cat.icon ?? '';
    this.fSortOrder = cat.sort_order;
    this.fFeatured  = cat.featured ?? false;
    this.fIsActive  = cat.is_active;
    this.editError.set('');
    this.showEdit.set(true);
  }

  closeEdit(): void { this.showEdit.set(false); }

  saveEdit(): void {
    const cat = this.editingCat();
    if (!cat) return;
    if (!this.fNameFr.trim() || !this.fNameAr.trim()) {
      this.editError.set('Les noms FR et AR sont requis.');
      return;
    }
    const payload: Partial<CategoryPayload> = {
      name_fr:    this.fNameFr.trim(),
      name_ar:    this.fNameAr.trim(),
      name_en:    this.fNameEn.trim(),
      icon:       this.fIcon.trim(),
      sort_order: this.fSortOrder,
      featured:   this.fFeatured,
      is_active:  this.fIsActive,
    };
    this.saving.set(true);
    this.adminApi.updateCategory(cat.ID, payload).subscribe({
      next: () => { this.saving.set(false); this.showEdit.set(false); this.load(); },
      error: () => { this.saving.set(false); this.editError.set('Erreur lors de la sauvegarde.'); },
    });
  }

  // ── Create ─────────────────────────────────────────────────────────────────

  openCreate(parent: Category | null = null): void {
    this.createParent.set(parent);
    this.cSlug = ''; this.cNameFr = ''; this.cNameAr = ''; this.cNameEn = '';
    this.cIcon = ''; this.cSortOrder = 0; this.cFeatured = false; this.cIsActive = true;
    this.createError.set('');
    this.showCreate.set(true);
  }

  closeCreate(): void { this.showCreate.set(false); }

  slugifyCreate(): void {
    this.cSlug = this.cNameFr
      .toLowerCase()
      .normalize('NFD').replace(/[\u0300-\u036f]/g, '')
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '');
  }

  saveCreate(): void {
    if (!this.cSlug.trim() || !this.cNameFr.trim() || !this.cNameAr.trim()) {
      this.createError.set('Slug, nom FR et nom AR sont requis.');
      return;
    }
    const parent = this.createParent();
    const payload: CategoryCreatePayload = {
      slug:       this.cSlug.trim(),
      name_fr:    this.cNameFr.trim(),
      name_ar:    this.cNameAr.trim(),
      name_en:    this.cNameEn.trim(),
      icon:       this.cIcon.trim(),
      sort_order: this.cSortOrder,
      featured:   parent ? false : this.cFeatured,
      is_active:  this.cIsActive,
      parent_id:  parent?.ID,
    };
    this.creating.set(true);
    this.adminApi.createCategory(payload).subscribe({
      next: () => { this.creating.set(false); this.showCreate.set(false); this.load(); },
      error: () => { this.creating.set(false); this.createError.set('Erreur lors de la création (slug déjà utilisé ?).'); },
    });
  }

  // ── Delete ─────────────────────────────────────────────────────────────────

  deleteCategory(cat: Category): void {
    const msg = cat.children?.length
      ? `Supprimer « ${cat.name_fr} » et ses ${cat.children.length} sous-catégorie(s) ?`
      : `Supprimer « ${cat.name_fr} » ?`;
    if (!confirm(msg)) return;
    this.adminApi.deleteCategory(cat.ID).subscribe({
      next:  () => this.load(),
      error: () => this.error.set('Impossible de supprimer cette catégorie.'),
    });
  }
}
