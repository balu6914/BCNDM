import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { UiSwitchModule } from 'ngx-ui-switch';
import { ModalModule } from 'ngx-bootstrap/modal';

// Sell routes
import { DashboardSellRoutingModule } from './dashboard.sell.routes';
// Sell components
import { DashboardSellComponent } from './main/dashboard.sell.component';
import { DashboardSellAddComponent } from './add/dashboard.sell.add.component';
import { DashboardSellEditComponent } from './edit/dashboard.sell.edit.component';
import { DashboardSellDeleteComponent } from './delete/dashboard.sell.delete.component';

import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';

@NgModule({
  imports: [
    CommonModule,
    UiSwitchModule.forRoot({
      size: 'small',
      color: '#007bff',
    }),
    ModalModule.forRoot(),
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
  ],
})
export class DashboardSellModule { }
