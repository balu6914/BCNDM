import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DatatableComponent } from '@swimlane/ngx-datatable';

import { StreamService } from '../services/stream.service';
import { SearchService } from '../services/search.service';
import { AuthService } from '../../../auth/services/auth.service';
import { Router } from '@angular/router';
import { TasPipe } from '../../../common/pipes/converter.pipe';

@Component({
  selector: 'dashboard-sell-map',
  templateUrl: './dashboard.sell.map.component.html',
  styleUrls: [ './dashboard.sell.map.component.scss' ]
})
export class DashboardSellMapComponent {
    temp = [];
    streams = [];
    user: any;


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
            this.streams.push(stream)
          }

        });
    },
    err => {
      console.log(err)
    });
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
