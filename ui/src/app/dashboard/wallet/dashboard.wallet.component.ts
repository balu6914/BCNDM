import { Component, OnInit } from '@angular/core';
import { AuthService } from 'app/auth/services/auth.service';
import { User } from 'app/common/interfaces/user.interface';
import { Page } from '../../common/interfaces/page.interface';
import { Subscription } from '../../common/interfaces';
import { StreamsType, StreamSection } from '../../shared/table/table-tabbed/section.streams';
import { SubscriptionService } from '../../common/services/subscription.service';
import { StreamService } from '../../common/services/stream.service';
import { TasPipe } from '../../common/pipes/converter.pipe';
import { forkJoin } from 'rxjs/observable/forkJoin';
import { TableType } from '../../shared/table/table';


@Component({
  selector: 'dpc-user-wallet-balance',
  templateUrl: './dashboard.wallet.component.html',
  styleUrls: ['./dashboard.wallet.component.scss']
})
export class DashboardWalletComponent implements OnInit {
  user: User;
  newBalance: number;
  subscription: any;
  // TODO: Remove this Mock of user balance its tmp hack for balance wallet widget
  mockBalance: any;
  sIncome: any[] = [];
  sOutcome: any[] = [];
  pageData: Page<Subscription>;
  activeStreamsSection: StreamSection = new StreamSection();
  page = 0;
  limit = 50;
  types = TableType;


  constructor(
    private authService: AuthService,
    private subscriptionService: SubscriptionService,
    private streamService: StreamService,
    private tasPipe: TasPipe,

  ) { }

  fetchSubscriptions() {
    if (this.activeStreamsSection.name === StreamsType.Bought) {
      this.subscriptionService.bought(this.page, this.limit).subscribe(
        (page: Page<Subscription>) => {
          page.content = page.content.map(sub => {
            sub.type = 'Outcome';
            return sub;
          });
          this.pageData = page;
          this.sOutcome = this.pageData.content;
        },
        err => console.log(err)
      );
    } else if (this.activeStreamsSection.name === StreamsType.Sold) {
      this.subscriptionService.owned(this.page, this.limit).subscribe(
        (page: Page<Subscription>) => {
          page.content.forEach(sub => {
            sub.type = 'Income';
            return sub;
          });
          this.pageData = page;
          this.sIncome = this.pageData.content;
        },
        err => console.log(err)
      );
    } else {
      const owned = this.subscriptionService.owned(this.page, this.limit / 2);
      const bought = this.subscriptionService.bought(this.page, this.limit / 2);
      forkJoin([owned, bought]).subscribe(results => {
        const result = results[0];
        const content = results[0].content.concat(results[1].content);
        content.forEach(sub => {
          sub.type = sub.stream_owner !== this.user.id ? 'Outcome' : 'Income';
          return sub;
        });
        result.content = content;
        result.total = results[0].total + results[1].total;
        result.limit = this.limit;
        this.pageData = result;
      });
    }
  }

  ngOnInit() {
    // Fetch current User
    this.authService.getCurrentUser().subscribe(
      data => {
        this.user = data;
      },
      err => {
        console.log(err);
      }
    );
    this.fetchSubscriptions();
  }

  onPageChanged(page: number) {
    this.page = page;
    this.fetchSubscriptions();
  }

  switchedTab(section) {
    this.activeStreamsSection = section;
    this.page = 0;
    this.fetchSubscriptions();
  }
}
