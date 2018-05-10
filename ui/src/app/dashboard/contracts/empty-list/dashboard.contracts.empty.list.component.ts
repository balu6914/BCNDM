import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'dashboard-contracts-empty-list',
  templateUrl: './dashboard.contracts.empty.list.component.html',
  styleUrls: [ './dashboard.contracts.empty.list.component.scss' ]
})
export class DashboardContractsEmptyListComponent {


    constructor(
        private router: Router
    ) { }

    ngOnInit() {
    }


}
