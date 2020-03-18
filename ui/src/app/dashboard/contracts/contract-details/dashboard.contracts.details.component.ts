import { Component, OnInit, ViewChild, Input } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'dpc-dashboard-contracts-details',
  templateUrl: './dashboard.contracts.details.component.html',
  styleUrls: [ './dashboard.contracts.details.component.scss' ]
})
export class DashboardContractsDetailsComponent implements OnInit {

    @Input() selectedContract: any;
    constructor(
      private router: Router
    ) { }

    ngOnInit() {
      console.log('Contract details: ', this.selectedContract);
    }

    calculateProfit() {
      console.log('Contract share / price', (this.selectedContract.share / this.selectedContract.stream.price) * 10000);
    }
}
