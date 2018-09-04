import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from '../../../environments/environment';
import { Page } from '../interfaces/page.interface';
import { Subscription } from '../interfaces/subscription.interface';

@Injectable()
export class SubscriptionService {
  constructor(
    private http: HttpClient
  ) { }

  get(): Observable<Subscription[]> {
    return this.http.get<Subscription[]>(`${environment.API_SUBSCRIPTIONS}`);
  }

  owned(page: number, limit: number): Observable<Page<Subscription>> {
    return this.http.get<Page<Subscription>>(`${environment.API_SUBSCRIPTIONS}/owned`, {
      params: new HttpParams()
        .set('page', page.toString())
        .set('limit', limit.toString())
    });
  }

  bought(page: number, limit: number): Observable<Page<Subscription>> {
    return this.http.get<Page<Subscription>>(`${environment.API_SUBSCRIPTIONS}/bought`, {
      params: new HttpParams()
        .set('page', page.toString())
        .set('limit', limit.toString())
    });
  }

  add(data) {
    return this.http.post(`${environment.API_SUBSCRIPTIONS}`, data);
  }

  update(id: string, data): Observable<Subscription> {
    return this.http.put<Subscription>(`${environment.API_SUBSCRIPTIONS}/${id}`, data);
  }

}
