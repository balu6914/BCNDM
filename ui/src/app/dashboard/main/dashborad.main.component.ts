import { Component, ViewChild, Input } from '@angular/core';
import { MdlDialogService } from '@angular-mdl/core';
import { DatatableComponent } from '@swimlane/ngx-datatable';

import { AuthService } from '../../auth/services/auth.service';
import { SearchService } from './services/search.service';
import { StreamService } from './services/stream.service';
import { SubscriptionService } from './services/subscription.service';

import { TasPipe } from '../../common/pipes/converter.pipe';
import { User } from '../../common/interfaces/user.interface';
import { Subscription } from '../../common/interfaces/subscription.interface';

import { Chart } from 'chart.js';
import { Table, TableType } from '../../shared/table';


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

    constructor(
        private AuthService: AuthService,
        private subscriptionService: SubscriptionService,
        private streamService: StreamService,
        private searchService: SearchService,
        private dialogService: MdlDialogService,
        private tasPipe: TasPipe,
    ) {}

    ngOnInit() {
      this.table.title = "Subscriptions";
      this.table.tableType = TableType.Dashboard;
      this.table.headers = ["Stream Name", "Price Paid","Start Date", "End Date", "URL"];

        // Fetch current User
        this.user = {};
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
                  const stream = result.Stream;

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
            // Set table content
            this.table.content = this.subscriptions;
          },
          err => {
            console.log(err);
          }
        );

          // TODO: Replace this coordinates with map bounds
          const southWestLng = -180;
          const southWestLat = -90;
          const northEastLng = 180;
          const northEastLat = 90;

          // Search streams on drawed region
          this.searchService.searchStreams(
            "geo", southWestLng, southWestLat, southWestLng, northEastLat,
            northEastLng, northEastLat, northEastLng, southWestLat).subscribe(
            (result: any) => {
              this.temp = [...result.Streams];
          },
          err => {
            console.log(err)
          });

          // Create mocked statistics graph
          this.createMyChart()
      }

      createMyChart() {
        let chartData = {
          type: "line",
          data: {
            labels: [
              "Jan 2017",
              "Apr 2017",
              "Sep 2017",
              "Dec 2017",
              "Mar 2018",
              "Jul 2018",
              "Oct 2018",
              "Feb 2019"
            ],
            datasets: [{
                type: "bar",
                label: "Dataset 1",
                data: [5, 10, 15, 7, 3, 10, 2, 45, 12, 3, 35, 2, 5],
                backgroundColor: "rgba(6, 210, 216, 1)",
                borderColor: "rgba(6, 210, 216, 1)",
                borderWidth: 1,
                barThickness: 1
              },
              {
                label: "Dataset 2",
                data: [25, 43, 38, 33, 52, 65, 62, 49],
                backgroundColor: "rgba(0, 125, 255, .1)",
                borderColor: "#007DFF",
                borderWidth: 4,
                pointBackgroundColor: "#ffffff",
                pointRadius: 3,
                pointBorderWidth: 1
              }
            ]
          },
          options: {
            maintainAspectRatio: false,
            scaleShowVerticalLines: false,
            tooltips: {
              backgroundColor: "#007DFF",
              xPadding: 15,
              yPadding: 5,
              titleMarginBottom: 0,
              bodySpacing: 2,
              cornerRadius: 0,
              displayColors: false,
              caretSize: 0,
              callbacks: {
                label: function (tooltipItem, data) {
                  return tooltipItem.yLabel + " TOK";
                },
                title: function (tooltipItem, data) {
                  return;
                }
              }
            },
            legend: {
              display: false
            },
            scales: {
              yAxes: [{
                afterTickToLabelConversion: function (q) {
                  for (var tick in q.ticks) {
                    let newLabel = q.ticks[tick] + " TOK ";
                    q.ticks[tick] = newLabel;
                  }
                },
                gridLines: {
                  color: "rgba(223,233,247,1)",
                  zeroLineColor: "rgba(223,233,247,1)",
                  borderDash: [15, 15],
                  drawBorder: false
                },
                ticks: {
                  fontColor: "rgba(158,175,200, 1)",
                  fontSize: 11,
                  stepSize: 25
                }
              }],
              xAxes: [{
                barPercentage: 10,
                categoryPercentage: 0.1,
                barThickness: 5,
                gridLines: {
                  lineWidth: 0,
                  color: "rgba(255,255,255,0)",
                  zeroLineColor: "rgba(255,255,255,0)"
                },
                ticks: {
                  fontColor: "rgba(158,175,200, 1)",
                  fontSize: 11
                }
              }]
            }
          }
        };

        let c: any = document.getElementById("myChart");
        let ctx = c.getContext("2d");
        let chart = new Chart(ctx, chartData);
      }
}
