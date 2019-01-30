import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { AppBootstrapModule } from 'app/app-bootstrap/app-bootstrap.module';
// Add wallet routes
import { DashboardWalletRoutingModule } from './dashboard.wallet.routes';
// Add components
import { DashboardWalletComponent } from './dashboard.wallet.component';
import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    AppBootstrapModule,
    CommonAppModule,
    SharedModule,
    DashboardWalletRoutingModule,
  ],
  declarations: [
    DashboardWalletComponent,
  ],
  schemas: [
    CUSTOM_ELEMENTS_SCHEMA
  ]
})
export class DashboardWalletModule { }
