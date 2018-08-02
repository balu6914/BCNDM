import { Component, Input, Output, EventEmitter, SimpleChanges, OnChanges } from '@angular/core';
import { Table, TableType } from '../table';

@Component({
  selector: 'dpc-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss']
})
export class TableComponent implements OnChanges {
  types = TableType;

  public loading = false;

  @Input() table: Table = new Table();

  @Output() pageChanged = new EventEmitter<number>();

  constructor() { }

  onPageChange(page: number) {
    this.loading = true;
    this.pageChanged.emit(page);
  }

  ngOnChanges(changes: SimpleChanges) {
    this.loading = false;
  }

}
