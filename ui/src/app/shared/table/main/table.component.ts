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
  selectedRow = Stream;

  @Input() table: Table = new Table();

  @Output() pageChanged = new EventEmitter<number>();

  constructor() { }


  showRowDetails(row) {
    this.selectedRow = row;
    if (this.table.hasDetails && this.selectedRow) {
      this.flip = (this.flip === 'inactive') ? 'active' : 'inactive';
    }
  }

  onPageChange(page: number) {
    this.pageChanged.emit(page);
  }

}
