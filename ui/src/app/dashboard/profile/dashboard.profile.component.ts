import { Component, OnInit } from '@angular/core';

import { User } from '../../common/interfaces/user.interface';
import { AuthService } from '../../auth/services/auth.service';

@Component({
  selector: 'dpc-dashboard-profile',
  templateUrl: 'dashboard.profile.component.html',
})
export class DashboardProfileComponent implements OnInit {
  user: User;
  constructor(
    private authService: AuthService,
  ) {}

  ngOnInit() {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
      },
      err => {
        console.log(err);
      });
  }
}
