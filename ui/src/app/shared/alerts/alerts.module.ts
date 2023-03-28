import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { AlertContainerComponent } from './alert-container/alert-container.component';
import { AlertComponent } from './alert/alert.component';
import { DomService } from './services/dom.service';
import { AlertService } from './services/alert.service';

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [AlertComponent, AlertContainerComponent],
  exports: [],
  providers: [DomService, AlertService],
  entryComponents: [AlertComponent, AlertContainerComponent]
})
export class AlertsModule { }
