import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// AI routes
import { DashboardAiRoutingModule } from './dashboard.ai.routes';
// AI components
import { DashboardAiComponent } from './main/dashboard.ai.component';
import { DashboardAiExecuteComponent } from './execute/dashboard.ai.execute.component';

import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';
import { AppBootstrapModule } from 'app/app-bootstrap/app-bootstrap.module';
import { DashboardAiAddComponent } from './add';
import { DashboardAiEditComponent } from './edit';

@NgModule({
  imports: [
    CommonModule,
    AppBootstrapModule,
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    CommonAppModule,
    DashboardAiRoutingModule,
    SharedModule,
  ],
  declarations: [
    DashboardAiComponent,
    DashboardAiExecuteComponent,
    DashboardAiAddComponent,
    DashboardAiEditComponent,
  ],
  entryComponents: [
    DashboardAiExecuteComponent,
    DashboardAiAddComponent,
    DashboardAiEditComponent,
  ],
})
export class DashboardAiModule { }
