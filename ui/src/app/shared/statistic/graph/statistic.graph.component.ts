import { Component, OnInit, Input } from '@angular/core';
import { Chart } from 'chart.js';

@Component({
  selector: 'dpc-statistic-graph',
  templateUrl: './statistic.graph.component.html',
  styleUrls: ['./statistic.graph.component.scss']
})
export class StatisticGraphComponent {
  @Input() dataIncome: number[];
  @Input() dataBalance: number[];

  constructor(
  ) { }

  ngOnInit() {
    this.createMyChart();
    console.log("what we have ? ", this.dataIncome)

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
          label: "Income",
          data: this.dataIncome,
          backgroundColor: "rgba(6, 210, 216, 1)",
          borderColor: "rgba(6, 210, 216, 1)",
          borderWidth: 1,
          barThickness: 1
        },
        {
          label: "Wallet balance",
          data: this.dataBalance,
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
              return tooltipItem.yLabel + " TOK";
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
          yAxes: [{
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
          }],
          xAxes: [{
            type: 'time',
            barPercentage: 10,
            categoryPercentage: 0.1,
            barThickness: 5,
            distribution: 'linear',
             time: {
             unit: 'month',
             unitStepSize: 1,
             displayFormats: {
                'month': 'MMM'
               }
             },
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
