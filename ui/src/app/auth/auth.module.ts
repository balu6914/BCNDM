import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule, HTTP_INTERCEPTORS } from "@angular/common/http";
import { NgProgressModule, NgProgressInterceptor } from 'ngx-progressbar';
import {
    FormsModule,
    ReactiveFormsModule
} from '@angular/forms';
// Interceptors
import { TokenInterceptor } from './services/token.http.interceptor.service';
import { UnauthorizedInterceptor } from './services/unauthorized.http.interceptor';
// Auth routes
import { AuthRoutingModule } from './auth.routes';
// Add services
import { AuthService } from './services/auth.service';
import { AuthGuardService } from './guardians/auth.guardian';
import { AuthLoggedinGuardService } from './guardians/auth.loggedin.guardin';
import  { UserService } from './services/user.service';
// Auth components
import { SignupComponent } from './signup';
import { LoginComponent } from './login';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    HttpClientModule,
    AuthRoutingModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  declarations: [
      SignupComponent,
      LoginComponent
  ],
  providers: [
      UserService,
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
