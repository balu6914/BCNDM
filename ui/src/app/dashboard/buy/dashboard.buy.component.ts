import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { StreamService } from '../../common/services/stream.service';
import { TasPipe } from '../../common/pipes/converter.pipe';
import { AuthService } from '../../auth/services/auth.service';
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

      // TODO: Replace this coordinates with map bounds
      const southWestLng = -180;
      const southWestLat = -90;
      const northEastLng = 180;
      const northEastLat = 90;

      // Search streams on drawed region
      this.streamService.searchStreams(
        "geo", southWestLng, southWestLat, southWestLng, northEastLat,
        northEastLng, northEastLat, northEastLng, southWestLat).subscribe(
        (result: any) => {
          this.temp = [...result.Streams];
          result.Streams.forEach(stream => {
            if (stream.owner != this.user.id) {
              this.streams.push(stream)
            }
          }
        );
        // Set table content
        this.table.content = this.streams;
      },
      err => {
        console.log(err)
      });
    }
}
