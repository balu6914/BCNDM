import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { SearchService } from '../services/search.service';
import { TasPipe } from '../../../common/pipes/converter.pipe';
import { AuthService } from '../../../auth/services/auth.service';

@Component({
  selector: 'buy-map-container',
  templateUrl: './dashboard.buy.map.component.html',
  styleUrls: [ './dashboard.buy.map.component.scss' ]
})

export class DashboardBuyMapComponent {
    temp = [];
    streams = [];
    user: any;

    constructor(
        private AuthService: AuthService,
        public searchService: SearchService,
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
            if (stream.owner != this.user.id) {
              this.streams.push(stream)
            }

          });
      },
      err => {
        console.log(err)
      });
    }
}
