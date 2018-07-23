import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ArchwizardModule } from 'ng2-archwizard';

import { LayoutModule } from '../layout'
// Auth module
import { AuthModule } from '../auth/auth.module';

import { DashboardRoutingModule } from './dashboard.routes';
import { CommonAppModule } from '../common/common.module';
import { SharedModule } from '../shared/shared.module';

// Dashboard components
import { DashboardComponent } from './dashboard.component';
import { DashboardMainComponent } from './main';
import { WalletModule } from './wallet/wallet.module';
import { SubscriptionModule } from './subscription/index';

import { NgxDatatableModule } from '@swimlane/ngx-datatable';
// Import subscription module
import { SubscriptionService } from './main/services/subscription.service';
import { StreamService } from './main/services/stream.service';
import { SearchService } from './main/services/search.service';

import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';
import { ClipboardModule } from 'ngx-clipboard';
import { MapComponent } from '../common/map/map.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    FormsModule,
    ReactiveFormsModule,
    ArchwizardModule,
    // App module
    AuthModule,
    CommonAppModule,
    SharedModule,
    LayoutModule,
    WalletModule,
    SubscriptionModule,
    DashboardRoutingModule,
    NgxDatatableModule,
    LeafletModule.forRoot(),
    LeafletDrawModule.forRoot(),
    ClipboardModule
  ],
  declarations: [
      DashboardComponent,
      DashboardMainComponent,
      MapComponent,
  ],
  providers: [
      SubscriptionService,
      StreamService,
      SearchService,
  ],
 schemas: [ CUSTOM_ELEMENTS_SCHEMA ]
})
export class DashboardModule { }
