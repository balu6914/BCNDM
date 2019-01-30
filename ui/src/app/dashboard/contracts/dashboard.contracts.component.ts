import { Component, OnInit, ViewChild } from '@angular/core';
import { BsModalService } from 'ngx-bootstrap/modal';

import { AuthService } from 'app/auth/services/auth.service';
import { ContractService } from 'app/dashboard/contracts/contract.service';
import { Contract } from 'app/common/interfaces/contract.interface';
import { Table, TableType } from 'app/shared/table/table';
import { AlertService } from 'app/shared/alerts/services/alert.service';
import { Page } from 'app/common/interfaces/page.interface';
import { DashboardContractsAddComponent } from './add/';

@Component({
  selector: 'dpc-dashboard-contracts-list',
  templateUrl: './dashboard.contracts.component.html',
  styleUrls: ['./dashboard.contracts.component.scss']
})
export class DashboardContractsComponent implements OnInit {
    user: any;
    subscription: any;
    temp = [];
    selectedContract = [];
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
        private contractService: ContractService,
        private alertService: AlertService,
    ) { }

  ngOnInit() {
    this.table.title = 'Contracts';
    this.table.tableType = TableType.Contract;
    this.table.headers = ['Stream Name', 'Share offered', 'Creation date',  'Expiration date', 'Status', ''];
    this.table.page = new Page<Contract>(0, 0, 0, []);

    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        this.fetchContracts();
      },
      err => {
        this.alertService.error(`Error: ${err.status} - ${err.statusText}`);
      }
    );
  }

  fetchContracts() {
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
  }

  onPageChange(page: number) {
    this.query.page = page;
    this.fetchContracts();
  }

  modalNewContract() {
    // Make sure  help container is closed
    this.openedHelp = false;
    // Show DashboardContractsAddComponent as Modal
    this.modalService.show(DashboardContractsAddComponent)
      .content.contractCreated.subscribe(
        response => {
          response.parties.forEach(partner => {
            const row: Contract = {
              stream_name: response.stream_name,
              start_time: response.start_time,
              end_time: response.end_time,
              partner_id: partner.partner_id,
              share: partner.share,
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
