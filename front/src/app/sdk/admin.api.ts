import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ApiClient } from './api-client';
import { Category, Permission, Role, User, UserListResponse } from './types';

export interface CategoryPayload {
  name_fr:    string;
  name_ar:    string;
  name_en?:   string;
  icon?:      string;
  sort_order?: number;
  featured?:  boolean;
  is_active?: boolean;
}

export interface CategoryCreatePayload extends CategoryPayload {
  slug:      string;
  parent_id?: number;
}

export interface RolePayload {
  name:           string;
  slug:           string;
  description:    string;
  permission_ids: number[];
}

/**
 * Thin wrapper around all /admin/** endpoints.
 * Requires the authenticated user to have role=admin.
 */
@Injectable({ providedIn: 'root' })
export class AdminApi {
  private readonly api = inject(ApiClient);

  // ── Permissions ────────────────────────────────────────────────────────────

  listPermissions(): Observable<Permission[]> {
    return this.api.get('/admin/permissions');
  }

  // ── Roles ──────────────────────────────────────────────────────────────────

  listRoles(): Observable<Role[]> {
    return this.api.get('/admin/roles');
  }

  createRole(payload: RolePayload): Observable<Role> {
    return this.api.post('/admin/roles', payload);
  }

  updateRole(id: number, payload: RolePayload): Observable<Role> {
    return this.api.put(`/admin/roles/${id}`, payload);
  }

  deleteRole(id: number): Observable<{ message: string }> {
    return this.api.delete(`/admin/roles/${id}`);
  }

  // ── Users ──────────────────────────────────────────────────────────────────

  listUsers(page = 1, pageSize = 20, userType?: 'external' | 'internal'): Observable<UserListResponse> {
    const type = userType ? `&user_type=${userType}` : '';
    return this.api.get(`/admin/users?page=${page}&page_size=${pageSize}${type}`);
  }

  setUserRoles(userId: number, roleIds: number[]): Observable<{ message: string }> {
    return this.api.put(`/admin/users/${userId}/roles`, { role_ids: roleIds });
  }

  updateUser(userId: number, payload: {
    display_name: string; phone: string; email?: string | null;
    address?: string | null; city?: string | null;
    postal_code?: string | null; country?: string | null;
    is_verified: boolean; role: string;
  }): Observable<User> {
    return this.api.put(`/admin/users/${userId}`, payload);
  }

  banUser(userId: number): Observable<{ message: string }> {
    return this.api.put(`/admin/users/${userId}/ban`, {});
  }

  unbanUser(userId: number): Observable<{ message: string }> {
    return this.api.put(`/admin/users/${userId}/unban`, {});
  }

  deleteUser(userId: number): Observable<{ message: string }> {
    return this.api.delete(`/admin/users/${userId}`);
  }

  resetUserPassword(userId: number, newPassword: string): Observable<{ message: string }> {
    return this.api.put(`/admin/users/${userId}/reset-password`, { new_password: newPassword });
  }

  // ── Categories ─────────────────────────────────────────────────────────────

  listAdminCategories(): Observable<Category[]> {
    return this.api.get<{ data: Category[] }>('/admin/categories').pipe(map(r => r.data));
  }

  createCategory(payload: CategoryCreatePayload): Observable<Category> {
    return this.api.post('/admin/categories', payload);
  }

  updateCategory(id: number, payload: Partial<CategoryPayload>): Observable<{ ok: boolean }> {
    return this.api.put(`/admin/categories/${id}`, payload);
  }

  deleteCategory(id: number): Observable<{ ok: boolean }> {
    return this.api.delete(`/admin/categories/${id}`);
  }
}
