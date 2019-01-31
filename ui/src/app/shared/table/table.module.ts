import { CommonModule } from '@angular/common';
import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import {RouterModule} from '@angular/router';
import { CommonAppModule } from 'app/common/common.module';
import { TableRowComponent } from './row/table.row.component';
import { TableRowDetailsComponent } from './details/table.row.details.component';
import { TableComponent } from './main/table.component';
import { TablePaginationComponent } from './pagination/table.pagination.component';
import { TableTabbedComponent } from './table-tabbed/table-tabbed.component';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    CommonAppModule
  ],
  declarations: [
    TableComponent,
    TableTabbedComponent,
    TableRowComponent,
    TableRowDetailsComponent,
    TablePaginationComponent
  ],
  exports: [
    TableComponent,
    TableTabbedComponent
  ],
  schemas: [
    CUSTOM_ELEMENTS_SCHEMA
  ],

})
export class TableModule { }
