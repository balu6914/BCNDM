import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Buy routes
import { DashboardBuyRoutingModule } from './dashboard.buy.routes';
// Buy components
import { DashboardBuyComponent } from './index';
// Import subscription coimponent
import { SubscriptionAddComponent } from '../../dashboard/buy/add';

import { CommonAppModule } from '../../common/common.module';
import { SharedModule } from '../../shared/shared.module';
import { AppBootstrapModule } from '../../app-bootstrap/app-bootstrap.module';

@NgModule({
  imports: [
    CommonModule,
    AppBootstrapModule,
    MdlModule,
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    CommonAppModule,
    DashboardBuyRoutingModule,
    SharedModule,
  ],
  declarations: [
    DashboardBuyComponent,
    SubscriptionAddComponent,
  ],
  entryComponents: [
    SubscriptionAddComponent,
  ],
})
export class DashboardBuyModule { }
