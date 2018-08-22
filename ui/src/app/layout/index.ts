import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgProgressModule } from 'ngx-progressbar';
import {RouterModule} from '@angular/router';
import { AppBootstrapModule } from '../app-bootstrap/app-bootstrap.module'
import { CommonAppModule } from '../common/common.module';

 // Layout Cmponents
import { HeaderComponent } from './header';

@NgModule({
  imports: [
    CommonModule,
    AppBootstrapModule,
    RouterModule,
    NgProgressModule,
    CommonAppModule,
  ],
  declarations: [
      HeaderComponent,
  ],
  exports: [
      HeaderComponent,
  ]
})
export class LayoutModule { }
