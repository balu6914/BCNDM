import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MdlModule } from '@angular-mdl/core';
import { NgProgressModule } from 'ngx-progressbar';
import {RouterModule} from '@angular/router';
// Layout Cmponents
import { HeaderComponent } from './header';
import { SidebarComponent } from './sidebar';

@NgModule({
  imports: [
    CommonModule,
    MdlModule,
    RouterModule,
    NgProgressModule,
  ],
  declarations: [
      HeaderComponent,
      SidebarComponent
  ],
  exports: [
      HeaderComponent,
      SidebarComponent
  ]
})
export class LayoutModule { }
