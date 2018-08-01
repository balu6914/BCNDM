import { Component, OnInit, Input } from '@angular/core';
import { Table, TableType } from '../table';

@Component({
  selector: 'dpc-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss']
})
export class TableComponent implements OnInit {
  types = TableType;

  @Input() table: Table = new Table();
  constructor() { }

  ngOnInit() {
  }
}
