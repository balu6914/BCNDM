import { Component, OnInit, Input } from '@angular/core';
import { Table } from '../table';

@Component({
  selector: 'dpc-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss']
})
export class TableComponent implements OnInit {

  @Input() table: Table = new Table();
  constructor() { }

  ngOnInit() {
  }
}
