import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Stream } from '../interfaces/stream.interface';
import { Query } from '../interfaces/query.interface';
import { Observable } from 'rxjs/Observable';

// Import RxJs required methods
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { Page } from '../interfaces/page.interface';


@Injectable()
export class StreamService {
  // Resolve HTTP using the constructor
  constructor(private http: HttpClient) { }

  getStream(id: string): Observable<Stream[]> {
    return this.http.get(`${environment.API_STREAMS}/${id}`)
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  addStream(data) {
    return this.http.post(`${environment.API_STREAMS}`, data)
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  addStreamBulk(csv) {
    return this.http.post(`${environment.API_STREAMS}/bulk`, csv)
      .map((res: Response) => { })
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  removeStream(id: string) {
    return this.http.delete(`${environment.API_STREAMS}/${id}`)
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  updateStream(id: string, data) {
    return this.http.put(`${environment.API_STREAMS}/${id}`, data)
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  searchStreams(q: Query): Observable<Page<Stream>> {
    return this.http.get(`${environment.API_STREAMS}`, {
      params: q.generateQuery()
    })
      .map((res: Response) => res)
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }
}
