import { Component, OnInit, OnDestroy } from '@angular/core';
import { ISubscription } from 'rxjs/Subscription';

import { AuthService } from 'app/auth/services/auth.service';
import { BalanceService } from 'app/shared/balance/balance.service';
import { Balance } from 'app/common/interfaces/balance.interface';
import { environment } from 'environments/environment';

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
  aiEnabled: Boolean;
  balance = new Balance();
  private userSubscription: ISubscription;
  private loggedInSubscription: ISubscription;
  private balanceSubscription: ISubscription;
  private balanceChangeSubscription: ISubscription;
  constructor(
    private authService: AuthService,
    private balanceService: BalanceService,
  ) {
    this.aiEnabled = environment.AI_ENABLED;
    console.log('IS AI Enabled ? ', this.aiEnabled);
    // Listen to balance changed events
    this.balanceChangeSubscription = this.balanceService.balance.subscribe(result => {
      this.balance = result;
    },
  );

   }

  ngOnInit() {
    this.loggedInSubscription = this.authService.loggedIn
      .subscribe(is => {
        this.isLoggedin = is;
        if (this.isLoggedin) {
          this.userSubscription = this.authService.getCurrentUser()
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
          // TODO remove this Mock of fiatAmount when we add this info on API side
          this.balance.fiatAmount = this.balance.amount;
          // Publish new balance data to balance message buss
          this.balanceService.changed(this.balance);
        },
        err => {
          console.error('Error fetching user balance ', err);
        });
  }

  logout() {
    this.authService.logout();
  }
}