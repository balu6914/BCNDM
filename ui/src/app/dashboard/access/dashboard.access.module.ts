import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ModalModule } from 'ngx-bootstrap/modal';

// Sell routes
import { DashboardAccessRoutingModule } from './dashboard.access.routes';
// Components
import { DashboardAccessComponent } from './main/dashboard.access.component';
import { DashboardAccessDetailsComponent } from './access-details/dashboard.access.details.component';
import { DashboardAccessAddComponent } from './add/dashboard.access.add.component';
import { DashboardAccessHelpComponent } from './help/dashboard.access.help.component';
import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';


@NgModule({
  imports: [
    CommonModule,
    ModalModule.forRoot(),
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    CommonAppModule,
    DashboardAccessRoutingModule,
    SharedModule,
  ],
  declarations: [
    DashboardAccessComponent,
    DashboardAccessAddComponent,
    DashboardAccessDetailsComponent,
    DashboardAccessHelpComponent,
  ],
  entryComponents: [
    DashboardAccessAddComponent,
  ]
})
export class DashboardAccessModule { }
