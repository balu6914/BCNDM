import {Injectable, EventEmitter, Output} from '@angular/core';
import {Observable, throwError, of} from 'rxjs';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';

import { environment } from 'environments/environment';
import { UserService } from 'app/common/services/user.service';
import * as jwt_decode from 'jwt-decode';

@Injectable()
export class AuthService {
  token: string;
  user: any;

  @Output() loggedIn: EventEmitter<Boolean> = new EventEmitter();
  constructor(
    private http: HttpClient,
    private router: Router,
    private userService: UserService
  ) {
    this.token = localStorage.getItem('token');
  }

  // Get user token
  login(username: string, password: string) {
    return this.http.post(`${environment.API_AUTH_TOKENS}`, JSON.stringify({
        email: username,
        password: password
      }), {
        headers: new HttpHeaders({
          'Content-Type': 'application/json'
        })
      })
      .map((res: any) => {
        const data = res;
        this.token = data.token;
        localStorage.setItem('token', this.token);
        this.fetchCurrentUser();
        this.loggedIn.emit(true);
    })
    .catch((error: any) => {
      return throwError(error);
    });
  }

  // Logout user, remove token from local storage
  logout() {
    localStorage.removeItem('token');
    this.loggedIn.emit(false);
    this.user = null;
    this.router.navigate(['login']);
  }

  // Check if user is logged in
  isLoggedin() {
    return !!localStorage.getItem('token');
  }

  setCurrentUser(data) {
    if (data) {
      this.user = data;
      this.loggedIn.emit(true);
    }
  }

  getCurrentUser() {
    if (this.user) {
      return of(this.user);
    } else {
      if (this.isLoggedin()) {
        return this.fetchCurrentUser();
      }
    }
  }

  fetchCurrentUser() {
    const jwtDecode = jwt_decode(this.token);
    const userID = jwtDecode.sub;

    return this.userService.getUser(userID).map(
      (data: any) =>  {
        this.setCurrentUser(data);
        return data;
      }
    );
  }

  getUserToken(): string {
    return localStorage.getItem('token');
  }
}
