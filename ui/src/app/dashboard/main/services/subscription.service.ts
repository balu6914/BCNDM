import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../../environments/environment';
import { Stream } from '../../../common/interfaces/stream.interface'
import { Subscription } from '../../../common/interfaces/subscription.interface'
import { Observable } from 'rxjs/Rx';

// Import RxJs required methods
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Injectable()
export class SubscriptionService {
     // Resolve HTTP using the constructor
     constructor (private http: HttpClient) {}

     // Fetch  stream
     get () : Observable<Subscription[]> {
         return this.http.get(`${environment.API_URL}/streams/purch`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
     }
     add (data): Observable<Stream[]> {
        return this.http.post(`${environment.API_URL}/streams/purch`, data)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
     remove (id:string): Observable<Stream[]> {
        return this.http.delete(`${environment.API_URL}/streams/purch/${id}`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
     update (id:string, data): Observable<Stream> {
        return this.http.put(`${environment.API_URL}/streams/purch/${id}`, data)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }

}
