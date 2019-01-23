import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Sell routes
import { DashboardContractsRoutingModule } from './dashboard.contracts.routes';
// Components
import { DashboardContractsComponent } from './dashboard.contracts.component';
import { DashboardContractsEmptyListComponent } from './empty-list/dashboard.contracts.empty.list.component';
import { DashboardContractsDetailsComponent } from './contract-details/dashboard.contracts.details.component';
import { DashboardContractsAddComponent } from './add/dashboard.contracts.add.component';
import { DashboardContractsSignComponent } from './sign/dashboard.contracts.sign.component';
import { DashboardContractsHelpComponent } from './help/dashboard.contracts.help.component';
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
    DashboardContractsRoutingModule,
    SharedModule,
  ],
  declarations: [
    DashboardContractsComponent,
    DashboardContractsAddComponent,
    DashboardContractsEmptyListComponent,
    DashboardContractsDetailsComponent,
    DashboardContractsHelpComponent,
    DashboardContractsSignComponent,
  ],
  entryComponents: [
    DashboardContractsAddComponent,
    DashboardContractsSignComponent,
  ]
})
export class DashboardContractsModule { }
