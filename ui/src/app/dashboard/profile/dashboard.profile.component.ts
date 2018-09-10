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
    console.log("ROOOOOOOOOOOOO")
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
        console.log("Here is a user, yeeeee", this.user)
      },
      err => {
        console.log(err);
      });
  }
}
