import { Component, OnInit, OnChanges, Input } from '@angular/core';
import { Chart } from 'chart.js';
import { WalletBalanceStatisticPipe } from 'app/common/pipes/balance.income.pipe';
import * as moment from 'moment';
@Component({
  selector: 'dpc-statistic-graph',
  templateUrl: './statistic.graph.component.html',
  styleUrls: ['./statistic.graph.component.scss']
})
export class StatisticGraphComponent implements OnInit, OnChanges {
  income: any[] = [];
  outcome: any[] = [];
  public chart: Chart;
  @Input() dataOutcome: any[];
  @Input() dataIncome: any[];

  constructor(
    private walletBalanceStatisticPipe: WalletBalanceStatisticPipe,
  ) { }

  ngOnChanges() {
    if (this.chart) {
      this.parseInputs();
      this.chart.destroy();
      this.createChart();
    }
  }

  ngOnInit() {
    this.parseInputs();
    this.createChart();
  }

  createChart() {
    const chartData = {
      type: 'line',
      data: {
        datasets: [{
          label: 'Outcome',
          data: this.outcome,
          backgroundColor: 'rgba(6, 210, 216, .1)',
          borderColor: 'rgba(6, 210, 216, 1)',
          borderWidth: 1,
          barThickness: 1,
        },
        {
          label: 'Income',
          data: this.income,
          backgroundColor: 'rgba(0, 125, 255, .1)',
          borderColor: '#007DFF',
          borderWidth: 4,
          pointBackgroundColor: '#ffffff',
          pointRadius: 3,
          pointBorderWidth: 1,
        }
        ]
      },
      options: {
        maintainAspectRatio: false,
        scaleShowVerticalLines: false,
        tooltips: {
          backgroundColor: '#007DFF',
          xPadding: 15,
          yPadding: 5,
          titleMarginBottom: 0,
          bodySpacing: 2,
          cornerRadius: 0,
          displayColors: false,
          caretSize: 0,
          callbacks: {
            label: function(tooltipItem, data) {
              return `${tooltipItem.yLabel} DPC on ${moment(tooltipItem.xLabel).format('MMM, D')}`;
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
              color: 'rgba(255,255,255,0)',
              zeroLineColor: 'rgba(255,255,255,0)'
            },
            ticks: {
              fontColor: 'rgba(158,175,200, 1)',
              fontSize: 11
            }
          }],
          yAxes: [{
            display: true,
            afterTickToLabelConversion: function(q) {
              q.ticks.forEach( t => {
                t = t + ' DPC ';
              });
            },
            gridLines: {
              color: 'rgba(223,233,247,1)',
              zeroLineColor: 'rgba(223,233,247,1)',
              borderDash: [15, 15],
              drawBorder: false
            },
            ticks: {
              fontColor: 'rgba(158,175,200, 1)',
              fontSize: 11,
              stepSize: this.calculteStepSize()
            }
          }]
        }
      }
    };
    this.chart = new Chart('statisticChart', chartData);
  }
  // Method parse inputed subscriptions array to chart.js data friendly objects array
  parseInputs() {
    this.outcome = this.walletBalanceStatisticPipe.transform(this.dataOutcome);
    this.income = this.walletBalanceStatisticPipe.transform(this.dataIncome);
  }
  // Method takes largest ammount and do the math to build a step size
  calculteStepSize() {
    const d = this.income.concat(this.outcome);
    return Math.max(...d.map(o => o.y)) / 5;
  }
}
