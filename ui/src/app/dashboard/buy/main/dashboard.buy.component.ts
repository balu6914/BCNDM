import { Component, OnInit, ViewChild } from '@angular/core';
import { AuthService } from 'app/auth/services/auth.service';
import { Query } from 'app/common/interfaces/query.interface';
import { StreamService } from 'app/common/services/stream.service';
import { Table, TableType } from 'app/shared/table/table';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { MapComponent } from 'app/shared/map/leaflet/map.leaflet.component';
import { TableComponent } from 'app/shared/table/main/table.component';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap';
import { DashboardBuyGroupComponent } from 'app/dashboard/buy/group/dashboard.buy.group.component';

@Component({
  selector: 'dpc-dashboard-buy',
  templateUrl: './dashboard.buy.component.html',
  styleUrls: ['./dashboard.buy.component.scss']
})
export class DashboardBuyComponent implements OnInit {
  user: any;
  table: Table = new Table();
  query = new Query();

  @ViewChild('map')
  private map: MapComponent;

  @ViewChild('tableComponent')
  private tableComponent: TableComponent;

  constructor(
    private AuthService: AuthService,
    public streamService: StreamService,
    public alertService: AlertService,
    private modalService: BsModalService,
    public midpcPipe: MidpcPipe,
  ) { }

  ngOnInit() {
    this.table.title = 'Streams';
    this.table.tableType = TableType.Buy;
    this.table.headers = ['Stream Name', 'Stream Type', 'Stream Price', ''];
    this.table.hasDetails = true;

    // Fetch current User
    this.user = {};
    this.AuthService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.query.owner = '-'.concat(this.user.id);
        this.fetchStreams();
      },
      err => {
        console.log(err);
      });

      this.map.viewChanged.subscribe(
        bounds => {
          this.query.setPoint('x0', bounds["_southWest" ]["lng"]);
          this.query.setPoint('y0', bounds["_southWest" ]["lat"]);
          this.query.setPoint('x1', bounds["_southWest" ]["lng"]);
          this.query.setPoint('y1', bounds["_northEast" ]["lat"]);

          this.query.setPoint('x2', bounds["_northEast" ]["lng"]);
          this.query.setPoint('y2', bounds["_northEast" ]["lat"]);
          this.query.setPoint('x3', bounds["_northEast" ]["lng"]);
          this.query.setPoint('y3', bounds["_southWest" ]["lat"]);
          this.fetchStreams();
        },
        err => {
          console.log(err);
        }
      );
  }

  onPageChange(page: number) {
    this.query.page = page;
    this.fetchStreams();
  }

  onBuyClick(){     
    const initialState = {
      streamsList: this.table.page.content,
    };
    this.modalService.show(DashboardBuyGroupComponent, { initialState });
  }
  
  onFiltersChange(filters: any) {
    filters.minPrice = this.midpcPipe.transform(filters.minPrice);
    filters.maxPrice = this.midpcPipe.transform(filters.maxPrice);
    Object.assign(this.query, filters);
    this.fetchStreams();
  }

  private fetchStreams() {
    this.streamService.searchStreams(this.query).subscribe(
      (result: any) => {
        const temp = Object.assign({}, this.table);
        temp.page = result;
        // Set table content
        this.table = temp;
      },
      err => {
        this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
      });
  }

  onHoverRow(row) {
    this.map.mouseHoverMarker(row);
  }

  onUnhoverRow(row) {
    this.map.mouseUnhoverMarker(row);
  }

  onHoverMarker(streamId) {
    this.tableComponent.activateRow(streamId);
  }

}
