import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service'; import { FormGroup, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms'; import { BalanceService } from '../balance.service';
import { Balance } from '../../../common/interfaces/balance.interface';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { MitasPipe } from '../../../common/pipes/converter.pipe';

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
    private mitasPipe: MitasPipe,
    public  alertService: AlertService,
  ){}

  ngOnInit() {
    this.form = this.formBuilder.group({
      'amount':        ['', Validators.required],
    });
    this.processing = false;
  }

  onSubmit(form, isValid: boolean) {
    this.errorMsg = null;
    if(isValid) {
      this.processing = true;
      // Convert to mTAS
      const toMiTasAmount =  {
        amount: this.mitasPipe.transform(form.amount),
      }
      this.balanceService.buy(toMiTasAmount).subscribe(
        response => {
          this.balanceUpdate.emit('update');
          this.alertService.success(` You successfully transfer ${this.form.value.amount} TAS to your account`);
        },
        err => {
          this.processing = false;
          this.alertService.error(`Something went wrong. Please try again later.`);
        }
      );
    }
  }
}
