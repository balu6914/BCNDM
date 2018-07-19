import { Component, Input, OnInit } from '@angular/core';
import { ngCopy } from 'angular-6-clipboard';
import { Stream, Subscription } from '../../../common/interfaces';
import { TableType } from '../../table';

@Component({
  selector: 'dpc-table-row',
  templateUrl: './table.row.component.html',
  styleUrls: ['./table.row.component.scss']
})

export class TableRowComponent implements OnInit {
  types = TableType

  @Input() row: Stream | Subscription;
  @Input() rowType: TableType
  constructor(
  ) { }

  private isStream(row: Stream | Subscription): row is Stream {
    return (<Stream>row).url !== undefined;
  }

  ngOnInit() {
  }

  public copyToClipboard() {
    if (this.isStream(this.row)) {
      ngCopy(this.row.url, null)
    }
  }
}
