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
      {'id' : "1", 'stream': {'name': 'WeIO board pressure', 'price':'100'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-12-15T12:14:56.806Z', 'share':'10','signed':false, 'expired': false},
      {'id' : "2", 'stream': {'name': 'WeIO board temperature', 'price':'10000'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-11-15T12:14:56.806Z', 'share':'15','signed':true, 'expired': false },
      {'id' : "3", 'stream': {'name': 'WeIO board humidity', 'price':'15'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-12-15T12:14:56.806Z', 'share':'50', 'signed':false, 'expired': false},
      {'id' : "4", 'stream': {'name': 'aqi', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},

      {'id' : "5", 'stream': {'name': 'vibration', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "6", 'stream': {'name': 'temperature', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "7", 'stream': {'name': 'humidity', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "8", 'stream': {'name': 'humidity_20', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "9", 'stream': {'name': 'radiation', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "10", 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "11", 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "12", 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "13", 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "14", 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
      {'id' : "15", 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'expirationDate': '2018-07-15T12:14:56.806Z', 'share':'5', 'signed':true, 'expired': true},
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
    openedHelp: boolean = false;

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

  toggleHelp() {
    this.openedHelp = !this.openedHelp;
    console.log("toggle trigger, opened: ", this.openedHelp);
  }

  modalNewContract() {
    // Make sure  help container is closed
    this.openedHelp = false;
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
