import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../auth/services/auth.service';
import { Page } from '../../common/interfaces/page.interface';
import { Subscription } from '../../common/interfaces/subscription.interface';
import { TasPipe } from '../../common/pipes/converter.pipe';
import { StreamService } from '../../common/services/stream.service';
import { SubscriptionService } from '../../common/services/subscription.service';
import { StreamSection, StreamsType } from './section-streams/section.streams';

@Component({
  selector: 'dpc-dashboard-main',
  templateUrl: './dashboard.main.component.html',
  styleUrls: ['./dashboard.main.component.scss']
})
export class DashboardMainComponent implements OnInit {
  user: any;
  streams = [];
  map: any;
  pageData: Page<Subscription>;
  streamTypes: StreamsType;
  activeStreamsSection: StreamSection = new StreamSection();
  page = 0;
  limit = 20;

  constructor(
    private AuthService: AuthService,
    private subscriptionService: SubscriptionService,
    private streamService: StreamService,
    private tasPipe: TasPipe,
  ) { }

  fetchSubscriptions() {
    if (this.activeStreamsSection.name === StreamsType.Bought) {
      this.subscriptionService.bought(this.page, this.limit).subscribe(
        (page: Page<Subscription>) => this.pageData = page,
        err => console.log(err)
      );
    } else {
      this.subscriptionService.owned(this.page, this.limit).subscribe(
        (page: Page<Subscription>) => this.pageData = page,
        err => console.log(err)
      );
    }
  }

  ngOnInit() {
    // Fetch current User
    this.AuthService.getCurrentUser().subscribe(
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
