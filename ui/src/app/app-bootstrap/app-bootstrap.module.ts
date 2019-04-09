import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
// Import third party modules
import { SidebarModule } from 'ng-sidebar';
// Import ngx-bootstrap modules
import { BsDropdownModule } from 'ngx-bootstrap/dropdown';
import { ModalModule } from 'ngx-bootstrap/modal';
import { TooltipModule } from 'ngx-bootstrap/tooltip';
import { NgxSelectModule } from 'ngx-select-ex';
import { UiSwitchModule } from 'ngx-ui-switch';

@NgModule({
  imports: [
    CommonModule,
    // Import third party modules
    SidebarModule.forRoot(),
    NgxSelectModule,
    UiSwitchModule.forRoot({
      size: 'small',
      color: '#007bff',
    }),
    // Import ngx-bootstrap modules
    BsDropdownModule.forRoot(),
    TooltipModule.forRoot(),
    ModalModule.forRoot(),
  ],
  exports: [
    BsDropdownModule,
    TooltipModule,
    ModalModule,
    SidebarModule,
    UiSwitchModule,
  ]
})
export class AppBootstrapModule {}
