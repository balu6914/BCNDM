import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { NgProgressModule, NgProgressInterceptor } from 'ngx-progressbar';
import { NgPipesModule } from 'ngx-pipes';

import { AppRoutingModule } from './app.routes';
import { AppBootstrapModule } from './app-bootstrap/app-bootstrap.module'
import { AppComponent } from './app.component';
import { NoContentComponent } from './no-content';

// Import our modules
import { CommonAppModule } from './common/common.module';
import { SharedModule } from './shared/shared.module';
import { LayoutModule } from './layout'
import { AuthModule } from './auth/auth.module';
import { DashboardModule } from './dashboard/dashboard.module';

@NgModule({
    imports: [
        NgProgressModule,
        BrowserModule,
        BrowserAnimationsModule,
        NgPipesModule,
        // App modules
        AppBootstrapModule,
        AuthModule,
        CommonAppModule,
        SharedModule,
        LayoutModule,
        DashboardModule,
        AppRoutingModule,
    ],
    declarations: [
        AppComponent,
        NoContentComponent,
    ],
    providers: [],
    bootstrap: [AppComponent]
})
export class AppModule { }
