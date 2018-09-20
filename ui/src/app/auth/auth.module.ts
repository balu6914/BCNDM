import { CommonModule } from '@angular/common';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgProgressInterceptor } from 'ngx-progressbar';

import { CommonAppModule } from 'app/common/common.module';
// Auth routes
import { AuthRoutingModule } from './auth.routes';
import { AuthGuardService } from './guardians/auth.guardian';
import { AuthLoggedinGuardService } from './guardians/auth.loggedin.guardin';
import { LoginComponent } from './login';
// Add services
import { AuthService } from './services/auth.service';
// Interceptors
import { TokenInterceptor } from './services/token.http.interceptor.service';
import { UnauthorizedInterceptor } from './services/unauthorized.http.interceptor';
// Auth components
import { SignupComponent } from './signup';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    HttpClientModule,
    AuthRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    CommonAppModule,
  ],
  declarations: [
    SignupComponent,
    LoginComponent
  ],
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
