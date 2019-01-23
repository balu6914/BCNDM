import { Component, OnInit, ViewChild } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';

import { AuthService } from 'app/auth/services/auth.service';
import { AccessService } from 'app/dashboard/access/access.service';
import { Access } from 'app/common/interfaces/access.interface';
import { Table, TableType } from 'app/shared/table/table';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { Page } from 'app/common/interfaces/page.interface';
import { DashboardAccessAddComponent } from './add/';

@Component({
  selector: 'dpc-dashboard-access-list',
  templateUrl: './dashboard.access.component.html',
  styleUrls: ['./dashboard.access.component.scss']
})
export class DashboardAccessComponent implements OnInit {
    user: any;
    subscription: any;
    temp = [];
    selectedAccess = [];
    table: Table = new Table();
    openedHelp = false;
    query = {
      isOwner: true,
      isPartner: true,
      page: 0
    };

    constructor(
        private modalService: BsModalService,
        private authService: AuthService,
        private accessService: AccessService,
        private alertService: AlertService,
    ) { }

  ngOnInit() {
    this.table.title = 'Access';
    this.table.tableType = TableType.Access;
    this.table.headers = ['Partners', 'Status'];
    this.table.page = new Page<Access>(0, 0, 0, [
      {partner_id: 'Air France', signed: false},
      {partner_id: 'Korean Air', signed: true},
    ]);

    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        // this.fetchAccess();
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );
  }

  fetchAccess() {
    this.accessService.get(this.query).subscribe(
      result => {
        const temp = Object.assign({}, this.table);
        temp.page = result;
        // Set table content
        this.table = temp;
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );
  }

  onPageChange(page: number) {
    this.query.page = page;
    this.fetchAccess();
  }

  modalNewAccess() {
    // Make sure  help container is closed
    this.openedHelp = false;
    // Show DashboardSellAddComponent as Modal
    this.modalService.show(DashboardAccessAddComponent)
      .content.accessCreated.subscribe(
        response => {
          response.items.forEach(partner => {
            console.log(partner);
            const row: Access = {
              partner_id: partner.partner_id,
              signed: false,
            };
            this.table.page.content.push(row);
          });
        },
        err => {
          console.log(err);
        }
      );
    }

    toggleHelp() {
      this.openedHelp = !this.openedHelp;
    }
}
