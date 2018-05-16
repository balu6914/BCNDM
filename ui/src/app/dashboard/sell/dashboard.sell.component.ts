import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { AuthService } from '../../auth/services/auth.service';
import { User } from '../../common/interfaces/user.interface';

@Component({
  selector: 'dashboard-sell',
  templateUrl: './dashboard.sell.component.html',
  styleUrls: [ './dashboard.sell.component.scss' ]
})
export class DashboardSellComponent {
    user:any;
    constructor(private AuthService: AuthService, public http: HttpClient) {
    }

    ngOnInit() {
    }

}
