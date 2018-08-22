import { Component, ViewChild, Input } from '@angular/core';
import { MdlDialogService } from '@angular-mdl/core';

import { AuthService } from '../../auth/services/auth.service';
import { StreamService } from '../../common/services/stream.service';
import { SubscriptionService } from '../../common/services/subscription.service';

import { TasPipe } from '../../common/pipes/converter.pipe';
import { User } from '../../common/interfaces/user.interface';
import { Subscription } from '../../common/interfaces/subscription.interface';

import { Table, TableType } from '../../shared/table/table';
import { Query } from '../../common/interfaces/query.interface';
import { Page } from '../../common/interfaces/page.interface';
import { Stream } from '../../common/interfaces';
import { StreamsType } from './section-streams/section.streams';
import { StreamSection } from './section-streams/section.streams';

@Component({
  selector: 'dashboard-main',
  templateUrl: './dashboard.main.component.html',
  styleUrls: [ './dashboard.main.component.scss' ]
})
export class DashboardMainComponent {
    user:any;
    subscriptions = [];
    streams = [];
    temp = [];
    map: any;
    table: Table = new Table();
    streamTypes: StreamsType;
    activeStreamsSection: StreamSection = new StreamSection();

    constructor(
        private AuthService: AuthService,
        private subscriptionService: SubscriptionService,
        private streamService: StreamService,
        private tasPipe: TasPipe,
    ) {}

    ngOnInit() {
      console.log(this.activeStreamsSection)
      this.table.tableType = TableType.Dashboard;
      this.table.headers = ["Stream Name", "Price Paid","Start Date", "End Date", "URL"];

        // Fetch current User
        this.AuthService.getCurrentUser().subscribe(
            data =>  {
                this.user = data;
            },
            err => {
                console.log(err)
            }
        );

        // Fetch all subscriptions
        this.subscriptionService.get().subscribe(
          (result: any) => {
            this.temp = [...result.Subscriptions];
            result.Subscriptions.forEach(subscription => {
              this.streamService.getStream(subscription["id"]).subscribe(
                (result: any) => {
                  const stream = result;

                  // Create name and price field in susbcription
                  subscription["stream_name"] = stream["name"];
                  const mitasPrice = this.tasPipe.transform(stream["price"]);
                  subscription["stream_price"] = mitasPrice;

                  // Set markers on the map
                  this.streams.push(stream);
                },
                err => {
                  console.log(err);
                }
              );

              // Push marker to the markers list
              this.subscriptions.push(subscription);
            });
            result.content = this.subscriptions;
            // Set table content
            this.table.page = result;
          },
          err => {
            console.log(err);
          }
        );

          const query = new Query();

          // Search streams on drawed region
          this.streamService.searchStreams(query).subscribe(
            (result: Page<Stream>) => {
              this.temp = result.content;
          },
          err => {
            console.log(err)
          });
      }


      switchedTab(section) {
        this.activeStreamsSection = section;
        //TODO: Call API here to fetch subscriptions and cast data to
        // this.subscriptions and this.streams (representation on map are streams)
        // Remove loging after subscriptions data implementation
        console.log("Catched tab change event:", this.activeStreamsSection);
      }

      fetchSubscriptions() {

      }

}
