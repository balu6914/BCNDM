import { Component, OnInit } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';

import { AlertService } from 'app/shared/alerts/services/alert.service';
import { AuthService } from 'app/auth/services/auth.service';
import { Query } from 'app/common/interfaces/query.interface';
import { StreamService } from 'app/common/services/stream.service';
import { ExecutionsService } from 'app/common/services/executions.service';
import { Table, TableType } from 'app/shared/table/table';
import { DashboardAiExecuteComponent } from 'app/dashboard/ai/execute/dashboard.ai.execute.component';
import { User } from 'app/common/interfaces/user.interface';
import { Stream } from 'app/common/interfaces/stream.interface';
import { TableComponent } from 'app/shared/table/main/table.component';
import { SubscriptionService } from 'app/common/services/subscription.service';
import { Page } from 'app/common/interfaces/page.interface';
import { Subscription } from 'app/common/interfaces/subscription.interface';
import { DashboardAiAddComponent } from 'app/dashboard/ai/add/dashboard.ai.add.component';
import { Execution } from 'app/common/interfaces/execution.interface';

@Component({
  selector: 'dpc-dashboard-ai',
  templateUrl: './dashboard.ai.component.html',
  styleUrls: ['./dashboard.ai.component.scss']
})
export class DashboardAiComponent implements OnInit {
  user: User;
  tableAlgorithms: Table = new Table();
  tableDatasets: Table = new Table();
  tableExecutions: Table = new Table();
  query = new Query();
  checkedDatasets = [];
  checkedAlgos = [];

  constructor(
    private modalService: BsModalService,
    private streamService: StreamService,
    private authService: AuthService,
    private executionsService: ExecutionsService,
    public alertService: AlertService,
    private subscriptionService: SubscriptionService,
  ) {
  }

  ngOnInit() {
    // Config tableDatasets
    this.tableDatasets.title = 'Datasets';
    this.tableDatasets.tableType = TableType.Ai;
    this.tableDatasets.headers = [' Name', ' Type', ' Price', 'Execute', ''];
    this.tableDatasets.hasDetails = true;

    // Config tableAlgorithms
    this.tableAlgorithms.title = 'Algorithms';
    this.tableAlgorithms.tableType = TableType.Ai;
    this.tableAlgorithms.headers = [' Name', ' Type', ' Price', 'Execute', ''];
    this.tableAlgorithms.hasDetails = true;

    // Config tableExecutions
    this.tableExecutions.title = 'Jobs Queue';
    this.tableExecutions.tableType = TableType.Executions;
    this.tableExecutions.headers = ['ID', 'Mode', 'Algo', 'Data', 'State', 'Result'];
    this.tableExecutions.hasDetails = true;

    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.query.owner = this.user.id;
        this.fetchStreams();
        this.fetchExecutions();
      },
      err => {
        console.log(err);
      }
    );
  }

  setDatasetsTable(page: any) {
    const tempDatas = Object.assign({}, this.tableDatasets);
    // Concat Page fields
    tempDatas.page.total = page.limit;
    tempDatas.page.total = page.total;
    tempDatas.page.content.push(...page.content);
    // Set tableDatasets content
    this.tableDatasets = tempDatas;
  }

  setAlgorithmsTable(page: any) {
    const tempAlgos = Object.assign({}, this.tableAlgorithms);
    // Concat Page fields
    tempAlgos.page.total = page.limit;
    tempAlgos.page.total = page.total;
    tempAlgos.page.content.push(...page.content);
    // Set tableDatasets content
    this.tableAlgorithms = tempAlgos;
  }

  addAiStreamToTable(stream: Stream) {
    if (stream.type === 'Algorithm') {
      this.tableAlgorithms.page.content.push(stream);
    } else if ((stream.type === 'Dataset')) {
      this.tableDatasets.page.content.push(stream);
    }
  }

  private fetchStreams() {
    // Fetch streams of type Dataset
    this.query.streamType = 'Dataset';
    this.streamService.searchStreams(this.query).subscribe(
      (page: Page<Stream>) => {
        this.setDatasetsTable(page);
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );

    // Fetch streams of type Algorithm
    this.query.streamType = 'Algorithm';
    this.streamService.searchStreams(this.query).subscribe(
      (page: Page<Stream>) => {
        this.setAlgorithmsTable(page);
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );

    this.subscriptionService.bought(0, 100).subscribe(
      (page: Page<Subscription>) => {
        page.content.forEach( sub => {
          const date = new Date();
          const subDate = new Date(sub.start_date);
          subDate.setHours(subDate.getHours() + Number(sub.hours));
          if (date < subDate) {
            this.streamService.getStream(sub.stream_id).subscribe(
              (stream: Stream) => {
                this.addAiStreamToTable(stream);
              },
              err => {
                this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
              }
            );
          }
        });
      },
      err => console.log(err)
    );
  }

  fetchExecutions() {
    this.executionsService.getExecutions().subscribe(
      (execResp: any) => {
        execResp.executions.forEach( exec => {
          // TODO: Add Algorithm and Dataset names in Execution structure
          this.streamService.getStream(exec.algo).subscribe(
            (stream: Stream) => {
              exec.algo = stream.name;
            },
            err => {
              this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
            }
          );
          this.streamService.getStream(exec.data).subscribe(
            (stream: Stream) => {
              exec.data = stream.name;
            },
            err => {
              this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
            }
          );
        });

        if (execResp.executions) {
          this.tableExecutions.page = {
            page: 0,
            limit: 5,
            total: execResp.executions.length,
            content: execResp.executions,
          };
        }
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );
  }

  fetchExecutionResult(execution: Execution) {
    this.executionsService.getExecutionResult(execution.id).subscribe(
      (execResult: any) => {
        execution.result = JSON.stringify(execResult.result);
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );
  }

  onPageChange(page: number) {
    this.query.page = page;
    this.fetchStreams();
  }

  modalNewExecution() {
    const initialState = {
      datasets: this.checkedDatasets,
      algos: this.checkedAlgos,
    };

    // Show DashboardAiExecuteComponent as Modal
    this.modalService.show(DashboardAiExecuteComponent, { initialState })
      .content.executionStarted.subscribe(
        response => {
          this.fetchExecutions();
        },
        err => {
          this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
        }
      );
    }

    onDataSelected(row: Stream) {
      const index = this.checkedDatasets.findIndex(
        stream =>  row.id === stream.id
      );
      if (index === -1) {
        this.checkedDatasets.push(row);
      } else {
        this.checkedDatasets.splice(index, 1);
      }
    }

    onAlgoSelected(row: Stream) {
      const index = this.checkedAlgos.findIndex(
        algo =>  row.id === algo.id
      );
      if (index === -1) {
        this.checkedAlgos.push(row);
      } else {
        this.checkedAlgos.splice(index, 1);
      }
    }

    openModalAdd(type: String) {
      const initialState = {
        streamType: type,
        ownerID: this.user.id,
      };

      // Show DashboardAilAddComponent as Modal
      this.modalService.show(DashboardAiAddComponent, { initialState })
        .content.aiStreamCreated.subscribe(
          stream => {
            this.addAiStreamToTable(stream);
          }
        );
    }

    deleteAlgo(algoID) {
      this.checkedAlgos = this.checkedAlgos.filter(a => a.id !== algoID);
    }

    deleteData(dataID) {
      this.checkedDatasets = this.checkedDatasets.filter(d => d.id !== dataID);
    }
}
