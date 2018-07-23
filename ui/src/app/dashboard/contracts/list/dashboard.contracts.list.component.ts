import { Component, ViewChild } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { AuthService } from '../../../auth/services/auth.service';
import { Router } from '@angular/router';

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
        {'id' : 1, 'stream': {'name': 'WeIO pressure', 'price':'1'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10','signed':false, 'expiered': false},
        {'id' : 2, 'stream': {'name': 'WeIO temperature', 'price':'10'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10','signed':true, 'expiered': false },
        {'id' : 3, 'stream': {'name': 'WeIO humidity', 'price':'15'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expiered': false},
        {'id' : 4, 'stream': {'name': 'WeIO water', 'price':'5'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':'true', 'expiered': false},
        {'id' : 5, 'stream': {'name': 'WeIO radiation', 'price':'50'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expiered': false},
        {'id' : 5, 'stream': {'name': 'Spark', 'price':'30'}, 'creationDate':"2018-02-15T12:14:56.806Z", 'exparationDate': '2018-02-15T12:14:56.806Z', 'share':'10', 'signed':false, 'expiered': true},
    ];
    user: any;
    subscription: any;
    temp = [];
    selectedContract = [];
    // ngx-table custom messages
    tableMessages =  {
        emptyMessage: "You don't have any smart contracts yet..."
    }

    constructor(
        private AuthService: AuthService,
        private router: Router
    ) { }

    ngOnInit() {
        this.subscription = this.AuthService.getCurrentUser();
              this.subscription
              .subscribe(data => {
                  this.user = data;
              });
            this.temp = [...this.myContractsList];
    }
}
