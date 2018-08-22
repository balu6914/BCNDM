import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Balance } from '../../common/interfaces/balance.interface'
import { Observable } from 'rxjs/Rx';
import { Subject } from 'rxjs/Subject';

// Import RxJs required methods
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Injectable()
export class BalanceService {
    // Resolve HTTP using the constructor
    private _balance = new Subject<Balance>();
    balance = this._balance.asObservable();

    constructor (
    private http: HttpClient,
  ) {}

    get() : Observable<Balance[]> {
        return this.http.get(`${environment.API_TOKENS}`)
                        .map((res:Response) => res)
                        .catch((error:any) => Observable.throw(error || 'Server error'));
    }
    buy(data: any) {
        return this.http.post(`${environment.API_TOKENS}/buy`, data)
                        .map((res:Response) => res)
                        .catch((error:any) => Observable.throw(error || 'Server error'));
    }
    // Balance Message Buss will brodcast notification about balance value changes
    changed(value: Balance) {
      this._balance.next(value);
    }
}
