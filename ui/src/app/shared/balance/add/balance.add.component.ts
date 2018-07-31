import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';
import { FormGroup, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { BalanceService } from '../balance.service';
import { Balance } from '../../../common/interfaces/balance.interface';

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
  ){}

  ngOnInit() {
    this.form = this.formBuilder.group({
      'amount':        ['', Validators.required],
    });
    this.processing = false;
  }

  onSubmit(model: Balance, isValid: boolean) {
    this.errorMsg = null;
    if(isValid) {
      this.processing = true;
      this.balanceService.buy(model).subscribe(
        response => {
          this.balanceUpdate.emit('update');
          this.processing = false;
        },
        err => {
          this.errorMsg = err;
          this.processing = false;
        }
      )
    }
  }
}
