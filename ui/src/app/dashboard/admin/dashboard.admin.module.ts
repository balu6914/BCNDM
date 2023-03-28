import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ModalModule } from 'ngx-bootstrap/modal';

import { DashboardAdminRoutingModule } from './dashboard.admin.routes';
import { DashboardAdminComponent } from './dashboard.admin.component';
import { DashboardAdminSignupComponent } from './signup/dashboard.admin.signup.component';
import { DashboardAdminEditComponent } from './edit/dashboard.admin.edit.component';
import { DashboardAdminDeleteComponent } from './delete/dashboard.admin.delete.component';
import { DashboardAdminLockComponent } from './lock/dashboard.admin.lock.component';

import { CommonAppModule } from 'app/common/common.module';
import { SharedModule } from 'app/shared/shared.module';


@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    CommonAppModule,
    DashboardAdminRoutingModule,
    SharedModule,
    ModalModule.forRoot(),
  ],
  declarations: [
    DashboardAdminComponent,
    DashboardAdminSignupComponent,
    DashboardAdminEditComponent,
    DashboardAdminDeleteComponent,
    DashboardAdminLockComponent,
  ],
  entryComponents: [
    DashboardAdminSignupComponent,
    DashboardAdminEditComponent,
    DashboardAdminDeleteComponent,
    DashboardAdminLockComponent,
  ],
})
export class DashboardAdminModule { }
