import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../../auth/services/auth.service';
import { Query } from '../../../common/interfaces/query.interface';
import { StreamService } from '../../../common/services/stream.service';
import { Table, TableType } from '../../../shared/table/table';
import { AlertService } from 'app/shared/alerts/services/alert.service';


@Component({
  selector: 'dpc-dashboard-buy',
  templateUrl: './dashboard.buy.component.html',
  styleUrls: ['./dashboard.buy.component.scss']
})
export class DashboardBuyComponent implements OnInit {
  user: any;
  table: Table = new Table();
  query = new Query();

  constructor(
    private AuthService: AuthService,
    public streamService: StreamService,
    public alertService: AlertService,
  ) { }

  ngOnInit() {
    this.table.title = 'Streams';
    this.table.tableType = TableType.Buy;
    this.table.headers = ['Stream Name', 'Stream Type', 'Stream Price'];
    this.table.hasDetails = true;

    // Fetch current User
    this.user = {};
    this.AuthService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.query.owner = '-'.concat(this.user.id);
        this.fetchStreams();
      },
      err => {
        console.log(err);
      });
  }

  onPageChange(page: number) {
    this.query.page = page;
    this.fetchStreams();
  }

  onFiltersChange(filters: any) {
    Object.assign(this.query, filters);
    this.fetchStreams();
  }

  private fetchStreams() {
    this.streamService.searchStreams(this.query).subscribe(
      (result: any) => {
        const temp = Object.assign({}, this.table);
        temp.page = result;
        // Set table content
        this.table = temp;
      },
      err => {
        console.log(err);
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
      });
  }

}
