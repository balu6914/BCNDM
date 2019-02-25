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
   streamsList= [];
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
    this.form = this.formBuilder.group({
      hours: ['', [Validators.required, Validators.min(1)]]
    });
  }

  ngOnInit() {
    let sum = 0;

    this.streamsList.forEach( stream => {
       sum += stream.price;
    })
    this.totalSum = sum;
  }
 
  onBuyAllClick(){
    if (this.form.invalid) {
      this.alertService.error('Please subscribe for at least one hour.');
      return;
    }
   
    let myRequests = []

    this.streamsList.forEach( stream => {

      const subsReq = {
        hours: this.form.value.hours,
        stream_id: stream.id,
      };

      myRequests.push(this.subscriptionService.add(subsReq));
    });

    console.log("Streams to buy:",myRequests);
    // TODO: Send the HTTP Request to the backend
  }
}
