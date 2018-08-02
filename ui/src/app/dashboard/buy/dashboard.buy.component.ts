import { Component } from '@angular/core';
import { AuthService } from '../../auth/services/auth.service';
import { Query } from '../../common/interfaces/query.interface';
import { TasPipe } from '../../common/pipes/converter.pipe';
import { StreamService } from '../../common/services/stream.service';
import { Table, TableType } from '../../shared/table/table';


@Component({
  selector: 'dashboard-buy',
  templateUrl: './dashboard.buy.component.html',
  styleUrls: [ './dashboard.buy.component.scss' ]
})
export class DashboardBuyComponent {
    temp = [];
    streams = [];
    user: any;
    table: Table = new Table();

    constructor(
        private AuthService: AuthService,
        public streamService: StreamService,
        private tasPipe: TasPipe
    ) { }

    ngOnInit() {
      this.table.title = "Streams";
      this.table.tableType =  TableType.Buy;
      this.table.headers = ["Stream Name", "Stream Type","Stream Price"];

      // Fetch current User
      this.user = {};
      this.AuthService.getCurrentUser().subscribe(
          data =>  {
              this.user = data;
          },
          err => {
              console.log(err)
          }
      );

      const query = new Query();

      this.streamService.searchStreams(query).subscribe(
        (result: any) => {
          this.temp = result.content;
          result.content.forEach(stream => {
            if (stream.owner !== this.user.id) {
              this.streams.push(stream)
            }
          }
        );
        result.content = this.streams;
        // Set table content
        this.table.page = result;
      },
      err => {
        console.log(err)
      });
    }

    onPageChanged(page: number) {
      const query = new Query();
      query.page = page;
      this.streamService.searchStreams(query).subscribe(
        (result: any) => {
          this.streams = [];
          result.content.forEach(stream => {
            if (stream.owner !== this.user.id) {
              this.streams.push(stream);
            }
          }
        );
        const temp = Object.assign({}, this.table);
        temp.page.content = this.streams;
        // Set table content
        this.table = temp;
      },
      err => {
        console.log(err);
      });

    }
}
