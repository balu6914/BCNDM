import { Component, ViewChild, Input } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'dashboard-contracts-details',
  templateUrl: './dashboard.contracts.details.component.html',
  styleUrls: [ './dashboard.contracts.details.component.scss' ]
})
export class DashboardContractsDetailsComponent {

    @Input('selectedContract') contract: any;

    constructor(
        private router: Router
    ) { }

    ngOnInit() {
        console.log("here it is ", this.contract);
    }

    calculateProfit() {
        console.log("here is a math", (this.contract.share/this.contract.stream.price)*10000)

    }


}
