import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/common/validators/patterns';
import { BsModalRef } from 'ngx-bootstrap';
import { BigQuery, Stream } from 'app/common/interfaces/stream.interface';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
import { StreamService } from 'app/common/services/stream.service';

@Component({
  selector: 'dpc-dashboard-sell-add',
  templateUrl: './dashboard.sell.add.component.html',
  styleUrls: ['./dashboard.sell.add.component.scss']
})
export class DashboardSellAddComponent implements OnInit {
  form: FormGroup;
  submitted = false;
  external = false;
  bqTouched = false;
  bqMail = false;

  @Output()
  streamCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private streamService: StreamService,
    private midpcPipe: MidpcPipe,
    private formBuilder: FormBuilder,
    public modalAddStream: BsModalRef,
    public alertService: AlertService
  ) { }

  ngOnInit() {
    const floatValidator = Validators.pattern(floatRegEx);
    const urlValidator = Validators.pattern(urlRegEx);
    this.form = this.formBuilder.group({
      visibility:  ['', [Validators.required]],
      name:        ['', [Validators.required, Validators.maxLength(256)]],
      type:        ['', [Validators.required, Validators.maxLength(32)]],
      description: ['', [Validators.required, Validators.maxLength(2048)]],
      url:         ['', [Validators.required, Validators.maxLength(2048)]],
      terms:       ['', [Validators.required, Validators.maxLength(2048), urlValidator]],
      price:       ['', [Validators.required, Validators.maxLength(9), floatValidator]],
      lat:         ['', [Validators.required, Validators.maxLength(11), floatValidator, Validators.min(-90), Validators.max(90)]],
      long:        ['', [Validators.required, Validators.maxLength(12), floatValidator, Validators.min(-180), Validators.max(180)]],
      snippet:     ['', [Validators.maxLength(2048)]],
      project: [{ value: '', disabled: true }, [this.requiredBQ.bind(this)]],
      dataset: [{ value: '', disabled: true }, [this.requiredBQ.bind(this)]],
      table:   [{ value: '', disabled: true }, [this.requiredBQ.bind(this)]],
      fields:  [{ value: '', disabled: true }, [this.requiredBQ.bind(this)]]
    });
  }
  requiredBQ(c: FormControl) {
    if (this.external && c.value === '') {
      return {required: true};
    }
    return null;
  }

  changeExt() {
    this.bqTouched = !this.bqTouched;
    if (!this.bqMail) {
      return;
    }

    this.external = !this.external;
    this.external
      ? this.form.get('project').enable()
      : this.form.get('project').disable();
    this.external
      ? this.form.get('dataset').enable()
      : this.form.get('dataset').disable();
    this.external
      ? this.form.get('table').enable()
      : this.form.get('table').disable();
    this.external
      ? this.form.get('fields').enable()
      : this.form.get('fields').disable();
     !this.external
      ? this.form.get('url').enable()
      : this.form.get('url').disable();
  }

  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const stream: Stream = {
        visibility: this.form.value.visibility,
        name: this.form.value.name,
        type: this.form.value.type,
        description: this.form.value.description,
        price: this.midpcPipe.transform(this.form.value.price),
        snippet: this.form.value.snippet,
        location: {
          type: 'Point',
          coordinates: [
            parseFloat(this.form.value.long),
            parseFloat(this.form.value.lat)
          ]
        },
        external: this.external,
        terms: encodeURIComponent(this.form.value.terms),
      };
      if (this.external) {
        stream.bq = new BigQuery(
          this.form.value.project,
          this.form.value.dataset,
          this.form.value.table,
          this.form.value.fields
        );
      } else {
        stream.url = encodeURIComponent(this.form.value.url);
      }

      this.streamService.addStream(stream).subscribe(
        res => {
          // Add ID from http response to stream
          stream.id = res['id'];
          this.streamCreated.emit(stream);
          this.alertService.success(`Stream successfully added!`);
          this.modalAddStream.hide();
        },
        err => {
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );
    }
  }
}
