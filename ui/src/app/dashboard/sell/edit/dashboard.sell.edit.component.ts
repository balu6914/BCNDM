import { Component } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { StreamService } from '../../../common/services/stream.service';
import { Stream } from '../../../common/interfaces/stream.interface';
import { MitasPipe } from '../../../common/pipes/converter.pipe';

@Component({
  selector: 'dashboard-sell-edit',
  templateUrl: './dashboard.sell.edit.component.html',
  styleUrls: [ './dashboard.sell.edit.component.scss' ]
})
export class DashboardSellEditComponent {
  user:any;
  form: FormGroup;
  modalMsg: string;
  submitted: boolean = false;
  formEdit: any;

  constructor(
      private streamService: StreamService,
      private mitasPipe: MitasPipe,
      private formBuilder: FormBuilder,
      public  modalNewStream: BsModalRef,
  ) {
    this.form = formBuilder.group({
      'name':        ['', [<any>Validators.required, <any>Validators.minLength(3)]],
      'type':        ['', Validators.required],
      'description': ['', [<any>Validators.required, <any>Validators.minLength(5)]],
      'url':         ['', Validators.required],
      'price':       ['', Validators.required],
      'lat':         ['', Validators.required],
      'long':        ['', Validators.required]
    })
  }

  ngOnInit() {
    this.form.controls.name.setValue(this.formEdit.name);
    this.form.controls.type.setValue(this.formEdit.type);
    this.form.controls.description.setValue(this.formEdit.description);
    this.form.controls.url.setValue(this.formEdit.url);
    this.form.controls.price.setValue(this.formEdit.price);
    this.form.controls.lat.setValue(this.formEdit.lat);
    this.form.controls.long.setValue(this.formEdit.long);
  }

  onSubmit() {
    const stream: Stream = {
      name: this.form.value.name,
      type: this.form.value.type,
      description: this.form.value.description,
      url: this.form.value.url,
      price: this.mitasPipe.transform(this.form.value.price),
      location: {
        "type": "Point",
        "coordinates": [parseFloat(this.form.value.long),
                        parseFloat(this.form.value.lat)]
      }
    };

    // Send addStream request
    this.streamService.updateStream(this.formEdit.id, stream).subscribe(
      response => {
        this.modalMsg = `Stream succesfully updated!`;
      },
      err => {
        console.log(err);
        this.modalMsg = `Status: ${err.status} - ${err.statusText}`;
    });

    // Hide modalNewStream and show modalResponse
    this.submitted = true;
  }
}
