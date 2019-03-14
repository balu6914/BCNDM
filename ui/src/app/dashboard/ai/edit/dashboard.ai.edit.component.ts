import { Component, Output, EventEmitter, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap';
import { StreamService } from 'app/common/services/stream.service';
import { Stream } from 'app/common/interfaces/stream.interface';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/shared/validators/patterns';

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
        url:         ['', [Validators.required, Validators.maxLength(2048), urlValidator]],
        price:       ['', [Validators.required, Validators.maxLength(9), floatValidator]],
        lat:         ['', [Validators.required, Validators.maxLength(11), Validators.min(-90), Validators.max(90), floatValidator]],
        long:        ['', [Validators.required, Validators.maxLength(12), Validators.min(-180), Validators.max(180), floatValidator]],
        snippet:     ['', [Validators.maxLength(2048)]],
        metadata:    ['', [Validators.maxLength(2048)]],
      },
      {
        validator: this.metadataValidator,
      },
    );
  }

  ngOnInit() {
    this.editData.snippet =  this.editData.snippet || null;
    this.editData.metadata =  this.editData.metadata || null;
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
        url: this.form.value.url,
        price: this.midpcPipe.transform(this.form.value.price),
        location: {
          'type': 'Point',
          'coordinates': [
            parseFloat(this.form.value.long),
            parseFloat(this.form.value.lat)
          ]
        },
        metadata: JSON.parse(this.form.value.metadata),
      };

      // Send addStream request
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