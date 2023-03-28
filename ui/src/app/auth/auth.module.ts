import { CommonModule } from '@angular/common';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgProgressInterceptor } from 'ngx-progressbar';

import { CommonAppModule } from 'app/common/common.module';

// Auth Guardians
import { AuthGuardService } from './guardians/auth.guardian';
import { AuthLoggedinGuardService } from './guardians/auth.loggedin.guardian';
// Auth Services
import { AuthService } from './services/auth.service';
// Auth Interceptors
import { TokenInterceptor } from './services/token.interceptor.service';
import { UnauthorizedInterceptor } from './services/unauthorized.interceptor.service';

@NgModule({
  providers: [
    AuthService,
    AuthGuardService,
    AuthLoggedinGuardService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: NgProgressInterceptor,
      multi: true
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: TokenInterceptor,
      multi: true
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: UnauthorizedInterceptor,
      multi: true
    },
  ],
})
export class AuthModule { }
