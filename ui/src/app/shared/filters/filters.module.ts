import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgxSelectModule } from 'ngx-select-ex';
import { SidebarModule } from 'ng-sidebar';

import { CommonAppModule } from 'app/common/common.module';
import { SidebarFiltersComponent } from './sidebar-filters/sidebar.filters.component';

@NgModule({
  imports: [
    CommonModule,
    NgxSelectModule,
    FormsModule,
    ReactiveFormsModule,
    SidebarModule,
    CommonAppModule
  ],
  declarations: [
    SidebarFiltersComponent
  ],
  exports: [
    SidebarFiltersComponent
  ],
})
export class FiltersModule { }
