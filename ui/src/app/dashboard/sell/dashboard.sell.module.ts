import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Sell routes
import { DashboardSellRoutingModule } from './dashboard.sell.routes';
// Sell components
import { DashboardSellComponent } from './index';
import { DashboardSellAddComponent } from './add';
import { DashboardSellEditComponent } from './edit';
import { DashboardSellMapComponent } from './map';

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
    DashboardSellRoutingModule,
    SharedModule,
  ],
  declarations: [
      DashboardSellComponent,
      DashboardSellAddComponent,
      DashboardSellEditComponent,
      DashboardSellMapComponent,
  ],
})
export class DashboardSellModule { }
