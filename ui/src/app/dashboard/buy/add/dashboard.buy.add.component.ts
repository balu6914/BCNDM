import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';
import { SubscriptionService } from '../../../common/services/subscription.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { BalanceService } from '../../../shared/balance/balance.service';
import { Balance } from '../../../common/interfaces/balance.interface';

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

  onSubmit() {
    if (this.form.invalid) {
      this.modalSubscription.hide();
      this.alertService.error('Please subscribe for at least one hour.');
      return;
    }
    const subsReq = {
      hours: this.form.value.hours,
      id: this.stream.id,
    };
    const streamName = this.stream.name;

    // Send subscription request
    this.subscriptionService.add(subsReq).subscribe(
      response => {
        // Fetch new user balance and publish new balance value to message buss
        this.balanceService.get().subscribe(
        (result: any) => {
          this.balance.amount = result.balance;
          this.balance.fiatAmount = this.balance.amount;
          this.balanceService.changed(this.balance);
          // Close modal and show success message
          this.modalSubscription.hide();
          this.alertService.success(`You now have access to ${this.stream.name} stream in next ${subsReq.hours} hours`);
        },
        err => {
          console.error("Error fetching user balance ", err)
        });

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
