import { Component, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { MitasPipe } from 'app/common/pipes/converter.pipe';

@Component({
  selector: 'dashboard-contracts-add',
  templateUrl: './dashboard.contracts.add.component.html',
  styleUrls: [ './dashboard.contracts.add.component.scss' ]
})
export class DashboardContractsAddComponent {
  form: FormGroup;
  modalMsg: string;
  submitted: boolean = false;

  @Output() contractCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private mitasPipe: MitasPipe,
    private formBuilder: FormBuilder,
    public  modalNewContract: BsModalRef,
  ) {
    this.form = this.formBuilder.group({
      streamName:     ['', [<any>Validators.required]],
      streamPrice:    ['', [<any>Validators.required]],
      parties:        ['', [<any>Validators.required]],
      shareOffered:   ['', [<any>Validators.required]],
      expirationDate: ['', [<any>Validators.required]]
    });
  }

  onSubmit() {
    const contract = {
      streamName: this.form.value.streamName,
      streamPrice: this.mitasPipe.transform(this.form.value.streamPrice),
      parties: this.form.value.parties,
      shareOffered: this.form.value.shareOffered,
      expirationDate: this.form.value.expirationDate
    }

    // TODO: Send request to transactions service to create th contract
    this.modalMsg = `Contract succesfully created!`;
    this.contractCreated.emit(contract);

    // Hide modalNewStream and show modalResponse
    this.submitted = true;
  }
}
