import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { User } from 'app/common/interfaces/user.interface';
import { Balance } from 'app/common/interfaces/balance.interface';
import { BalanceService } from './balance.service';
import { BalanceAddComponent } from './add/balance.add.component';
import { BalanceWithdrawComponent } from './withdraw/balance.withdraw.component';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
import { AuthService } from 'app/auth/services/auth.service';

@Component({
  selector: 'dpc-balance-widget',
  templateUrl: './balance.component.html',
  styleUrls: ['./balance.component.scss']
})
export class BalanceComponent implements OnInit {
  balance = new Balance();
  modalRef: BsModalRef;
  user: User;

  @Input() showWalletKey: boolean;
  @Output() balanceUpdate = new EventEmitter();

  constructor(
    private authService: AuthService,
    private modalService: BsModalService,
    private balanceService: BalanceService,
    private midpcPipe: MidpcPipe,
  ) { }

  // Open BUY tokens dialog
  onBuyTokensClick() {
    const initialState = {
      user: this.user,
    };

    this.modalRef = this.modalService.show(BalanceAddComponent, {initialState});
    // Listen to balance update event
    this.modalRef.content.balanceUpdate.subscribe(e => {
      // Fetch updated user balance
      this.getBalance(this.user.id);
      this.modalRef.hide();
    });
  }
  // Open Withdraw tokens dialog
  onWithdrawTokensClick() {
    const initialState = {
      user: this.user,
    };

    this.modalRef = this.modalService.show(BalanceWithdrawComponent, {initialState});
    // Listen to balance update event
    this.modalRef.content.balanceUpdate.subscribe(e => {
      // Fetch updated user balance
      this.getBalance(this.user.id);
      this.modalRef.hide();
    });
  }

  ngOnInit() {
    // Get balance if used as modal or fetch user if not
    if (this.user !== undefined) {
      this.getBalance(this.user.id);
    } else {
      this.authService.getCurrentUser().subscribe(
        data => {
          this.user = data;
          this.getBalance(this.user.id);
        },
        err => {
          console.log(err);
        },
      );
    }
  }

  getBalance(userID) {
    this.balanceService.getBalance(userID).subscribe(
      (result: any) => {
        this.balance.amount = result.balance;
        // TODO remove this Mock of fiatAmount when we add this info on API side
        this.balance.fiatAmount = this.balance.amount;
        // Publish new balance data to balance message buss
        this.balanceService.changed(this.balance);

        this.balanceUpdate.emit(result.balance);
      },
      err => {
        console.error('Error fetching user balance ', err);
      });
  }
}
