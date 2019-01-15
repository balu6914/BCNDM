import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';
import { BalanceService } from '../balance.service';
import { Balance } from '../../../common/interfaces/balance.interface';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { MitasPipe } from '../../../common/pipes/converter.pipe';

@Component({
  selector: 'dpc-balance-withdraw',
  templateUrl: './balance.withdraw.component.html',
  styleUrls: ['./balance.withdraw.component.scss']
})
export class BalanceWithdrawComponent implements OnInit {
  form: FormGroup;
  processing: Boolean;
  @Output()
  // Emit event when we successfully sold some token , to get updated balance.
  balanceUpdate = new EventEmitter();

  constructor(
    public  modalAddTokens: BsModalRef,
    private balanceService: BalanceService,
    private formBuilder: FormBuilder,
    private mitasPipe: MitasPipe,
    public  alertService: AlertService,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      'amount':        ['', Validators.required],
    });
    this.processing = false;
  }

  onSubmit(form, isValid: boolean) {
    if(isValid) {
      this.processing = true;
      // Convert to mTAS
      const toMiTasAmount =  {
        amount: this.mitasPipe.transform(form.amount),
      }
      this.balanceService.withdraw(toMiTasAmount).subscribe(
        response => {
          this.balanceUpdate.emit('update');
          this.alertService.success(` You successfully Withdraw ${this.form.value.amount} TAS from your account`);
        },
        err => {
          this.processing = false;
          this.alertService.error(`Something went wrong. Please try again later.`);
        }
      );
    }
  }



}
