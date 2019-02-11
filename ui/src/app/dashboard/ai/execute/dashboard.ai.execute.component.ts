import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { ExecutionsService } from 'app/common/services/executions.service';
import { floatRegEx, urlRegEx } from 'app/shared/validators/patterns';
import { BsModalRef } from 'ngx-bootstrap';
import { BigQuery, Stream } from 'app/common/interfaces/stream.interface';
import { ExecutionReq, StartExecutionReq } from 'app/common/interfaces/execution.interface';

@Component({
  selector: 'dpc-dashboard-ai-execute',
  templateUrl: './dashboard.ai.execute.component.html',
  styleUrls: ['./dashboard.ai.execute.component.scss']
})
export class DashboardAiExecuteComponent implements OnInit {
  form: FormGroup;
  streams = [];
  algos = [];

  executionReqList: ExecutionReq[] = [];
  startExecutionReqList: StartExecutionReq[] = [];
  selectedMode: string;

  @Output() executionStarted: EventEmitter<any> = new EventEmitter();
  constructor(
    public modalExecute: BsModalRef,
    private formBuilder: FormBuilder,
    private executionsService: ExecutionsService,
    public alertService: AlertService,
  ) {
  }

  ngOnInit() {
  }

  execute() {
    this.streams.forEach ( (stream, i) => {
      const executionReq: ExecutionReq = {
        data: stream.name,
        mode: this.selectedMode,
      };
      this.executionReqList.push(executionReq);
    });

    this.algos.forEach( algo => {
      const startExecutionReq: StartExecutionReq = {
        algo: algo.name,
        executions: this.executionReqList,
      };
      this.startExecutionReqList.push(startExecutionReq);
    });

    this.startExecutionReqList.forEach( startExecReq => {
      this.executionsService.startExecution(startExecReq).subscribe(
        (result: any) => {
          this.alertService.success(`Execution succesfully started!`);
          this.executionStarted.emit(result);
          this.modalExecute.hide();
        },
        err => {
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
          this.modalExecute.hide();
        }
      );
    });
  }

  selectMode(mode: string) {
    this.selectedMode = mode;
  }
}
