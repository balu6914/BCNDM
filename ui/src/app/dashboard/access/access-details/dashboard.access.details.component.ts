import { Component, ViewChild, Input, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'dpc-dashboard-access-details',
  templateUrl: './dashboard.access.details.component.html',
  styleUrls: [ './dashboard.access.details.component.scss' ]
})
export class DashboardAccessDetailsComponent implements OnInit {

    @Input() selectedAccess: any;
    constructor(
        private router: Router
    ) { }

    ngOnInit() {
      console.log('Access details: ', this.selectedAccess);
    }

    calculateProfit() {
      console.log('Access share / price', (this.selectedAccess.share / this.selectedAccess.stream.price) * 10000);
    }
}
