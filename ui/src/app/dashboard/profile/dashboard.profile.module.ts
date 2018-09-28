import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { AppBootstrapModule } from 'app/app-bootstrap/app-bootstrap.module';
// Add  routes
import { DashboardProfileRoutingModule } from './dashboard.profile.routes';
// Add components
import { DashboardProfileComponent } from './dashboard.profile.component';
import { DashboardProfilePasswordUpdateComponent  } from './password/dashboard.profile.change.password.component';
import { DashboardProfileUpdateComponent } from './update/dashboard.profile.update.component';
import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';
import { ClipboardModule } from 'ngx-clipboard';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    AppBootstrapModule,
    CommonAppModule,
    SharedModule,
    DashboardProfileRoutingModule,
    ClipboardModule,
  ],
  declarations: [
    DashboardProfileComponent,
    DashboardProfileUpdateComponent,
    DashboardProfilePasswordUpdateComponent,
  ],
  schemas: [
  ]
})
export class DashboardProfileModule { }
