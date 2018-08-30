import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Sell routes
import { DashboardSellRoutingModule } from './dashboard.sell.routes';
// Sell components
import { DashboardSellComponent } from './index';
import { DashboardSellAddComponent } from './add';
import { DashboardSellEditComponent } from './edit';
import { DashboardSellDeleteComponent } from './delete';

import { CommonAppModule } from '../../common/common.module';
import { SharedModule } from '../../shared/shared.module';
import { AppBootstrapModule } from '../../app-bootstrap/app-bootstrap.module';

@NgModule({
  imports: [
    CommonModule,
    AppBootstrapModule,
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    CommonAppModule,
    DashboardSellRoutingModule,
    SharedModule,
  ],
  declarations: [
    DashboardSellComponent,
    DashboardSellAddComponent,
    DashboardSellEditComponent,
    DashboardSellDeleteComponent,
  ],
  entryComponents: [
    DashboardSellAddComponent,
    DashboardSellEditComponent,
    DashboardSellDeleteComponent,
  ]
})
export class DashboardSellModule { }
