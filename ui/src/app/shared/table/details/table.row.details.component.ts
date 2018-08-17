import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';

import { Stream } from '../../../common/interfaces/stream.interface';

@Component({
  selector: 'dpc-table-row-details',
  templateUrl: './table.row.details.component.html',
  styleUrls: ['./table.row.details.component.scss']
})

export class TableRowDetailsComponent implements OnInit {

  @Input() stream: Stream;
  @Output() backClicked = new EventEmitter<String>();

  constructor(
  ) {}

  ngOnInit() {
  }

  close() {
    this.backClicked.emit('trigger');
  }


}
