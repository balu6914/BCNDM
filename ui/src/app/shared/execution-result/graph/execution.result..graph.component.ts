import { Component, OnInit, OnChanges, Input } from '@angular/core';
import { Chart } from 'chart.js';
import { WalletBalanceStatisticPipe } from 'app/common/pipes/balance.income.pipe';
import * as moment from 'moment';
@Component({
  selector: 'dpc-execution-result-graph',
  templateUrl: './execution.result.graph.component.html',
  styleUrls: ['./execution.result.graph.component.scss']
})
export class ExecutionResultGraphComponent implements OnInit, OnChanges {
  data: any = {};
  preview: any = {};
  sequence: number[];
  public chart: Chart;
  @Input() resultData: string;
  constructor(
  ) { }

  ngOnChanges() {
    if (this.chart) {
      this.chart.destroy();
      this.createChart();
    }
  }

  ngOnInit() {
    this.data = JSON.parse(this.resultData);
    // We are showing just first 1% of execution result on graph
    // take 1% of results array length
    if (this.data.content.length) {
      const max = 10000;
      const n = Math.round(1 / 100 * this.data.content.length);
      // Check if 1% of result is bigger then max allowed size for preview
      const t = n <= max ? n : max;
      // Take first n elements of array result
      this.preview = this.data.content.slice(0, t);
      // Create new array
      this.sequence = Array.apply(null, {length: this.preview.length}).map(Number.call, Number);
      this.createChart();
    }
  }

  createChart() {
    const chartData = {
      type: 'line',
      data: {
        datasets: [{
          data: this.preview,
          backgroundColor: 'rgba(6, 210, 216, .1)',
          borderColor: 'rgba(6, 210, 216, 1)',
          borderWidth: 1,
          barThickness: 1
        }]
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
        },
        legend: {
          display: false
        },
        scales: {
          xAxes: [{
            display: true,
            labels: this.sequence,
            ticks: {
                fontColor: 'rgba(158,175,200, 1)',
                fontSize: 11
            }
          }],
          yAxes: [{
            display: true,
            gridLines: {
              color: 'rgba(223,233,247,1)',
              zeroLineColor: 'rgba(223,233,247,1)',
              borderDash: [15, 15],
              drawBorder: false
            },
            ticks: {
              fontColor: 'rgba(158,175,200, 1)',
              fontSize: 11,
            }
          }]
        },
        title: {
          display: true,
          text: `Preview of first ${this.preview.length} records, the 1% of your total result set`,
          position: 'top',
        }
      },
    };
    this.chart = new Chart('statisticChart', chartData);
  }

  downloadAsCSV() {
    const csvContent = this.data.content;
    const f = new Blob([csvContent], {
      type: 'text/csv'
    });
    const url = window.URL.createObjectURL(f);
    window.open(url);
  }
}
