import { Component, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { BsModalService } from 'ngx-bootstrap/modal';

import { DashboardSellAddComponent } from './add/dashboard.sell.add.component';
import { StreamService } from '../../common/services/stream.service';
import { AuthService } from '../../auth/services/auth.service';
import { Table, TableType } from '../../shared/table/table';
import { Query } from '../../common/interfaces/query.interface';
import { Page } from '../../common/interfaces/page.interface';
import { Stream } from '../../common/interfaces';
import { MapComponent } from '../../shared/map/leaflet/map.leaflet.component';

@Component({
  selector: 'dashboard-sell',
  templateUrl: './dashboard.sell.component.html',
  styleUrls: ['./dashboard.sell.component.scss']
})
export class DashboardSellComponent {
  user: any;
  temp = [];
  streams = [];
  table: Table = new Table();

  @ViewChild('map')
  private map: MapComponent;

  constructor(
    private streamService: StreamService,
    private AuthService: AuthService,
    private router: Router,
    private modalService: BsModalService,
  ) {
  }

  ngOnInit() {
    // Fetch current User
    this.user = {};
    this.AuthService.getCurrentUser().subscribe(
      data => {
        this.user = data;
      },
      err => {
        console.log(err)
      }
    );

    // Config table
    this.table.title = 'Streams';
    this.table.tableType = TableType.Sell;
    this.table.headers = ['Stream Name', 'Stream Type', 'Stream Price'];
    this.table.hasDetails = true;

    const query = new Query();

    this.streamService.searchStreams(query).subscribe(
      (result: Page<Stream>) => {
        this.temp = result.content;
        result.content.forEach(stream => {
          if (stream.owner === this.user.id) {
            this.streams.push(stream);
          }
        });
        result.content = this.streams;
        // Set table content
        this.table.page = result;
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
        data => this.router.navigate(['/dashboard/sell']),
        error => console.log(error),
      );
    }
  }

  openModalAdd() {
    // Show DashboardSellAddComponent as Modal
    this.modalService.show(DashboardSellAddComponent)
      .content.streamCreated.subscribe(
        stream => {
          // Push new stream to table
          this.streams.push(stream);
          // Set MArker on the map
          this.map.addMarker(stream);
        },
        err => {
          console.log(err)
        }
    );
  }

  editStream(stream) {
    this.map.editMarker(stream);
  }

  deleteStream(id) {
    this.map.removeMarker(id);
  }
}
