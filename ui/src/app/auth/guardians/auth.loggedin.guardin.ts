import { Injectable } from '@angular/core';
import { Router, CanActivate } from '@angular/router';
import { AuthService } from 'app/auth/services/auth.service';

@Injectable()
export class AuthLoggedinGuardService implements CanActivate {
  constructor(public auth: AuthService, public router: Router) {}

  canActivate(): boolean {
    if (this.auth.isLoggedin()) {
      this.router.navigate(['dashboard']);
      return false;
    }
    return true;
  }
}
