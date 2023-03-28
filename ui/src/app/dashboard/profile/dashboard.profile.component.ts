import { Component, OnInit } from '@angular/core';

import { User } from 'app/common/interfaces/user.interface';
import { AuthService } from 'app/auth/services/auth.service';

@Component({
  selector: 'dpc-dashboard-profile',
  templateUrl: 'dashboard.profile.component.html',
  styleUrls: ['dashboard.profile.component.scss'],
})
export class DashboardProfileComponent implements OnInit {
  user: User;

  constructor(
    private authService: AuthService,
  ) {}

  ngOnInit() {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      resp => {
        this.user = resp;
      },
      err => {
        console.log(err);
      });
  }
}
