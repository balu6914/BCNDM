import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TableModule } from './table/table.module';
import { BalanceModule } from './balance/balance.module';
import { MapModule } from './map/map.module';
import { FiltersModule } from './filters/filters.module';

@NgModule({
  imports: [
    CommonModule,
    TableModule,
    BalanceModule,
    MapModule,
    FiltersModule
  ],
  exports: [
    TableModule,
    BalanceModule,
    MapModule,
    FiltersModule
  ]
})
export class SharedModule { }
