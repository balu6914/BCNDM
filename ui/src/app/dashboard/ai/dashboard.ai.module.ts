import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ModalModule } from 'ngx-bootstrap/modal';

// AI routes
import { DashboardAiRoutingModule } from './dashboard.ai.routes';
// AI components
import { DashboardAiComponent } from './main/dashboard.ai.component';
import { DashboardAiExecuteComponent } from './execute/dashboard.ai.execute.component';

import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';

import { DashboardAiAddComponent } from './add/dashboard.ai.add.component';
import { DashboardAiEditComponent } from './edit/dashboard.ai.edit.component';
import { DashboardAiDeleteComponent } from './delete/dashboard.ai.delete.component';

@NgModule({
  imports: [
    CommonModule,
    ModalModule.forRoot(),
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
    DashboardAiDeleteComponent,
  ],
  entryComponents: [
    DashboardAiExecuteComponent,
    DashboardAiAddComponent,
    DashboardAiEditComponent,
    DashboardAiDeleteComponent,
  ],
})
export class DashboardAiModule { }
