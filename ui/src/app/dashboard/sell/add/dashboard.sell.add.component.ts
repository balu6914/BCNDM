import { Component, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { StreamService } from '../../../common/services/stream.service';
import { Stream } from '../../../common/interfaces/stream.interface';
import { MitasPipe } from '../../../common/pipes/converter.pipe';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-sell-add',
  templateUrl: './dashboard.sell.add.component.html',
  styleUrls: ['./dashboard.sell.add.component.scss']
})
export class DashboardSellAddComponent {
  form: FormGroup;
  submitted: boolean = false

  @Output() streamCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private streamService: StreamService,
    private mitasPipe: MitasPipe,
    private formBuilder: FormBuilder,
    public modalNewStream: BsModalRef,
    public alertService: AlertService,
  ) {
    const floatValidator = Validators.pattern('[-+]?([0-9]\.[0-9]+|[0-9]+)');
    const urlValidator = Validators.pattern('(https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9]\.[^\s]{2,})');

    this.form = this.formBuilder.group({
      'name':        ['', [Validators.required, Validators.maxLength(16)]],
      'type':        ['', [Validators.required, Validators.maxLength(32)]],
      'description': ['', [Validators.required, Validators.maxLength(256)]],
      'url':         ['', [Validators.required, Validators.maxLength(128), urlValidator]],
      'price':       ['', [Validators.required, Validators.maxLength(9), floatValidator]],
      'lat':         ['', [Validators.required, Validators.maxLength(9), floatValidator, Validators.min(-90), Validators.max(90)]],
      'long':        ['', [Validators.required, Validators.maxLength(9), floatValidator, Validators.min(-180), Validators.max(180)]],
      'snippet':     ['', [Validators.maxLength(256)]]
    });
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const stream: Stream = {
        name:        this.form.value.name,
        type:        this.form.value.type,
        description: this.form.value.description,
        url:         this.form.value.url,
        price:       this.mitasPipe.transform(this.form.value.price),
        snippet:     this.form.value.snippet || "",
        location: {
          'type': 'Point',
          'coordinates': [
            parseFloat(this.form.value.long),
            parseFloat(this.form.value.lat)
          ]
        }
      };

      // Send addStream request
      this.streamService.addStream(stream).subscribe(
        res => {
          // Add ID from http response to stream
          stream.id = res['id'];
          this.streamCreated.emit(stream);
          this.modalNewStream.hide();
          this.alertService.success(`Stream successfully added!`);
        },
        err => {
          this.modalNewStream.hide();
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );
    }
  }
}
