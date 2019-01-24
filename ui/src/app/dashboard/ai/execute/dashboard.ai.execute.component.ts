import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { floatRegEx, urlRegEx } from 'app/shared/validators/patterns';
import { BsModalRef } from 'ngx-bootstrap';
import { BigQuery, Stream } from 'app/common/interfaces/stream.interface';

@Component({
  selector: 'dpc-dashboard-ai-execute',
  templateUrl: './dashboard.ai.execute.component.html',
  styleUrls: ['./dashboard.ai.execute.component.scss']
})
export class DashboardAiExecuteComponent implements OnInit {
  form: FormGroup;
  submitted = false;
  streams = [];
  algos = [];

  @Output() executionStarted: EventEmitter<any> = new EventEmitter();
  constructor(
    public modalExecute: BsModalRef,
    private formBuilder: FormBuilder,
    public alertService: AlertService,
  ) { }

  ngOnInit() {
    // TODO. create request form
  }

  execute() {
    // TODO: Send request to execution  endpoint
    console.log('Execute Algorithm: ', this.algos);
    console.log('On Dataset: ', this.streams);
    this.alertService.success(`Execution succesfully started!`);
    this.modalExecute.hide();
  }
}
