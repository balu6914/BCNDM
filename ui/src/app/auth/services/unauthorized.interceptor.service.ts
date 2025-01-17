import { Injectable, Injector } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import 'rxjs/add/operator/do';

import { AuthService } from './auth.service';

@Injectable()
export class UnauthorizedInterceptor implements HttpInterceptor {

  constructor(
    private inj: Injector,
  ) {}

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const auth = this.inj.get(AuthService);
    return next.handle(req).catch(
      event =>  {
        if (event instanceof HttpErrorResponse && event.status === 403) {
          console.error('403 Forbiden!');
          // handle 403 errors
          auth.logout();
        }
        return throwError(event);
      }
    );
  }
}
