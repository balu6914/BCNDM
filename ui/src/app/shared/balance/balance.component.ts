import { Component, OnInit, Input } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { Balance } from './balance';
import { BalanceService } from './balance.service';
import { BalanceAddComponent } from './add/balance.add.component';

@Component({
  selector: 'dpc-balance-widget',
  templateUrl: './balance.component.html',
  styleUrls: ['./balance.component.scss']
})
export class BalanceComponent implements OnInit {
  balance: Balance;
  modalRef: BsModalRef;

  constructor(
    private modalService: BsModalService,
    private balanceService: BalanceService,
  ) {}

    // Open BUY tokens dialog
    onBuyTokensClick() {
        this.modalRef = this.modalService.show(BalanceAddComponent);
        // Listen to balance update event
        this.modalRef.content.balanceUpdate.subscribe(e => {
          // Fetch updated user balance
          this.getBalance().then(
          (response) => {
              this.modalRef.hide()
          },
        )
        });
    }

  ngOnInit() {
    // TODO: Remove this Mock of user balance its tmp hack until we add this on API side.
    let mockBalance = {
      amount: 0,
      symbol: 'TAS',
      fiatAmount: 1200,
      fiatSymbol: 'USD'
    }
    this.balance = mockBalance;
    this.getBalance();
  }

  getBalance() {
    let promise = new Promise((resolve, reject) => {
      this.balanceService.get().subscribe(
        (result: any) => {
          this.balance.amount = result.balance;
          resolve();
        },
        err => {
          console.error("Error fetching user balance ", err)
          reject();
        });
    });
    return promise;
  }

}
