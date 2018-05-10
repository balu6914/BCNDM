import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import {NgPipesModule} from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
// Buy routes
import { DashboardBuyRoutingModule } from './dashboard.buy.routes';
// // Add services
import { SearchService } from './services/search.service';
import { TokenInterceptor } from '../../auth/services/token.http.interceptor.service';
// Buy components
import { DashboardBuyComponent } from './index';
import { DashboardBuyMapComponent } from './map';
// Import subscription module
import { SubscriptionModule } from '../subscription';

import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';
import { CommonAppModule } from '../../common/common.module';


@NgModule({
  imports: [
    CommonModule,
    MdlModule,
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    NgxDatatableModule,
    CommonAppModule,
    DashboardBuyRoutingModule,
    SubscriptionModule,
    LeafletModule.forRoot(),
    LeafletDrawModule.forRoot()
  ],
  declarations: [
      DashboardBuyComponent,
      DashboardBuyMapComponent,
  ],
  providers: [
      SearchService,
  ]
})
export class DashboardBuyModule { }
