import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { NgPipesModule } from 'ngx-pipes';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
// Sell routes
import { DashboardSellRoutingModule } from './dashboard.sell.routes';
// // Add services
import { SearchService } from './services/search.service';
import { StreamService } from './services/stream.service';
import { TokenInterceptor } from '../../auth/services/token.http.interceptor.service';
// Sell components
import { DashboardSellComponent } from './index';
import { DashboardSellAddComponent } from './add';
import { DashboardSellEditComponent } from './edit';
import { DashboardSellMapComponent } from './map';

import { CommonAppModule } from '../../common/common.module';
import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';

@NgModule({
  imports: [
    CommonModule,
    MdlModule,
    FormsModule,
    ReactiveFormsModule,
    NgPipesModule,
    NgxDatatableModule,
    CommonAppModule,
    DashboardSellRoutingModule,
    LeafletModule.forRoot(),
    LeafletDrawModule.forRoot()
  ],
  declarations: [
      DashboardSellComponent,
      DashboardSellAddComponent,
      DashboardSellEditComponent,
      DashboardSellMapComponent,
  ],
  providers: [
      StreamService,
      SearchService,
  ]
})
export class DashboardSellModule { }
