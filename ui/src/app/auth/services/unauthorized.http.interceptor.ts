import { Injectable, Injector } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent, HttpErrorResponse } from '@angular/common/http'
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/do';

import { AuthService } from './auth.service';

@Injectable()
export class UnauthorizedInterceptor implements HttpInterceptor {

  constructor(private inj: Injector) {}

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        const auth = this.inj.get(AuthService);
        return next.handle(req).catch(event =>  {
            if (event instanceof HttpErrorResponse && event.status == 401) {
                console.error("403 Forbiden!")
                // handle 401 errors
                auth.logout()
            } else {
                return Observable.throw(event);
            }
        });
    }
}
