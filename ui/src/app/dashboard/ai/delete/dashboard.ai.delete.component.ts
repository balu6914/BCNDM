import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap';
import { StreamService } from 'app/common/services/stream.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-ai-delete',
  templateUrl: './dashboard.ai.delete.component.html',
  styleUrls: ['./dashboard.ai.delete.component.scss']
})
export class DashboardAiDeleteComponent implements OnInit {
  stream: any;

  @Output() streamDeleted: EventEmitter<any> = new EventEmitter();
  constructor(
    private streamService: StreamService,
    public  modalDeleteAi: BsModalRef,
    public  alertService: AlertService,
  ) {}


  confirm(): void {
    this.streamService.removeStream(this.stream.id).subscribe(
      res => {
        this.streamDeleted.emit(this.stream.id);
        this.modalDeleteAi.hide();
        this.alertService.success(`AI succesfully removed!`);
      },
      err => {
        this.modalDeleteAi.hide();
        this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
    });
  }

  ngOnInit() {
  }

}
