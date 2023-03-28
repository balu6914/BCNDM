import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Subject } from 'rxjs/Subject';
import { environment } from 'environments/environment';
import { Balance } from 'app/common/interfaces/balance.interface';

@Injectable({
  providedIn: 'root'
})
export class BalanceService {
  private _balance = new Subject<Balance>();
  balance = this._balance.asObservable();

  constructor(
    private http: HttpClient,
  ) { }

  get(): Observable<Balance> {
    return this.http.get<Balance>(`${environment.API_TOKENS}`);
  }

  getBalance(userID: string) {
    return this.http.get<Balance>(`${environment.API_TOKENS}/${userID}`);
  }

  buy(data: any) {
    return this.http.post(`${environment.API_TOKENS}/buy`, data);
  }
  withdraw(data: any) {
    return this.http.post(`${environment.API_TOKENS}/withdraw`, data);
  }

  // Balance Message Buss will broadcast notification about balance value changes
  changed(value: Balance) {
    this._balance.next(value);
  }
}
