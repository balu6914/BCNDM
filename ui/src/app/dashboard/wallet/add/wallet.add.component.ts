import { Component, OnInit, HostListener } from '@angular/core';

import { MitasPipe } from '../../../common/pipes/converter.pipe';

declare let paypal: any;

@Component({
  selector: 'wallet-add',
  templateUrl: './wallet.add.component.html',
  styleUrls: ['./wallet.add.component.scss']
})
export class WalletAddComponent implements OnInit {

    public showLoading: boolean = false;
    public showButton: boolean = true;

    constructor(
        private mitasPipe: MitasPipe
    ) {}

    ngOnInit() {
    }

   @HostListener('keydown.esc')
   public onEsc(): void {
   }

   public didPaypalScriptLoad: boolean = false;
   public loading: boolean = true;
   public paymentAmount: number = 0;
   public paypalConfig: any = {
       env: 'sandbox', // 'production' or 'sandbox'
       commit: true, // Show a 'Pay Now' button

       payment: () => {
           this.showLoading = true;
           if (!this.paymentAmount || this.paymentAmount > 10000) {
               this.showLoading = false;
               const msg = "Payment amount must be between 1 and 10000 TAS";
               // this.dialogRef.hide();
               // this.mdlDialogService.confirm(msg, "", "OK");
           } else {
           return paypal.request({
               method: 'post',
               url: '/api/payments/create',
               headers: {'Content-Type': 'application/json',
                         'Authorization': "Bearer " + localStorage.getItem('token'),
               },
               json: {
                   "price": this.paymentAmount.toString(),
                   "quantity": 1,
                   "total": this.paymentAmount.toString(),
                   "currency": "USD"
               }
           }).then( (data) => {
               return data.paymentID;
           });
       }
       },
       onAuthorize: (data) => {
           return paypal.request({
               method: 'post',
               url: '/api/payments/execute',
               headers: {'Content-Type': 'application/json',
                         'Authorization': "Bearer " + localStorage.getItem('token'),
               },
               json: {
                   "paymentID": data.paymentID,
                   "payerID":   data.payerID
               }
           }).then( () => {
               // Hide loading
               this.showLoading = false;
               const mitas = this.mitasPipe.transform(this.paymentAmount.toString());
               // this.dialogRef.hide(mitas);
               // Display msg with tranfer value in TAS
               let successMsg = `Success! We just transfered ` +
                                this.paymentAmount.toString() +
                                ` TAS to your account balance!`
               // this.mdlDialogService.confirm(successMsg, "", "OK");
           });
       }
   };

    public ngAfterViewChecked(): void {
        if(!this.didPaypalScriptLoad) {
            this.loadPaypalScript().then(() => {
                paypal.Button.render(this.paypalConfig, '#paypal-button');
                this.loading = false;
            });
        }
    }

    public loadPaypalScript(): Promise<any> {
        this.didPaypalScriptLoad = true;
        return new Promise((resolve, reject) => {
            const scriptElement = document.createElement('script');
            scriptElement.src = 'https://www.paypalobjects.com/api/checkout.js';
            scriptElement.onload = resolve;
            document.body.appendChild(scriptElement);
        });
    }
}
