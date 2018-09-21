import { Component, OnInit, ViewChild } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';

import { AuthService } from 'app/auth/services/auth.service';
import { ContractService } from 'app/dashboard/contracts/contract.service';
import { Contract } from 'app/common/interfaces/contract.interface';
import { Table, TableType } from 'app/shared/table/table';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { Page } from 'app/common/interfaces/page.interface';
import { DashboardContractsAddComponent } from './add/';
import { Query } from 'app/common/interfaces/query.interface';

@Component({
  selector: 'dpc-dashboard-contracts-list',
  templateUrl: './dashboard.contracts.component.html',
  styleUrls: [ './dashboard.contracts.component.scss' ]
})
export class DashboardContractsComponent implements OnInit {
    user: any;
    subscription: any;
    temp = [];
    selectedContract = [];
    table: Table = new Table();
    openedHelp = false;
    query = new Query();

    constructor(
        private modalService: BsModalService,
        private authService: AuthService,
        private contractService: ContractService,
        private alertService: AlertService,
    ) { }

  ngOnInit() {
    this.table.title = 'Contracts';
    this.table.tableType = TableType.Contract;
    this.table.headers = ['Stream Name', 'Share offered', 'Creation date',  'Expiration date', 'Status', ''];
    this.table.page = new Page<Contract>(0, 0, 0, []);

    // Fetch contracts with user as query owner
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.query.owner = this.user.id;
        this.contractService.get(this.query).subscribe(
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
      },
      err => {
        console.log(err);
      }
    );
  }

  toggleHelp() {
    this.openedHelp = !this.openedHelp;
  }

  modalNewContract() {
    // Make sure  help container is closed
    this.openedHelp = false;
    // Show DashboardSellAddComponent as Modal
    this.modalService.show(DashboardContractsAddComponent)
      .content.contractCreated.subscribe(
        response => {
          this.table.page.content.push(response);
        },
        err => {
          console.log(err);
        }
    );
  }
}
