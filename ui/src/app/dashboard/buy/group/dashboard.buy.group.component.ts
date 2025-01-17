import { Component, EventEmitter, OnInit } from '@angular/core';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/common/validators/patterns';
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
  submitted = false;

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
    });
    this.totalSum = sum;

    this.form = this.formBuilder.group({
      hours: ['', [Validators.required, Validators.min(1)]]
    }, {
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
        this.alertService.error(`Error fetching user balance.`);
      });
  }

  hoursValidator(fg: FormGroup) {
    if (this.balance.amount < this.totalSum * fg.get('hours').value) {
      fg.controls.hours.setErrors({ 'insufficientFunds': true });
    } else {
      return null;
    }
  }

  onBuyAllClick() {
    this.submitted = true;

    if (this.form.valid) {
      const subsReq = [];
      this.streamsList.forEach(stream => {
        const subReq = {
          hours: Number(this.form.value.hours),
          stream_id: stream.id,
        };
        subsReq.push(subReq);
      });

      // Send subscription request
      this.subscriptionService.add(subsReq).subscribe(
        (response: any) => {
          this.modalBuyGroupStream.hide();
          // Display the results
          response.Responses.forEach(resp => {
            if (resp.errorMessage == null) {
              this.alertService.success(`You now have access to ${resp.subscriptionID} for ${this.form.value.hours} hours.`);
            } else {
              this.alertService.error(`Error: Subscription to Stream ${resp.streamID}failed due to ${resp.errorMessage}.`);
            }
          });
          // Fetch new user balance and publish new balance value to message bus
          this.getBalance();
        },
        err => {
          this.modalBuyGroupStream.hide();
          this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
        });
    }
  }
}
