import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { NgPipesModule } from 'ngx-pipes';
import { MdlSelectModule } from '@angular2-mdl-ext/select';
import { MdlDatePickerModule } from '@angular-mdl/datepicker';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Sell routes
import { DashboardContractsRoutingModule } from './dashboard.contracts.routes';
// Components
import { DashboardContractsListComponent } from './list/dashboard.contracts.list.component';
import { DashboardContractsEmptyListComponent } from './empty-list/dashboard.contracts.empty.list.component';
import { DashboardContractsDetailsComponent } from './contract-details/dashboard.contracts.details.component'
import { DashboardContractsAddComponent } from './add/dashboard.contracts.add.component';
import { CommonAppModule } from '../../common/common.module';

@NgModule({
  imports: [
    MdlModule,
    MdlSelectModule,
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    CommonModule,
    CommonAppModule,
    DashboardContractsRoutingModule
  ],
  declarations: [
      DashboardContractsListComponent,
      DashboardContractsEmptyListComponent,
      DashboardContractsDetailsComponent,
      DashboardContractsAddComponent
  ],
  providers: [
      // StreamService,
  ]
})
export class DashboardContractsModule { }
