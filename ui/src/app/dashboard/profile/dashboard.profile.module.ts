import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { AppBootstrapModule } from '../../app-bootstrap/app-bootstrap.module';
// Add  routes
import { DashboardProfileRoutingModule } from './dashboard.profile.routes';
// Add components
import { DashboardProfileComponent } from './dashboard.profile.component';
import { DashboardProfilePasswordUpdateComponent  } from './password/dashboard.profile.change.password.component';
import { DashboardProfileUpdateComponent } from './update/dashboard.profile.update.component';
import { CommonAppModule } from '../../common/common.module';
import { SharedModule } from '../../shared/shared.module';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    AppBootstrapModule,
    CommonAppModule,
    SharedModule,
    DashboardProfileRoutingModule,
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
