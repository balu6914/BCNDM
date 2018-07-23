import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TableModule } from './table/table.module';
import { BalanceModule } from './balance/balance.module';

@NgModule({
  imports: [
    CommonModule,
    TableModule,
    BalanceModule
  ],
  exports: [TableModule, BalanceModule]
})
export class SharedModule { }
