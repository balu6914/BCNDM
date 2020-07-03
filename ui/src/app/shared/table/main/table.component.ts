import { Component, Input, Output, EventEmitter, SimpleChanges, OnChanges, ViewChild } from '@angular/core';
import { trigger, state, style, transition, animate } from '@angular/animations';
import { Table, TableType } from 'app/shared/table/table';
import { Stream } from 'app/common/interfaces/stream.interface';
import { TableRowComponent } from 'app/shared/table/row/table.row.component';
import { Access } from 'app/common/interfaces/access.interface';
import { Execution } from 'app/common/interfaces/execution.interface';

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
  rowToMark: any;

  @Input() table: Table = new Table();
  @Output() deleteEvt: EventEmitter<any> = new EventEmitter();
  @Output() editEvt: EventEmitter<any> = new EventEmitter();
  @Output() pageChanged = new EventEmitter<number>();
  @Output() hoverRow: EventEmitter<any> = new EventEmitter();
  @Output() unhoverRow: EventEmitter<any> = new EventEmitter();
  @Output() contractSigned: EventEmitter<any> = new EventEmitter();
  @Output() checkboxChangedEvt: EventEmitter<any> = new EventEmitter();
  @Output() accessApproved: EventEmitter<any> = new EventEmitter();
  @Output() accessRevoked: EventEmitter<any> = new EventEmitter();
  @Output() fetchExecResult: EventEmitter<any> = new EventEmitter();

  constructor() { }


  showRowDetails(row: any) {
    this.selectedRow = row;
    if (this.table.hasDetails && this.selectedRow) {
      this.flip = (this.flip === 'inactive') ? 'active' : 'inactive';
    }
  }

  rowDeleted(rowDeleted: any) {
    // Remove row
    const list: any = this.table.page.content;
    list.forEach( (row, i) => {
      // TODO: Remove thiss check and emit full row for all componets
      if (row.id === rowDeleted.id) {
        if (rowDeleted.email !== undefined) {
          this.deleteEvt.emit(rowDeleted);
        } else {
          // Remove row from table
          this.table.page.content.splice(i, 1);
          this.deleteEvt.emit(row.id);
        }

        // Emit event to DashboardSellComponent
        // If its delted from details page, go back to list
        if (this.selectedRow && row.id === this.selectedRow.id) {
          this.flip = 'inactive';
        }
      }
    });
  }

  rowEdited(rowEdited: any) {
    // Update row values
    const rows: any = this.table.page.content;
    rows.forEach( (row, i) => {
      if (row.id === rowEdited.id) {
        // Update row table
        this.table.page.content[i] = rowEdited;
        // Emit event to DashboardSellComponent
        this.editEvt.emit(rowEdited);
        this.selectedRow = rowEdited;
      }
    });
  }

  onPageChange(page: number) {
    this.pageChanged.emit(page);
  }

  onHoveringRow(row) {
    this.hoverRow.emit(row);
  }

  onUnhoveringRow(row) {
    this.unhoverRow.emit(row);
  }

  activateRow(streamId) {
    this.rowToMark = streamId;
  }

  onContractSigned(contract: any) {
    // Update signed field
    const rows: any = this.table.page.content;
    rows.forEach( (row, i) => {
      if (row.stream_id === contract.stream_id) {
        contract.signed = true;
        this.table.page.content[i] = contract;
      }
    });
  }

  onCheckboxChangedEvt(row: Stream) {
    this.checkboxChangedEvt.emit(row);
  }

  onAccessApproved(row: Access) {
    this.accessApproved.emit(row);
  }

  onAccessRevoked(row: Access) {
    this.accessRevoked.emit(row);
  }

  onFetchExecResult(row: Execution) {
    this.fetchExecResult.emit(row);
  }
}
