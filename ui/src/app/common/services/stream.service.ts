import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'environments/environment';
import { Page } from '../interfaces/page.interface';
import { Query } from '../interfaces/query.interface';
import { Stream } from '../interfaces/stream.interface';

@Injectable()
export class StreamService {
  constructor(
    private http: HttpClient
  ) { }

  getStream(id: string): Observable<Stream> {
    return this.http.get<Stream>(`${environment.API_STREAMS}/${id}`);
  }

  addStream(data) {
    return this.http.post(`${environment.API_STREAMS}`, data);
  }

  addStreamBulk(csv) {
    return this.http.post(`${environment.API_STREAMS}/bulk`, csv);
  }

  removeStream(id: string) {
    return this.http.delete(`${environment.API_STREAMS}/${id}`);
  }

  updateStream(id: string, data) {
    return this.http.put(`${environment.API_STREAMS}/${id}`, data);
  }

  searchStreams(q: Query): Observable<Page<Stream>> {
    return this.http.get<Page<Stream>>(`${environment.API_STREAMS}`, {
      params: q.generateQuery()
    });
  }

}
