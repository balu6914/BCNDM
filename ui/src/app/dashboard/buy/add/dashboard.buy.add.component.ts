import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';
import { SubscriptionService } from 'app/common/services/subscription.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { BalanceService } from 'app/shared/balance/balance.service';
import { Balance } from 'app/common/interfaces/balance.interface';
import 'rxjs/add/observable/throw';
import { Stream } from 'app/common/interfaces/stream.interface';

@Component({
  selector: 'dpc-dashboard-buy-add',
  templateUrl: './dashboard.buy.add.component.html',
  styleUrls: ['./dashboard.buy.add.component.scss']
})
export class DashboardBuyAddComponent implements OnInit {
  user: any;
  form: FormGroup;
  stream: Stream;
  modalMsg: string;
  balance = new Balance();
  submitted = false;

  constructor(
    private subscriptionService: SubscriptionService,
    private formBuilder: FormBuilder,
    public modalSubscription: BsModalRef,
    private alertService: AlertService,
    private balanceService: BalanceService,
  ) {
    this.form = this.formBuilder.group({
      hours: ['', [Validators.required, Validators.min(1)]]
    }, {
      validator: this.hoursValidator.bind(this)
    });
  }

  ngOnInit() {
    this.getBalance();
  }

  hoursValidator(fg: FormGroup) {
    if (this.stream) {
      if (this.balance.amount < (this.stream.price * fg.get('hours').value)) {
        fg.controls.hours.setErrors({ 'insufficientFunds': true });
      } else {
        return null;
      }
    }
  }

  getBalance() {
    this.balanceService.get().subscribe(
      (result: any) => {
        this.balance.amount = result.balance;
        this.balance.fiatAmount = this.balance.amount;
        this.balanceService.changed(this.balance);
      },
      err => {
        this.alertService.error('Error fetching user balance.');
      });
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const subsReq = [{
        hours: Number(this.form.value.hours),
        stream_id: this.stream.id,
      }];

      // Send subscription request
      this.subscriptionService.add(subsReq).subscribe(
        (response: any) => {
          this.modalSubscription.hide();
          // Display the results
          response.Responses.forEach(resp => {
            if (resp.errorMessage == null) {
              this.alertService.success(`You now have access to ${resp.subscriptionID} for ${this.form.value.hours} hours.`);
            } else {
              this.alertService.error(`Error trying to subscribe to Stream ${resp.streamID} due to ${resp.errorMessage}.`);
            }
          });
          // Fetch new user balance and publish new balance value to message buss
          this.getBalance();
        },
        err => {
          this.modalSubscription.hide();
          this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
        }
      );
    }
  }

}
