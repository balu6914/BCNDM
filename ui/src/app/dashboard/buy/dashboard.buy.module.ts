import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ModalModule } from 'ngx-bootstrap/modal';

// Buy components
import { DashboardBuyComponent } from './main/dashboard.buy.component';
import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';
// Import add subscription componen
import { DashboardBuyAddComponent } from './add/dashboard.buy.add.component';
// Buy routes
import { DashboardBuyRoutingModule } from './dashboard.buy.routes';
import { DashboardBuyGroupComponent } from './group/dashboard.buy.group.component';

@NgModule({
  imports: [
    CommonModule,
    ModalModule.forRoot(),
    FormsModule,
    ReactiveFormsModule,
    CommonAppModule,
    SharedModule,
    DashboardBuyRoutingModule,
  ],
  declarations: [
    DashboardBuyComponent,
    DashboardBuyAddComponent,
    DashboardBuyGroupComponent,
  ],
  entryComponents: [
    DashboardBuyAddComponent,
    DashboardBuyGroupComponent,
  ],
})
export class DashboardBuyModule { }
