import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/shared/validators/patterns';
import { BsModalRef } from 'ngx-bootstrap';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BigQuery, Stream } from 'app/common/interfaces/stream.interface';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
import { SubscriptionService } from 'app/common/services/subscription.service';
import { BalanceService } from 'app/shared/balance/balance.service';
import { Balance } from 'app/common/interfaces/balance.interface';
import { StreamService } from 'app/common/services/stream.service';

@Component({
  selector: 'dpc-dashboard-buy-group',
  templateUrl: './dashboard.buy.group.component.html',
  styleUrls: ['./dashboard.buy.group.component.scss']
})
export class DashboardBuyGroupComponent implements OnInit {
  form: FormGroup;
  streamsList = [];
  totalSum: number;
  balance = new Balance();

  @Output()
  streamCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private subscriptionService: SubscriptionService,
    private formBuilder: FormBuilder,
    private alertService: AlertService,
    public modalBuyGroupStream: BsModalRef,
    private balanceService: BalanceService,
  ) {
  }

  ngOnInit() {
    this.getBalance();
    let sum = 0;

    this.streamsList.forEach(stream => {
      sum += stream.price;
    })
    this.totalSum = sum;

    this.form = this.formBuilder.group({
      hours: ['', [Validators.required, Validators.min(1)]]
    },
      {
        validator: this.hoursValidator.bind(this)
      });
  }

  getBalance() {
    this.balanceService.get().subscribe(
      (result: any) => {
        this.balance.amount = result.balance;
        this.balance.fiatAmount = this.balance.amount;
        this.balanceService.changed(this.balance);
      },
      err => {
        console.error("Error fetching user balance ", err)
      });
  }

  hoursValidator(fg: FormGroup) {
    if (this.balance.amount < this.totalSum * fg.get('hours').value) {
      fg.controls.hours.setErrors({ 'insufficientFunds': true });
    }
    else {
      return null;
    }
  }

  onBuyAllClick() {
    if (this.form.valid) {

      let myRequests = []

      this.streamsList.forEach(stream => {
        const subsReq = {
          hours: this.form.value.hours,
          stream_id: stream.id,
        };

        myRequests.push(subsReq);
      });

      // Send subscription request
      this.subscriptionService.add(myRequests).subscribe(
        (response: any) => {
          this.modalBuyGroupStream.hide();
          //Display the results
          response.Responses.forEach(resp => {
            if (resp.errorMessage == null) {
              this.alertService.success(`You now have access to ${resp.subscriptionID} for ${myRequests[0].hours} hours.`);
            }
            else {
              this.alertService.error(`Error trying to subscribe to Stream ${resp.streamID} due to ${resp.errorMessage}.`);
            }
          });
          // Fetch new user balance and publish new balance value to message bus
          this.getBalance();
        },
        err => {
          if (err.status === 402) {
            this.alertService.error(`No enough funds.`);
          } else {
            this.alertService.error(`Something went wrong. Please try again later.`);
          }
        });
    }
  }
}
