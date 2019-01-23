import { Component, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'dpc-dashboard-access-help',
  templateUrl: './dashboard.access.help.component.html',
  styleUrls: [ './dashboard.access.help.component.scss' ]
})
export class DashboardAccessHelpComponent {

  @Output() closeTrigger = new EventEmitter();
  @Output() addAccessTrigger = new EventEmitter();

  constructor(
  ) {}

  close() {
    this.closeTrigger.emit();
  }

  create() {
    this.addAccessTrigger.emit();
  }

}
