import { CommonModule } from '@angular/common';
import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonAppModule } from '../../common/common.module';
import { TableRowComponent } from './row/table.row.component';
import { TableComponent } from './main/table.component';
import { DeleteButtonComponent } from './delete-button/delete-button.component';
import { TablePaginationComponent } from './pagination/table.pagination.component';

@NgModule({
  imports: [
    CommonModule,
    CommonAppModule
  ],
  declarations: [
    TableComponent,
    TableRowComponent,
    DeleteButtonComponent,
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
