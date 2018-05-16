import { Component } from '@angular/core';
import { AuthService } from '../../auth/services/auth.service';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'sidebar-component',
  templateUrl: './sidebar.component.html',
  styleUrls: [ './sidebar.component.scss' ],
  providers: [
  ],
})
export class SidebarComponent {
    isLoggedin: Boolean;
    loggedSubscription: any;
    toggledBar = false;
    subscription: any;

  constructor(private AuthService: AuthService) {

  }

  ngOnInit() {
      this.subscription = this.AuthService.loggedIn;
      this.subscription
      .subscribe(is => {
          this.isLoggedin = is;
      });
  }

  logout() {
      this.AuthService.logout();
      this.isLoggedin = false;
  }

  toggleSideBar() {
      this.toggledBar = !this.toggledBar;

  }
}
