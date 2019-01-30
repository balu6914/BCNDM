import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import {
  FormsModule,
  ReactiveFormsModule
} from '@angular/forms';
import { NgxSelectModule } from 'ngx-select-ex';
import { AppBootstrapModule } from 'app/app-bootstrap/app-bootstrap.module';
import { CommonAppModule } from 'app/common/common.module';
import { SidebarFiltersComponent } from './sidebar-filters/sidebar.filters.component';



@NgModule({
  imports: [
    CommonModule,
    NgxSelectModule,
    FormsModule,
    ReactiveFormsModule,
    AppBootstrapModule,
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
