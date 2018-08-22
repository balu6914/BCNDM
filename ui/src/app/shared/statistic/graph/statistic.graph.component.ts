import { Component, OnInit } from '@angular/core';
import { Chart } from 'chart.js';

@Component({
  selector: 'dpc-statistic-graph',
  templateUrl: './statistic.graph.component.html',
  styleUrls: ['./statistic.graph.component.scss']
})
export class StatisticGraphComponent {

  constructor(
  ) { }

  ngOnInit() {
    this.createMyChart();

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
