import { Component, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { StreamService } from '../../../common/services/stream.service';
import { Stream } from '../../../common/interfaces/stream.interface';
import { MitasPipe } from '../../../common/pipes/converter.pipe';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/shared/validators/patterns';

@Component({
  selector: 'dpc-dashboard-sell-add',
  templateUrl: './dashboard.sell.add.component.html',
  styleUrls: ['./dashboard.sell.add.component.scss']
})
export class DashboardSellAddComponent {
  form: FormGroup;
  submitted = false;
  bigQuery = false;

  @Output() streamCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private streamService: StreamService,
    private mitasPipe: MitasPipe,
    private formBuilder: FormBuilder,
    public modalAddStream: BsModalRef,
    public alertService: AlertService,
  ) {
    const floatValidator = Validators.pattern(floatRegEx);
    const urlValidator = Validators.pattern(urlRegEx);

    this.form = this.formBuilder.group({
      'name':        ['name', [Validators.required, Validators.maxLength(32)]],
      'type':        ['type', [Validators.required, Validators.maxLength(32)]],
      'description': ['desc', [Validators.required, Validators.maxLength(256)]],
      'url':         ['http://www.com.com', [Validators.required, Validators.maxLength(128), urlValidator]],
      'price':       ['22', [Validators.required, Validators.maxLength(9), floatValidator]],
      'lat':         ['22', [Validators.required, Validators.maxLength(11), floatValidator, Validators.min(-90), Validators.max(90)]],
      'long':        ['45.566', [Validators.required, Validators.maxLength(12), floatValidator, Validators.min(-180), Validators.max(180)]],
      'snippet':     ['{my:stream}', [Validators.maxLength(256)]],
      'project':     [{value: 'digisense-196722', disabled: true}, [Validators.required, Validators.maxLength(256)]],
      'dataset':     [{value: 'demo', disabled: true}, [Validators.required, Validators.maxLength(256)]],
      'table':       [{value: 'aqi', disabled: true}, [Validators.required, Validators.maxLength(256)]],
      'fields':      [{value: 'Time,Value', disabled: true}, [Validators.required, Validators.maxLength(256)]]
    });
  }


  changeBQ() {
    this.bigQuery = !this.bigQuery;
    this.bigQuery ? this.form.get('project').enable() : this.form.get('project').disable();
    this.bigQuery ? this.form.get('dataset').enable() : this.form.get('dataset').disable();
    this.bigQuery ? this.form.get('table').enable() : this.form.get('table').disable();
    this.bigQuery ? this.form.get('fields').enable() : this.form.get('fields').disable();
  }


  onSubmit() {
    this.submitted = true;
    if (this.form.valid) {
      const stream: Stream = {
        name:        this.form.get('name').value,
        type:        this.form.get('type').value,
        description: this.form.get('description').value,
        url:         this.form.get('url').value,
        price:       this.mitasPipe.transform(this.form.get('price').value),
        snippet:     this.form.get('snippet').value,
        location: {
          'type': 'Point',
          'coordinates': [
            parseFloat(this.form.get('long').value),
            parseFloat(this.form.get('lat').value)
          ]
        },
        bq: this.bigQuery
      };
      if (this.bigQuery) {
        stream.project = this.form.get('project').value;
        stream.dataset = this.form.get('dataset').value;
        stream.table = this.form.get('table').value;
        stream.fields = this.form.get('fields').value;
      }

      console.log(stream);
      // Send addStream request
      this.streamService.addStream(stream).subscribe(
        res => {
          // Add ID from http response to stream
          stream.id = res['id'];
          this.streamCreated.emit(stream);
          this.alertService.success(`Stream successfully added!`);
        },
        err => {
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );

      this.modalAddStream.hide();
    }
  }
}
