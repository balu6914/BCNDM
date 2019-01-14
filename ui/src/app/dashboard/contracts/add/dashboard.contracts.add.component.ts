import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormArray, FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';

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
  parties: FormArray;
  streamsName = [];
  streams = [];
  submitted = false;
  invalidShare = false;
  date = new Date();
  minDate: string;

  @Output() contractCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private authService: AuthService,
    private streamService: StreamService,
    private contractService: ContractService,
    private formBuilder: FormBuilder,
    public modalNewContract: BsModalRef,
    public alertService: AlertService,
  ) {
    // Add one day to current date and set it as min date
    this.date.setDate(this.date.getDate() + 1);
    this.minDate = this.date.toISOString().split('T')[0];

    this.form = this.formBuilder.group({
      streamName: ['', [Validators.required]],
      endTime:    ['', [Validators.required]],
      parties:    this.formBuilder.array([this.createPartner()])
    }, {
      validator: [this.shareValidator, this.partnerValidator.bind(this)]
    });

  }

  createPartner(): FormGroup {
    return this.formBuilder.group({
      partner: ['', [Validators.required]],
      share:   ['', [Validators.required, Validators.min(1)]],
    });
  }

  partnerValidator(fg: FormGroup) {
    // Verify that partner ID is not the owner ID
    fg.value.parties.forEach( (item, i) => {
        if (item.partner === this.user.id) {
          // Create a custom error field used as *ngIf condition for style
          fg.controls.parties['controls'][i].controls.partner.setErrors({
            'ownerID': true
          });
        }
    });
  }

  shareValidator(fg: FormGroup) {
    // Check if partner sharing sum is < 100
    let sum = 0;
    fg.value.parties.forEach( (item, i) => {
        sum += parseInt(item.share, 10);
        if (sum > 100) {
          // Create a custom error field used as *ngIf condition for style
          fg.controls.parties['controls'][i].controls.share.setErrors({
            'shareSum': true
          });
        }
    });

    return null;
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
        items: [],
      };

      this.form.value.parties.forEach( item => {
        const partner = {
          partner_id: item.partner,
          share: parseInt(item.share, 10)
        };
        createContractReq.items.push(partner);
      });

      this.contractService.create(createContractReq).subscribe(
        data => {
          const contractRow = {
            stream_name: this.form.value.streamName,
            start_time: new Date().toISOString(),
            end_time: createContractReq.end_time,
            parties: createContractReq.items,
          };
          this.contractCreated.emit(contractRow);
          this.alertService.success(`Contract succesfully created!`);
        },
        err => {
          if (err.status === 500) {
            this.alertService.error(`Not created. An active contract already exist for this stream.`);
          } else {
            this.alertService.error(`Error: ${err.status} - ${err.statusText}.`);
          }
        }
      );

      this.modalNewContract.hide();
    }
  }

  onAddPartner() {
    this.parties = this.form.get('parties') as FormArray;
    this.parties.push(this.createPartner());
  }

  onDeletePartner(index: number) {
    this.parties.controls.splice(index, 1);
  }
}
