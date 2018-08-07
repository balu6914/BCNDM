import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
// Import third party modules
import { SidebarModule } from 'ng-sidebar';
// Import ngx-bootstrap modules
import { BsDropdownModule } from 'ngx-bootstrap/dropdown';
import { TooltipModule } from 'ngx-bootstrap/tooltip';
import { ModalModule } from 'ngx-bootstrap/modal';

@NgModule({
  imports: [
    CommonModule,
    // Import third party modules
    SidebarModule.forRoot(),
    // Import ngx-bootstrap modules
    BsDropdownModule.forRoot(),
    TooltipModule.forRoot(),
    ModalModule.forRoot()
  ],
  exports: [
    BsDropdownModule,
    TooltipModule,
    ModalModule,
    SidebarModule,
  ]
})
export class AppBootstrapModule {}
