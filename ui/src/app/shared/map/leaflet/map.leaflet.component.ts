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
      L.tileLayer('https://api.mapbox.com/styles/v1/gesaleh/cjk0yrl8kaavm2smu4bx3okrh/tiles/256/{z}/{x}/{y}?access_token=pk.eyJ1IjoiZ2VzYWxlaCIsImEiOiJjamQ4bXFuZ3kybDZiMnhxcjl6Mjlmc3hmIn0.RVdSuXXmCgZJubeCAncjJQ', {})
    ],
    zoom: 5,
    center: L.latLng({ lat: 48, lng: 2 })
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

  @Input() streamList: any;
  constructor(
    private tasPipe: TasPipe,
  ) {}

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
    const that = this;

    map.on('move', function(e) {
      that.streamList.forEach(stream => {
        // Create marker and set icon
        const newMarker = L.marker(
          [stream.location.coordinates[1],
           stream.location.coordinates[0]], {
             icon: L.icon({
                 iconUrl:  that.getIcon(stream.type),
                 iconSize: [50, 50]
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
                ${that.tasPipe.transform(stream.price)} TAS
              </span>
              </p>
            </div>
          </div>
        `
        newMarker.bindPopup(msg);

        // Add marker as a map layer
        map.addLayer(newMarker);
      });
    });

    // TODO: Fix map.on('load') event (don't set view to force move event)
    map.setView([48.864716, 2.349014], 5);
  }
}
