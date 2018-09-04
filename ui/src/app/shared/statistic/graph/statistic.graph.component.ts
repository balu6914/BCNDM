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

  }

  createMyChart() {
    let chartData = {
      type: "line",
      data: {
        datasets: [{
          type: "bar",
          label: "Income",
          data: this.income,
          backgroundColor: "rgba(6, 210, 216, 1)",
          borderColor: "rgba(6, 210, 216, 1)",
          borderWidth: 1,
          barThickness: 1
        },
        {
          label: "Wallet balance",
          data: this.balance,
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
              stepSize: 2500
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
