import { CommonModule } from '@angular/common';
import { NgModule, ModuleWithProviders } from '@angular/core';
import { TableModule } from './table/table.module';
import { BalanceModule } from './balance/balance.module';
import { MapModule } from './map/map.module';
import { FiltersModule } from './filters/filters.module';
import { AlertsModule } from './alerts/alerts.module';
import { StatisticModule } from './statistic/statistic.module';

@NgModule({
  imports: [
    CommonModule,
    TableModule,
    BalanceModule,
    MapModule,
    FiltersModule,
    AlertsModule,
    StatisticModule,
  ],
  exports: [
    TableModule,
    BalanceModule,
    MapModule,
    FiltersModule,
    AlertsModule,
    StatisticModule,
  ],
})
export class SharedModule {}
