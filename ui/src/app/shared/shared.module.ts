import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TableModule } from './table/table.module';
import { BalanceModule } from './balance/balance.module';

// Import map component and ngx-leaflet
import { MapComponent } from './map/leaflet/map.leaflet.component';
import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';

@NgModule({
  imports: [
    CommonModule,
    TableModule,
    BalanceModule,
    LeafletModule.forRoot(),
    LeafletDrawModule.forRoot(),
  ],
  declarations: [
    MapComponent,
  ],
  exports: [
    TableModule,
    BalanceModule,
    MapComponent,
  ]
})
export class SharedModule { }
