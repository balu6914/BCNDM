import { Component, OnInit, ViewContainerRef } from '@angular/core';
import { MdlDialogService, MdlDialogReference, MdlDialogOutletService } from '@angular-mdl/core';

import { AuthService } from '../../../auth/services/auth.service';
import { User } from '../../../common/interfaces/user.interface';
import {WalletAddComponent} from '../add'

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
    private dialogService: MdlDialogService,
    private mdlDialogService: MdlDialogOutletService,
    private vcRef: ViewContainerRef
) {
    this.mdlDialogService.setDefaultViewContainerRef(this.vcRef);
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
    public onBuyTokensClick($event: MouseEvent) {
          let pDialog = this.dialogService.showCustomDialog({
            component: WalletAddComponent,
            isModal: true,
            styles: {'width': '350px'},
            clickOutsideToClose: true,
            enterTransitionDuration: 400,
            leaveTransitionDuration: 400
          });
          pDialog.subscribe( dialogRef => {
              dialogRef.onHide().subscribe(data => {
                  // Check if balance is updated
                  if(data) {
                      this.user.balance = this.user.balance + data;
                  }
            });
          });
    }
  }
