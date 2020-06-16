import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ArchwizardModule } from 'ng2-archwizard';
import { ClipboardModule } from 'ngx-clipboard';

// Auth module
import { AuthModule } from 'app/auth/auth.module';

// Dashboard components
import { DashboardComponent } from './dashboard.component';
import { DashboardMainComponent } from './main/dashboard.main.component';
import { DashboardMainStreamsComponent } from './main/section-streams/dashboard.main.streams.component';
import { DashboardRoutingModule } from './dashboard.routes';
import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    FormsModule,
    ReactiveFormsModule,
    ArchwizardModule,
    // App module
    AuthModule,
    CommonAppModule,
    SharedModule,
    DashboardRoutingModule,
    ClipboardModule,
  ],
  declarations: [
    DashboardComponent,
    DashboardMainComponent,
    DashboardMainStreamsComponent,
  ],
 schemas: [ CUSTOM_ELEMENTS_SCHEMA ]
})
export class DashboardModule { }
