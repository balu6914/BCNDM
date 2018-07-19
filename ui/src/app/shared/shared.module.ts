import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TableModule } from './table/table.module';

@NgModule({
  imports: [
    CommonModule,
    TableModule
  ],
  exports: [TableModule]
})
export class SharedModule { }
