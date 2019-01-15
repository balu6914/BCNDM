import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap';
import { StreamService } from '../../../common/services/stream.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dashboard-sell-delete',
  templateUrl: './dashboard.sell.delete.component.html',
  styleUrls: ['./dashboard.sell.delete.component.scss']
})
export class DashboardSellDeleteComponent implements OnInit {
  stream: any;

  @Output() streamDeleted: EventEmitter<any> = new EventEmitter();
  constructor(
    private streamService: StreamService,
    public  modalDeleteStream: BsModalRef,
    public  alertService: AlertService,
  ) {}


  confirm(): void {
    // Send addStream request
    this.streamService.removeStream(this.stream.id).subscribe(
      res => {
        this.streamDeleted.emit(this.stream.id)
        this.modalDeleteStream.hide();
        this.alertService.success(`Stream succesfully removed!`);
      },
      err => {
        this.modalDeleteStream.hide();
        this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
    });
  }

  ngOnInit() {
  }

}
