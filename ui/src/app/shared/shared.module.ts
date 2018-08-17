import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TableModule } from './table/table.module';
import { BalanceModule } from './balance/balance.module';
import { MapModule } from './map/map.module';
import { FiltersModule } from './filters/filters.module';
import { AlertsModule } from './alerts/alerts.module';

@NgModule({
  imports: [
    CommonModule,
    TableModule,
    BalanceModule,
    MapModule,
    FiltersModule,
    AlertsModule
  ],
  exports: [
    TableModule,
    BalanceModule,
    MapModule,
    FiltersModule,
    AlertsModule
  ]
})
export class SharedModule { }
