import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { AuthService } from '../../../auth/services/auth.service';
import { Router } from '@angular/router';
import { Table, TableType } from '../../../shared/table/table';
import { Contract } from '../../../common/interfaces'

@Component({
  selector: 'dashboard-contracts-list',
  templateUrl: './dashboard.contracts.list.component.html',
  styleUrls: [ './dashboard.contracts.list.component.scss' ]
})
export class DashboardContractsListComponent {
    tableColumns = [
        { prop: 'name' },
        { name: 'type' },
        { name: 'description' },
        { name: 'price' }
    ];
    // TODO: Remove this Mock contracts list for demo purpose
    myContractsList = [
      {'id' : "1", 'stream': {'name': 'WeIO pressure', 'price':'100'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-02-15T12:14:56.806Z', 'share':'10','signed':false, 'expired': false},
      {'id' : "2", 'stream': {'name': 'WeIO temperature', 'price':'10000'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-02-15T12:14:56.806Z', 'share':'10','signed':true, 'expired': false },
      {'id' : "3", 'stream': {'name': 'WeIO humidity', 'price':'15'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expired': false},
      {'id' : "4", 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':true, 'expired': false},
      {'id' : "5", 'stream': {'name': 'WeIO radiation', 'price':'50'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expired': false},
      {'id' : "5", 'stream': {'name': 'Spark', 'price':'30'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expired': true}
    ];
    user: any;
    subscription: any;
    temp = [];
    selectedContract = [];
    // ngx-table custom messages
    tableMessages =  {
        emptyMessage: "You don't have any smart contracts yet..."
    }
    table: Table = new Table();

    constructor(
        private AuthService: AuthService,
        private router: Router
    ) { }

  ngOnInit() {
    this.table.title = "Contracts";
    this.table.tableType = TableType.Contract;
    this.table.headers = ["Stream Name", "Tokens per hour", "Share offered", "Expiration date", "Status"];
    this.table.content = this.myContractsList;

    this.AuthService.getCurrentUser().subscribe(
      data => {
        this.user = data;
      },
      err => {
        console.log(err)
      }
    );

    this.temp = [...this.myContractsList];
  }
}
