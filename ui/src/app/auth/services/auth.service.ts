import {Injectable, EventEmitter, Output} from '@angular/core';
import {Observable, throwError, of} from 'rxjs';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';

import { environment } from 'environments/environment';
import { UserService } from 'app/common/services/user.service';
import * as jwt_decode from 'jwt-decode';

const COOKIE_NAME: string = "token";
const COOKIE_EXPIRY_IN_HOURS: number = 2;

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
		this.token = this.getCookie(COOKIE_NAME);
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
				this.setCookie(COOKIE_NAME, this.token, COOKIE_EXPIRY_IN_HOURS)
        this.fetchCurrentUser();
        this.loggedIn.emit(true);
    })
    .catch((error: any) => {
      return throwError(error);
    });
  }

  // Logout user, remove token from local storage
  logout() {
		this.deleteCookie(COOKIE_NAME)
    this.loggedIn.emit(false);
    this.user = null;
    this.router.navigate(['login']);
  }

  // Check if user is logged in
  isLoggedin() {
		return !!this.getCookie(COOKIE_NAME);
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
		return this.getCookie('token');
  }

	private getCookie(name: string) {
		let ca: Array<string> = document.cookie.split(';');
		let caLen: number = ca.length;
		let cookieName = `${name}=`;
		let c: string;

		for (let i: number = 0; i < caLen; i += 1) {
			c = ca[i].replace(/^\s+/g, '');
			if (c.indexOf(cookieName) == 0) {
				return c.substring(cookieName.length, c.length);
			}
		}
		return '';
	}
	private deleteCookie(name) {
		this.setCookie(name, '', -1);
	}

	private setCookie(name: string, value: string, expireDays: number, path: string = '') {
		let d: Date = new Date();
		d.setTime(d.getTime() + expireDays * 60 * 60 * 1000);
		console.log(d);
		let expires: string = `expires=${d.toString()}`;
		let cpath: string = path ? `; path=${path}` : '';
		document.cookie = `${name}=${value}; ${expires}${cpath}`;
	}

}
