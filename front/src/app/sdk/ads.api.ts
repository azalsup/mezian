import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ApiClient } from './api-client';
import type { Ad, AdsQuery, AdsResponse, CreateAdPayload } from './types';

@Injectable({ providedIn: 'root' })
export class AdsApi {
  private readonly api = inject(ApiClient);

  getAds(query: AdsQuery = {}): Observable<AdsResponse> {
    const p = new URLSearchParams();
    if (query.q)             p.set('q',         query.q);
    if (query.cat)           p.set('cat',        query.cat);
    if (query.sub)           p.set('sub',        query.sub);
    if (query.city)          p.set('city',       query.city);
    if (query.minPrice)      p.set('min_price',  String(query.minPrice));
    if (query.maxPrice)      p.set('max_price',  String(query.maxPrice));
    if (query.page)          p.set('page',       String(query.page));
    if (query.sort)          p.set('sort',       query.sort);
    const qs = p.toString();
    return this.api.get<AdsResponse>(`/ads${qs ? '?' + qs : ''}`);
  }

  getBySlug(slug: string): Observable<Ad> {
    return this.api.get<{ data: Ad }>(`/ads/${slug}`).pipe(map(r => r.data));
  }

  createAd(payload: CreateAdPayload): Observable<Ad> {
    return this.api.post<{ data: Ad }>('/ads', payload).pipe(map(r => r.data));
  }
}
