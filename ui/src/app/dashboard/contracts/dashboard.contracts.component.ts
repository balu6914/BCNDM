import { Component, ViewChild } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';

import { Table, TableType } from 'app/shared/table/table';
import { Contract } from 'app/common/interfaces/contract.interface'
import { Page } from 'app/common/interfaces/page.interface';
import { DashboardContractsAddComponent } from './add/';

@Component({
  selector: 'dashboard-contracts-list',
  templateUrl: './dashboard.contracts.component.html',
  styleUrls: [ './dashboard.contracts.component.scss' ]
})
export class DashboardContractsComponent {
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
        private modalService: BsModalService,
    ) { }

  ngOnInit() {
    this.table.title = "Contracts";
    this.table.tableType = TableType.Contract;
    this.table.headers = ["Stream Name", "Tokens per hour", "Share offered", "Expiration date", "Status"];
    this.table.page = new Page<Contract>(0, 20, 10, this.myContractsList);

    this.temp = [...this.myContractsList];
  }

  modalNewContract() {
    // Show DashboardSellAddComponent as Modal
    this.modalService.show(DashboardContractsAddComponent)
      .content.contractCreated.subscribe(
        res => {
          // TODO: Use all values from response to create the contract
          const contract = {
            id: 'myID',
            stream: {
              name: res.streamName,
              price: res.streamPrice
            },
            creationDate: "2018-02-15T12:14:56.806Z",
            expirationDate: '2018-02-15T12:14:56.806Z',
            share: res.shareOffered,
            signed: false,
            expired: true
          }
          this.myContractsList.push(contract);
        },
        err => {
          console.log(err);
        }
    );
  }

}
