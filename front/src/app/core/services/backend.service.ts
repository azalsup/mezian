import { Injectable, inject } from '@angular/core';
import { ApiClient } from '../../sdk';
import { Observable, map, catchError, of } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class BackendService {
  private readonly api = inject(ApiClient);

  checkHealth(): Observable<boolean> {
    return this.api.get('/health').pipe(
      map(() => true),
      catchError(() => of(false))
    );
  }
}
