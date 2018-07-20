import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DatatableComponent } from '@swimlane/ngx-datatable';

import { StreamService } from '../services/stream.service';
import { SearchService } from '../services/search.service';
import { AuthService } from '../../../auth/services/auth.service';
import { Router } from '@angular/router';
import { TasPipe } from '../../../common/pipes/converter.pipe';

import * as L from 'leaflet';
import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';
import { icon, latLng, Layer, marker, tileLayer } from 'leaflet';

@Component({
  selector: 'dashboard-sell-map',
  templateUrl: './dashboard.sell.map.component.html',
  styleUrls: [ './dashboard.sell.map.component.scss' ]
})
export class DashboardSellMapComponent {
    options = {
		layers: [
            L.tileLayer('https://api.mapbox.com/styles/v1/gesaleh/cjdbxg3f6c6sq2smdj7cp4wwa/tiles/256/{z}/{x}/{y}?access_token=pk.eyJ1IjoiZ2VzYWxlaCIsImEiOiJjamQ4bXFuZ3kybDZiMnhxcjl6Mjlmc3hmIn0.RVdSuXXmCgZJubeCAncjJQ', {
                attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="http://mapbox.com">Mapbox</a>',
                maxZoom: 18,
                id: 'mapbox.dark',
                accessToken: 'pk.eyJ1IjoiZ2VzYWxlaCIsImEiOiJjamQ4bXFuZ3kybDZiMnhxcjl6Mjlmc3hmIn0.RVdSuXXmCgZJubeCAncjJQ'
            })
		],
		zoom: 5,
		center: L.latLng({ lat: 48.864716, lng: 2.349014 })
	};

	drawOptions = {
        position: 'topright',
        draw: {
            marker: false,
            polygon: false,
            polyline: false,
            circle: false,
            circlemarker: false,

        },
        edit: {
            remove: false,
            edit: false
        }
	};

    tableColumns = [
        { prop: 'name' },
        { name: 'type' },
        { name: 'description' },
        { name: 'price'}
    ];

    streams = [];
    user: any;
    subscription: any;
    temp = [];
    temp2 = [];
    userEventFlag = 0;
    markerInd: any;
    markers: Layer[] = [];

    @ViewChild(DatatableComponent) table: DatatableComponent;

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
