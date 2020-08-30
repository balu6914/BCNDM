import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from 'environments/environment';
import { Page } from 'app/common/interfaces/page.interface';
import { Subscription } from 'app/common/interfaces/subscription.interface';

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

 report(page: number, limit: number, reqType: string): Observable<any> {
   var options = {
     headers: new HttpHeaders({
        'Accept':'application/pdf'
     }),
      params: new HttpParams()
        .set('page', page.toString())
        .set('limit', limit.toString())
        .set('type', reqType),
     'responseType': 'blob' as 'json'
  }
    return this.http.get<any>(`${environment.API_SUBSCRIPTIONS}/report`, options);
  }

  add(data) {
    return this.http.post(`${environment.API_SUBSCRIPTIONS}`, data);
  }

  update(id: string, data): Observable<Subscription> {
    return this.http.put<Subscription>(`${environment.API_SUBSCRIPTIONS}/${id}`, data);
  }

}
