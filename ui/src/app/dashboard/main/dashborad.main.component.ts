import { Component, ViewChild } from '@angular/core';
import { MdlDialogService } from '@angular-mdl/core';
import { HttpClient } from '@angular/common/http';
import { DatatableComponent } from '@swimlane/ngx-datatable';

import { AuthService } from '../../auth/services/auth.service';
import { SearchService } from './services/search.service';
import { StreamService } from './services/stream.service';
import { SubscriptionService } from './services/subscription.service';
import { TasPipe } from '../../common/pipes/converter.pipe';
import { User } from '../../common/interfaces/user.interface';

import * as L from 'leaflet';
import { icon, latLng, Layer, marker } from 'leaflet';
import { LeafletModule } from '@asymmetrik/ngx-leaflet';
import { LeafletDrawModule } from '@asymmetrik/ngx-leaflet-draw';

import { Subscription } from '../../common/interfaces/subscription.interface';

import { Chart } from 'chart.js';
import {} from '@types/googlemaps';

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
          (result: any) => {
              this.temp2 = [...result.Subscriptions];

              result.Subscriptions.forEach(subscription => {
                  this.streamService.getStream(subscription["id"]).subscribe(
                    (result: any) => {
                        const stream = result.Stream
                        subscription["stream_name"] = stream["name"]
                        subscription["stream_price"] = stream["price"]

                        // Create marker with stream coordinate
                        const newMarker = L.marker(
                        [stream["location"]["coordinates"][1],
                         stream["location"]["coordinates"][0]], {}
                        );
                        // Use yellow color for owner streams and blue for others
                        var defIcon = L.icon({
                            iconUrl:  '/assets/images/green-marker.png',
                            iconSize: [45, 45]
                        });
                        newMarker.setIcon(defIcon);
                        // Popup Msg
                        const msg = "<b>" + stream["name"] + "</b>"
                        // Push marker to the markers list
                        this.mySubscriptions.push(subscription);
                        this.markersSubs.push(newMarker);

                        // Set markers on the map
                        this.setTabMarkers();
                    },
                    err => {
                        console.log(err)
                    });
                })
          },
          err => {
              console.log(err)
          });

          let chartData = {
              type: "line",
              data: {
                  labels: [
                      "Jan 2017",
                      "Apr 2017",
                      "Sep 2017",
                      "Dec 2017",
                      "Mar 2018",
                      "Jul 2018",
                      "Oct 2018",
                      "Feb 2019"
                  ],
                  datasets: [
                      {
                          type: "bar",
                          label: "Dataset 1",
                          data: [5, 10, 15, 7, 3, 10, 2, 45, 12, 3, 35, 2, 5],
                          backgroundColor: "rgba(6, 210, 216, 1)",
                          borderColor: "rgba(6, 210, 216, 1)",
                          borderWidth: 1,
                          barThickness: 1
                      },
                      {
                          label: "Dataset 2",data: [25, 43, 38, 33, 52, 65, 62, 49],
                          backgroundColor: "rgba(0, 125, 255, .1)",
                          borderColor: "#007DFF",
                          borderWidth: 4,
                          pointBackgroundColor: "#ffffff",
                          pointRadius: 3,
                          pointBorderWidth: 1
                      }
                  ]
              },
              options: {
                  maintainAspectRatio: false,
                  scaleShowVerticalLines: false,
                  tooltips: {
                      backgroundColor: "#007DFF",
                      xPadding: 15,
                      yPadding: 5,
                      titleMarginBottom: 0,
                      bodySpacing: 2,
                      cornerRadius: 0,
                      displayColors: false,
                      caretSize: 0,
                      callbacks: {
                          label: function(tooltipItem, data) {
                              return tooltipItem.yLabel + " TOK";
                          },
                          title: function(tooltipItem, data) {
                              return;
                          }
                      }
                  },
                  legend: {
                      display: false
                  },
                  scales: {
                      yAxes: [
                          {
                              afterTickToLabelConversion: function(q) {
                                  for (var tick in q.ticks) {
                                      let newLabel = q.ticks[tick] + " TOK ";
                                      q.ticks[tick] = newLabel;
                                  }
                              },
                              gridLines: {
                                  color: "rgba(223,233,247,1)",
                                  zeroLineColor: "rgba(223,233,247,1)",
                                  borderDash: [15, 15],
                                  drawBorder: false
                              },
                              ticks: {
                                  fontColor: "rgba(158,175,200, 1)",
                                  fontSize: 11,
                                  stepSize: 25
                              }
                          }
                      ],
                      xAxes: [
                          {
                              barPercentage: 10,
                              categoryPercentage: 0.1,
                              barThickness: 5,
                              gridLines: {
                                  lineWidth: 0,
                                  color: "rgba(255,255,255,0)",
                                  zeroLineColor: "rgba(255,255,255,0)"
                              },
                              ticks: {
                                  fontColor: "rgba(158,175,200, 1)",
                                  fontSize: 11
                              }
                          }
                      ]
                  }
              }
          };

          let c: any = document.getElementById("myChart");
          let ctx = c.getContext("2d");
          let chart = new Chart(ctx, chartData);

          // When the window has finished loading create our google map below

          // Basic options for a simple Google Map
          // For more options see: https://developers.google.com/maps/documentation/javascript/reference#MapOptions

          var bounds = new google.maps.LatLngBounds();
          let mapOptions : any =  {

              // The latitude and longitude to center the map (always required)
              center: new google.maps.LatLng(40.67, -73.94), // New York

              // TODO: Fix this
              // How zoomed in you want the map to start at (always required)
              zoom: 11,
              disableDefaultUI: true,
              zoomControl: true,

              // How you would like to style the map.
              // This is where you would paste any style found on Snazzy Maps.
              styles: [
                  {
                      featureType: "administrative",
                      elementType: "labels.text.fill",
                      stylers: [
                          {
                              color: "#444444",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.country",
                      elementType: "geometry.stroke",
                      stylers: [
                          {
                              color: "#85a9c1",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.country",
                      elementType: "labels",
                      stylers: [
                          {
                              visibility: "on",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.country",
                      elementType: "labels.text",
                      stylers: [
                          {
                              visibility: "on",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.country",
                      elementType: "labels.text.fill",
                      stylers: [
                          {
                              color: "#9eafc8",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.province",
                      elementType: "geometry.stroke",
                      stylers: [
                          {
                              color: "#9eafc8",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.province",
                      elementType: "labels",
                      stylers: [
                          {
                              visibility: "on",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.province",
                      elementType: "labels.text",
                      stylers: [
                          {
                              visibility: "on",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.locality",
                      elementType: "labels",
                      stylers: [
                          {
                              visibility: "on",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.locality",
                      elementType: "labels.text",
                      stylers: [
                          {
                              visibility: "on"
                          },
                          {
                              color: "#9eafc8",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.locality",
                      elementType: "labels.text.fill",
                      stylers: [
                          {
                              color: "#9eafc8",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.locality",
                      elementType: "labels.text.stroke",
                      stylers: [
                          {
                              color: "#ffffff",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.locality",
                      elementType: "labels.icon",
                      stylers: [
                          {
                              lightness: "66",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.neighborhood",
                      elementType: "labels",
                      stylers: [
                          {
                              visibility: "off",
                          }
                      ]
                  },
                  {
                      featureType: "administrative.land_parcel",
                      elementType: "labels",
                      stylers: [
                          {
                              visibility: "off",
                          }
                      ]
                  },
                  {
                      featureType: "landscape",
                      elementType: "all",
                      stylers: [
                          {
                              color: "#f2f2f2",
                          }
                      ]
                  },
                  {
                      featureType: "landscape",
                      elementType: "geometry.fill",
                      stylers: [
                          {
                              color: "#f3f7fa",
                          }
                      ]
                  },
                  {
                      featureType: "poi",
                      elementType: "all",
                      stylers: [
                          {
                              visibility: "off",
                          }
                      ]
                  },
                  {
                      featureType: "road",
                      elementType: "all",
                      stylers: [
                          {
                              saturation: -100,
                          },
                          {
                              lightness: 45,
                          }
                      ]
                  },
                  {
                      featureType: "road",
                      elementType: "geometry.fill",
                      stylers: [
                          {
                              color: "#ffffff",
                          }
                      ]
                  },
                  {
                      featureType: "road",
                      elementType: "labels.text.fill",
                      stylers: [
                          {
                              color: "#9eafc8",
                          }
                      ]
                  },
                  {
                      featureType: "road.highway",
                      elementType: "all",
                      stylers: [
                          {
                              visibility: "simplified",
                          }
                      ]
                  },
                  {
                      featureType: "road.highway",
                      elementType: "labels",
                      stylers: [
                          {
                              visibility: "off",
                          }
                      ]
                  },
                  {
                      featureType: "road.arterial",
                      elementType: "labels.icon",
                      stylers: [
                          {
                              visibility: "off",
                          }
                      ]
                  },
                  {
                      featureType: "transit",
                      elementType: "all",
                      stylers: [
                          {
                              visibility: "off",
                          }
                      ]
                  },
                  {
                      featureType: "water",
                      elementType: "all",
                      stylers: [
                          {
                              color: "#c0e4f3",
                          },
                          {
                              visibility: "on",
                          }
                      ]
                  }
              ]
          };

          let infoWindowContent = [
              [
                  '<div class="map-tooltip">' +
                      '<p class="map-tooltip__title">Luftdaten Hum 8806</p>' +
                      '<div id="bodyContent" class="map-tooltip__content">' +
                      '<p class="map-tooltip__subtitle">Humidity</p>' +
                      '<p class="map-tooltip__stake">Stake: <span class="map-tooltip__stake-amount">0.043256 TAS</span></p>' +
                      "</div>" +
                      "</div>"
              ],
              [
                  '<div class="map-tooltip">' +
                      '<p class="map-tooltip__title">Lorem Ipsum</p>' +
                      '<div id="bodyContent" class="map-tooltip__content">' +
                      '<p class="map-tooltip__subtitle">Humidity</p>' +
                      '<p class="map-tooltip__stake">Stake: <span class="map-tooltip__stake-amount">0.043256 TAS</span></p>' +
                      "</div>" +
                      "</div>"
              ],
              [
                  '<div class="map-tooltip">' +
                      '<p class="map-tooltip__title">Sold Luftdaten Hum 8806</p>' +
                      '<div id="bodyContent" class="map-tooltip__content">' +
                      '<p class="map-tooltip__subtitle">Humidity</p>' +
                      '<p class="map-tooltip__stake">Stake: <span class="map-tooltip__stake-amount">0.043256 TAS</span></p>' +
                      "</div>" +
                      "</div>"
              ],
              [
                  '<div class="map-tooltip">' +
                      '<p class="map-tooltip__title">Air Cuality</p>' +
                      '<div id="bodyContent" class="map-tooltip__content">' +
                      '<p class="map-tooltip__subtitle">Humidity</p>' +
                      '<p class="map-tooltip__stake">Stake: <span class="map-tooltip__stake-amount">0.043256 TAS</span></p>' +
                      "</div>" +
                      "</div>"
              ]
          ];

          // Get the HTML DOM element that will contain your map
          // We are using a div with id="map" seen below in the <body>
          let mapElement = document.getElementById("map");

          // Create the Google Map using our element and options defined above
          let map = new google.maps.Map(mapElement, mapOptions);

          // Multiple Markers
          let markers = [
              ["Datapace 1", 40.67, -73.94, "assets/img/icons/map-co2.svg"],
              ["Datapace 2", 40.64, -73.93, "assets/img/icons/map-chart.svg"],
              ["Datapace 3", 40.68, -73.91, "assets/img/icons/map-temp.svg"],
              ["Datapace 4", 40.62, -73.91, "assets/img/icons/map-water.svg"]
          ];

          // Display multiple markers on a map
          let infoWindow = new google.maps.InfoWindow(),
              marker,
              i;

          // Loop through our array of markers & place each one on the map
          for (i = 0; i < markers.length; i++) {
              let position = new google.maps.LatLng(Number(markers[i][1]), Number(markers[i][2]));
              bounds.extend(position);
              let marker = new google.maps.Marker({
                  position: position,
                  map: map,
                  title: String(markers[i][0]),
                  icon: String(markers[i][3]),
              });

              // Allow each marker to have an info window
              google.maps.event.addListener(
                  marker,
                  "click",
                  (function(marker, i) {
                      return function() {
                          infoWindow.setContent(infoWindowContent[i][0]);
                          infoWindow.open(map, marker);
                      };
                  })(marker, i)
              );

              // Automatically center the map fitting all markers on the screen
              map.fitBounds(bounds);
          }

          var styles = {
              default: null,
              hide: [
                  {
                      featureType: "poi.business",
                      stylers: [{ visibility: "off" }]
                  },
                  {
                      featureType: "transit",
                      elementType: "labels.icon",
                      stylers: [{ visibility: "off" }]
                  }
              ]
          };
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
        this.tabIndex = event.index;
        this.drawnItems.clearLayers();
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
