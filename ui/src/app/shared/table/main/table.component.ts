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

  rowDelete(id) {
    // Remove row
    const list: any = this.table.content;
    list.forEach( (row, i) => {
      if (row.id == id) {
        this.table.content.splice(i,1);
      }
    });
  }

  rowEdited(stream) {
    // Update row values
    const list: any = this.table.content;
    list.forEach( (row, i) => {
      if (row.id == stream.id) {
        this.table.content[i] = stream;
      }
    });
  }
}
