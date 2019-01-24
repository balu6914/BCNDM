import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// AI routes
import { DashboardAIRoutingModule } from './dashboard.ai.routes';
// AI components
import { DashboardAiComponent } from './dashboard.ai.component';
import { DashboardAiExecuteComponent } from './execute/dashboard.ai.execute.component';

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
    DashboardAIRoutingModule,
    SharedModule,
  ],
  declarations: [
    DashboardAiComponent,
    DashboardAiExecuteComponent,
  ],
  entryComponents: [
    DashboardAiExecuteComponent,
  ],
})
export class DashboardAiModule { }
