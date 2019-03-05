import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { floatRegEx, urlRegEx } from 'app/shared/validators/patterns';
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
    this.form = this.formBuilder.group({
      visibility: ['', [Validators.required]],
      name: ['', [Validators.required, Validators.maxLength(32)]],
      description: ['', [Validators.required, Validators.maxLength(256)]],
      snippet: ['', [Validators.maxLength(256)]],
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
    });
  }


  onSubmit() {
    this.submitted = true;

    if (this.form.valid) {
      const stream: Stream = {
        visibility: this.form.get('visibility').value,
        name: this.form.get('name').value,
        type: this.streamType,
        description: this.form.get('description').value,
        url: this.form.get('url').value,
        price: this.midpcPipe.transform(this.form.get('price').value),
        snippet: this.form.get('snippet').value,
        location: {
          type: 'Point',
          coordinates: [
            parseFloat(this.form.get('long').value),
            parseFloat(this.form.get('lat').value)
          ]
        },
        external: false,
      };

      // Send addStream request
      this.streamService.addStream(stream).subscribe(
        res => {
          // Add ID from http response to stream
          stream.id = res['id'];
          this.aiStreamCreated.emit(stream);
          this.alertService.success(`${this.streamType} successfully added!`);
        },
        err => {
          if (err.status === 401) {
            this.alertService.error(`You don't have permission to add this ${this.streamType}.`);
            return;
          }
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );

      this.modalAddAiStream.hide();
    }
  }
}
