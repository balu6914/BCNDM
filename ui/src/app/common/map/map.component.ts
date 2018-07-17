import { Component, OnInit } from "@angular/core";

import { TasPipe } from "../../common/pipes/converter.pipe";
import mapStyle from './map-style.json';

import {} from "@types/googlemaps";

@Component({
  selector: "dpc-map",
  templateUrl: "./map.component.html",
  styleUrls: ["./map.component.css"]
})
export class MapComponent implements OnInit {
  streamList = [];
  subscriptionList = [];
  temp = [];
  map: any;

  constructor(private tasPipe: TasPipe) {}

  // Create Google Map
  create(mapElement: HTMLElement) {
    // Basic options for a simple Google Map
    // For more options see: https://developers.google.com/maps/documentation/javascript/reference#MapOptions
    let mapOptions: any = {
      // The latitude and longitude to center the map (always required)
      center: new google.maps.LatLng(48.86, 2.34), // Paris

      // How zoomed in you want the map to start at (always required)
      zoom: 6,
      disableDefaultUI: true,
      zoomControl: true,

      // How you would like to style the map.
      // This is where you would paste any style found on Snazzy Maps.
      styles: mapStyle,
    };
    // Create the Google Map using our element and options defined above
    this.map = new google.maps.Map(mapElement, mapOptions);

    const that = this;
    google.maps.event.addListener(this.map, "idle", function(ev) {
      let bounds = that.map.getBounds();
      var southWestLng = bounds.getSouthWest().lng();
      var southWestLat = bounds.getSouthWest().lat();
      var northEastLng = bounds.getNorthEast().lng();
      var northEastLat = bounds.getNorthEast().lat();

      that.streamList.forEach(stream => {
        // Set markers on the map
        that.addMarker(stream);
      });
    });

    // Automatically center the map fitting all markers on the screen
    //map.fitBounds(bounds);
  }

  // Get stream type icon
  getIcon(type) {
    var icons = {
      temperature: "assets/img/icons/map-temp.svg",
      humidity: "assets/img/icons/map-water.svg",
      air: "assets/img/icons/map-co2.svg",
      default: "assets/img/icons/map.svg"
    };
    return icons[type] || icons["default"];
  }

  // Display stream marker on a map
  addMarker(stream) {
    const name = stream["name"];
    const lng = stream["location"]["coordinates"][1];
    const lat = stream["location"]["coordinates"][0];
    const position = new google.maps.LatLng(lng, lat);
    const mitasPrice = this.tasPipe.transform(stream["price"]);
    const type = stream["type"];
    const icon = this.getIcon(type);

    // Create new marker on the map
    let marker = new google.maps.Marker({
      position: position,
      map: this.map,
      title: name,
      icon: icon
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
    marker.addListener("click", function() {
      infowindow.open(this.map, marker);
    });
  }

  setStreamList(list: any) {
    this.streamList = list;
  }

  ngOnInit() {}
}
