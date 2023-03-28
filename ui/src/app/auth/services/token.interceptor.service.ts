import { Injectable, Injector } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { AuthService } from './auth.service';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class TokenInterceptor implements HttpInterceptor {

  constructor(
    private inj: Injector,
  ) {}

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const auth = this.inj.get(AuthService);
    request = request.clone({
      setHeaders: {
        Authorization: `${auth.getUserToken()}`
      }
    });

    return next.handle(request);
  }
}
