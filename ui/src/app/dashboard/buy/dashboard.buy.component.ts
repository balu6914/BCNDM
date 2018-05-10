import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { AuthService } from '../../auth/services/auth.service';
import { User } from '../../common/interfaces/user.interface';

@Component({
  selector: 'dashboard-buy',
  templateUrl: './dashboard.buy.component.html',
  styleUrls: [ './dashboard.buy.component.scss' ]
})
export class DashboardBuyComponent {
    user:any;
    constructor(private AuthService: AuthService, public http: HttpClient) {
    }

    ngOnInit() {
    }

}
