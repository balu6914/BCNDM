import { Component, OnInit, ViewChild } from '@angular/core';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { BsModalService } from 'ngx-bootstrap/modal';
import { AuthService } from 'app/auth/services/auth.service';
import { Query } from 'app/common/interfaces/query.interface';
import { Table, TableType } from 'app/shared/table/table';
import { DashboardAdminSignupComponent } from 'app/dashboard/admin/signup/dashboard.admin.signup.component';
import { User } from 'app/common/interfaces/user.interface';
import { TableComponent } from 'app/shared/table/main/table.component';

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

  @ViewChild('tableComponent')
  private tableComponent: TableComponent;

  constructor(
    private authService: AuthService,
    private modalService: BsModalService,
    public alertService: AlertService,
  ) {
  }

  ngOnInit() {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.admin = data;
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
    this.table.headers = ['Email', 'Name', 'Last Name', 'Phone', 'Address', 'Company', ''];
    this.table.hasDetails = true;
  }

  openModalSignup() {
    // Open DashboardAdminSignupComponent as Modal
    this.modalService.show(DashboardAdminSignupComponent)
      .content.userCreated.subscribe(
        user => {
          // Add created user to the table
          this.table.page.content.push(user);
        },
        err => {
          console.log(err);
        }
      );
  }

  editUser(user: User) {
    // TODO: edit User
  }

  deleteUser(userID: string) {
    // TODO: delete User
  }

  fetchUsers() {
    // TODO: fetch list of Users
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
