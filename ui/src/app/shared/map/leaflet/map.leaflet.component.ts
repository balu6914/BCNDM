import { Component, Input } from "@angular/core";

import { TasPipe } from "app/common/pipes/converter.pipe";

import * as L from 'leaflet';

@Component({
  selector: "dpc-map-leaflet",
  templateUrl: "./map.leaflet.component.html",
  styleUrls: ["./map.leaflet.component.scss"]
})
export class MapComponent {
  options = {
    layers: [
      L.tileLayer('https://api.mapbox.com/styles/v1/gesaleh/cjk0yt1dj8snd2sk6xqz29gha/tiles/256/{z}/{x}/{y}@2x?access_token=pk.eyJ1IjoiZ2VzYWxlaCIsImEiOiJjamQ4bXFuZ3kybDZiMnhxcjl6Mjlmc3hmIn0.RVdSuXXmCgZJubeCAncjJQ', {})
    ],
    center: L.latLng({ lat: 48, lng: 2 }),
    zoom: 5
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
  map: L.Map;
  layerGroup = new L.LayerGroup();
  markers = [];
  activeMarkers = [];

  @Input() streamList: any;
  constructor(
    private tasPipe: TasPipe,
  ) {
  }

  // Get stream type icon
  getIcon(type) {
    const icons = {
      temperature: "assets/img/icons/map-temp.svg",
      humidity: "assets/img/icons/map-water.svg",
      air: "assets/img/icons/map-co2.svg",
      default: "assets/img/icons/map.svg"
    };
    return icons[type] || icons["default"];
  }

  onMapReady(map: L.Map) {
    this.map = map;
    this.map.addLayer(this.layerGroup);
    this.addMarkers();
  }

  ngOnChanges() {
    this.addMarkers();
  }

  focusMap() {
    // Get Markers range Bounts and set map to this view in order to always
    // see all markers in focus on map Init.
    const that = this;
    this.activeMarkers = this.markers.map(m => {
      const a = [];
      const lat = m.location.coordinates[1];
      const long = m.location.coordinates[0];
      if(lat && long) {
        a.push(L.latLng(lat, long));
      }
      return a;
    });
    Promise.all(this.activeMarkers).then(result => {
      if(this.activeMarkers.length && this.map) {
        var b = new L.LatLngBounds(this.activeMarkers);
        this.map.fitBounds(b);
      }
    });
  }

  addMarkers() {
    if(this.map) {
      this.map.removeLayer(this.layerGroup);
      this.layerGroup = new L.LayerGroup();
      this.map.addLayer(this.layerGroup);
      this.markers = [];
      this.streamList.forEach( stream => {
        this.addMarker(stream);
      });
      this.focusMap();
    }
  }

  addMarker(stream: any) {
    // Create marker and set icon
    const newMarker = L.marker(
      [stream.location.coordinates[1],
       stream.location.coordinates[0]], {
         icon: L.icon({
             iconUrl:  this.getIcon(stream.type),
             iconSize: [45, 45]
         })
       }
    );
    // Create popup message and add it to marker
    const msg = `
      <div class="map-tooltip">
        <p class="map-tooltip__title"> ${stream.name} </p>
        <div id="bodyContent" class="map-tooltip__content">
          <p class="map-tooltip__subtitle"> ${stream.type} </p>
          <p class="map-tooltip__stake">
            Stake: <span class="map-tooltip__stake-amount">
            ${this.tasPipe.transform(stream.price)} TAS
          </span>
          </p>
        </div>
      </div>
    `
    newMarker.bindPopup(msg);

    // Add marker as a map layer
    this.layerGroup.addLayer(newMarker);

    // Get ID of the marker and add it to stream object
    stream.mapId = this.layerGroup.getLayerId(newMarker);

    // Add stream to markers list
    this.markers.push(stream);
  }

  // Callback of editEvt from TableComponent
  editMarker(stream) {
    this.removeMarker(stream.id);
    this.addMarker(stream);
  }

  // Callback of deleteEvt from TableComponent
  removeMarker(id) {
    this.markers.forEach( (stream, i) => {
      if (stream.id == id) {
        const markerLayer = this.layerGroup.getLayer(stream.mapId);
        this.layerGroup.removeLayer(markerLayer);
        this.markers.splice(i,1);
      }
    });
  }
}
