import { Component, Input, Output, EventEmitter, SimpleChanges, OnChanges } from '@angular/core';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { Table, TableType } from '../table';
import { Stream } from '../../../common/interfaces/stream.interface';

@Component({
  selector: 'dpc-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
  // Flip component animation
  animations: [
    trigger('flipState', [
      state('active', style({
        transform: 'rotateY(179deg)'
      })),
      state('inactive', style({
        transform: 'rotateY(0)'
      })),
      transition('active => inactive', animate('500ms ease-out')),
      transition('inactive => active', animate('500ms ease-in'))
    ])
  ]
})
export class TableComponent {
  types = TableType;
  flip = 'inactive';
  selectedRow: any;

  @Input() table: Table = new Table();
  @Output() deleteEvt: EventEmitter<any> = new EventEmitter();
  @Output() editEvt: EventEmitter<any> = new EventEmitter();
  @Output() pageChanged = new EventEmitter<number>();

  constructor() { }


  showRowDetails(row: any) {
    this.selectedRow = row;
    if (this.table.hasDetails && this.selectedRow) {
      this.flip = (this.flip === 'inactive') ? 'active' : 'inactive';
    }
  }

  rowDeleted(id) {
    // Remove row
    const list: any = this.table.page.content;
    list.forEach( (row, i) => {
      if (row.id == id) {
        // Remove row from table
        this.table.page.content.splice(i,1);
        // Emit event to DashboardSellComponent
        // If its delted from details page, go back to list
        if(this.selectedRow && row.id === this.selectedRow.id) {
          this.flip = 'inactive';
        }
        this.deleteEvt.emit(row.id);
      }
    });
  }

  rowEdited(stream) {
    // Update row values
    let rows: any = this.table.page.content;
    rows.forEach( (row, i) => {
      if (row.id == stream.id) {
        // Update row table
        this.table.page.content[i] = stream;
        // Emit event to DashboardSellComponent
        this.editEvt.emit(stream);
        this.selectedRow = stream;
      }
    });
  }

  onPageChange(page: number) {
    this.pageChanged.emit(page);
  }
}
