import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Buy components
import { DashboardBuyComponent } from '.';
import { AppBootstrapModule } from '../../app-bootstrap/app-bootstrap.module';
import { CommonAppModule } from '../../common/common.module';
import { SharedModule } from '../../shared/shared.module';
// Import add subscription component
import { DashboardBuyAddComponent } from './add';
// Buy routes
import { DashboardBuyRoutingModule } from './dashboard.buy.routes';


@NgModule({
  imports: [
    CommonModule,
    AppBootstrapModule,
    FormsModule,
    ReactiveFormsModule,
    CommonAppModule,
    DashboardBuyRoutingModule,
    SharedModule,
  ],
  declarations: [
    DashboardBuyComponent,
    DashboardBuyAddComponent,
  ],
  entryComponents: [
    DashboardBuyAddComponent,
  ],
})
export class DashboardBuyModule { }
