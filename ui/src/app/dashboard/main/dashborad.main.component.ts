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

import { Subscription } from '../../common/interfaces/subscription.interface';

import { Chart } from 'chart.js';
import {} from '@types/googlemaps';

import {html, render} from 'lit-html';

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
    map: any;

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

        // Fetch all subscriptions
        this.SubscriptionService.get().subscribe(
          (result: any) => {
              this.temp2 = [...result.Subscriptions];
              result.Subscriptions.forEach(subscription => {
                  this.streamService.getStream(subscription["id"]).subscribe(
                    (result: any) => {
                        const stream = result.Stream

                        // Create name and price field in susbcription
                        subscription["stream_name"] = stream["name"]
                        const mitasPrice = this.tasPipe.transform(stream["price"])
                        subscription["stream_price"] = mitasPrice

                        // Set markers on the map
                        this.setMarkers(stream);
                    },
                    err => {
                        console.log(err)
                    });

                    // Push marker to the markers list
                    this.mySubscriptions.push(subscription);
                });
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
          let mapOptions : any =  {

              // The latitude and longitude to center the map (always required)
              center: new google.maps.LatLng(48.86, 2.34), // Paris

              // How zoomed in you want the map to start at (always required)
              zoom: 6,
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

          // Get the HTML DOM element that will contain your map
          // We are using a div with id="map" seen below in the <body>
          let mapElement = document.getElementById("map");

          // Create the Google Map using our element and options defined above
          this.map = new google.maps.Map(mapElement, mapOptions);

          const that = this;
          google.maps.event.addListener(this.map, 'idle', function(ev){
              let bounds = that.map.getBounds();
              var southWestLng = bounds.getSouthWest().lng();
              var southWestLat = bounds.getSouthWest().lat();
              var northEastLng = bounds.getNorthEast().lng();
              var northEastLat = bounds.getNorthEast().lat();

              // Search streams on drawed region
              that.searchService.searchStreams(
                "geo", southWestLng, southWestLat, southWestLng, northEastLat,
                       northEastLng, northEastLat, northEastLng, southWestLat).subscribe(
                  (result: any) => {
                      that.temp = [...result.Streams];

                      result.Streams.forEach(stream => {
                          // TODO: set Sell Streams markers
                          //that.setMarkers(stream);
                      });
                  },
                  err => {
                      console.log(err)
                  });
          });

          // Automatically center the map fitting all markers on the screen
          //map.fitBounds(bounds);

      }

    // Display stream marker on a map
    setMarkers(stream) {
        const name = stream["name"];
        const lng = stream["location"]["coordinates"][1];
        const lat = stream["location"]["coordinates"][0];
        const position = new google.maps.LatLng(lng, lat);
        const mitasPrice = this.tasPipe.transform(stream["price"])
        const type = stream["type"]

        // Check the stream type and set proper icon
        let icon: string;
        var icons = {
            'temperature': 'Coke',
            'humidity': 'Pepsi',
            'air': 'Lemonade',
            'default': 'Default item'
        };
        switch(type) {
            case "temperature": {
                icon = "assets/img/icons/map-temp.svg";
                break;
            }
            case "humidity": {
                icon = "assets/img/icons/map-water.svg";
                break;
            }
            case "air": {
                icon = "assets/img/icons/map-co2.svg";
                break;
            }
            default: {
                icon = "assets/img/icons/map.svg";
                break;
            }
        }

        // Create new marker on the map
        let marker = new google.maps.Marker({
            position: position,
            map: this.map,
            title: name,
            icon: icon,
        });

        // Create new marker infowindow
        var infowindow = new google.maps.InfoWindow({
            content: `
                <div class="map-tooltip">
                  <p class="map-tooltip__title"> ${name} </p>
                  <div id="bodyContent" class="map-tooltip__content">
                    <p class="map-tooltip__subtitle"> ${type} </p>
                    <p class="map-tooltip__stake">
                      Stake: <span class="map-tooltip__stake-amount">
                      ${mitasPrice} TAS
                      </span>
                    </p>
                  </div>
                </div>
          `
        });

        // Set infowindow to marker
        marker.addListener('click', function() {
          infowindow.open(this.map, marker);
        });

    }
}
