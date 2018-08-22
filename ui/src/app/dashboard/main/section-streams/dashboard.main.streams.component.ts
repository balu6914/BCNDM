import { Component, ViewChild, Input, Output, EventEmitter } from '@angular/core';
import { SubscriptionService } from '../../../common/services/subscription.service';
import { StreamService } from '../../../common/services/stream.service';

import { TasPipe } from '../../../common/pipes/converter.pipe';
import { Subscription } from '../../../common/interfaces/subscription.interface';
import { Stream } from '../../../common/interfaces/stream.interface';

import { Table, TableType } from '../../../shared/table/table';
import { Query } from '../../../common/interfaces/query.interface';
import { Page } from '../../../common/interfaces/page.interface';
import { StreamSection } from './section.streams';
import { StreamsType } from './section.streams';

@Component({
  selector: 'dpc-dashboard-main-section-streams',
  templateUrl: './dashboard.main.streams.component.html',
  styleUrls: [ './dashboard.main.streams.component.scss' ]
})

export class DashboardMainStreamsComponent {
    @Input()  data: Subscription[];
    @Output() tabChanged: EventEmitter<any> = new EventEmitter();
    streams = [];
    temp = [];
    table: Table = new Table();
    section: StreamSection = new StreamSection();
    types = StreamsType;

    constructor(
        private subscriptionService: SubscriptionService,
        private streamService: StreamService,
        private tasPipe: TasPipe,
    ) {}

    ngOnInit() {
      this.table.tableType = TableType.Dashboard;
      this.table.headers = ["Stream Name", "Price Paid","Start Date", "End Date", "URL"];
      this.table.page.content = this.data;

      }

      onTabSwitch(t) {
        this.section.name = t;
        this.tabChanged.emit(this.section);
      }


}
