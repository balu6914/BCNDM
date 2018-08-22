import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import 'rxjs/add/operator/catch';
// Import RxJs required methods
import 'rxjs/add/operator/map';
import { Observable } from 'rxjs/Observable';
import { environment } from '../../../environments/environment';
import { Page } from '../interfaces/page.interface';
import { Subscription } from '../interfaces/subscription.interface';
import { Query } from '../interfaces/query.interface';


@Injectable()
export class SubscriptionService {
  // Resolve HTTP using the constructor
  constructor(private http: HttpClient) { }

  get(): Observable<Subscription[]> {
    return this.http.get(`${environment.API_SUBSCRIPTIONS}`)
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  owned(page: number, limit: number): Observable<Page<Subscription>> {
    return this.http.get(`${environment.API_SUBSCRIPTIONS}/owned`, {
      params: new HttpParams()
        .set('page', page.toString())
        .set('limit', limit.toString())
    })
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  bought(page: number, limit: number): Observable<Page<Subscription>> {
    return this.http.get(`${environment.API_SUBSCRIPTIONS}/bought`, {
      params: new HttpParams()
        .set('page', page.toString())
        .set('limit', limit.toString())
    })
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  add(data) {
    return this.http.post(`${environment.API_SUBSCRIPTIONS}`, data)
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  update(id: string, data): Observable<Subscription> {
    return this.http.put(`${environment.API_SUBSCRIPTIONS}/${id}`, data)
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }
}
