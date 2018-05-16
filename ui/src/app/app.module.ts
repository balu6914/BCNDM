import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { MdlModule } from '@angular-mdl/core';
import { MdlDatePickerModule } from '@angular-mdl/datepicker';
import { NgProgressModule, NgProgressInterceptor } from 'ngx-progressbar';
import {NgPipesModule} from 'ngx-pipes';


import { AppRoutingModule } from './app.routes';
import { AppComponent } from './app.component';
import { NoContentComponent } from './no-content';

import  { TokenInterceptor } from './auth/services/token.http.interceptor.service';
// Import our modules
import { CommonAppModule } from './common/common.module';
import { LayoutModule } from './layout'
import { AuthModule } from './auth/auth.module';
import { DashboardModule } from './dashboard/dashboard.module';
import { WalletModule } from './dashboard/wallet/wallet.module';

@NgModule({
    imports: [
        NgProgressModule,
        BrowserModule,
        NgPipesModule,
        MdlModule,
        MdlDatePickerModule,
        // App modules
        AuthModule,
        AppRoutingModule,
        CommonAppModule,
        LayoutModule,
        DashboardModule,
        WalletModule
    ],
    declarations: [
        AppComponent,
        NoContentComponent,
    ],
    providers: [],
    bootstrap: [AppComponent]
})
export class AppModule { }
