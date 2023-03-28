import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { ExecutionsService } from 'app/common/services/executions.service';
import { floatRegEx, urlRegEx } from 'app/common/validators/patterns';
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
  datasets = [];
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
    this.datasets.forEach ( (dataset, i) => {
      const executionReq: ExecutionReq = {
        data: dataset.id,
        // TODO: Make this fields configurable.
        local_args: ['--isShuffle', '--nEstimators', '100',
        '--oobScore', '--maxFeatures', '0.3', '--minSamplesLeaf',
        '20', '--minSamplesSplit', '4', '--criterion', 'gini',
        '--classWeight', 'balanced', '--bootstrap', '--nJobs', '4',
        '--randomState', '0'],
        type: 'ITERATIVE',
        global_timeout: 3000000,
        local_timeout: 1000000,
        preprocess_args: [],
        mode: 'FEDERATED',
        global_args: [],
        files: [],
      };
      this.executionReqList.push(executionReq);
    });

    this.algos.forEach( algo => {
      const startExecutionReq: StartExecutionReq = {
        algo: algo.id,
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
