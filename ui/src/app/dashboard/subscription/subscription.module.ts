import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { RouterModule } from "@angular/router";
import {
    FormsModule,
    ReactiveFormsModule
} from '@angular/forms';
// Add services
import { TokenInterceptor } from '../../auth/services/token.http.interceptor.service';
import { SubscriptionSrvice } from './services/subscription.service';
// Add components
import { SubscriptionAddComponent } from './add'
import { CommonAppModule } from '../../common/common.module';

@NgModule({
  imports: [
    RouterModule,
    CommonModule,
    MdlModule,
    FormsModule,
    ReactiveFormsModule,
    CommonAppModule
  ],
  declarations: [
      SubscriptionAddComponent,
  ],
  providers: [
    SubscriptionSrvice,
  ],
  entryComponents:[SubscriptionAddComponent]
})
export class SubscriptionModule { }
