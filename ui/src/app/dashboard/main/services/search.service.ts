import { Injectable }     from '@angular/core';
import { HttpClient, HttpResponse, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../../environments/environment';
import { Stream } from '../../../common/interfaces/stream.interface'
import { Observable } from 'rxjs/Rx';

// Import RxJs required methods
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';


@Injectable()
export class SearchService {
     // Resolve HTTP using the constructor
     constructor (private http: HttpClient) {}

     // serach streams
     searchStreams(type, x0, y0, x1, y1, x2, y2, x3 ,y3): Observable<Stream[]> {
        return this.http.get(`${environment.API_URL}/streams/search?type=${type}&x0=${x0}&y0=${y0}&x1=${x1}&y1=${y1}&x2=${x2}&y2=${y2}&x3=${x3}&y3=${y3}`)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error || 'Server error'));
    }
}
