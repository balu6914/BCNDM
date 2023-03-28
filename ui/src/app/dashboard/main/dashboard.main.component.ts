import { Component, OnInit } from '@angular/core';
import { AuthService } from 'app/auth/services/auth.service';
import { Page } from 'app/common/interfaces/page.interface';
import { Subscription } from 'app/common/interfaces/subscription.interface';
import { DpcPipe } from 'app/common/pipes/converter.pipe';
import { StreamService } from 'app/common/services/stream.service';
import { SubscriptionService } from 'app/common/services/subscription.service';
import { StreamSection, StreamsType } from './section-streams/section.streams';

@Component({
  selector: 'dpc-dashboard-main',
  templateUrl: './dashboard.main.component.html',
  styleUrls: ['./dashboard.main.component.scss']
})
export class DashboardMainComponent implements OnInit {
  user: any;
  streams = [];
  sIncome: any[] = [];
  sOutcome: any[] = [];
  map: any;
  pageData: Page<Subscription>;
  streamTypes: StreamsType;
  activeStreamsSection: StreamSection = new StreamSection();
  page = 0;
  limit = 50;

  constructor(
    private authService: AuthService,
    private subscriptionService: SubscriptionService,
    private streamService: StreamService,
    private dpcPipe: DpcPipe,
  ) { }

fetchStreams(page: Page<Subscription>) {
    this.streams = [];
    page.content.forEach(sub => {
      this.streamService.getStream(sub.stream_id).subscribe(
        (stream: any) => {
          // Update map markers.
          this.streams = this.streams.concat(stream);
        },
        err => {
          console.log(err);
        }
      );
    });
  }

  fetchSubscriptions() {
    if (this.activeStreamsSection.name === StreamsType.Bought) {
      this.subscriptionService.bought(this.page, this.limit).subscribe(
        (page: Page<Subscription>) => {
          this.pageData = page;
          this.sOutcome = this.pageData.content;
          this.fetchStreams(page);
        },
        err => console.log(err)
      );
    } else {
      this.subscriptionService.owned(this.page, this.limit).subscribe(
        (page: Page<Subscription>) => {
          this.pageData = page;
          this.sIncome = this.pageData.content;
          this.fetchStreams(page);
        },
        err => console.log(err)
      );
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

  onHoverRow(row) {
    this.map.mouseHoverMarker(row);
  }

  onUnhoverRow(row) {
    this.map.mouseUnhoverMarker(row);
  }
}
