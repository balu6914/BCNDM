import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from 'environments/environment';
import { Access } from 'app/common/interfaces/access.interface';
import { Page } from 'app/common/interfaces/page.interface';
import { Query } from 'app/common/interfaces/query.interface';

@Injectable({
  providedIn: 'root'
})
export class AccessService {
  constructor(
    private http: HttpClient,
  ) { }

  requestAccess(data: any) {
    return this.http.post(`${environment.API_ACCESS_CONTROL}`, data);
  }

  getAccessSent(state: string) {
    let params = new HttpParams();
    params = params.set('state', state);
    return this.http.get(`${environment.API_ACCESS_CONTROL}/sent`, {
      params: params
    });
  }

  getAccessReceived(state: string) {
    let params = new HttpParams();
    params = params.set('state', state);
    return this.http.get(`${environment.API_ACCESS_CONTROL}/received`, {
      params: params
    });
  }

  approveAccessRequest(requestID: string) {
    return this.http.put(`${environment.API_ACCESS_CONTROL}/${requestID}/approve`, {});
  }

  revokeAccessRequest(requestID: string) {
    return this.http.put(`${environment.API_ACCESS_CONTROL}/${requestID}/revoke`, {});
  }

}
