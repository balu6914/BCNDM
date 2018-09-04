import { Component, OnInit, Input } from '@angular/core';
import { Chart } from 'chart.js';
import { WalletBalanceStatisticPipe } from '../../../common/pipes/balance.income.pipe';
import * as moment from 'moment';
@Component({
  selector: 'dpc-statistic-graph',
  templateUrl: './statistic.graph.component.html',
  styleUrls: ['./statistic.graph.component.scss']
})
export class StatisticGraphComponent {
  income: any[] = [];
  balance: any[] = [];
  @Input() dataIncome: number[];
  @Input() dataBalance: number[];

  constructor(
    private walletBalanceStatisticPipe: WalletBalanceStatisticPipe,
  ) { }

  ngOnInit() {
    this.income = this.walletBalanceStatisticPipe.transform(this.dataIncome);
    this.createMyChart();
    console.log("what we have ? ", this.income);


  }

  createMyChart() {
    let chartData = {
      type: "line",
      data: {
        // labels: [
        //   "Jan 2017",
        //   "Apr 2017",
        //   "Sep 2017",
        //   "Dec 2017",
        //   "Mar 2018",
        //   "Jul 2018",
        //   "Oct 2018",
        //   "Feb 2019"
        // ],
        datasets: [{
          type: "bar",
          label: "Income",
          data: this.income,
          // data: [
          //   {
          //     x: this.newDate(12),
          //     y: 2.44
          //   },
          //   {
          //     x: this.newDate(13),
          //     y: 31
          //   },
          //   {
          //     x: this.newDate(2),
          //     y: 13.33
          //   }
          // ],
          backgroundColor: "rgba(6, 210, 216, 1)",
          borderColor: "rgba(6, 210, 216, 1)",
          borderWidth: 1,
          barThickness: 1
        },
        {
          label: "Wallet balance",
          data: [
            {
              x: this.newDate(0),
              y: 34.11
            },
            {
              x: this.newDate(1),
              y: 34
            },
            {
              x: this.newDate(228),
              y: 100
            }
          ],
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
            label: function(tooltipItem, data) {
              console.log("tooltipItem", tooltipItem)
              return `${tooltipItem.yLabel} TAS on ${moment(tooltipItem.xLabel).format('MMM, D')}`
            },
            title: function(tooltipItem, data) {
              return;
            }
          }
        },
        legend: {
          display: false
        },
        scales: {
          xAxes: [{
            type: 'time',
            barPercentage: 10,
            categoryPercentage: 0.1,
            barThickness: 5,
            distribution: 'linear',
            gridLines: {
              lineWidth: 0,
              color: "rgba(255,255,255,0)",
              zeroLineColor: "rgba(255,255,255,0)"
            },
            ticks: {
              fontColor: "rgba(158,175,200, 1)",
              fontSize: 11
            }
          }],
          yAxes: [{
            display: true,
            afterTickToLabelConversion: function(q) {
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
          }]
        }
      }
    };
    let c: any = document.getElementById("myChart");
    let ctx = c.getContext("2d");
    let chart = new Chart(ctx, chartData);
  }

  newDate(days) {
    console.log("here is how it looks like", moment().add(days, 'd').toDate())
    return moment().add(days, 'd').toDate();
  }
}
