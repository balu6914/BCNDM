import { Component, Input, Output, EventEmitter, OnChanges } from '@angular/core';

import { DpcPipe } from 'app/common/pipes/converter.pipe';

import * as L from 'leaflet';

@Component({
  selector: 'dpc-map-leaflet',
  templateUrl: './map.leaflet.component.html',
  styleUrls: ['./map.leaflet.component.scss']
})
export class MapComponent implements OnChanges {
  options = {
    layers: [
      L.tileLayer('https://api.mapbox.com/styles/v1/gesaleh/cjk0yt1dj8snd2sk6xqz29gha/tiles/256/{z}/{x}/{y}@2x?' +
      'access_token=pk.eyJ1IjoiZ2VzYWxlaCIsImEiOiJjamQ4bXFuZ3kybDZiMnhxcjl6Mjlmc3hmIn0.RVdSuXXmCgZJubeCAncjJQ', {})
    ],
    center: L.latLng({ lat: 48, lng: 2 }),
    zoom: 5,
    minZoom: 2,
    maxBounds: L.latLngBounds(L.latLng(-90, -180),  // southWest
                              L.latLng(90, 180)),   // northEast
    maxBoundsViscosity: 1
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
  firstPageLoad = true;
  tempId: any;
  drawLayerGroup = new L.LayerGroup();
  drawCreated = false;

  @Input() streamList: any;
  @Output() viewChanged: EventEmitter<any> = new EventEmitter();
  @Output() hoverMarker: EventEmitter<any> = new EventEmitter();
  constructor(
    private dpcPipe: DpcPipe,
  ) {
  }

  ngOnChanges() {
    this.addMarkers();
  }

  // Get stream type icon
  getIcon(type) {
    const icons = {
      chart: 'assets/img/icons/map-chart.svg',
      temperature: 'assets/img/icons/map-temp.svg',
      humidity: 'assets/img/icons/map-water.svg',
      air: 'assets/img/icons/map-co2.svg',
      default: 'assets/img/icons/map.svg'
    };
    return icons[type] || icons['default'];
  }

  // Get stream type icon
  getIconHover(type) {
    const icons = {
      chart: 'assets/img/icons/map-chart-red.svg',
      temperature: 'assets/img/icons/map-temp-red.svg',
      humidity: 'assets/img/icons/map-water-red.svg',
      air: 'assets/img/icons/map-co2-red.svg',
      default: 'assets/img/icons/map-red.svg'
    };
    return icons[type] || icons['default'];
  }

  onMapReady(map: L.Map) {
    this.map = map;
    this.map.addLayer(this.layerGroup);
    this.map.addLayer(this.drawLayerGroup);

    // Set markers and polygon layer inside the created polygon
    map.on('draw:created', (e: any) => {
      const layer = e.layer as L.Rectangle;
      this.drawLayerGroup.addLayer(layer);
      const drawBounds = layer.getBounds();
      this.viewChanged.emit(drawBounds);
      this.map.flyToBounds(drawBounds, {});
      this.drawCreated = true;
    });

    map.on('draw:drawstart', e => {
      this.drawLayerGroup.clearLayers();
    });
  }

  // This event is triggered when move or zoom is detected
  onMapMoveEnd() {
    if (!this.drawCreated) {
      this.drawLayerGroup.clearLayers();
      const bounds = this.map.getBounds();
      this.viewChanged.emit(bounds);
    }
  }

  // On click fetch always all streams in current view.
  onMapClick() {
    this.drawLayerGroup.clearLayers();
    const bounds = this.map.getBounds();
    this.viewChanged.emit(bounds);
    this.drawCreated = false;
  }

  onMapMouseMove(event: any) {
    this.hoverMarker.emit(this.tempId);
  }

  // Get marker bounds to apply flyToBounds effect on new map view
  focusMap() {
    const markers = this.layerGroup.getLayers();

    if (markers && markers.length) {
      const featureGroup = new L.FeatureGroup(markers);
      const bounds = featureGroup.getBounds();
      const options = {
        maxZoom: 8
      };

      this.map.flyToBounds(bounds, options);
    }
  }

  addMarkers() {
    if (this.map) {
      this.layerGroup.clearLayers();
      this.markers = [];

      this.streamList.forEach( stream => {
        this.addMarkerNoFocus(stream);
      });

      if (this.firstPageLoad) {
        // Fit map bounds
        this.focusMap();
        this.firstPageLoad = false;
      }
    }
  }

  addMarkerNoFocus(stream: any) {
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
            ${this.dpcPipe.transform(stream.price)} DPC
          </span>
          </p>
        </div>
      </div>
    `;

    newMarker.bindPopup(msg);

    // Add marker as a map layer
    this.layerGroup.addLayer(newMarker);

    // Get ID of the marker and add it to stream object
    stream.mapId = this.layerGroup.getLayerId(newMarker);

    // TODO: Use (leafletMouseMove)="onMapMouseMove($event)"
    newMarker.on('mouseover', () => {
      this.tempId = stream.mapId;
    });
    newMarker.on('mouseout', () => {
      this.tempId = 0;
    });

    // Add stream to markers list
    this.markers.push(stream);
  }

  // Callback of addEvt from TableComponent
  addMarker(stream) {
    // Do map focus only when adding streams from ui
    this.addMarkerNoFocus(stream);
    this.focusMap();
  }

  // Callback of editEvt from TableComponent
  editMarker(stream) {
    this.removeMarker(stream.id);
    this.addMarker(stream);
  }

  // Callback of deleteEvt from TableComponent
  removeMarker(id) {
    this.markers.forEach( (stream, i) => {
      if (stream.id === id) {
        const markerLayer = this.layerGroup.getLayer(stream.mapId);
        this.layerGroup.removeLayer(markerLayer);
        this.markers.splice(i, 1);
      }
    });

    // Fit map bounds
    this.focusMap();
  }

  mouseHoverMarker(row: any) {
    const markerLayer = <L.Marker>this.layerGroup.getLayer(row.mapId);
    const icon = L.icon({
        iconUrl:  this.getIconHover(row.type),
        iconSize: [45, 45]
    }) ;
    markerLayer.setIcon(icon);
  }

  mouseUnhoverMarker(row: any) {
    const markerLayer = <L.Marker>this.layerGroup.getLayer(row.mapId);
    const icon = L.icon({
        iconUrl:  this.getIcon(row.type),
        iconSize: [45, 45]
    }) ;
    markerLayer.setIcon(icon);
  }

}
