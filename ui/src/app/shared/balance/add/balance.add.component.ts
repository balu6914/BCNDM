import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';
import { BalanceService } from 'app/shared/balance/balance.service';
import { Balance } from 'app/common/interfaces/balance.interface';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';

@Component({
  selector: 'dpc-balance-add',
  templateUrl: './balance.add.component.html',
  styleUrls: ['./balance.add.component.scss']
})
export class BalanceAddComponent implements OnInit {
  form: FormGroup;
  errorMsg: String;
  processing: Boolean;
  @Output()
  // Emit event when we successfully buy more token , to get updated balance.
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
      'amount': ['', Validators.required],
    });
    this.processing = false;
  }

  onSubmit(form, isValid: boolean) {
    this.errorMsg = null;
    if (isValid) {
      this.processing = true;
      // Convert to miDPC
      const toMidpcAmount =  {
        amount: this.midpcPipe.transform(form.amount),
      };
      this.balanceService.buy(toMidpcAmount).subscribe(
        response => {
          this.modalAddTokens.hide();
          this.balanceUpdate.emit('update');
          this.alertService.success(` You successfully transfer ${this.form.value.amount} DPC to your account`);
        },
        err => {
          this.modalAddTokens.hide();
          this.processing = false;
          this.alertService.error(`Something went wrong. Please try again later.`);
        }
      );
    }
  }
}
