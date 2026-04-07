import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiClient } from './api-client';
import { Permission, Role, User, UserListResponse } from './types';

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

  listUsers(page = 1, pageSize = 20): Observable<UserListResponse> {
    return this.api.get(`/admin/users?page=${page}&page_size=${pageSize}`);
  }

  setUserRoles(userId: number, roleIds: number[]): Observable<{ message: string }> {
    return this.api.put(`/admin/users/${userId}/roles`, { role_ids: roleIds });
  }
}
