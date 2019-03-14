import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'environments/environment';
import { Page } from 'app/common/interfaces/page.interface';
import { Query } from 'app/common/interfaces/query.interface';
import { Execution } from 'app/common/interfaces/execution.interface';

@Injectable()
export class ExecutionsService {
  constructor(
    private http: HttpClient
  ) { }

  startExecution(data) {
    return this.http.post(`${environment.API_EXECUTIONS}`, data);
  }

  getExecutions():  Observable<Page<Execution>> {
    return this.http.get<Page<Execution>>(`${environment.API_EXECUTIONS}`);
  }

  getExecution(id: string): Observable<Execution> {
    return this.http.get<Execution>(`${environment.API_EXECUTIONS}/${id}`);
  }

  getExecutionResult(id: string) {
    return this.http.get(`${environment.API_EXECUTIONS}/${id}/result`);
  }
}
