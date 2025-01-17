import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { NgProgressModule, NgProgressInterceptor } from 'ngx-progressbar';
import { NgPipesModule } from 'ngx-pipes';
import { SidebarModule } from 'ng-sidebar';
import { BsDropdownModule } from 'ngx-bootstrap/dropdown';
import { ModalModule } from 'ngx-bootstrap/modal';

// Import our modules
import { AppRoutingModule } from './app.routes';
import { AppComponent } from './app.component';
import { NoContentComponent } from './no-content/no-content.component';
import { HeaderComponent } from './header/header.component';

import { CommonAppModule } from './common/common.module';
import { SharedModule } from './shared/shared.module';
import { AuthModule } from './auth/auth.module';
import { DashboardModule } from './dashboard/dashboard.module';
import { BalanceComponent } from './shared/balance/balance.component';
import { TermsComponent } from './dashboard/login/terms/terms.component';

@NgModule({
  imports: [
    NgProgressModule,
    BrowserModule,
    BrowserAnimationsModule,
    NgPipesModule,
    SidebarModule.forRoot(),
    BsDropdownModule.forRoot(),
    ModalModule.forRoot(),
    AuthModule,
    CommonAppModule,
    SharedModule,
    DashboardModule,
    AppRoutingModule,
    HttpClientModule,
  ],
  declarations: [
    AppComponent,
    NoContentComponent,
    HeaderComponent,
		TermsComponent
  ],
  exports: [
    BalanceComponent,
  ],
	entryComponents: [
		TermsComponent
	],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
