import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../auth/services/auth.service';
import { BalanceService } from '../../shared/balance/balance.service';
import { Balance } from '../../common/interfaces/balance.interface';

@Component({
  selector: 'dpc-header-component',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
  providers: [
  ],
})

export class HeaderComponent implements OnInit {
  isLoggedin: Boolean;
  user: any;
  balance = new Balance();

  constructor(
    private AuthService: AuthService,
    private balanceService: BalanceService,
  ) {
    // Listen to balance changed events
    this.balanceService.balance.subscribe(result => {
      this.balance = result;
    },
  );

   }

  ngOnInit() {
    this.AuthService.loggedIn
      .subscribe(is => {
        this.isLoggedin = is;

        this.AuthService.getCurrentUser()
          .subscribe(data => {
            console.log("I have data", data)
              this.user = data;
              // Get user balance
              if (data && data.id) {
                this.getBalance();
              }
          });
      });
    }

  getBalance() {
      this.balanceService.get().subscribe(
        (result: any) => {
          this.balance.amount = result.balance;
          //TODO remove this Mock of fiatAmount when we add this info on API side
          this.balance.fiatAmount = this.balance.amount;
          // Publish new balance data to balance message buss
          this.balanceService.changed(this.balance);
        },
        err => {
          console.error("Error fetching user balance ", err)
        });
  }

  logout() {
    this.AuthService.logout();
  }
}
