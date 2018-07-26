import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Stream } from '../interfaces/stream.interface'
import { Observable } from 'rxjs/Rx';

// Import RxJs required methods
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';


@Injectable()
export class StreamService {
     // Resolve HTTP using the constructor
     constructor (private http: HttpClient) {}

     getStream(id:string) : Observable<Stream[]> {
         return this.http.get(`${environment.API_STREAMS}/${id}`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
     }
     addStream (data) {
        return this.http.post(`${environment.API_STREAMS}`, data)
                         .map((res:Response) => {})
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
    addStreamBulk (csv) {
        return this.http.post(`${environment.API_STREAMS}/bulk`, csv)
                        .map((res:Response) => {})
                        .catch((error:any) => Observable.throw(error || 'Server error'));
    }
     removeStream (id:string) {
        return this.http.delete(`${environment.API_STREAMS}/${id}`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
     updateStream (id:string, data) {
        return this.http.put(`${environment.API_STREAMS}/${id}`, data)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
     // serach streams
     searchStreams(type, x0, y0, x1, y1, x2, y2, x3 ,y3): Observable<Stream[]> {
        return this.http.get(`${environment.API_STREAMS}/search?type=${type}&x0=${x0}&y0=${y0}&x1=${x1}&y1=${y1}&x2=${x2}&y2=${y2}&x3=${x3}&y3=${y3}`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }

}
