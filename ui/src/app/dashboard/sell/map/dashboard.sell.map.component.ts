import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { StreamService } from '../services/stream.service';
import { SearchService } from '../services/search.service';
import { AuthService } from '../../../auth/services/auth.service';
import { Router } from '@angular/router';
import { TasPipe } from '../../../common/pipes/converter.pipe';
import { Table, TableType } from '../../../shared/table/table';

@Component({
  selector: 'dashboard-sell-map',
  templateUrl: './dashboard.sell.map.component.html',
  styleUrls: [ './dashboard.sell.map.component.scss' ]
})
export class DashboardSellMapComponent {
    temp = [];
    streams = [];
    user: any;
    table: Table = new Table();

    constructor(
        private streamService: StreamService,
        private searchService: SearchService,
        private AuthService: AuthService,
        private router: Router,
        private tasPipe: TasPipe
    ) { }

  ngOnInit() {
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

    // Config table
    this.table.title = "Streams";
    this.table.tableType = TableType.Sell;
    this.table.headers = ["Stream Name", "Stream Type","Stream Price"];

    // TODO: Replace this coordinates with map bounds
    const southWestLng = -180;
    const southWestLat = -90;
    const northEastLng = 180;
    const northEastLat = 90;

    // Search streams on drawed region
    this.searchService.searchStreams(
      "geo", southWestLng, southWestLat, southWestLng, northEastLat,
      northEastLng, northEastLat, northEastLng, southWestLat).subscribe(
      (result: any) => {
        this.temp = [...result.Streams];
        result.Streams.forEach(stream => {
          if (stream.owner == this.user.id) {
            this.streams.push(stream);
          }
        });
        // Set table content
        this.table.content = this.streams;
      },
      err => {
        console.log(err);
      }
    );
  }

  // Add Bulk event
  onFileChange(event): void {
    const fileList: FileList = event.target.files;
    if (fileList.length > 0) {
      const file = fileList[0];

      const formData = new FormData();
      formData.append('csv', file, file.name);

      this.streamService.addStreamBulk(formData).subscribe(
        data => this.router.navigate(['/dashboard/sell/map']),
        error => console.log(error),
    );
    }
  }
}
