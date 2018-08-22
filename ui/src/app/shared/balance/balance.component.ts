import { Component, OnInit, Input } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { Balance } from '../../common/interfaces/balance.interface';
import { BalanceService } from './balance.service';
import { BalanceAddComponent } from './add/balance.add.component';
import { MitasPipe } from '../../common/pipes/converter.pipe';

@Component({
  selector: 'dpc-balance-widget',
  templateUrl: './balance.component.html',
  styleUrls: ['./balance.component.scss']
})
export class BalanceComponent implements OnInit {
  balance = new Balance();
  modalRef: BsModalRef;

  constructor(
    private modalService: BsModalService,
    private balanceService: BalanceService,
    private mitasPipe: MitasPipe,
  ) {}

    // Open BUY tokens dialog
    onBuyTokensClick() {
        this.modalRef = this.modalService.show(BalanceAddComponent);
        // Listen to balance update event
        this.modalRef.content.balanceUpdate.subscribe(e => {
          // Fetch updated user balance
          this.getBalance().then(
          (response) => {
              this.modalRef.hide();
          },
        )
        });
    }

  ngOnInit() {
    this.getBalance();
  }

  getBalance() {
    let promise = new Promise((resolve, reject) => {
      this.balanceService.get().subscribe(
        (result: any) => {
          this.balance.amount = result.balance;
          //TODO remove this Mock of fiatAmount when we add this info on API side
          this.balance.fiatAmount = this.balance.amount;
          // Publish new balance data to balance message buss
          this.balanceService.changed(this.balance);
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
