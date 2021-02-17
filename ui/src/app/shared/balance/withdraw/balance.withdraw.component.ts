import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';
import { BalanceService } from 'app/shared/balance/balance.service';
import { Balance } from 'app/common/interfaces/balance.interface';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
import { User } from 'app/common/interfaces/user.interface';

@Component({
  selector: 'dpc-balance-withdraw',
  templateUrl: './balance.withdraw.component.html',
  styleUrls: ['./balance.withdraw.component.scss']
})
export class BalanceWithdrawComponent implements OnInit {
  form: FormGroup;
  processing: Boolean;
  user: User;

  @Output()
  // Emit event when we successfully sold some token , to get updated balance.
  balanceUpdate = new EventEmitter();

  constructor(
    public  modalAddTokens: BsModalRef,
    private balanceService: BalanceService,
    private formBuilder: FormBuilder,
    private midpcPipe: MidpcPipe,
    public  alertService: AlertService,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      'amount':        ['', Validators.required],
    });
    this.processing = false;
  }

  onSubmit(form, isValid: boolean) {
    if (isValid) {
      this.processing = true;
      // Convert to miDPC
      const toMidpcAmount =  {
        amount: this.midpcPipe.transform(form.amount),
        fund_id: this.user.id,
      };
      this.balanceService.withdraw(toMidpcAmount).subscribe(
        response => {
          this.balanceUpdate.emit('update');
          this.alertService.success(` You successfully Withdraw ${this.form.value.amount} DPC from your account`);
        },
        err => {
          this.processing = false;
          this.alertService.error(`Something went wrong. Please try again later.`);
        }
      );
    }
  }
}
