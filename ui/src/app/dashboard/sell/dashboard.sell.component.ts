import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { BsModalService } from 'ngx-bootstrap/modal';

import { DashboardSellAddComponent } from './add/dashboard.sell.add.component';
import { StreamService } from '../../common/services/stream.service';
import { AuthService } from '../../auth/services/auth.service';
import { Table, TableType } from '../../shared/table/table';
import { Query } from '../../common/interfaces/query.interface';
import { Page } from '../../common/interfaces/page.interface';
import { Stream } from '../../common/interfaces';

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

  openModal() {
    // Open DashboardSellAddComponent Modal
    this.modalService.show(DashboardSellAddComponent);
  }
}
