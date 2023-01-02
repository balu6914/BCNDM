import { Component, Output, EventEmitter, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';
import { StreamService } from 'app/common/services/stream.service';
import { Stream } from 'app/common/interfaces/stream.interface';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/common/validators/patterns';

@Component({
  selector: 'dpc-dashboard-ai-edit',
  templateUrl: './dashboard.ai.edit.component.html',
  styleUrls: ['./dashboard.ai.edit.component.scss']
})
export class DashboardAiEditComponent implements OnInit {
  form: FormGroup;
  editData: Stream;
  streamID: any;
  ownerID: any;
  submitted = false;

  @Output() streamEdited: EventEmitter<any> = new EventEmitter();
  constructor(
    private streamService: StreamService,
    private midpcPipe: MidpcPipe,
    private formBuilder: FormBuilder,
    public modalEditStream: BsModalRef,
    public alertService: AlertService,
  ) {
    const floatValidator = Validators.pattern(floatRegEx);
    const urlValidator = Validators.pattern(urlRegEx);

    this.form = this.formBuilder.group(
      {
        visibility:  ['', [Validators.required]],
        type:        ['', [Validators.required]],
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

  ngOnInit() {
    this.editData.snippet = this.editData.snippet || '';
    this.editData.metadata = this.editData.metadata || null;
    this.form.setValue(this.editData);
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

    if (this.form.valid) {
      const stream: Stream = {
        visibility: this.form.value.visibility,
        name: this.form.value.name,
        type: this.editData.type,
        description: this.form.value.description,
        snippet: this.form.value.snippet,
        url: "",
        encodedURL: this.streamService.encodeURL(this.form.value.url),
        terms: "",
        encodedTerms: this.streamService.encodeURL(this.form.value.terms),
        price: this.midpcPipe.transform(this.form.value.price),
        location: {
          'type': 'Point',
          'coordinates': [
            parseFloat(this.form.value.long || '0'),
            parseFloat(this.form.value.lat || '0')
          ]
        },
        metadata: JSON.parse(this.form.value.metadata),
      };

      this.streamService.updateStream(this.streamID, stream).subscribe(
        res => {
          // Set this parameters to configure table rows
          stream.id = this.streamID;
          stream.owner = this.ownerID;
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
