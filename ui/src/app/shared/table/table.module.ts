import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
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
  ]
})
export class TableModule { }
