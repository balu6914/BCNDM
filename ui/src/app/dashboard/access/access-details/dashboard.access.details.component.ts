import { Component, ViewChild, Input, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'dpc-dashboard-access-details',
  templateUrl: './dashboard.access.details.component.html',
  styleUrls: [ './dashboard.access.details.component.scss' ]
})
export class DashboardAccessDetailsComponent implements OnInit {

    @Input('selectedAccess') access: any;

    constructor(
        private router: Router
    ) { }

    ngOnInit() {
        console.log("here it is ", this.access);
    }

    calculateProfit() {
        console.log("here is a math", (this.access.share/this.access.stream.price)*10000)
    }

}
