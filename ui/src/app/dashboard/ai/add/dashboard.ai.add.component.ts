import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { floatRegEx, urlRegEx } from 'app/common/validators/patterns';
import { BsModalRef } from 'ngx-bootstrap';
import { Stream } from 'app/common/interfaces/stream.interface';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
import { StreamService } from 'app/common/services/stream.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dpc-dashboard-ai-add',
  templateUrl: './dashboard.ai.add.component.html',
  styleUrls: ['./dashboard.ai.add.component.scss']
})
export class DashboardAiAddComponent implements OnInit {
  form: FormGroup;
  submitted = false;
  streamType = '';
  ownerID: any;

  @Output()
  aiStreamCreated: EventEmitter<any> = new EventEmitter();
  constructor(
    private streamService: StreamService,
    private midpcPipe: MidpcPipe,
    private formBuilder: FormBuilder,
    public modalAddAiStream: BsModalRef,
    public alertService: AlertService
  ) { }

  ngOnInit() {
    const floatValidator = Validators.pattern(floatRegEx);
    const urlValidator = Validators.pattern(urlRegEx);
    this.form = this.formBuilder.group(
      {
        visibility:  ['', [Validators.required]],
        name:        ['', [Validators.required, Validators.maxLength(32)]],
        description: ['', [Validators.required, Validators.maxLength(256)]],
        url:         ['', [Validators.maxLength(2048)]],
        terms:       ['', [Validators.required, Validators.maxLength(2048), urlValidator]],
        price:       ['', [Validators.required, Validators.maxLength(9), floatValidator]],
        lat:         ['', [Validators.maxLength(11), Validators.min(-90), Validators.max(90), floatValidator]],
        long:        ['', [Validators.maxLength(12), Validators.min(-180), Validators.max(180), floatValidator]],
        snippet:     ['', [Validators.maxLength(2048)]],
        metadata:    [''],
      },
      {
        validator: this.metadataValidator,
      },
    );
  }

  metadataValidator(fg: FormGroup) {
    try {
      JSON.parse(fg.value.metadata);
    } catch (e) {
      fg.controls.metadata.setErrors({'invalid': true});
      return;
    }

    return null;
  }

  onSubmit() {
    this.submitted = true;
    const d = new Date;
    const ts = d.valueOf();

    if (this.form.valid) {
      const stream: Stream = {
        visibility: this.form.value.visibility,
        name: this.form.value.name,
        type: this.streamType,
        description: this.form.value.description,
        snippet: this.form.value.snippet,
        url: encodeURIComponent(this.form.value.url || `http://www.localhost.com/${ts}`),
        price: this.midpcPipe.transform(this.form.value.price),
        location: {
          type: 'Point',
          coordinates: [
            parseFloat(this.form.value.long || '0'),
            parseFloat(this.form.value.lat || '0')
          ]
        },
        metadata: JSON.parse(this.form.value.metadata),
        terms: encodeURIComponent(this.form.value.terms),
      };

      this.streamService.addStream(stream).subscribe(
        res => {
          // Set this parameters to configure table rows
          stream.id = res['id'];
          stream.owner = this.ownerID;
          this.aiStreamCreated.emit(stream);
          this.alertService.success(`${this.streamType} successfully added!`);
        },
        err => {
          console.log(err);
          this.alertService.error(`Status: ${err.status} - ${err.error.error}`);
        }
      );

      this.modalAddAiStream.hide();
    }
  }
}
