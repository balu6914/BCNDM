import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { AuthService } from 'app/auth/services/auth.service';
import { User } from 'app/common/interfaces/user.interface';
import { StreamService } from 'app/common/services/stream.service';
import { Query } from 'app/common/interfaces/query.interface';
import { ContractService } from 'app/dashboard/contracts/contract.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-contracts-add',
  templateUrl: './dashboard.contracts.add.component.html',
  styleUrls: [ './dashboard.contracts.add.component.scss' ]
})
export class DashboardContractsAddComponent implements OnInit {
  user: User;
  query = new Query();
  form: FormGroup;
  streamsName = [];
  streams = [];

  @Output() contractCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private authService: AuthService,
    private streamService: StreamService,
    private  contractService: ContractService,
    private formBuilder: FormBuilder,
    public  modalNewContract: BsModalRef,
    public  alertService: AlertService,
  ) {
    this.form = this.formBuilder.group({
      streamName:   ['', [<any>Validators.required]],
      parties:      ['', [<any>Validators.required]],
      shareOffered: ['', [<any>Validators.required]],
      endTime:      ['', [<any>Validators.required]],
      endTimeHour:  ['', [<any>Validators.required]]
    });
  }

  ngOnInit () {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.query.owner = this.user.id;
      },
      err => {
        console.log(err);
      }
    );

    // Fetch streams
    this.streamService.searchStreams(this.query).subscribe(
      (result: any) => {
        this.streams = result.content;
        this.streamsName = this.streams.map(stream => stream.name);
      },
      err => {
        this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
      }
    );

  }

  onSubmit() {
    const index = this.streamsName.indexOf(this.form.value.streamName);
    const streamId = this.streams[index].id;

    const createcontractReq = {
      stream_id: streamId,
      end_time: `${this.form.value.endTime}T${this.form.value.endTimeHour}:00Z`,
      items: [
        {
          partner_id: this.form.value.parties,
          share: parseInt(this.form.value.shareOffered, 10)
        }
      ],
    };

    this.contractService.create(createcontractReq).subscribe(
      data => {
        const contractRow = {
          stream_id: streamId,
          start_time:  new Date().toISOString(),
          end_time: `${this.form.value.endTime}T${this.form.value.endTimeHour}:00Z`,
          share: this.form.value.shareOffered
        };
        this.contractCreated.emit(contractRow);
        this.alertService.success(`Contract succesfully created!`);
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );

    this.modalNewContract.hide();
  }
}
