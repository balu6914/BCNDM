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

  fetchStreams(page: Page<Subscription>) {
    this.streams = [];
    page.content.forEach(sub => {
      this.streamService.getStream(sub.stream_id).subscribe(
        (stream: any) => {
          // Create name and price field in the Subscription
          sub['stream_name'] = stream['name'];
          const mitasPrice = this.tasPipe.transform(stream['price']);
          sub['stream_price'] = mitasPrice;
          // Set markers on the map
          this.streams.push(stream);
        },
        err => {
          console.log(err);
        }
      );
    });
    this.pageData = page;
  }

  fetchSubscriptions() {
    if (this.activeStreamsSection.name === StreamsType.Bought) {
      this.subscriptionService.bought(this.page, this.limit).subscribe(
        (page: Page<Subscription>) => this.fetchStreams(page),
        err => console.log(err)
      );
    } else {
      this.subscriptionService.owned(this.page, this.limit).subscribe(
        (page: Page<Subscription>) => this.fetchStreams(page),
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
