import { Component, Output, EventEmitter, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';

import { StreamService } from '../../../common/services/stream.service';
import { Stream } from '../../../common/interfaces/stream.interface';
import { MitasPipe } from '../../../common/pipes/converter.pipe';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/shared/validators/patterns';

@Component({
  selector: 'dpc-dashboard-sell-edit',
  templateUrl: './dashboard.sell.edit.component.html',
  styleUrls: ['./dashboard.sell.edit.component.scss']
})
export class DashboardSellEditComponent implements OnInit {
  form: FormGroup;
  editData: any;
  streamID: any;
  submitted = false;

  @Output() streamEdited: EventEmitter<any> = new EventEmitter();
  constructor(
    private streamService: StreamService,
    private mitasPipe: MitasPipe,
    private formBuilder: FormBuilder,
    public modalEditStream: BsModalRef,
    public alertService: AlertService,
  ) {
    const floatValidator = Validators.pattern(floatRegEx);
    const urlValidator = Validators.pattern(urlRegEx);

    this.form = this.formBuilder.group({
      'name':        ['', [Validators.required, Validators.maxLength(32)]],
      'type':        ['', [Validators.required, Validators.maxLength(32)]],
      'description': ['', [Validators.required, Validators.maxLength(256)]],
      'url':         ['', [Validators.required, Validators.maxLength(2048), urlValidator]],
      'price':       ['', [Validators.required, Validators.maxLength(9), floatValidator]],
      'lat':         ['', [Validators.required, Validators.maxLength(11), floatValidator, Validators.min(-90), Validators.max(90)]],
      'long':        ['', [Validators.required, Validators.maxLength(12), floatValidator, Validators.min(-180), Validators.max(180)]],
      'snippet':     ['', [Validators.maxLength(256)]]
    });
  }

  ngOnInit() {
    this.editData.snippet =  this.editData.snippet || '';
    this.form.setValue(this.editData);
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const stream: Stream = {
        name: this.form.value.name,
        type: this.form.value.type,
        description: this.form.value.description,
        snippet: this.form.value.snippet,
        url: this.form.value.url,
        price: this.mitasPipe.transform(this.form.value.price),
        location: {
          'type': 'Point',
          'coordinates': [
            parseFloat(this.form.value.long),
            parseFloat(this.form.value.lat)
          ]
        }
      };

      // Send addStream request
      this.streamService.updateStream(this.streamID, stream).subscribe(
        res => {
          stream.id = this.streamID;
          this.streamEdited.emit(stream);
          this.alertService.success(`Stream successfully updated!`);
        },
        err => {
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );

      this.modalEditStream.hide();
    }
  }
}
