// Imports
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'environments/environment';
import { User } from 'app/common/interfaces/user.interface';

@Injectable()
export class UserService {
  constructor(
    private http: HttpClient
  ) { }

  getUser(userID: string): Observable<User> {
    return this.http.get<User>(`${environment.API_AUTH}/${userID}`);
  }

  getAllUsers(): Observable<User[]> {
    return this.http.get<User[]>(`${environment.API_AUTH}`);
  }

  getNonPartners(): Observable<User[]> {
      return this.http.get<User[]>(`${environment.API_AUTH}/non-partners`);
  }

  addUser(data): Observable<User> {
    return this.http.post<User>(`${environment.API_AUTH}`, data);
  }

  removeUser(id: string) {
    return this.http.delete(`${environment.API_AUTH}`);
  }

  updateUser(user: User): Observable<User> {
    return this.http.patch<User>(`${environment.API_AUTH}/${user.id}`, user);
  }
}
