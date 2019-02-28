import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';
import { SubscriptionService } from 'app/common/services/subscription.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { BalanceService } from 'app/shared/balance/balance.service';
import { Balance } from 'app/common/interfaces/balance.interface';
import 'rxjs/add/observable/throw';

@Component({
  selector: 'dpc-dashboard-buy-add',
  templateUrl: './dashboard.buy.add.component.html',
  styleUrls: ['./dashboard.buy.add.component.scss']
})
export class DashboardBuyAddComponent {
  user: any;
  form: FormGroup;
  stream: any;
  modalMsg: string;
  balance = new Balance();

  constructor(
    private subscriptionService: SubscriptionService,
    private formBuilder: FormBuilder,
    public modalSubscription: BsModalRef,
    private alertService: AlertService,
    private balanceService: BalanceService,
  ) {
    this.form = this.formBuilder.group({
      hours: ['', [Validators.required, Validators.min(1)]]
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

  onSubmit() {
    if (this.form.invalid) {
      this.modalSubscription.hide();
      this.alertService.error('Please subscribe for at least one hour.');
      return;
    }
    const subsReq = [{
      hours: this.form.value.hours,
      stream_id: this.stream.id,
    }];

    // Send subscription request
    this.subscriptionService.add(subsReq).subscribe(
      (response: any) => {
        this.modalSubscription.hide();
        //Display the results
        response.Responses.forEach(resp => {
          if (resp.errorMessage == null) {
            this.alertService.success(`You now have access to ${resp.subscriptionID} for ${subsReq[0].hours} hours.`);
          }
          else {
            this.alertService.error(`Error trying to subscribe to Stream ${resp.streamID} due to ${resp.errorMessage}.`);
          }
        });
        // Fetch new user balance and publish new balance value to message buss
        this.getBalance()
      },
      err => {
        this.modalSubscription.hide();
        if (err.status === 402) {
          this.alertService.error(`No enough funds.`);
        } else {
          this.alertService.error(`Something went wrong. Please try again later.`);
        }
      });
  }
}
