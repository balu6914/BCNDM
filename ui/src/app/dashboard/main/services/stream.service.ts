import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../../environments/environment';
import { Stream } from '../../../common/interfaces/stream.interface'
import { Observable } from 'rxjs/Rx';

// Import RxJs required methods
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';


@Injectable()
export class StreamService {
     // Resolve HTTP using the constructor
     constructor (private http: HttpClient) {}

     // Fetch  stream
     getAll() : Observable<Stream[]> {
         return this.http.get(`${environment.API_URL}/streams`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
     }
     getStream(id:string) : Observable<Stream[]> {
         return this.http.get(`${environment.API_URL}/streams/${id}`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
     }
     addStream (data): Observable<Stream[]> {
        return this.http.post(`${environment.API_URL}/streams`, data)
                         .map((res:Response) => {})
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
     removeStream (id:string): Observable<Stream[]> {
        return this.http.delete(`${environment.API_URL}/streams/${id}`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
     updateStream (id:string, data): Observable<Stream> {
        return this.http.put(`${environment.API_URL}/streams/${id}`, data)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
}
