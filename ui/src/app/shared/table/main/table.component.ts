import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Table, TableType } from '../table';

@Component({
  selector: 'dpc-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss']
})
export class TableComponent implements OnInit {
  types = TableType;

  @Input() table: Table = new Table();
  @Output() deleteEvt: EventEmitter<any> = new EventEmitter();
  @Output() editEvt: EventEmitter<any> = new EventEmitter();
  constructor() { }

  ngOnInit() {
  }

  rowDeleted(id) {
    // Remove row
    const list: any = this.table.content;
    list.forEach( (row, i) => {
      if (row.id == id) {
        // Remove row from table
        this.table.content.splice(i,1);
        // Emit event to DashboardSellComponent
        this.deleteEvt.emit(row.id);
      }
    });
  }

  rowEdited(stream) {
    // Update row values
    let rows: any = this.table.content;
    rows.forEach( (row, i) => {
      if (row.id == stream.id) {
        // Update row table
        this.table.content[i] = stream;
        // Emit event to DashboardSellComponent
        this.editEvt.emit(stream);
      }
    });
  }
}
