// Imports
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { User } from '../../common/interfaces/user.interface';

@Injectable()
export class UserService {
  constructor(
    private http: HttpClient
  ) { }

  getUser(): Observable<User> {
    return this.http.get<User>(`${environment.API_AUTH}`);
  }

  addUser(data): Observable<User> {
    return this.http.post<User>(`${environment.API_AUTH}`, data);
  }

  removeUser(id: string) {
    return this.http.delete(`${environment.API_AUTH}`);
  }
  updateUser(id: string, data): Observable<User> {
    return this.http.put<User>(`${environment.API_AUTH}`, data);
  }

}
