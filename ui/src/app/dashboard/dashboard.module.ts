import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ArchwizardModule } from 'ng2-archwizard';

import { LayoutModule } from '../layout'
// Auth module
import { AuthModule } from '../auth/auth.module';

import { DashboardRoutingModule } from './dashboard.routes';
// Dashboard components
import { DashboardComponent } from './dashboard.component';
import { DashboardMainComponent } from './main';
import { WalletModule } from './wallet/wallet.module';
import { SubscriptionModule } from './subscription/index';

import { CommonAppModule } from '../common/common.module';
// Import subscription module
import { SubscriptionService } from './main/services/subscription.service';
import { StreamService } from './main/services/stream.service';
import { SearchService } from './main/services/search.service';

import { ClipboardModule } from 'ngx-clipboard';
import { SharedModule } from '../shared/shared.module';

@NgModule({
  imports: [
    CommonModule,
    MdlModule,
    FormsModule,
    FormsModule,
    ReactiveFormsModule,
    ArchwizardModule,
    // App module
    AuthModule,
    CommonAppModule,
    LayoutModule,
    WalletModule,
    SubscriptionModule,
    DashboardRoutingModule,
    ClipboardModule,
    SharedModule,
  ],
  declarations: [
      DashboardComponent,
      DashboardMainComponent,
  ],
  providers: [
      SubscriptionService,
      StreamService,
      SearchService,
  ],
 schemas: [ CUSTOM_ELEMENTS_SCHEMA ]
})
export class DashboardModule { }
