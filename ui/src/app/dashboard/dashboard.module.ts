import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ArchwizardModule } from 'ng2-archwizard';
import { ClipboardModule } from 'ngx-clipboard';

// Dashboard components
import { DashboardComponent } from './dashboard.component';
import { DashboardMainComponent } from './main/dashboard.main.component';
import { DashboardMainStreamsComponent } from './main/section-streams/dashboard.main.streams.component';
import { DashboardRoutingModule } from './dashboard.routes';
import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';

import { SignupComponent } from './signup/signup.component';
import { LoginComponent } from './login/login.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    ArchwizardModule,
    // App module
    CommonAppModule,
    SharedModule,
    DashboardRoutingModule,
    ClipboardModule,
  ],
  declarations: [
    SignupComponent,
    LoginComponent,
    DashboardComponent,
    DashboardMainComponent,
    DashboardMainStreamsComponent,
  ],
 schemas: [ CUSTOM_ELEMENTS_SCHEMA ]
})
export class DashboardModule { }
