import { Component, OnInit, ViewChild } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';

import { AccessService } from 'app/dashboard/access/access.service';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { AuthService } from 'app/auth/services/auth.service';
import { UserService } from 'app/common/services/user.service';
import { User } from 'app/common/interfaces/user.interface';
import { Access } from 'app/common/interfaces/access.interface';
import { Page } from 'app/common/interfaces/page.interface';
import { Table, TableType } from 'app/shared/table/table';
import { DashboardAccessAddComponent } from 'app/dashboard/access/add/dashboard.access.add.component';

@Component({
  selector: 'dpc-dashboard-access-list',
  templateUrl: './dashboard.access.component.html',
  styleUrls: ['./dashboard.access.component.scss']
})
export class DashboardAccessComponent implements OnInit {
    user: User;
    users: User[];
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
        private userService: UserService,
        private accessService: AccessService,
        private alertService: AlertService,
    ) { }

  ngOnInit() {
    this.table.title = 'Access';
    this.table.tableType = TableType.Access;
    this.table.headers = ['Partner', 'Origin', 'Status', ''];
    this.table.page = new Page<Access>(0, 0, 0, []);

    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      resp => {
        this.user = resp;
        this.fetchAccessRequests();
      },
      err => {
        this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
      }
    );
  }

  fetchAccessRequests() {
    this.userService.getAllUsers().subscribe(
      (respUsers: any) => {
        this.users = respUsers.users;

        // access requests sent
        this.accessService.getAllAccessSent().subscribe(
          (resp: any) => {
            this.reqsToTable(resp.Requests, 'sent');
          },
          err => {
            this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
          }
        );
        // access requests received
        this.accessService.getAllAccessReceived().subscribe(
          (resp: any) => {
            this.reqsToTable(resp.Requests, 'received');
          },
          err => {
            this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
          }
        );
      },
    );
  }

  reqsToTable(requests: any, origin: string) {
    // modify receiver field to show name instead of ID
    requests.forEach( req => {
      const index = this.users.findIndex(
        user => {
          if (origin === 'sent') {
            return req.receiver === user.id;
          } else {
            return req.sender === user.id;
          }
        }
      );

      // Set partner name origin
      req.partner = `${this.users[index].first_name} ${this.users[index].last_name}`;
      req.origin = origin;

      this.table.page.content.push(req);
    });
  }

  onPageChange(page: number) {
    this.query.page = page;
    this.fetchAccessRequests();
  }

  modalNewAccess() {
    // Make sure  help container is closed
    this.openedHelp = false;
    // Show DashboardAccessAddComponent as Modal
    this.modalService.show(DashboardAccessAddComponent)
      .content.accessCreated.subscribe(
        resp => {
          const row: any = {
            partner: `${resp.receiver.first_name} ${resp.receiver.last_name}`,
            state: 'pending',
            origin: 'sent',
          };
          this.table.page.content.push(row);
        },
      );
    }

    toggleHelp() {
      this.openedHelp = !this.openedHelp;
    }

    onAccessApproved(row: Access) {
      this.accessService.approveAccessRequest(row.id).subscribe(
        (resp: any) => {
          const index = this.table.page.content.findIndex(
            (access: Access) => row.id === access.id
          );
          const rowToUpdate = <Access> this.table.page.content[index];
           rowToUpdate.state = 'approved';
        },
        err => {
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );
    }

    onAccessRevoked(row: Access) {
      this.accessService.revokeAccessRequest(row.id).subscribe(
        (resp: any) => {
          const index = this.table.page.content.findIndex(
            (access: Access) => row.id === access.id
          );
          const rowToUpdate = <Access> this.table.page.content[index];
          rowToUpdate.state = 'revoked';
        },
        err => {
          this.alertService.error(`Status: ${err.status} - ${err.statusText}`);
        }
      );
    }
}
