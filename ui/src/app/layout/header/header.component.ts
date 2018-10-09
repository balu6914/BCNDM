import { Component, OnInit, OnDestroy } from '@angular/core';
import { ISubscription } from 'rxjs/Subscription';

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

export class HeaderComponent implements OnInit, OnDestroy {
  isLoggedin: Boolean;
  user: any;
  balance = new Balance();
  private userSubscription: ISubscription;
  private loggedInSubscription: ISubscription;
  private balanceSubscription: ISubscription;
  private balanceChangeSubscription: ISubscription;
  constructor(
    private AuthService: AuthService,
    private balanceService: BalanceService,
  ) {
    // Listen to balance changed events
    this.balanceChangeSubscription = this.balanceService.balance.subscribe(result => {
      this.balance = result;
    },
  );

   }

  ngOnInit() {
    this.loggedInSubscription = this.AuthService.loggedIn
      .subscribe(is => {
        this.isLoggedin = is;
        if (this.isLoggedin) {
          this.userSubscription = this.AuthService.getCurrentUser()
          .subscribe(data => {
            this.user = data;
            // Get user balance
            this.getBalance();
          });
        }
      });
    }

    ngOnDestroy() {
      this.userSubscription.unsubscribe();
      this.loggedInSubscription.unsubscribe();
      this.balanceSubscription.unsubscribe();
      this.balanceChangeSubscription.unsubscribe();
    }

  getBalance() {
      this.balanceSubscription = this.balanceService.get().subscribe(
        (result: any) => {
          this.balance.amount = result.balance;
          //TODO remove this Mock of fiatAmount when we add this info on API side
          this.balance.fiatAmount = this.balance.amount;
          // Publish new balance data to balance message buss
          this.balanceService.changed(this.balance);
        },
        err => {
          console.error('Error fetching user balance ', err);
        });
  }

  logout() {
    this.AuthService.logout();
  }
}
