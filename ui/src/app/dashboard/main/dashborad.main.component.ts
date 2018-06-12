import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DatatableComponent } from '@swimlane/ngx-datatable';
import { AuthService } from '../../auth/services/auth.service';
import { User } from '../../common/interfaces/user.interface';
import { StreamService } from './services/stream.service';
import { SubscriptionService } from './services/subscription.service';
import { SearchService } from './services/search.service';
import { MdlDialogService } from '@angular-mdl/core';

import * as L from 'leaflet';
import { icon, latLng, Layer, marker } from 'leaflet';
import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';
import { TasPipe } from '../../common/pipes/converter.pipe';

@Component({
  selector: 'dashboard-main',
  templateUrl: './dashboard.main.component.html',
  styleUrls: [ './dashboard.main.component.scss' ]
})
export class DashboardMainComponent {
    user:any;
    mySubscriptions = [];
    myStreamsList = [];
    temp2 = [];
    temp = [];
    userEventFlag = 0;
    markers: Layer[] = [];
    markersSubs: Layer[] = [];
    markerInd: any;

    // ngx-table custom messages
    tableStreamsMessages =  {
        emptyMessage: "You don't have any streams yet..."
    }
    // ngx-table custom messages
    tableSubsMessages =  {
        emptyMessage: "You don't have any subscription yet..."
    }
    // Leaflet options
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

    drawnItems: any;
    tabIndex = 0;

    @ViewChild(DatatableComponent) table: DatatableComponent;
    constructor(
        private AuthService: AuthService,
        private SubscriptionService: SubscriptionService,
        private streamService: StreamService,
        public searchService: SearchService,
        public http: HttpClient,
        private dialogService: MdlDialogService,
        private tasPipe: TasPipe
    ) {}

    ngOnInit() {
        this.user = {};
        this.AuthService.getCurrentUser().subscribe(
            data =>  {
                this.user = data;
            },
            err => {
                console.log(err)
            }
        );

        // Fetch all subscriptions
        this.SubscriptionService.get().subscribe(
          result => {
              this.temp2 = [...result];

              // Create list of user subscriptions
              this.mySubscriptions = this.temp2;

              for (var i = 0; i < this.mySubscriptions.length; i++) {
                  // Create marker with stream coordinates
                  const newMarker = L.marker(
                  [this.mySubscriptions[i]["stream_data"]["coordinates"][1],
                   this.mySubscriptions[i]["stream_data"]["coordinates"][0]], {}
                  );
                  // Use yellow color for owner streams and blue for others
                  var defIcon = L.icon({
                      iconUrl:  '/assets/images/green-marker.png',
                      iconSize: [45, 45]
                  });
                  newMarker.setIcon(defIcon);
                  // Popup Msg
                  const msg = "<b>" + this.mySubscriptions[i]["stream_data"]["name"] + "</b>"
                  // Push marker to the markers list
                  this.markersSubs.push(newMarker);
              }
              // Set markers on the map
              this.setTabMarkers();
          },
          err => {
              console.log(err)
          }
      );
    }

    onMapReady(map: L) {
        const search = this.searchService;
        const that = this;

        that.drawnItems = new L.FeatureGroup();
        map.addLayer(that.drawnItems);

        // Set Markers on initial charged map
        const bounds = map.getBounds();
        setMapMarkers(bounds["_southWest" ]["lng"], bounds["_southWest" ]["lat"],
                      bounds["_southWest" ]["lng"], bounds["_northEast" ]["lat"],
                      bounds["_northEast" ]["lng"], bounds["_northEast" ]["lat"],
                      bounds["_northEast" ]["lng"], bounds["_southWest" ]["lat"]);

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

            this.drawnItems.addLayer(layer);
        });

        // Generic function to draw map markers from coordinates
        function setMapMarkers(x1, y1, x2, y2, x3, y3, x4, y4) {
            // Remove all layers and reset markers and color flag
            var layers = that.drawnItems.getLayers();
            that.drawnItems.clearLayers();
            that.userEventFlag = 0;
            that.markers = [];
            that.myStreamsList = [];

            // Search streams on drawed region
            that.searchService.searchStreams("geo",x1, y1, x2, y2, x3, y3, x4, y4).subscribe(
                (result: any) => {
                    that.temp = [...result.Streams];
                    // Add stream markers on the map (Name, Description and price)
                    result.Streams.forEach(stream => {
                    if (stream["owner"] == that.user["id"]) {
                        // Create marker with stream coordinates
                        const newMarker = L.marker(
                        [stream["location"]["coordinates"][1],
                         stream["location"]["coordinates"][0]], {}
                        );
                        // Use yellow color for owner streams and blue for others
                        var defIcon = L.icon({
                            iconUrl:  '/assets/images/yellow-marker.png',
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
                        that.myStreamsList.push(stream);
                        that.markers.push(newMarker);
                    }
                });

                    // Set markers on the map
                    that.setTabMarkers();
                    // Refresh ngx-datatable list
                    that.myStreamsList = [...that.myStreamsList];
                    // Update temp list to use with search filter
                    that.temp = that.myStreamsList;
                },
                err => { console.log(err) }
            );
        }
    }

    setMarkerColor(i: number){
        if (this.myStreamsList[i]["owner"] != this.user["email"]) {
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

    // Subscriptions table fileter
    updateSubscriptions(event) {
        const val = event.target.value.toLowerCase();
        // filter our data
        const temp = this.temp2.filter(function(d) {
            const n  = d.streamUrl.toLowerCase().indexOf(val) !== -1 || !val;
            return n;
        });

        // update the rows
        this.mySubscriptions = temp;
        // Whenever the filter changes, always go back to the first page
        this.table.offset = 0;
    }
    // Mystreams table fileter
    updateMyStreams(event) {
        const val = event.target.value.toLowerCase();
        // filter our data
        const temp = this.temp.filter(function(d) {
            const n =  d.name.toLowerCase().indexOf(val) !== -1 || !val;
            const t =  d.type.toLowerCase().indexOf(val) !== -1 || !val;
            return n || t;
        });

        // update the rows
        this.myStreamsList = temp;
        // Whenever the filter changes, always go back to the first page
        this.table.offset = 0;
    }

    // Draw subscriptions or myStreams markers
    setTabMarkers() {
        if (this.tabIndex == 0) {
            for (var i = 0; i < this.markersSubs.length; i++) {
                this.drawnItems.addLayer(this.markersSubs[i]);
            }
        } else {
            for (var i = 0; i < this.markers.length; i++) {
                this.drawnItems.addLayer(this.markers[i]);
            }
        }
    }

    tabChanged(event) {
        var layers = this.drawnItems.getLayers();
        this.drawnItems.clearLayers();
        this.tabIndex = event.index;
        this.setTabMarkers();
    }

    // Open BUY tokens dialog
    onUrlClick(url) {
        const urlMsg = '<a href="' + url + '"> LINK </a>'
        let resultCopy = this.dialogService.confirm(urlMsg, "OK", "Copy");
        resultCopy.subscribe(
            copy => {
                console.log("Copy");
            },
            err => {}
         );
    }
}
