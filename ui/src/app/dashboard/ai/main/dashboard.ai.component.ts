import { Component, OnInit, ViewChild } from '@angular/core';
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

@Component({
  selector: 'dpc-dashboard-ai',
  templateUrl: './dashboard.ai.component.html',
  styleUrls: ['./dashboard.ai.component.scss']
})
export class DashboardAiComponent implements OnInit {
  user: User;
  tableStreams: Table = new Table();
  tableAlgos: Table = new Table();
  tableExecutions: Table = new Table();
  query = new Query();
  checkedStreams = [];
  checkedAlgos = [];

  @ViewChild('tableStreamsComponent')
  private tableStreamsComponent: TableComponent;

  @ViewChild('tableAlgosComponent')
  private tableAlgosComponent: TableComponent;

  @ViewChild('tableExecutionComponent')
  private tableExecutionComponent: TableComponent;

  constructor(
    private modalService: BsModalService,
    private streamService: StreamService,
    private authService: AuthService,
    private executionsService: ExecutionsService,
    public alertService: AlertService,
  ) {
  }

  ngOnInit() {
    // Config tableStreams
    this.tableStreams.title = 'Datasets';
    this.tableStreams.tableType = TableType.Ai;
    this.tableStreams.headers = ['Name', 'Type', 'Price', ''];
    this.tableStreams.hasDetails = true;

    // Config tableStreams
    this.tableAlgos.title = 'Algorithms';
    this.tableAlgos.tableType = TableType.Ai;
    this.tableAlgos.headers = ['Name', 'Type', 'Price', ''];
    this.tableAlgos.hasDetails = true;

    // Config tableStreams
    this.tableExecutions.title = 'Jobs Queue';
    this.tableExecutions.tableType = TableType.Executions;
    this.tableExecutions.headers = ['ID', 'Mode', 'Algo', 'Data', 'State'];
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

  private fetchStreams() {
    this.query.streamType = 'Dataset';
    this.streamService.searchStreams(this.query).subscribe(
      (result: any) => {
        const tempStreams = Object.assign({}, this.tableStreams);
        tempStreams.page = result;
        // Set tableStreams content
        this.tableStreams = tempStreams;
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      });

      this.query.streamType = 'Algorithm';
      this.streamService.searchStreams(this.query).subscribe(
        (result: any) => {
          const tempAlgos = Object.assign({}, this.tableAlgos);
          tempAlgos.page = result;

          // Set tableAlgos content
          this.tableAlgos = tempAlgos;
        },
        err => {
          this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
        });
  }

  fetchExecutions() {
    this.executionsService.getExecutions().subscribe(
      (result: any) => {
        this.tableExecutions.page = {
          page: 0,
          total: 5,
          limit: 50,
          content: result.executions,
        };
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
      streams: this.checkedStreams,
      algos: this.checkedAlgos,
    };

    // Show DashboardAiExecuteComponent as Modal
    this.modalService.show(DashboardAiExecuteComponent, { initialState })
      .content.executionStarted.subscribe(
        response => {
          this.fetchExecutions();
        },
        err => {
        }
      );
    }

    onStreamSelected(row: Stream) {
      const index = this.checkedStreams.findIndex(
        stream =>  row.id === stream.id
      );
      if (index === -1) {
        this.checkedStreams.push(row);
      } else {
        this.checkedStreams.splice(index, 1);
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
}