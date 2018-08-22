import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { NgPipesModule } from 'ngx-pipes';
import { MdlDatePickerModule } from '@angular-mdl/datepicker';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Sell routes
import { DashboardContractsRoutingModule } from './dashboard.contracts.routes';
// Components
import { DashboardContractsComponent } from './dashboard.contracts.component';
import { DashboardContractsEmptyListComponent } from './empty-list/dashboard.contracts.empty.list.component';
import { DashboardContractsDetailsComponent } from './contract-details/dashboard.contracts.details.component'
import { DashboardContractsAddComponent } from './add/dashboard.contracts.add.component';
import { DashboardContractsHelpComponent } from './help/dashboard.contracts.help.component';
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
    DashboardContractsRoutingModule,
    SharedModule,
  ],
  declarations: [
    DashboardContractsComponent,
    DashboardContractsAddComponent,
    DashboardContractsEmptyListComponent,
    DashboardContractsDetailsComponent,
    DashboardContractsHelpComponent,
  ],
  entryComponents: [
    DashboardContractsAddComponent,
  ]
})
export class DashboardContractsModule { }
