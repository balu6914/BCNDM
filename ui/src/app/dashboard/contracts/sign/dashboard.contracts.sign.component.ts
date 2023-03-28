import { Component, Output, EventEmitter } from '@angular/core';
import { BsModalRef } from 'ngx-bootstrap';

import { ContractService } from 'app/dashboard/contracts/contract.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-contracts-sign',
  templateUrl: './dashboard.contracts.sign.component.html',
  styleUrls: [ './dashboard.contracts.sign.component.scss' ]
})
export class DashboardContractsSignComponent {
  contract: any;

  @Output() contractSigned: EventEmitter<any> = new EventEmitter();
  constructor(
    private  contractService: ContractService,
    public  modalSignContract: BsModalRef,
    public  alertService: AlertService,
  ) { }

  confirm() {
    const signReq = {
      stream_id: this.contract.stream_id,
      end_time: this.contract.end_time
    };

    this.contractService.sign(signReq).subscribe(
      res => {
        this.contractSigned.emit(this.contract);
        this.alertService.success(`Contract succesfully signed!`);
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );

    this.modalSignContract.hide();
  }
}
