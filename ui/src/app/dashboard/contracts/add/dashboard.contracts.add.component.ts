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
  user = <User>{};
  query = new Query();
  form: FormGroup;
  streamsName = [];
  streams = [];
  submitted = false;
  date = new Date();
  minDate: string;

  @Output() contractCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private authService: AuthService,
    private streamService: StreamService,
    private  contractService: ContractService,
    private formBuilder: FormBuilder,
    public  modalNewContract: BsModalRef,
    public  alertService: AlertService,
  ) {
    // Add one day to current date and set it as min date
    this.date.setDate(this.date.getDate() + 1);
    this.minDate = this.date.toISOString().split('T')[0];

    this.form = this.formBuilder.group({
      streamName:   ['', [Validators.required]],
      parties:      ['', [Validators.required]],
      shareOffered: ['', [Validators.required]],
      endTime:      ['', [Validators.required]]
    });
  }

  ngOnInit () {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.query.owner = this.user.id;
        // Fetch streams as owner
        this.streamService.searchStreams(this.query).subscribe(
          (result: any) => {
            this.streams = result.content;
            this.streamsName = this.streams.map(stream => stream.name);
          },
          err => {
            this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
          }
        );
      },
      err => {
        console.log(err);
      }
    );
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const index = this.streamsName.indexOf(this.form.value.streamName);
      const streamID = this.streams[index].id;

      const createContractReq = {
        stream_id: streamID,
        end_time: `${this.form.value.endTime}T00:00:00Z`,
        items: [
          {
            partner_id: this.form.value.parties,
            share: parseInt(this.form.value.shareOffered, 10)
          }
        ],
      };

      this.contractService.create(createContractReq).subscribe(
        data => {
          const contractRow = {
            stream_id: streamID,
            start_time:  new Date().toISOString(),
            end_time: this.form.value.endTime,
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
}
