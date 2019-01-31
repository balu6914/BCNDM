import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Sell routes
import { DashboardAccessRoutingModule } from './dashboard.access.routes';
// Components
import { DashboardAccessComponent } from './dashboard.access.component';
import { DashboardAccessDetailsComponent } from './access-details/dashboard.access.details.component';
import { DashboardAccessAddComponent } from './add/dashboard.access.add.component';
import { DashboardAccessSignComponent } from './sign/dashboard.access.sign.component';
import { DashboardAccessHelpComponent } from './help/dashboard.access.help.component';
import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';
import { AppBootstrapModule } from 'app/app-bootstrap/app-bootstrap.module';

@NgModule({
  imports: [
    CommonModule,
    AppBootstrapModule,
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
    DashboardAccessSignComponent,
  ],
  entryComponents: [
    DashboardAccessAddComponent,
    DashboardAccessSignComponent,
  ]
})
export class DashboardAccessModule { }
