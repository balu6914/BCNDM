import { Component, ViewContainerRef } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MdlDialogService, MdlSnackbarService, MdlDialogOutletService } from '@angular-mdl/core';
import { FormGroup, FormControl, ReactiveFormsModule, FormBuilder, Validators  } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';

import { AuthService } from '../../../auth/services/auth.service';
import { SubscriptionSrvice } from '../services/subscription.service';
import { StreamService } from '../../sell/services/stream.service';
import { TasPipe } from '../../../common/pipes/converter.pipe';

@Component({
  selector: 'subscription-add',
  templateUrl: './subscription.add.component.html',
  styleUrls: [ './subscription.add.component.scss' ]
})
export class SubscriptionAddComponent {
    user:any;
    stream: any;
    isDataAvailable:boolean = false;
    public showLoading: boolean;
    public form: FormGroup;
    constructor(
        private authService: AuthService,
        private subscriptionService: SubscriptionSrvice,
        private streamService: StreamService,
        private router: Router,
        private route: ActivatedRoute,
        private fb: FormBuilder,
        public http: HttpClient,
        private dialogService: MdlDialogService,
        private mdlSnackbarService: MdlSnackbarService,
        private mdlDialogService: MdlDialogOutletService,
        private vcRef: ViewContainerRef,
        private tasPipe: TasPipe
     ) {
         this.mdlDialogService.setDefaultViewContainerRef(this.vcRef);
    }

    ngOnInit() {
        this.form = this.fb.group({
                hours: ['', [<any>Validators.required]]
        });
        // Fetch stream
        let id = this.route.snapshot.params['id'];
        this.streamService.getStream(id).subscribe(
                (result: any) => {
                    this.stream = result.data;
                    this.isDataAvailable = true;
                },
                err => { console.log(err) }
              );

    }

    onSubmit(form, isValid: Boolean) {
        if(isValid) {
            // Confirm dialog
            const mitasPrice = this.stream.price * form.hours;
            const tasPrice = this.tasPipe.transform(mitasPrice)
            let confirmMsg = `Your will be charged ${(tasPrice)} TAS. Do you Agree ?`
            let result = this.dialogService.confirm(confirmMsg, 'Cancel', 'Yes, charge it!');
            result.subscribe( () => {
                this.showLoading = true;
                form.id = this.route.snapshot.params['id'];
                this.subscriptionService.add(form).subscribe(
                    response => {
                        let successMsg = `Success! You now have access to ${this.stream.name} in next ${form.hours} hours`
                        let result = this.dialogService.confirm(successMsg);
                        this.showLoading = false;
                        // Get current User and update balance
                        this.authService.getCurrentUser().subscribe(
                            data =>  {
                                data.balance = data.balance - mitasPrice;
                                this.authService.setCurrentUser(data);
                                this.router.navigate(['/dashboard/me']);
                            },
                            err => {
                                console.log(err);
                            }
                        );
                    },
                    err => {
                        this.showLoading = false;
                        let error = err.error
                        let errMsg = `Ups...Error was occured: <br> <strong>${error.error}</strong>`
                        let resultError = this.dialogService.confirm(errMsg, "");
                        resultError.subscribe( () => {
                            this.router.navigate(['/dashboard/buy/map']);
                        });
                     });
                },
                (err: any) => {
                    console.log('declined');
                    this.showLoading = false;
                }
            );
        }
    }
}
