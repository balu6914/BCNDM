import { Component, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'dpc-dashboard-contracts-help',
  templateUrl: './dashboard.contracts.help.component.html',
  styleUrls: [ './dashboard.contracts.help.component.scss' ]
})
export class DashboardContractsHelpComponent {

  @Output() closeTrigger = new EventEmitter();
  @Output() addContractTrigger = new EventEmitter();

  constructor(
  ) {}

  close() {
    this.closeTrigger.emit();
  }

  create() {
    this.addContractTrigger.emit();
  }

}
