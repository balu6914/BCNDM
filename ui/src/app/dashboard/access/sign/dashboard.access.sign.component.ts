import { Component, Output, EventEmitter } from '@angular/core';
import { BsModalRef } from 'ngx-bootstrap';

import { AccessService } from 'app/dashboard/access/access.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-access-sign',
  templateUrl: './dashboard.access.sign.component.html',
  styleUrls: [ './dashboard.access.sign.component.scss' ]
})
export class DashboardAccessSignComponent {
  access: any;

  @Output() accessSigned: EventEmitter<any> = new EventEmitter();
  constructor(
    private  accessService: AccessService,
    public  modalSignAccess: BsModalRef,
    public  alertService: AlertService,
  ) { }

  confirm() {
    const signReq = {
      stream_id: this.access.stream_id,
      end_time: this.access.end_time
    };

    this.accessService.sign(signReq).subscribe(
      res => {
        this.accessSigned.emit(this.access);
        this.alertService.success(`Access succesfully signed!`);
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );

    this.modalSignAccess.hide();
  }
}
