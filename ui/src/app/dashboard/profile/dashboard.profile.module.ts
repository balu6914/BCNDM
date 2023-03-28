import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ClipboardModule } from 'ngx-clipboard';

// Add  routes
import { DashboardProfileRoutingModule } from './dashboard.profile.routes';
// Add components
import { DashboardProfileComponent } from './dashboard.profile.component';
import { DashboardProfilePasswordUpdateComponent  } from './password/dashboard.profile.change.password.component';
import { DashboardProfileUpdateComponent } from './update/dashboard.profile.update.component';
import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
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
