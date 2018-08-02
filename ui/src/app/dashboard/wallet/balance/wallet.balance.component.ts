import { Component, OnInit } from '@angular/core';

import { AuthService } from '../../../auth/services/auth.service';
import { User } from '../../../common/interfaces/user.interface';
import {WalletAddComponent} from '../add'
import { BsModalService } from 'ngx-bootstrap/modal';

@Component({
  selector: 'user-wallet-balance',
  templateUrl: './wallet.balance.component.html',
  styleUrls: ['./wallet.balance.component.scss']
})
export class WalletBalanceComponent implements OnInit {

  user: any;
  newBalance: number;
  subscription: any;
  constructor(
    private AuthService: AuthService,
    private modalService: BsModalService,

) {
}

    ngOnInit() {
        this.user = {};
        this.AuthService.getCurrentUser().subscribe(
            data =>  {
                this.user = data
            },
            err => {
                console.log(err);
            });
    }

    // Open BUY tokens dialog
    public onBuyTokensClick() {
      // Open DashboardSellAddComponent Modal
        this.modalService.show(WalletAddComponent);

          // let pDialog = this.dialogService.showCustomDialog({
          //   component: WalletAddComponent,
          //   isModal: true,
          //   styles: {'width': '350px'},
          //   clickOutsideToClose: true,
          //   enterTransitionDuration: 400,
          //   leaveTransitionDuration: 400
          // });
          // pDialog.subscribe( dialogRef => {
          //     dialogRef.onHide().subscribe(data => {
          //         // Check if balance is updated
          //         if(data) {
          //             this.user.balance = this.user.balance + data;
          //         }
          //   });
          // });
    }
  }
