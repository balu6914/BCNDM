import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges } from '@angular/core';
import { Page } from 'app/common/interfaces/page.interface';
import { Subscription } from 'app/common/interfaces/subscription.interface';
import { Table, TableType } from 'app/shared/table/table';
import { StreamSection, StreamsType } from './section.streams';



@Component({
  selector: 'dpc-dashboard-main-section-streams',
  templateUrl: './dashboard.main.streams.component.html',
  styleUrls: ['./dashboard.main.streams.component.scss']
})

export class DashboardMainStreamsComponent implements OnInit, OnChanges {
  @Input() page: Page<Subscription>;
  @Output() tabChanged: EventEmitter<any> = new EventEmitter();
  @Output() pageChanged = new EventEmitter<number>();
  table: Table = new Table();
  section: StreamSection = new StreamSection();
  types = StreamsType;

  constructor() { }

  ngOnChanges(changes: SimpleChanges) {
    this.table.page = this.page;
  }

  ngOnInit() {
    this.table.tableType = TableType.Dashboard;
    this.table.headers = ['Stream Name', 'Price Paid', 'Start Date', 'End Date', 'URL'];
    this.table.page = new Page<Subscription>(0, 0, 0, []);
  }

  onPageChange(page: number) {
    this.pageChanged.emit(page);
  }

  onTabSwitch(t) {
    this.section.name = t;
    this.tabChanged.emit(this.section);
  }

}
