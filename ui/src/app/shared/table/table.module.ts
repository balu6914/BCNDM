import { CommonModule } from '@angular/common';
import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonAppModule } from '../../common/common.module';
import { TableRowComponent } from './row/table.row.component';
import { TableRowDetailsComponent } from './details/table.row.details.component';
import { TableComponent } from './main/table.component';
import { TablePaginationComponent } from './pagination/table.pagination.component';

@NgModule({
  imports: [
    CommonModule,
    CommonAppModule
  ],
  declarations: [
    TableComponent,
    TableRowComponent,
    TableRowDetailsComponent,
    TablePaginationComponent
  ],
  exports: [
    TableComponent
  ],
  schemas: [
    CUSTOM_ELEMENTS_SCHEMA
  ],

})
export class TableModule { }
