import { Component, OnInit } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { StreamService } from '../../../common/services/stream.service';

@Component({
  selector: 'dashboard-sell-delete',
  templateUrl: './dashboard.sell.delete.component.html',
  styleUrls: ['./dashboard.sell.delete.component.scss']
})
export class DashboardSellDeleteComponent implements OnInit {
  stream: any;
  modalMsg: string;
  submitted: boolean = false;

  constructor(
    private streamService: StreamService,
    public  modalDeleteStream: BsModalRef,
  ) {}


  confirm(): void {
    // Send addStream request
    this.streamService.removeStream(this.stream.id).subscribe(
      response => {
        this.modalMsg = `Stream succesfully removed!`;
      },
      err => {
        this.modalMsg = `Status: ${err.status} - ${err.statusText}`;
    });

    // Hide modalDeleteStream and show modalResponse
    this.submitted = true;
  }

  ngOnInit() {
  }

}
