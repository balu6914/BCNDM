import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/shared/validators/patterns';
import { BsModalRef } from 'ngx-bootstrap'
import { BigQuery, Stream } from '../../../common/interfaces/stream.interface';
import { MitasPipe } from '../../../common/pipes/converter.pipe';
import { StreamService } from '../../../common/services/stream.service';

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
    private mitasPipe: MitasPipe,
    private formBuilder: FormBuilder,
    public modalAddStream: BsModalRef,
    public alertService: AlertService
  ) { }

  ngOnInit() {
    const floatValidator = Validators.pattern(floatRegEx);
    const urlValidator = Validators.pattern(urlRegEx);
    this.form = this.formBuilder.group({
      name: ['', [Validators.required, Validators.maxLength(32)]],
      type: ['', [Validators.required, Validators.maxLength(32)]],
      description: ['', [Validators.required, Validators.maxLength(256)]],
      url: ['', [Validators.required, Validators.maxLength(2048), urlValidator]],
      price: ['', [Validators.required, Validators.maxLength(9), floatValidator]],
      lat: [
        '',
        [
          Validators.required,
          Validators.maxLength(11),
          floatValidator,
          Validators.min(-90),
          Validators.max(90)
        ]
      ],
      long: [
        '',
        [
          Validators.required,
          Validators.maxLength(12),
          floatValidator,
          Validators.min(-180),
          Validators.max(180)
        ]
      ],
      snippet: ['', [Validators.maxLength(256)]],
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
    this.bqTouched = true;
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
        name: this.form.get('name').value,
        type: this.form.get('type').value,
        description: this.form.get('description').value,
        price: this.mitasPipe.transform(this.form.get('price').value),
        snippet: this.form.get('snippet').value,
        location: {
          type: 'Point',
          coordinates: [
            parseFloat(this.form.get('long').value),
            parseFloat(this.form.get('lat').value)
          ]
        },
        external: this.external
      };
      if (this.external) {
        stream.bq = new BigQuery(
          this.form.get('project').value,
          this.form.get('dataset').value,
          this.form.get('table').value,
          this.form.get('fields').value
        );
      } else {
        stream.url = this.form.get('url').value;
      }

      // Send addStream request
      this.streamService.addStream(stream).subscribe(
        res => {
          // Add ID from http response to stream
          stream.id = res['id'];
          this.streamCreated.emit(stream);
          this.alertService.success(`Stream successfully added!`);
        },
        err => {
          if (err.status === 401) {
            this.alertService.error('You don\'t have permission to add this Stream.');
            return;
          }
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );

      this.modalAddStream.hide();
    }
  }
}
