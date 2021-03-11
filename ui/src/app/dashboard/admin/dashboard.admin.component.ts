import { Component, OnInit, ViewChild } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';
import { BsModalRef } from 'ngx-bootstrap/modal/bs-modal-ref.service';

import { UserService } from 'app/common/services/user.service';
import { AuthService } from 'app/auth/services/auth.service';
import { Query } from 'app/common/interfaces/query.interface';
import { Table, TableType } from 'app/shared/table/table';
import { DashboardAdminSignupComponent } from 'app/dashboard/admin/signup/dashboard.admin.signup.component';
import { User } from 'app/common/interfaces/user.interface';
import { TableComponent } from 'app/shared/table/main/table.component';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { BalanceService } from 'app/shared/balance/balance.service';

import { DashboardAdminEditComponent } from 'app/dashboard/admin/edit/dashboard.admin.edit.component';
import { DashboardAdminDeleteComponent } from 'app/dashboard/admin/delete/dashboard.admin.delete.component';
import { DashboardAdminLockComponent } from 'app/dashboard/admin/lock/dashboard.admin.lock.component';

@Component({
  selector: 'dpc-dashboard-admin',
  templateUrl: './dashboard.admin.component.html',
  styleUrls: ['./dashboard.admin.component.scss']
})
export class DashboardAdminComponent implements OnInit {
  admin: User;
  temp = [];
  users = [];
  table: Table = new Table();
  query = new Query();
  bsModalRef: BsModalRef;

  @ViewChild('tableComponent')
  private tableComponent: TableComponent;

  constructor(
    private authService: AuthService,
    private userService: UserService,
    private modalService: BsModalService,
    public alertService: AlertService,
    private balanceService: BalanceService,
  ) {
  }

  ngOnInit() {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      resp => {
        this.admin = resp;
        this.query.owner = this.admin.id;
        this.fetchUsers();
      },
      err => {
        console.log(err);
      }
    );

    // Config table
    this.table.title = 'Users';
    this.table.tableType = TableType.Users;
    this.table.headers = ['Role', 'Email', 'Name', 'Phone', 'Address', 'Company', 'Account Login', 'Status', 'Balance', ''];
    this.table.hasDetails = true;
  }

  openModalSignup() {
    // Open DashboardAdminSignupComponent as Modal
    this.modalService.show(DashboardAdminSignupComponent)
      .content.userCreated.subscribe(
        user => {
          user.balance = 0;
          // Add created user to the table
          this.table.page.content.push(user);
        },
        err => {
          console.log(err);
        }
      );
  }

  editUser(row: User) {
    const initialState = {
      user: {
        id: row.id,
        email: row.email,
        first_name: row.first_name,
        last_name: row.last_name,
        phone: row.phone,
        address: row.address,
        company: row.company,
        balance: row.balance,
        role: row.role,
      },
    };

    // Open DashboardSellDeleteComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardAdminEditComponent, {initialState})
    .content.userEdited.subscribe(
      resp => {
        const itemIndex = this.table.page.content.findIndex((item: any) => item.id === row.id);
        this.table.page.content[itemIndex] = resp;
      }
    );
  }

  lockUser(row: User) {
    const initialState = {
      user: {
        id: row.id,
        email: row.email,
        first_name: row.first_name,
        last_name: row.last_name,
        phone: row.phone,
        address: row.address,
        company: row.company,
        disabled: row.disabled,
        locked: row.locked,
        balance: row.balance,
        role: row.role,
      },
    };

    // Open DashboardAdminLockComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardAdminLockComponent, {initialState})
    .content.userLocked.subscribe(
      resp => {
        const itemIndex = this.table.page.content.findIndex((item: any) => item.id === resp.id);
        row.locked = resp;
        this.table.page.content[itemIndex] = row;
      }
    );
  }

  deleteUser(row: User) {
    const initialState = {
      user: {
        id: row.id,
        email: row.email,
        first_name: row.first_name,
        last_name: row.last_name,
        phone: row.phone,
        address: row.address,
        company: row.company,
        disabled: row.disabled,
        balance: row.balance,
        role: row.role,
      },
    };

    // Open DashboardAdminDeleteComponent as Modal
    this.bsModalRef = this.modalService.show(DashboardAdminDeleteComponent, {initialState})
    .content.userDeleted.subscribe(
      resp => {
        const itemIndex = this.table.page.content.findIndex((item: any) => item.id === resp.id);
        row.disabled = resp;
        this.table.page.content[itemIndex] = row;
      }
    );
  }

  fetchUsers() {
    this.userService.getAllUsers().subscribe(
      (resp: any) => {
        if (resp.users.length > 0) {
          resp.users.forEach((u: User) => {
            this.balanceService.getBalance(u.id).subscribe(
              (respBalance: any) => {
                u.balance = respBalance.balance;
                this.table.page.content.push(u);
              },
            );
          });
        }
      },
      err => {
        console.log(err);
      }
    );
  }

  onPageChange(page: number) {
    this.query.page = page;
    this.fetchUsers();
  }

  onHoverRow(row) {
  }

  onUnhoverRow(row) {
  }

}
