import { Component, OnInit, ViewChild } from '@angular/core';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { BsModalService } from 'ngx-bootstrap/modal';
import { AuthService } from 'app/auth/services/auth.service';
import { Query } from 'app/common/interfaces/query.interface';
import { StreamService } from 'app/common/services/stream.service';
import { MapComponent } from 'app/shared/map/leaflet/map.leaflet.component';
import { Table, TableType } from 'app/shared/table/table';
import { DashboardSellAddComponent } from 'app/dashboard/sell/add/dashboard.sell.add.component';
import { User } from 'app/common/interfaces/user.interface';
import { TableComponent } from 'app/shared/table/main/table.component';
import { MidpcPipe } from 'app/common/pipes/converter.pipe';
@Component({
  providers: [MapComponent],
  selector: 'dpc-dashboard-sell',
  templateUrl: './dashboard.sell.component.html',
  styleUrls: ['./dashboard.sell.component.scss']
})
export class DashboardSellComponent implements OnInit {
  user: User;
  temp = [];
  streams = [];
  table: Table = new Table();
  query = new Query();

  @ViewChild('map')
  private map: MapComponent;

  @ViewChild('tableComponent')
  private tableComponent: TableComponent;

  constructor(
    private streamService: StreamService,
    private AuthService: AuthService,
    private modalService: BsModalService,
    public alertService: AlertService,
    public midpcPipe: MidpcPipe,
  ) {
  }

  ngOnInit() {
    // Fetch current User
    this.AuthService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.query.owner = this.user.id;
        this.fetchStreams();
      },
      err => {
        console.log(err);
      }
    );

    // Config table
    this.table.title = 'Streams';
    this.table.tableType = TableType.Sell;
    this.table.headers = ['Stream Name', 'Stream Type', 'Stream Price', ''];
    this.table.hasDetails = true;

    this.map.viewChanged.subscribe(
      bounds => {
        this.query.setPoint('x0', bounds['_southWest' ]['lng']);
        this.query.setPoint('y0', bounds['_southWest' ]['lat']);
        this.query.setPoint('x1', bounds['_southWest' ]['lng']);
        this.query.setPoint('y1', bounds['_northEast' ]['lat']);

        this.query.setPoint('x2', bounds['_northEast' ]['lng']);
        this.query.setPoint('y2', bounds['_northEast' ]['lat']);
        this.query.setPoint('x3', bounds['_northEast' ]['lng']);
        this.query.setPoint('y3', bounds['_southWest' ]['lat']);
        this.fetchStreams();
      },
      err => {
        console.log(err);
      }
    );
  }

  // Add Bulk event
  onFileChange(event): void {
    const fileList: FileList = event.target.files;
    if (fileList.length > 0) {
      const file = fileList[0];

      const formData = new FormData();
      formData.append('csv', file, file.name);

      this.streamService.addStreamBulk(formData).subscribe(
        data => {
          this.fetchStreams();
          this.alertService.success(`CSV successfully uploaded`);
        },
        err => {
          if (err.status === 400) {
            this.alertService.error(`Error with CSV file format`);
          } else {
            this.alertService.error(`Status: ${err} - ${err.statusText}`);
          }
        }
      );
    }
  }

  openModalAdd() {
    const gmailSuffix = '@gmail.com';
    const initialState = {
      bqMail: this.user.email.toLowerCase().endsWith(gmailSuffix) ||
        this.user.contact_email.toLowerCase().endsWith(gmailSuffix)
    };

    // Show DashboardSellAddComponent as Modal
    this.modalService.show(DashboardSellAddComponent, { initialState })
      .content.streamCreated.subscribe(
        stream => {
          // Push new stream to table
          this.table.page.content.push(stream);
          // Set MArker on the map
          this.map.addMarker(stream);
        },
        err => {
          console.log(err);
        }
      );
  }

  editStream(stream) {
    this.map.editMarker(stream);
  }

  deleteStream(id) {
    this.map.removeMarker(id);
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

  onPageChange(page: number) {
    this.query.page = page;
    this.fetchStreams();
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
