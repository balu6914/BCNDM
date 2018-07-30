import { Component } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { SubscriptionService } from '../../../common/services/subscription.service';

@Component({
  selector: 'subscription-add',
  templateUrl: './subscription.add.component.html',
  styleUrls: [ './subscription.add.component.scss' ]
})
export class SubscriptionAddComponent {
    user:any;
    form: FormGroup;
    stream: any;
    modalMsg: string;
    submitted: boolean = false;

    constructor(
      private subscriptionService: SubscriptionService,
      private formBuilder: FormBuilder,
      public  modalSubscription: BsModalRef,
     ) {
       this.form = formBuilder.group({
         hours: ['', Validators.required],
       })
    }

    ngOnInit() {
    }

    onSubmit() {
      const subsReq = {
        hours: this.form.value.hours,
        id:    this.stream.id,
      }
      const streamName = this.stream.name;

      // Send subscription request
      this.subscriptionService.add(subsReq).subscribe(
        response => {
          this.modalMsg = `Stream succesfully added!`;`Success! You now have access to ${streamName} in next ${subsReq.hours} hours`
        },
        err => {
          this.modalMsg = `Status: ${err.status} - ${err.statusText}`;
      });

      // Hide modalSubscription and show modalResponse
      this.submitted = true;
    }
}
