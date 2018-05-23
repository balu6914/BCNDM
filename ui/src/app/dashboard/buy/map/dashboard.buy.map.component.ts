import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import * as L from 'leaflet';

import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';
import { SearchService } from '../services/search.service';
import { TasPipe } from '../../../common/pipes/converter.pipe';
import { icon, latLng, Layer, marker, tileLayer } from 'leaflet';
import { AuthService } from '../../../auth/services/auth.service';
import { DatatableComponent } from '@swimlane/ngx-datatable';

@Component({
  selector: 'buy-map-container',
  templateUrl: './dashboard.buy.map.component.html',
  styleUrls: [ './dashboard.buy.map.component.scss' ]
})

export class DashboardBuyMapComponent {
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

    temp = [];
    streamList = [];

    tableColumns = [
        { prop: 'name' },
        { name: 'type' },
        { name: 'description' },
        { name: 'price'}
    ];

	markers: Layer[] = [];
    markerInd: any;
    userEventFlag = 0;
    subscription: any;
    user: any;

    @ViewChild(DatatableComponent) table: DatatableComponent;
    constructor(
        private AuthService: AuthService,
        public searchService: SearchService,
        private tasPipe: TasPipe
    ) { }

    onMapReady(map: L) {
        const search = this.searchService;
        const that = this;

        var drawnItems = new L.FeatureGroup();
        map.addLayer(drawnItems);

        // Set Markers on initial charged map
        const bounds = map.getBounds();
        setMapMarkers(bounds["_southWest" ]["lng"], bounds["_southWest" ]["lat"],
                      bounds["_southWest" ]["lng"], bounds["_northEast" ]["lat"],
                      bounds["_northEast" ]["lng"], bounds["_northEast" ]["lat"],
                      bounds["_northEast" ]["lng"], bounds["_southWest" ]["lat"]);

        //map.on('click', function (e){
        //});

        // Set original marker color if mouse is over the map and was over the list
        map.on('mouseover', function (e){
            if (that.userEventFlag) {
                that.setMarkerColor(that.markerInd);
            }
        });

        // Set markers on new view
        map.on('moveend',function(e){
            const bounds = map.getBounds();
            setMapMarkers(bounds["_southWest" ]["lng"], bounds["_southWest" ]["lat"],
                          bounds["_southWest" ]["lng"], bounds["_northEast" ]["lat"],
                          bounds["_northEast" ]["lng"], bounds["_northEast" ]["lat"],
                          bounds["_northEast" ]["lng"], bounds["_southWest" ]["lat"]);
        });

        // Set markers and polygon layer inside the created polygon
        map.on('draw:created', function (e) {
            var layer = e.layer;
            var latLngs = layer.getLatLngs()[0];
            setMapMarkers(latLngs[0]["lng"], latLngs[0]["lat"],
                          latLngs[1]["lng"], latLngs[1]["lat"],
                          latLngs[2]["lng"], latLngs[2]["lat"],
                          latLngs[3]["lng"], latLngs[3]["lat"]);

            drawnItems.addLayer(layer);
        });

        // Generic function to draw map markers from coordinates
        function setMapMarkers(x1, y1, x2, y2, x3, y3, x4, y4) {
            // Remove all layers and reset markers and color flag
            var layers = drawnItems.getLayers();
            drawnItems.clearLayers();
            that.userEventFlag = 0;
            that.markers = [];
            that.streamList = [];

            // Search streams on drawed region
            that.searchService.searchStreams("geo",x1, y1, x2, y2, x3, y3, x4, y4).subscribe(
                (result: any) => {
                    that.temp = [...result.Streams];
                    // Add stream markers on the map (Name, Description and price)
                    result.Streams.forEach(stream => {
                        if (stream["owner"] != that.user["id"]) {
                            // Create marker with stream coordinates
                            const newMarker = L.marker(
                            [stream["location"]["coordinates"][1],
                             stream["location"]["coordinates"][0]], {}
                            );
                            // Use yellow color for owner streams and blue for others
                            var defIcon = L.icon({
                                iconUrl:  '/assets/images/blue-marker.png',
                                iconSize: [45, 45]
                            });
                            newMarker.setIcon(defIcon);

                            // Popup Msg
                            const name = stream["name"]
                            const description = stream["description"]
                            const price = that.tasPipe.transform(stream["price"])
                            const msg = `<b>${name}</b> <br> ${description} <br> ${price} TAS`
                            newMarker.bindPopup(msg);

                            // Push marker to the markers list
                            that.streamList.push(stream);
                            that.markers.push(newMarker);

                            drawnItems.addLayer(newMarker);
                        }
                    });

                    // Refresh streamList for
                    that.streamList = [...that.streamList];
                },
                err => { console.log(err) }
            );
        }
    }

    ngOnInit() {
        this.subscription = this.AuthService.getCurrentUser();
              this.subscription
              .subscribe(data => {
                  this.user = data;
              });
    }

    onUserEvent ( e ) {
        if ( e.type == "mouseenter" ) {
            for (var i = 0; i < this.streamList.length; i++) {
                // Set red-marker if mouse is over corresponding row
                if (this.streamList[i]["id"] == e.row["id"]) {
                    var redIcon = L.icon({
                        iconUrl:  '/assets/images/red-marker.png',
                        iconSize: [45, 45]
                    });
                    this.markers[i].setIcon(redIcon);

                    // set original color if marker changed at least one time
                    if (this.markerInd != i && this.userEventFlag) {
                        this.setMarkerColor(this.markerInd);
                    }
                    this.markerInd = i;

                    // Set flag for first mouse over
                    this.userEventFlag = 1;
                }
            }
        }
    }

    setMarkerColor(i: number){
        if (this.streamList[i]["owner"] != this.user["id"]) {
                var defIcon = L.icon({
                    iconUrl:  '/assets/images/blue-marker.png',
                    iconSize: [45, 45]
                });
                this.markers[i].setIcon(defIcon);
        } else {
            var defIcon = L.icon({
                iconUrl:  '/assets/images/yellow-marker.png',
                iconSize: [45, 45]
            });
            this.markers[i].setIcon(defIcon);
        }
    }

    // Mystreams table filter
    updateMyStreams(event) {
        const val = event.target.value.toLowerCase();
        // filter our data
        const temp2 = this.temp.filter(function(d) {
            const n =  d.name.toLowerCase().indexOf(val) !== -1 || !val;
            const t =  d.type.toLowerCase().indexOf(val) !== -1 || !val;
            return n || t;
        });

        // update the rows
        this.streamList = temp2;
        // Whenever the filter changes, always go back to the first page
        this.table.offset = 0;
    }
}