import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges } from '@angular/core';
import { Page } from 'app/common/interfaces/page.interface';
import { Subscription } from 'app/common/interfaces/subscription.interface';
import { Table, TableType } from 'app/shared/table/table';
import { StreamSection, StreamsType } from './section.streams';
import { SubscriptionService } from 'app/common/services/subscription.service';



@Component({
  selector: 'dpc-dashboard-table-tabbed',
  templateUrl: './table-tabbed.component.html',
  styleUrls: ['./table-tabbed.component.scss']
})

export class TableTabbedComponent implements OnInit, OnChanges {
  @Input() page: Page<Subscription>;
  @Input() type: TableType;

  @Output() tabChanged: EventEmitter<any> = new EventEmitter();
  @Output() pageChanged = new EventEmitter<number>();
  table: Table = new Table();
  section: StreamSection = new StreamSection();
  types = StreamsType;

  constructor(
        private subscriptionService: SubscriptionService,
  ) { }

  ngOnChanges() {
    const table = Object.assign({}, this.table);
    table.page = this.page;
    this.table = table;
  }

  ngOnInit() {
    if (this.type === TableType.Dashboard) {
      this.table.tableType = TableType.Dashboard;
      this.table.headers = ['Stream Name', 'Price Paid', 'Start time', 'End time', 'URL'];
    } else {
      this.table.tableType = TableType.Wallet;
      this.table.headers = ['Stream Name', 'Type', 'Date and time', 'Value'];
    }
    this.table.page = new Page<Subscription>(0, 0, 0, []);
  }

  onPageChange(page: number) {
    this.pageChanged.emit(page);
  }

  onTabSwitch(t) {
    this.section.name = t;
    this.tabChanged.emit(this.section);
  }

  report() {
    this.subscriptionService.report(0, 100, this.section.name.toLowerCase()).subscribe(
      (data: any) => {
        console.log(data)
        const a = document.createElement('a')
        const objectUrl = URL.createObjectURL(data)
        a.href = objectUrl
        a.download = 'report.pdf';
        a.click();
      },
      err => console.log(err)
    );
  }
}
