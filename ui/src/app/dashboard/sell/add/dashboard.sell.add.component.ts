import { Component, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { StreamService } from '../../../common/services/stream.service';
import { Stream } from '../../../common/interfaces/stream.interface';
import { MitasPipe } from '../../../common/pipes/converter.pipe';
import { AlertService } from 'app/shared/alerts/services/alert.service';

@Component({
  selector: 'dashboard-sell-add',
  templateUrl: './dashboard.sell.add.component.html',
  styleUrls: [ './dashboard.sell.add.component.scss' ]
})
export class DashboardSellAddComponent {
    form: FormGroup;

    @Output() streamCreated: EventEmitter<any> = new EventEmitter();
    constructor(
        private streamService: StreamService,
        private mitasPipe: MitasPipe,
        private formBuilder: FormBuilder,
        public  modalNewStream: BsModalRef,
        public  alertService: AlertService,
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
      this.streamService.addStream(stream).subscribe(
        res => {
          // Add ID from http response to stream
          stream.id = res.id;
          this.streamCreated.emit(stream);
          this.modalNewStream.hide();
          this.alertService.success(`Stream succesfully added!`);
        },
        err => {
          this.modalNewStream.hide();
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
      });
    }
}
