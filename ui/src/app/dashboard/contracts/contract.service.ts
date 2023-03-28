import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from 'environments/environment';
import { Contract } from 'app/common/interfaces/contract.interface';
import { Page } from 'app/common/interfaces/page.interface';
import { Query } from 'app/common/interfaces/query.interface';

@Injectable({
  providedIn: 'root'
})
export class ContractService {
  constructor(
    private http: HttpClient,
  ) { }

  create(data: any) {
    return this.http.post(`${environment.API_CONTRACTS}`, data);
  }

  sign(data: any) {
    return this.http.patch(`${environment.API_CONTRACTS}/sign`, data);
  }

  get(query: any): Observable<Page<Contract>>  {
    let params = new HttpParams();
    params = params.set('owner', query.isOwner.toString());
    params = params.set('partner', query.isPartner.toString());
    params = params.set('page', query.page.toString());

    return this.http.get<Page<Contract>>(`${environment.API_CONTRACTS}`, {
      params: params
    });
  }

}
