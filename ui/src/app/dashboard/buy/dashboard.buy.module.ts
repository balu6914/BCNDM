import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Buy routes
import { DashboardBuyRoutingModule } from './dashboard.buy.routes';
// Buy components
import { DashboardBuyComponent } from './index';
import { DashboardBuyMapComponent } from './map';
// Import subscription module
import { SubscriptionModule } from '../subscription';

import { CommonAppModule } from '../../common/common.module';
import { SharedModule } from '../../shared/shared.module';


@NgModule({
  imports: [
    CommonModule,
    MdlModule,
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    CommonAppModule,
    SubscriptionModule,
    DashboardBuyRoutingModule,
    SharedModule,
  ],
  declarations: [
      DashboardBuyComponent,
      DashboardBuyMapComponent,
  ],
})
export class DashboardBuyModule { }
